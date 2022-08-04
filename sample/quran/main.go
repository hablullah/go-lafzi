package main

import (
	"fmt"

	"github.com/hablullah/go-lafzi"
)

func main() {
	// Open storage
	storage, err := lafzi.OpenStorage("quran.lafzi")
	checkError(err)

	// Prepare storage
	goto bench
	err = prepareStorage(storage)
	checkError(err)

	// Run benchmark
bench:
	err = runBenchmark(storage)
	checkError(err)
}

func prepareStorage(st *lafzi.Storage) error {
	// Prepare documents
	var docs []lafzi.Document
	for i, ayah := range listAyah {
		docs = append(docs, lafzi.Document{
			ID:     i + 1,
			Arabic: ayah,
		})
	}

	// Save documents to storage
	return st.AddDocuments(docs...)
}

func runBenchmark(st *lafzi.Storage) error {
	for _, sc := range scenarios {
		fmt.Printf("SCENARIO %q\n", sc.Name)

		// Prepare variables to score this scenario
		var nMatched float64
		nDocuments := float64(len(sc.Documents) * len(sc.Queries))

		// Process each query
		for _, query := range sc.Queries {
			// Search this query
			results, err := st.Search(query)
			if err != nil {
				return err
			}

			// Convert result to map of "surah:ayah"
			mapResult := map[string]struct{}{}
			if len(results) > 0 {

				for _, r := range results {
					if surah := getSurah(r.DocumentID); surah != nil {
						ayah := r.DocumentID - surah.Start + 1
						key := fmt.Sprintf("%d:%d", surah.ID, ayah)
						mapResult[key] = struct{}{}
					}
				}
			}

			// Count how many relevant documents that found or not found
			var undetected []string
			for _, doc := range sc.Documents {
				if _, exist := mapResult[doc]; exist {
					nMatched++
				} else {
					undetected = append(undetected, doc)
				}
			}

			// Print the undetected document
			if len(undetected) > 0 {
				fmt.Printf("\tQUERY %q NOT FOUND IN:\n", query)
				for _, u := range undetected {
					fmt.Printf("\t\t%s\n", u)
				}
			}
		}

		scScore := nMatched / nDocuments
		fmt.Printf("SCORE %q => %f\n\n", sc.Name, scScore)
	}

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
