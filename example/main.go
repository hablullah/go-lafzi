// +build ignore

package main

import (
	"fmt"

	"github.com/hablullah/go-lafzi"
)

func main() {
	// Open database
	db, err := createDatabase(lafzi.SQLite)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Try each search
	for _, transliteration := range searchTexts {
		ids, err := db.Search(transliteration)
		if err != nil {
			panic(err)
		}

		fmt.Println(transliteration, ids)
	}
}

func createDatabase(storageType lafzi.StorageType) (*lafzi.Database, error) {
	db, err := lafzi.OpenDatabase("sample-db.lafzi", storageType)
	if err != nil {
		return nil, err
	}

	dbEntries := make([]lafzi.DatabaseEntry, len(arabicTexts))
	for i, text := range arabicTexts {
		dbEntries[i] = lafzi.DatabaseEntry{
			ID:         int64(i + 1),
			ArabicText: text,
		}
	}

	err = db.AddEntries(dbEntries...)
	if err != nil {
		return nil, err
	}

	return db, nil
}

var arabicTexts = []string{
	"بِسْمِ اللَّهِ الرَّحْمَـٰنِ الرَّحِيمِ",
	"الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ",
	"الرَّحْمَـٰنِ الرَّحِيمِ",
	"مَالِكِ يَوْمِ الدِّينِ",
	"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ",
	"اهْدِنَا الصِّرَاطَ الْمُسْتَقِيمَ",
	"صِرَاطَ الَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ الْمَغْضُوبِ عَلَيْهِمْ وَلَا الضَّالِّينَ",
	"الم",
	"ذَٰلِكَ الْكِتَابُ لَا رَيْبَ ۛ فِيهِ ۛ هُدًى لِّلْمُتَّقِينَ",
	"الَّذِينَ يُؤْمِنُونَ بِالْغَيْبِ وَيُقِيمُونَ الصَّلَاةَ وَمِمَّا رَزَقْنَاهُمْ يُنفِقُونَ",
	"وَالَّذِينَ يُؤْمِنُونَ بِمَا أُنزِلَ إِلَيْكَ وَمَا أُنزِلَ مِن قَبْلِكَ وَبِالْآخِرَةِ هُمْ يُوقِنُونَ",
	"أُولَـٰئِكَ عَلَىٰ هُدًى مِّن رَّبِّهِمْ ۖ وَأُولَـٰئِكَ هُمُ الْمُفْلِحُونَ",
	"إِنَّ الَّذِينَ كَفَرُوا سَوَاءٌ عَلَيْهِمْ أَأَنذَرْتَهُمْ أَمْ لَمْ تُنذِرْهُمْ لَا يُؤْمِنُونَ",
	"خَتَمَ اللَّهُ عَلَىٰ قُلُوبِهِمْ وَعَلَىٰ سَمْعِهِمْ ۖ وَعَلَىٰ أَبْصَارِهِمْ غِشَاوَةٌ ۖ وَلَهُمْ عَذَابٌ عَظِيمٌ",
	"وَمِنَ النَّاسِ مَن يَقُولُ آمَنَّا بِاللَّهِ وَبِالْيَوْمِ الْآخِرِ وَمَا هُم بِمُؤْمِنِينَ",
	"يُخَادِعُونَ اللَّهَ وَالَّذِينَ آمَنُوا وَمَا يَخْدَعُونَ إِلَّا أَنفُسَهُمْ وَمَا يَشْعُرُونَ",
	"فِي قُلُوبِهِم مَّرَضٌ فَزَادَهُمُ اللَّهُ مَرَضًا ۖ وَلَهُمْ عَذَابٌ أَلِيمٌ بِمَا كَانُوا يَكْذِبُونَ",
}

var searchTexts = []string{
	// Text exist in source
	"bismi",
	"alhamdulillah",
	"rabb",
	"rahman",
	"rahim",
	"malik",
	"yaumiddin",
	"iyyaka",
	"dollin",
	"kitab",
	"muttakin",
	"muflihun",
	"aandzartahum",
	"minan nas",
	"yaumil akhir",
	"qulubihim",
	"marodun",
	"yakdzibun",

	// Entire sentence
	"bismillahirrahmanirrahim",

	// Wrong transliteration
	"basma",
	"dooollliiin",
	"yaumal akhir",

	// Not exist in source
	"waylul",
	"istigfar",
}
