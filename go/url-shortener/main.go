package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	
)

const (
	shortCodeLength = 6
	chars           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	dataFileName    = "urls.json"
)

// URLMap stores the mapping between short codes and long URLs
type URLMap struct {
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
}

// loadURLMaps loads URL mappings from a JSON file
func loadURLMaps(filePath string) ([]URLMap, error) {
	var urlMaps []URLMap
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return urlMaps, nil // File does not exist, return empty slice
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading data file: %w", err)
	}

	if len(data) == 0 {
		return urlMaps, nil // File is empty, return empty slice
	}

	err = json.Unmarshal(data, &urlMaps)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return urlMaps, nil
}

// saveURLMaps saves URL mappings to a JSON file
func saveURLMaps(filePath string, urlMaps []URLMap) error {
	data, err := json.MarshalIndent(urlMaps, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing data file: %w", err)
	}
	return nil
}

// generateShortCode generates a unique short code
func generateShortCode(existingCodes map[string]bool) (string, error) {
	for {
		shortCode := make([]byte, shortCodeLength)
		for i := 0; i < shortCodeLength; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
			if err != nil {
				return "", fmt.Errorf("error generating random number: %w", err)
			}
			shortCode[i] = chars[num.Int64()]
		}
		code := string(shortCode)
		if !existingCodes[code] {
			return code, nil
		}
	}
}

// shortenURL shortens a given long URL
func shortenURL(longURL string, dataFilePath string) (string, error) {
	// Validate URL
	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	urlMaps, err := loadURLMaps(dataFilePath)
	if err != nil {
		return "", err
	}

	existingCodes := make(map[string]bool)
	for _, um := range urlMaps {
		if um.LongURL == longURL {
			return um.ShortCode, nil // URL already shortened, return existing short code
		}
		existingCodes[um.ShortCode] = true
	}

	shortCode, err := generateShortCode(existingCodes)
	if err != nil {
		return "", err
	}

	newURLMap := URLMap{ShortCode: shortCode, LongURL: longURL}
	urlMaps = append(urlMaps, newURLMap)

	err = saveURLMaps(dataFilePath, urlMaps)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

// retrieveURL retrieves the long URL for a given short code
func retrieveURL(shortCode string, dataFilePath string) (string, error) {
	urlMaps, err := loadURLMaps(dataFilePath)
	if err != nil {
		return "", err
	}

	for _, um := range urlMaps {
		if um.ShortCode == shortCode {
			return um.LongURL, nil
		}
	}
	return "", fmt.Errorf("short code '%s' not found", shortCode)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  url-shortener shorten <long-url>")
		fmt.Println("  url-shortener retrieve <short-code>")
		os.Exit(1)
	}

	command := os.Args[1]
	dataFilePath := filepath.Join(os.Getenv("HOME"), ".url-shortener", dataFileName)

	// Ensure the data directory exists
	dataDir := filepath.Dir(dataFilePath)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			log.Fatalf("Error creating data directory: %v", err)
		}
	}

	switch command {
	case "shorten":
		if len(os.Args) < 3 {
			fmt.Println("Usage: url-shortener shorten <long-url>")
			os.Exit(1)
		}
		longURL := os.Args[2]
		shortCode, err := shortenURL(longURL, dataFilePath)
		if err != nil {
			log.Fatalf("Error shortening URL: %v", err)
		}
		fmt.Printf("Shortened URL: %s\n", shortCode)
	case "retrieve":
		if len(os.Args) < 3 {
			fmt.Println("Usage: url-shortener retrieve <short-code>")
			os.Exit(1)
		}
		shortCode := os.Args[2]
		longURL, err := retrieveURL(shortCode, dataFilePath)
		if err != nil {
			log.Fatalf("Error retrieving URL: %v", err)
		}
		fmt.Printf("Original URL: %s\n", longURL)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage:")
		fmt.Println("  url-shortener shorten <long-url>")
		fmt.Println("  url-shortener retrieve <short-code>")
		os.Exit(1)
	}
}