package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func sortMapKeys(inputMap map[string]interface{}) map[string]interface{} {
	keys := make([]string, 0, len(inputMap))

	for key := range inputMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sortedMap := make(map[string]interface{})
	for _, key := range keys {
		value := inputMap[key]
		if subMap, isSubMap := value.(map[string]interface{}); isSubMap {
			sortedMap[key] = sortMapKeys(subMap)
		} else {
			sortedMap[key] = value
		}
	}

	return sortedMap
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <path_to_directory", os.Args[0])
		return
	}


	path := os.Args[1]
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}

	for _, file := range files {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".json") {
			continue
		}

		fullPath := filepath.Join(path, filename)

		contents, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", fullPath, err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(contents, &data); err != nil {
			log.Fatalf("Error unmarshalling JSON from %s: %v", fullPath, err)
		}

		sortedData := sortMapKeys(data)

		sortedJSON, err := json.MarshalIndent(sortedData, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling sorted JSON: %v", err)
		}

		if err := os.WriteFile(fullPath, sortedJSON, 0644); err != nil {
			log.Fatalf("Error writing sorted JSON to %s: %v", fullPath, err)
		}

		fmt.Println("Successfully sorted keys lexicographically:", filename)
	}
}
