package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	pattern := flag.String("p", "", "Search pattern (string or regex)")
	path := flag.String("path", ".", "Path to file or directory to search")
	recursive := flag.Bool("r", false, "Search recursively in directories")
	caseSensitive := flag.Bool("c", false, "Perform case-sensitive search")
	regex := flag.Bool("regex", false, "Treat pattern as a regular expression")

	flag.Parse()

	if *pattern == "" {
		fmt.Println("Error: Search pattern cannot be empty.")
		flag.Usage()
		os.Exit(1)
	}

	searchConfig := SearchConfig{
		Pattern:       *pattern,
		CaseSensitive: *caseSensitive,
		IsRegex:       *regex,
	}

	err := searchFiles(*path, *recursive, searchConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// SearchConfig holds the configuration for the search operation
type SearchConfig struct {
	Pattern       string
	CaseSensitive bool
	IsRegex       bool
}

// searchFiles searches for the pattern in files starting from the given path
func searchFiles(rootPath string, recursive bool, config SearchConfig) error {
	info, err := os.Stat(rootPath)
	if err != nil {
		return fmt.Errorf("invalid path %q: %w", rootPath, err)
	}

	if info.IsDir() {
		return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != rootPath && !recursive {
				return filepath.SkipDir // Skip subdirectories if not recursive
			}
			if !info.IsDir() {
				return searchFile(path, config)
			}
			return nil
		})
	} else {
		return searchFile(rootPath, config)
	}
}

// searchFile searches for the pattern within a single file
func searchFile(filePath string, config SearchConfig) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file %q: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		match := false
		if config.IsRegex {
			var re *regexp.Regexp
			var compileErr error
			if config.CaseSensitive {
				re, compileErr = regexp.Compile(config.Pattern)
			} else {
				re, compileErr = regexp.Compile("(?i)" + config.Pattern)
			}
			if compileErr != nil {
				return fmt.Errorf("invalid regex pattern %q: %w", config.Pattern, compileErr)
			}
			match = re.MatchString(line)
		} else {
			if config.CaseSensitive {
				match = strings.Contains(line, config.Pattern)
			} else {
				match = strings.Contains(strings.ToLower(line), strings.ToLower(config.Pattern))
			}
		}

		if match {
			fmt.Printf("%s:%d: %s\n", filePath, lineNumber, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %q: %w", filePath, err)
	}

	return nil
}
