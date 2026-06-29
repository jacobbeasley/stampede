package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	schemaPath := filepath.Join("migrations", "schema.sql")

	// Read the schema.sql file
	content, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Warning: schema.sql not found at %s, skipping cleanup", schemaPath)
			return
		}
		log.Fatalf("Error reading schema.sql: %v", err)
	}

	// Split by newline and filter out \restrict and \unrestrict lines
	lines := bytes.Split(content, []byte("\n"))
	var filtered [][]byte
	modified := false

	for _, line := range lines {
		// Check if line contains \restrict or \unrestrict
		if bytes.Contains(line, []byte(`\restrict`)) || bytes.Contains(line, []byte(`\unrestrict`)) {
			modified = true
			continue
		}
		filtered = append(filtered, line)
	}

	if modified {
		newContent := bytes.Join(filtered, []byte("\n"))
		err = ioutil.WriteFile(schemaPath, newContent, 0644)
		if err != nil {
			log.Fatalf("Error writing cleaned schema.sql: %v", err)
		}
		log.Println("Successfully cleaned migrations/schema.sql (removed \\restrict / \\unrestrict lines).")
	} else {
		log.Println("migrations/schema.sql was already clean.")
	}
}
