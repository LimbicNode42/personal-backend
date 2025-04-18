package datasets

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"log"
)

// AppendToJSONL is the resolver for the appendToJSONL field.
func (r *mutationResolver) AppendToJSONL(ctx context.Context, fileName string, record map[string]any) (bool, error) {
	file := strings.Join([]string{"datasets", fmt.Sprintf("%s", fileName)}, "/")

	// Open the file for appending (create if not exists)
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Convert record to JSON
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return false, err
	}

	// Append as a new line
	if _, err := f.WriteString(string(jsonBytes) + "\n"); err != nil {
		return false, err
	}

	return true, nil
}

// CreateDataset is the resolver for the createDataset field.
func (r *mutationResolver) CreateDataset(ctx context.Context, fileName string) (bool, error) {
	file := strings.Join([]string{"datasets", fmt.Sprintf("%s", fileName)}, "/")

	// Ensure the filename ends with `.jsonl` (optional, but helpful)
	if !strings.HasSuffix(file, ".jsonl") {
		file += ".jsonl"
	}

	// Create or open the file (no overwrite)
	stream, err := os.OpenFile(file, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		// If file exists, return false but not fatal
		if os.IsExist(err) {
			return false, nil // Already exists
		}
		return false, err // Other errors
	}
	defer stream.Close()

	return true, nil // File created
}

// ReadJSONL is the resolver for the readJSONL field.
func (r *queryResolver) ReadJSONL(ctx context.Context, fileName string) ([]map[string]any, error) {
	file := strings.Join([]string{"datasets", fmt.Sprintf("%s", fileName)}, "/")

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []map[string]interface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue // or return nil, err if you want strict validation
		}
		results = append(results, obj)
	}

	return results, scanner.Err()
}

// GetFiles is the resolver for the getFiles field.
func (r *queryResolver) GetFiles(ctx context.Context) ([]string, error) {
	// Define the directory path you want to scan
	dirPath := "./datasets" // Change this to your actual dataset directory

	// Read the directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() { // Only include files (not sub-directories)
			fileNames = append(fileNames, entry.Name())
			log.Println(entry.Name())
		} else {
			log.Println("no files found")
		}
	}

	return fileNames, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
