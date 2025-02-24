package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectGenerator bertanggung jawab untuk membuat struktur proyek baru
type ProjectGenerator struct {
	projectName  string
	basePath     string
	databaseType string
}

func (g *ProjectGenerator) SetDatabaseType(dbType string) {
	g.databaseType = dbType
}

// NewProjectGenerator membuat instance baru ProjectGenerator
func NewProjectGenerator(projectName string) *ProjectGenerator {
	return &ProjectGenerator{
		projectName: projectName,
		basePath:    projectName,
	}
}

// Generate membuat struktur folder dan file dasar untuk proyek baru
func (g *ProjectGenerator) Generate() error {
	// Create base directory
	if err := os.MkdirAll(g.basePath, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori proyek: %w", err)
	}

	// Create all required directories
	dirs := []string{
		"cmd",
		"internal/delivery/http",
		"internal/repository",
		"internal/usecase",
		"internal/entity",
		"pkg/database",
		"pkg/middleware",
	}

	for _, dir := range dirs {
		path := filepath.Join(g.basePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("gagal membuat direktori %s: %w", dir, err)
		}
	}

	// Generate go.mod
	if err := g.generateGoMod(); err != nil {
		return fmt.Errorf("gagal generate go.mod: %w", err)
	}

	// Generate main.go
	if err := g.generateMainFile(); err != nil {
		return fmt.Errorf("gagal generate main.go: %w", err)
	}

	// Generate database.go
	if err := g.generateDatabaseFile(); err != nil {
		return fmt.Errorf("gagal generate database.go: %w", err)
	}

	// Generate .env and .env.example
	if err := g.generateEnvFiles(); err != nil {
		return fmt.Errorf("gagal generate env files: %w", err)
	}

	// Generate Makefile
	if err := g.generateMakefile(); err != nil {
		return fmt.Errorf("gagal generate Makefile: %w", err)
	}

	// Generate .gitignore
	if err := g.generateGitignore(); err != nil {
		return fmt.Errorf("gagal generate .gitignore: %w", err)
	}

	// Generate README.md
	if err := g.generateReadme(); err != nil {
		return fmt.Errorf("gagal generate README.md: %w", err)
	}

	// Generate Dockerfile
	if err := g.generateDockerfile(); err != nil {
		return fmt.Errorf("gagal generate Dockerfile: %w", err)
	}

	return nil
}

func (g *ProjectGenerator) generateMainFile() error {
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"PROJECT_NAME/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Setup database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup router
	r := mux.NewRouter()
	
	// Setup middleware
	r.Use(loggingMiddleware)

	// Get port from env or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}`

	mainContent = strings.Replace(mainContent, "PROJECT_NAME", g.projectName, -1)
	return os.WriteFile(filepath.Join(g.basePath, "cmd", "main.go"), []byte(mainContent), 0644)
}

func (g *ProjectGenerator) generateDatabaseFile() error {
	var dbContent string

	if g.databaseType == "mysql" {
		dbContent = `package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config menyimpan konfigurasi database
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig membuat instance Config dari environment variables
func NewConfig() *Config {
	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

// Connect membuat koneksi ke database
func Connect() (*gorm.DB, error) {
	config := NewConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Setup connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// AutoMigrate menjalankan migrasi database untuk semua model
func AutoMigrate(db *gorm.DB) error {
	// Daftar model untuk auto-migrate akan ditambahkan saat generate modul
	return nil
}`
	} else {
		// Logika untuk database lain (misalnya PostgreSQL)
		dbContent = `package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config menyimpan konfigurasi database
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig membuat instance Config dari environment variables
func NewConfig() *Config {
	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

// Connect membuat koneksi ke database
func Connect() (*gorm.DB, error) {
	config := NewConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Setup connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// AutoMigrate menjalankan migrasi database untuk semua model
func AutoMigrate(db *gorm.DB) error {
	// Daftar model untuk auto-migrate akan ditambahkan saat generate modul
	return nil
}`
	}

	return os.WriteFile(filepath.Join(g.basePath, "pkg", "database", "db.go"), []byte(dbContent), 0644)
}

func (g *ProjectGenerator) generateEnvFiles() error {
	envContent := fmt.Sprintf(`# Application
APP_NAME=%s
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=%s
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0`, g.projectName, g.projectName)

	if err := os.WriteFile(filepath.Join(g.basePath, ".env"), []byte(envContent), 0644); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(g.basePath, ".env.example"), []byte(envContent), 0644)
}

func (g *ProjectGenerator) generateMakefile() error {
	makefileContent := fmt.Sprintf(`.PHONY: build run test clean deps lint mock dev migrate-create migrate-up migrate-down

# Build the application
build:
	go build -o bin/%s cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter
lint:
	go vet ./...
	golangci-lint run

# Generate mock files
mock:
	mockgen -source=internal/repository/repository.go -destination=internal/mocks/repository_mock.go
	mockgen -source=internal/usecase/usecase.go -destination=internal/mocks/usecase_mock.go

# Build and run in development mode
dev: build
	./bin/%s

# Create database migrations
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Run database migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down`, g.projectName, g.projectName)

	return os.WriteFile(filepath.Join(g.basePath, "Makefile"), []byte(makefileContent), 0644)
}

func (g *ProjectGenerator) generateGitignore() error {
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# Environment variables
.env
.env.local

# Logs
*.log

# OS specific
.DS_Store
Thumbs.db`

	return os.WriteFile(filepath.Join(g.basePath, ".gitignore"), []byte(gitignoreContent), 0644)
}

func (g *ProjectGenerator) generateReadme() error {
	readmeContent := fmt.Sprintf(`# %s

Proyek ini dibuat menggunakan [Capy](https://github.com/arraniry/capy) - Generator proyek Go dengan Clean Architecture.

## Struktur Proyek

`+"```"+`
.
├── cmd/                    # Entry points aplikasi
├── internal/               # Private application code
│   ├── delivery/          # Layer interface (HTTP handlers)
│   ├── repository/        # Layer data persistence
│   ├── usecase/           # Layer business logic
│   └── entity/            # Enterprise business rules
└── pkg/                   # Public libraries
    ├── database/          # Database utilities
    └── middleware/        # HTTP middleware
`+"```"+`

## Cara Menjalankan

1. Install dependencies:
`+"```"+`bash
go mod download
`+"```"+`

2. Setup environment variables:
`+"```"+`bash
cp .env.example .env
# Edit .env sesuai kebutuhan
`+"```"+`

3. Jalankan aplikasi:
`+"```"+`bash
make run
# atau
go run cmd/main.go
`+"```"+`

## Development

`+"```"+`bash
# Install dependencies
make deps

# Run tests
make test

# Run linter
make lint

# Build binary
make build
`+"```"+`

## API Endpoints

### Users
- `+"`GET /users`"+` - Get all users
- `+"`GET /users/{id}`"+` - Get user by ID
- `+"`POST /users`"+` - Create new user
- `+"`PUT /users/{id}`"+` - Update user
- `+"`DELETE /users/{id}`"+` - Delete user`, g.projectName)

	return os.WriteFile(filepath.Join(g.basePath, "README.md"), []byte(readmeContent), 0644)
}

func (g *ProjectGenerator) generateGoMod() error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)
`, g.projectName)
	return os.WriteFile(filepath.Join(g.basePath, "go.mod"), []byte(content), 0644)
}

func (g *ProjectGenerator) generateDockerfile() error {
	dockerfileContent := `# Gunakan image Go resmi sebagai base image
FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Salin go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Salin semua file ke dalam container
COPY . .

# Build aplikasi
RUN go build -o main ./cmd/main.go

# Gunakan image yang lebih kecil untuk menjalankan aplikasi
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Salin binary dari builder
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
`

	return os.WriteFile(filepath.Join(g.basePath, "Dockerfile"), []byte(dockerfileContent), 0644)
}
