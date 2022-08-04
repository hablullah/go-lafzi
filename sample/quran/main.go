package main

import (
	"fmt"
	"time"

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
	// Prepare variables to score all scenario
	var nDoc, nQuery int
	var nTruePos, nFalseNeg int
	start := time.Now()

	for _, sc := range scenarios {
		fmt.Printf("SCENARIO %q\n", sc.Name)

		// Prepare variables to score this scenario
		var nsTruePos, nsFalseNeg int
		nsQuery := len(sc.Queries)
		nsDoc := nsQuery * len(sc.Documents)

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
					nsTruePos++
				} else {
					undetected = append(undetected, doc)
					nsFalseNeg++
				}
			}

			// TODO: Print the undetected document
			// if len(undetected) > 0 {
			// 	fmt.Printf("\tQUERY %q NOT FOUND IN:\n", query)
			// 	for _, u := range undetected {
			// 		fmt.Printf("\t\t%s\n", u)
			// 	}
			// }
		}

		// Print scenario info
		fmt.Printf("\tN QUERY    : %d\n", nsQuery)
		fmt.Printf("\tN EXPECTED : %d\n", nsDoc)
		fmt.Printf("\tN TRUE POS : %d\n", nsTruePos)
		fmt.Printf("\tN FALSE NEG: %d\n", nsFalseNeg)
		fmt.Printf("\tRECALL     : %f\n", float64(nsTruePos)/float64(nsDoc))
		fmt.Println()

		nDoc += nsDoc
		nQuery += nsQuery
		nTruePos += nsTruePos
		nFalseNeg += nsFalseNeg
	}

	duration := time.Since(start).Seconds()
	speed := 1 / (float64(nQuery) / (duration * 1000))

	fmt.Printf("N QUERY    : %d\n", nQuery)
	fmt.Printf("N EXPECTED : %d\n", nDoc)
	fmt.Printf("N TRUE POS : %d\n", nTruePos)
	fmt.Printf("N FALSE NEG: %d\n", nFalseNeg)
	fmt.Printf("RECALL     : %f\n", float64(nTruePos)/float64(nDoc))
	fmt.Printf("DURATION   : %f s (%f ms/query)\n", duration, speed)

	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
