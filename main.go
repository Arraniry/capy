package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arraniry/capy/internal/generator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "capy",
	Short: "Capy is a CLI tool for generating Go projects with Clean Architecture",
	Long: `Capy adalah sebuah CLI tool yang membantu Anda membuat proyek Go 
dengan struktur Clean Architecture secara cepat dan mudah.
Anda dapat membuat proyek baru dan generate berbagai komponen seperti
controller, repository, dan usecase.`,
}

var newCmd = &cobra.Command{
	Use:   "new [nama-proyek] [database]",
	Short: "Membuat proyek Go baru dengan Clean Architecture",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		databaseType := args[1]
		fmt.Printf("Membuat proyek baru: %s dengan database: %s\n", projectName, databaseType)

		projectGen := generator.NewProjectGenerator(projectName)
		projectGen.SetDatabaseType(databaseType)
		if err := projectGen.Generate(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Proyek %s berhasil dibuat!\n", projectName)

		// Generate a default module after project creation
		moduleGen := generator.NewModuleGenerator("defaultModule")
		moduleGen.SetProjectPath(projectName)

		if err := moduleGen.Generate(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Modul default berhasil dibuat!\n")
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate [tipe] [nama]",
	Short: "Generate komponen (controller/repository/usecase)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		componentType := args[0]
		componentName := args[1]
		fmt.Printf("Generate %s: %s\n", componentType, componentName)

		compGen := generator.NewComponentGenerator(componentType, componentName)
		if err := compGen.Generate(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Komponen %s_%s.go berhasil dibuat!\n", componentName, componentType)
	},
}

var moduleCmd = &cobra.Command{
	Use:   "module [nama-modul]",
	Short: "Generate modul lengkap (model, controller, repository, dan usecase)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		fmt.Printf("Membuat modul baru: %s\n", moduleName)

		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		projectName := filepath.Base(wd)

		moduleGen := generator.NewModuleGenerator(moduleName)
		moduleGen.SetProjectPath(projectName)

		if err := moduleGen.Generate(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Modul %s berhasil dibuat!\n", moduleName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(moduleCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
