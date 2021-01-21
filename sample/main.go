// +build ignore

package main

import (
	"fmt"

	"github.com/hablullah/go-lafzi"
)

func main() {
	// Create and open database
	db, err := lafzi.OpenDatabase("sample-db.lafzi", lafzi.Bolt)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Add documents
	documents := []lafzi.Document{
		{ID: 1, ArabicText: "بِسْمِ اللَّهِ الرَّحْمَـٰنِ الرَّحِيمِ"},
		{ID: 2, ArabicText: "الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ"},
		{ID: 3, ArabicText: "الرَّحْمَـٰنِ الرَّحِيمِ"},
		{ID: 4, ArabicText: "مَالِكِ يَوْمِ الدِّينِ"},
		{ID: 5, ArabicText: "إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"},
		{ID: 6, ArabicText: "اهْدِنَا الصِّرَاطَ الْمُسْتَقِيمَ"},
		{ID: 7, ArabicText: "صِرَاطَ الَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ الْمَغْضُوبِ عَلَيْهِمْ وَلَا الضَّالِّينَ"},
	}

	err = db.AddDocuments(documents...)
	if err != nil {
		panic(err)
	}

	// Search for "rahman"
	docIDs, err := db.Search("rahman")
	if err != nil {
		panic(err)
	}

	fmt.Printf("'rahman' found in %v\n", docIDs)
}
