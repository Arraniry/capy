package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ModuleGenerator struct {
	moduleName  string
	projectPath string
}

func NewModuleGenerator(moduleName string) *ModuleGenerator {
	return &ModuleGenerator{
		moduleName:  moduleName,
		projectPath: "", // Kosongkan nilai awal
	}
}

// SetProjectPath mengatur path proyek
func (g *ModuleGenerator) SetProjectPath(projectPath string) {
	g.projectPath = projectPath
}

func (g *ModuleGenerator) Generate() error {
	// Generate model
	if err := g.generateModel(); err != nil {
		return fmt.Errorf("gagal generate model: %w", err)
	}

	// Generate controller
	if err := g.generateController(); err != nil {
		return fmt.Errorf("gagal generate controller: %w", err)
	}

	// Generate repository
	if err := g.generateRepository(); err != nil {
		return fmt.Errorf("gagal generate repository: %w", err)
	}

	// Generate usecase
	if err := g.generateUsecase(); err != nil {
		return fmt.Errorf("gagal generate usecase: %w", err)
	}

	return nil
}

func (g *ModuleGenerator) generateModel() error {
	template := `package entity

import (
	"time"
)

type {{.Name}} struct {
	ID        uint      ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
	Username  string    ` + "`json:\"username\" gorm:\"unique;not null\"`" + `
	Email     string    ` + "`json:\"email\" gorm:\"unique;not null\"`" + `
	Password  string    ` + "`json:\"password,omitempty\" gorm:\"not null\"`" + `
	FullName  string    ` + "`json:\"full_name\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// TODO: Tambahkan field sesuai kebutuhan
}
`
	return g.generateFile("internal/entity", g.moduleName+".go", template)
}

func (g *ModuleGenerator) generateController() error {
	template := fmt.Sprintf(`package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"%s/internal/entity"
	"github.com/gorilla/mux"
)

type {{.Name}}Handler struct {
	usecase {{.Name}}Usecase
}

type {{.Name}}Usecase interface {
	GetAll() ([]entity.{{.Name}}, error)
	GetByID(id uint) (*entity.{{.Name}}, error)
	Create({{.LowerName}} *entity.{{.Name}}) error
	Update({{.LowerName}} *entity.{{.Name}}) error
	Delete(id uint) error
}

func New{{.Name}}Handler(usecase {{.Name}}Usecase) *{{.Name}}Handler {
	return &{{.Name}}Handler{
		usecase: usecase,
	}
}

func (h *{{.Name}}Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/{{.LowerName}}s", h.GetAll).Methods("GET")
	r.HandleFunc("/{{.LowerName}}s/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/{{.LowerName}}s", h.Create).Methods("POST")
	r.HandleFunc("/{{.LowerName}}s/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/{{.LowerName}}s/{id}", h.Delete).Methods("DELETE")
}

func (h *{{.Name}}Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.usecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *{{.Name}}Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := h.usecase.GetByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *{{.Name}}Handler) Create(w http.ResponseWriter, r *http.Request) {
	var item entity.{{.Name}}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.Create(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *{{.Name}}Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item entity.{{.Name}}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item.ID = uint(id)
	if err := h.usecase.Update(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *{{.Name}}Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
`, g.projectPath)
	return g.generateFile("internal/delivery/http", g.moduleName+"_handler.go", template)
}

func (g *ModuleGenerator) generateRepository() error {
	template := fmt.Sprintf(`package repository

import (
	"%s/internal/entity"
	"gorm.io/gorm"
)

type {{.Name}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) *{{.Name}}Repository {
	return &{{.Name}}Repository{
		db: db,
	}
}

func (r *{{.Name}}Repository) GetAll() ([]entity.{{.Name}}, error) {
	var items []entity.{{.Name}}
	result := r.db.Find(&items)
	return items, result.Error
}

func (r *{{.Name}}Repository) GetByID(id uint) (*entity.{{.Name}}, error) {
	var item entity.{{.Name}}
	result := r.db.First(&item, id)
	return &item, result.Error
}

func (r *{{.Name}}Repository) Create({{.LowerName}} *entity.{{.Name}}) error {
	return r.db.Create({{.LowerName}}).Error
}

func (r *{{.Name}}Repository) Update({{.LowerName}} *entity.{{.Name}}) error {
	return r.db.Save({{.LowerName}}).Error
}

func (r *{{.Name}}Repository) Delete(id uint) error {
	return r.db.Delete(&entity.{{.Name}}{}, id).Error
}
`, g.projectPath)
	return g.generateFile("internal/repository", g.moduleName+"_repository.go", template)
}

func (g *ModuleGenerator) generateUsecase() error {
	template := fmt.Sprintf(`package usecase

import (
	"%s/internal/entity"
)

type {{.Name}}Usecase struct {
	repo {{.Name}}Repository
}

type {{.Name}}Repository interface {
	GetAll() ([]entity.{{.Name}}, error)
	GetByID(id uint) (*entity.{{.Name}}, error)
	Create({{.LowerName}} *entity.{{.Name}}) error
	Update({{.LowerName}} *entity.{{.Name}}) error
	Delete(id uint) error
}

func New{{.Name}}Usecase(repo {{.Name}}Repository) *{{.Name}}Usecase {
	return &{{.Name}}Usecase{
		repo: repo,
	}
}

func (u *{{.Name}}Usecase) GetAll() ([]entity.{{.Name}}, error) {
	return u.repo.GetAll()
}

func (u *{{.Name}}Usecase) GetByID(id uint) (*entity.{{.Name}}, error) {
	return u.repo.GetByID(id)
}

func (u *{{.Name}}Usecase) Create({{.LowerName}} *entity.{{.Name}}) error {
	// TODO: Add validation
	return u.repo.Create({{.LowerName}})
}

func (u *{{.Name}}Usecase) Update({{.LowerName}} *entity.{{.Name}}) error {
	// TODO: Add validation
	return u.repo.Update({{.LowerName}})
}

func (u *{{.Name}}Usecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
`, g.projectPath)
	return g.generateFile("internal/usecase", g.moduleName+"_usecase.go", template)
}

func (g *ModuleGenerator) generateFile(dir, filename, tmpl string) error {
	fullDir := filepath.Join(g.projectPath, dir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori %s: %w", fullDir, err)
	}

	t, err := template.New("module").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("gagal parse template: %w", err)
	}

	file, err := os.Create(filepath.Join(fullDir, filename))
	if err != nil {
		return fmt.Errorf("gagal membuat file: %s: %w", filepath.Join(fullDir, filename), err)
	}
	defer file.Close()

	data := struct {
		Name        string
		LowerName   string
		ProjectPath string
	}{
		Name:        strings.Title(g.moduleName),
		LowerName:   strings.ToLower(g.moduleName),
		ProjectPath: g.projectPath,
	}

	return t.Execute(file, data)
}
