package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: static-code-analyzer <path_to_go_file_or_directory>")
		return
	}

	targetPath := os.Args[1]

	info, err := os.Stat(targetPath)
	if err != nil {
		log.Fatalf("Error stating path %s: %v", targetPath, err)
	}

	fset := token.NewFileSet() // FileSet provides source position information

	if info.IsDir() {
		// Analyze all .go files in the directory
		err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".go" {
				fmt.Printf("\nAnalyzing file: %s\n", path)
				analyzeFile(fset, path)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking directory %s: %v", targetPath, err)
		}
	} else if filepath.Ext(targetPath) == ".go" {
		// Analyze a single .go file
		fmt.Printf("Analyzing file: %s\n", targetPath)
		analyzeFile(fset, targetPath)
	} else {
		log.Fatalf("Invalid target: %s. Must be a .go file or a directory.", targetPath)
	}
}

func analyzeFile(fset *token.FileSet, filePath string) {
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments) // Parse the file
	if err != nil {
		log.Printf("Error parsing file %s: %v", filePath, err)
		return
	}

	// --- Rule: Find empty functions ---
	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil && len(funcDecl.Body.List) == 0 {
				position := fset.Position(funcDecl.Pos())
				fmt.Printf("  [WARNING] Empty function: %s at %s:%d:%d\n", funcDecl.Name.Name, position.Filename, position.Line, position.Column)
			}
		}
		return true
	})
}
