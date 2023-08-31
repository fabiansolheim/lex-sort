package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
		fmt.Println("Usage: i18n-sort <path_to_json_file>")
		return
	}

	path := os.Args[1]
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		panic(err)
	}

	sortMapKeys(data)

	sortedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path, sortedJSON, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully sorted your keys alphabetically:", path)
}