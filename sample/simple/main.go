package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hablullah/go-lafzi"
)

var arabicTexts = []string{
	"بِسْمِ اللَّهِ الرَّحْمَـٰنِ الرَّحِيمِ",
	"الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ",
	"الرَّحْمَـٰنِ الرَّحِيمِ",
	"مَالِكِ يَوْمِ الدِّينِ",
	"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ",
	"اهْدِنَا الصِّرَاطَ الْمُسْتَقِيمَ",
	"صِرَاطَ الَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ الْمَغْضُوبِ عَلَيْهِمْ وَلَا الضَّالِّينَ",
}

func main() {
	// Open storage
	os.RemoveAll("sample.lafzi")
	storage, err := lafzi.OpenStorage("sample.lafzi")
	checkError(err)

	// Prepare documents
	var docs []lafzi.Document
	for i, arabicText := range arabicTexts {
		docs = append(docs, lafzi.Document{
			Identifier: fmt.Sprintf("%d", i+1),
			Arabic:     arabicText},
		)
	}

	// Save documents to storage
	err = storage.AddDocuments(docs...)
	checkError(err)

	// Search in storage
	results, err := storage.Search("rahman")
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
