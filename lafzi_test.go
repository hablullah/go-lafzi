package lafzi_test

import (
	"os"
	"testing"

	"github.com/hablullah/go-lafzi"
	"github.com/stretchr/testify/assert"
)

func TestSearchBolt(t *testing.T) {
	testSearch(t, lafzi.Bolt)
}

func TestSearchSQLite(t *testing.T) {
	testSearch(t, lafzi.SQLite)
}

func testSearch(t *testing.T, storageType lafzi.StorageType) {
	// Prepare database name
	dbName := "test-bolt.lafzi"
	if storageType == lafzi.SQLite {
		dbName = "test-sqlite.lafzi"
	}

	// Open database
	db, err := createDatabase(dbName, storageType)
	if err != nil {
		panic(err)
	}

	defer func() {
		db.Close()
		os.Remove(dbName)
	}()

	// Test each search
	for transliteration, expectedIDs := range searchTexts {
		ids, err := db.Search(transliteration)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, expectedIDs, ids, transliteration)
	}
}

func createDatabase(name string, storageType lafzi.StorageType) (*lafzi.Database, error) {
	db, err := lafzi.OpenDatabase(name, storageType)
	if err != nil {
		return nil, err
	}

	documents := make([]lafzi.Document, len(arabicTexts))
	for i, text := range arabicTexts {
		documents[i] = lafzi.Document{
			ID:         int64(i + 1),
			ArabicText: text,
		}
	}

	err = db.AddDocuments(documents...)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func removeDatabase() error {
	return os.RemoveAll("test-db.lafzi")
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

var searchTexts = map[string][]int64{
	// Text exist in source
	"bismi":         {1},
	"alhamdulillah": {2},
	"rabb":          {2, 12},
	"rahman":        {1, 3},
	"rahim":         {1, 3},
	"malik":         {4},
	"yaumiddin":     {4},
	"iyyaka":        {5},
	"dollin":        {7},
	"kitab":         {9},
	"muttakin":      {9},
	"muflihun":      {12},
	"aandzartahum":  {},
	"minan nas":     {15},
	"yaumil akhir":  {},
	"qulubihim":     {14, 17},
	"marodun":       {17},
	"yakdzibun":     {17},

	// Entire sentence
	"bismillahirrahmanirrahim": {1},

	// Wrong transliteration
	"basma":        {},
	"dooollliiin":  {7},
	"yaumal akhir": {},

	// Not exist in source
	"waylul":   {},
	"istigfar": {},
}
