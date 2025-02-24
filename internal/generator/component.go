package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ComponentGenerator bertanggung jawab untuk generate komponen
type ComponentGenerator struct {
	componentType string
	componentName string
}

// NewComponentGenerator membuat instance baru ComponentGenerator
func NewComponentGenerator(componentType, componentName string) *ComponentGenerator {
	return &ComponentGenerator{
		componentType: componentType,
		componentName: componentName,
	}
}

// Generate membuat file komponen baru
func (g *ComponentGenerator) Generate() error {
	switch strings.ToLower(g.componentType) {
	case "controller":
		return g.generateController()
	case "repository":
		return g.generateRepository()
	case "usecase":
		return g.generateUsecase()
	default:
		return fmt.Errorf("tipe komponen tidak valid: %s", g.componentType)
	}
}

func (g *ComponentGenerator) generateController() error {
	template := `package http

import (
	"net/http"
)

type {{.Name}}Handler struct {
	usecase {{.Name}}Usecase
}

type {{.Name}}Usecase interface {
	// TODO: Define usecase methods
}

func New{{.Name}}Handler(usecase {{.Name}}Usecase) *{{.Name}}Handler {
	return &{{.Name}}Handler{
		usecase: usecase,
	}
}

func (h *{{.Name}}Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement handler
}
`
	data := struct {
		Name string
	}{
		Name: strings.Title(g.componentName),
	}

	return g.generateFile("internal/delivery/http", template, data)
}

func (g *ComponentGenerator) generateRepository() error {
	template := `package repository

type {{.Name}}Repository struct {
	db interface{} // TODO: Replace with your database client
}

func New{{.Name}}Repository(db interface{}) *{{.Name}}Repository {
	return &{{.Name}}Repository{
		db: db,
	}
}

// TODO: Implement repository methods
`
	data := struct {
		Name string
	}{
		Name: strings.Title(g.componentName),
	}

	return g.generateFile("internal/repository", template, data)
}

func (g *ComponentGenerator) generateUsecase() error {
	template := `package usecase

type {{.Name}}Usecase struct {
	repo {{.Name}}Repository
}

type {{.Name}}Repository interface {
	// TODO: Define repository methods
}

func New{{.Name}}Usecase(repo {{.Name}}Repository) *{{.Name}}Usecase {
	return &{{.Name}}Usecase{
		repo: repo,
	}
}

// TODO: Implement usecase methods
`
	data := struct {
		Name string
	}{
		Name: strings.Title(g.componentName),
	}

	return g.generateFile("internal/usecase", template, data)
}

func (g *ComponentGenerator) generateFile(dir, tmpl string, data interface{}) error {
	t, err := template.New("component").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("gagal parse template: %w", err)
	}

	filename := fmt.Sprintf("%s_%s.go", strings.ToLower(g.componentName), strings.ToLower(g.componentType))
	path := filepath.Join(dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %w", err)
	}
	defer file.Close()

	return t.Execute(file, data)
}
