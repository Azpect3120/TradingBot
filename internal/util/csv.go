package util

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// parseCSV reads a CSV file and returns the data in a slice of
// slices of strings. The first slice is the rows of the CSV file.
// The first column (header) is not included in the data. The first
// row is the first row of the CSV file. The caller is expected to
// know the structure of the CSV file and how to parse the data.
func parseCSV(path string) ([][]string, error) {
	// Convert path to absolute path for security
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Open the file
	file, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Store the data
	var data [][]string

	// Read file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, strings.Split(scanner.Text(), ","))
	}

	return data[1:], nil
}

// GetNamesFromCSV reads a CSV file and returns the names in the
// first column. The first row is the header and is not included
// in the data. The names are returned in a slice of strings.
// The caller should have a CSV file that is properly formatted.
func GetNamesFromCSV(path string) ([]string, error) {
	// Convert path to absolute path for security
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Open the file
	file, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Store the names
	var names []string

	// Read file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, strings.Split(scanner.Text(), ",")[0])
	}

	return names[1:], nil
}
