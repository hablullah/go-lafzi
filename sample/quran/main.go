package main

import (
	"encoding/json"
	"fmt"

	"github.com/hablullah/go-lafzi"
)

func main() {
	// Open storage
	storage, err := lafzi.OpenStorage("quran.lafzi")
	checkError(err)

	// Prepare documents
	var docs []lafzi.Document
	for i, ayah := range listAyah {
		docs = append(docs, lafzi.Document{
			ID:     i + 1,
			Arabic: ayah,
		})
	}

	// Save documents to storage
	err = storage.AddDocuments(docs...)
	checkError(err)

	// Search in storage
	results, err := storage.Search("ulul albaab")
	checkError(err)

	// Print search result
	bt, _ := json.MarshalIndent(&results, "", "  ")
	fmt.Println(string(bt))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
