package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/hablullah/go-lafzi"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Open storage
	os.RemoveAll("quran.lafzi")
	storage, err := lafzi.OpenStorage("quran.lafzi")
	checkError(err)

	// Prepare storage
	err = prepareStorage(storage)
	checkError(err)
	goto bench

	// Run benchmark
bench:
	err = runBenchmark(storage)
	checkError(err)
}

func prepareStorage(st *lafzi.Storage) error {
	start := time.Now()
	fmt.Println("START INDEXING")

	// Prepare documents
	docs := make([]lafzi.Document, len(listAyah))
	for i, ayah := range listAyah {
		id := i + 1
		var identifier string
		if surah := getSurah(id); surah != nil {
			ayah := id - surah.Start + 1
			identifier = fmt.Sprintf("%d:%d", surah.ID, ayah)
		}

		docs[i] = lafzi.Document{
			Identifier: identifier,
			Arabic:     ayah,
		}
	}

	// Save documents to storage
	err := st.AddDocuments(docs...)
	if err != nil {
		return err
	}

	duration := time.Since(start).Seconds()
	fmt.Printf("INDEXING FINISHED IN %f s\n\n", duration)
	return nil
}

func runBenchmark(st *lafzi.Storage) error {
	// Prepare variables to score all scenario
	var nDoc, nQuery int
	var nTruePos, nFalseNeg int
	start := time.Now()

	runScenario := func(sc Scenario) error {
		start := time.Now()

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
					mapResult[r.Identifier] = struct{}{}
				}
			}

			// Count how many relevant documents that found or not found
			for _, doc := range sc.Documents {
				if _, exist := mapResult[doc]; exist {
					nsTruePos++
				} else {
					nsFalseNeg++
				}
			}
		}

		// Print scenario info
		fmt.Printf("SCENARIO %q\n", sc.Name)
		fmt.Printf("\tN QUERY    : %d\n", nsQuery)
		fmt.Printf("\tN EXPECTED : %d\n", nsDoc)
		fmt.Printf("\tN TRUE POS : %d\n", nsTruePos)
		fmt.Printf("\tN FALSE NEG: %d\n", nsFalseNeg)
		fmt.Printf("\tRECALL     : %f\n", float64(nsTruePos)/float64(nsDoc))
		fmt.Printf("\tDURATION   : %f s\n", time.Since(start).Seconds())
		fmt.Println()

		nDoc += nsDoc
		nQuery += nsQuery
		nTruePos += nsTruePos
		nFalseNeg += nsFalseNeg
		return nil
	}

	var wg sync.WaitGroup
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(runtime.GOMAXPROCS(0))

	wg.Add(len(scenarios))
	for _, sc := range scenarios {
		g.Go(func() error {
			defer wg.Done()
			return runScenario(sc)
		})
	}

	wg.Wait()

	if err := g.Wait(); err != nil {
		return err
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
