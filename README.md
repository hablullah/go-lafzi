# Go-Lafzi

[![Go Report Card][report-badge]][report-url]
[![Go Reference][doc-badge]][doc-url]

Go-Lafzi is a Go package for searching Arabic text using its transliteration (phonetic search). It works by
using indexed trigram for approximate string matching. For storing the indexes, it uses either [Bolt][bolt]
or [SQLite][sqlite] database.

If using Bolt database, the indexing process is slower and the index file is bigger than SQLite. However,
the search process is far faster. For example, here is the result when used for indexing an Al-Quran dataset
(6,236 documents) and searching using short keyword (in this case, I'm searching for "aamanna billah"):

| Storage | Indexing Time | Output Size  | Search Time |
|---------|---------------|--------------|-------------|
| SQLite  | 6,265 ms      | 20,389,888 B | 58 ms       |
| Bolt    | 9,648 ms      | 86,282,240 B | 19 ms       |

As you can see, there are trade offs between those two. Another important difference is SQLite uses CGo.
With that said, if you want to use pure Go language, you will need to use Bolt as the storage.

## Usage

For example, we want to find the word "rahman" within surah [Al-Fatiha][al-fatiha]:

```go
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
```

Which will give us following results :

```
'rahman' found in [1 3]
```

## Resources

1. Istiadi, MA. 2012. _Sistem Pencarian Ayat Al-Quran Berbasis Kemiripan Fonetis_. 
	Bachelor’s Thesis, Institut Pertanian Bogor ([PDF][istiadi-pdf], [university repo][istiadi-univ-repo])

## License

Go-Lafzi is distributed using [MIT] license.

[report-badge]: https://goreportcard.com/badge/github.com/hablullah/go-lafzi
[report-url]: https://goreportcard.com/report/github.com/hablullah/go-lafzi
[doc-badge]: https://pkg.go.dev/badge/github.com/hablullah/go-lafzi.svg
[doc-url]: https://pkg.go.dev/github.com/hablullah/go-lafzi

[bolt]: https://github.com/etcd-io/bbolt
[sqlite]: https://github.com/mattn/go-sqlite3
[al-fatiha]: http://tanzil.net/#1:1

[istiadi-pdf]: doc/2012-istiadi-ma.pdf
[istiadi-univ-repo]: http://repository.ipb.ac.id:8080/handle/123456789/56060?show=full

[MIT]: http://choosealicense.com/licenses/mit/