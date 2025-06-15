# Go-Lafzi [![Go Report Card][report-badge]][report-url] [![Go Reference][doc-badge]][doc-url]

Go-Lafzi is a Go package for searching Arabic text using its transliteration (phonetic search). Loosely based on research by Istiadi (2012) and multiple papers related with it.

It works by using indexed trigram for approximate string matching, while the search result is ranked using [_Longest Common Sequence_][lcs] with [Myers Diff Algorithm][myers]. For storing the indexes, it uses [SQLite][sqlite] database which brings several advantages:

- The indexing and lookup process is pretty fast, around 3 seconds for indexing entire Al-Quran and 90 ms per lookup. For more detail, checkout the code in [`sample/quran`](sample/quran).
- Can be safely used concurrently.

## Usage

For example, we want to find the word "rahman" within surah [Al-Fatiha][al-fatiha]:

```go
package main

import (
	"encoding/json"
	"fmt"

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
	storage, err := lafzi.OpenStorage("sample.lafzi")
	checkError(err)

	// Prepare documents
	var docs []lafzi.Document
	for i, arabicText := range arabicTexts {
		docs = append(docs, lafzi.Document{
			ID:     i + 1,
			Arabic: arabicText},
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
```

Which will give us following results :

```
[
  {
    "DocumentID": 1,
    "Text": "بِسْمِ اللَّهِ الرَّحْمَـٰنِ الرَّحِيمِ",
    "Confidence": 1,
    "Start": 17,
    "End": 27
  },
  {
    "DocumentID": 3,
    "Text": "الرَّحْمَـٰنِ الرَّحِيمِ",
    "Confidence": 1,
    "Start": 2,
    "End": 12
  }
]
```

## Resources

All resources mentioned here is also available in `doc` folder. This is done to prevent case where the university decided to close public access to these research. For example, paper by Istiadi was publicly available back in 2014, however now in 2022 it can only downloaded by member of its university.

By the way, the algorithm that implemented in this package is not exactly the same as in these papers. There are also some papers that I ignored, i.e. the papers to find Arabic text cross-verse in Qur'an, which I believe not really useful for general Arabic texts. There are also many parts that I've changed to make implementation easier and to increase performance in testing.

- Istiadi, Muhammad Abrar. "Sistem pencarian ayat al-quran berbasis kemiripan fonetis." (2012). ([PDF][istiadi-pdf], [link][istiadi-url])
- Zafran, Aidil, Moch Arif Bijaksana, and Kemas M. Lhaksmana. "Truncated query of phonetic search for al qur’an." 2019 7th International Conference on Information and Communication Technology (ICoICT). IEEE, 2019. ([PDF][zafran-pdf], [link][zafran-url])
- Rifaldi, Eki, Moch Arif Bijaksana, and Kemas Muslim Lhaksamana. "Sistem Pencarian Lintas Ayat Al-Qur'an Berdasarkan Kesamaan Fonetis." Indonesia Journal on Computing (Indo-JC) 4.2 (2019): 177-188. ([PDF][rifaldi-pdf], [link][rifaldi-url])
- Rasyad, Naufal, Moch Arif Bijaksana, and Kemas Muslim Lhaksmana. "Pencarian Potongan Ayat Al-Qur'an dengan Perbedaan Bunyi pada Tanda Berhenti Berdasarkan Kemiripan Fonetis." Jurnal Linguistik Komputasional 2.2 (2019): 56-61. ([PDF][rasyad-pdf], [link][rasyad-url])
- Satriady, Wildhan, Moch Arif Bijaksana, and Kemas M. Lhaksmana. "Quranic Latin Query Correction as a Search Suggestion." Procedia Computer Science 157 (2019): 183-190. ([PDF][satriady-pdf], [link][satriady-url])
- Octavia, Agni, Moch Arif Bijaksana, and Kemas Muslim Lhaksmana. "Verse Search System for Sound Differences in the Qur’an Based on the Text of Phonetic Similarities." Jurnal Sisfokom (Sistem Informasi dan Komputer) 9.3 (2020): 317-322. ([PDF][octavia-pdf], [link][octavia-url])
- Fitriani, Intan Khairunnisa, Moch Arif Bijaksana, and Kemas Muslim Lhaksmana. "Qur’an Search System for Handling Cross Verse Based on Phonetic Similarity." Jurnal Sisfokom (Sistem Informasi dan Komputer) 10.1 (2021): 46-51. ([PDF][fitriani-pdf], [link][fitriani-url])
- Purwita, Naila Iffah, et al. "Typo handling in searching of Quran verse based on phonetic similarities." Register: Jurnal Ilmiah Teknologi Sistem Informasi 6.2 (2020): 130-140. ([PDF][purwita-pdf], [link][purwita-url])
- Cendikia, Putri, Moch Arif Bijaksana, and Kemas M. Lhaksmana. "Pencarian Ayat Al-Qur'an Yang Tidak Utuh Berdasarkan Kemiripan Fonetis." eProceedings of Engineering 7.2 (2020). ([PDF][cendekia-pdf], [link][cendekia-url])
- Elder, Robert. "Myers Diff Algorithm - Code &amp; Interactive Visualization." (2017) ([archive][elder-archive], [link][elder-url])

## License

Go-Lafzi is distributed using [MIT] license.

[report-badge]: https://goreportcard.com/badge/github.com/hablullah/go-lafzi
[report-url]: https://goreportcard.com/report/github.com/hablullah/go-lafzi
[doc-badge]: https://pkg.go.dev/badge/github.com/hablullah/go-lafzi.svg
[doc-url]: https://pkg.go.dev/github.com/hablullah/go-lafzi
[lcs]: https://en.wikipedia.org/wiki/Longest_common_subsequence_problem
[myers]: https://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.4.6927
[sqlite]: https://github.com/mattn/go-sqlite3
[al-fatiha]: http://tanzil.net/#1:1
[istiadi-pdf]: doc/2012-ma-istiadi.pdf
[istiadi-url]: http://repository.ipb.ac.id:8080/handle/123456789/56060?show=full
[zafran-pdf]: doc/2019-a-zafran.pdf
[zafran-url]: https://ieeexplore.ieee.org/abstract/document/8835336/
[rifaldi-pdf]: doc/2019-e-rifaldi.pdf
[rifaldi-url]: http://socj.telkomuniversity.ac.id/ojs/index.php/indojc/article/view/342
[rasyad-pdf]: doc/2019-n-rasyad.pdf
[rasyad-url]: http://inacl.id/journal/index.php/jlk/article/view/25
[satriady-pdf]: doc/2019-w-satriady.pdf
[satriady-url]: https://www.sciencedirect.com/science/article/pii/S1877050919310749
[octavia-pdf]: doc/2020-a-octavia.pdf
[octavia-url]: http://jurnal.atmaluhur.ac.id/index.php/sisfokom/article/view/935
[fitriani-pdf]: doc/2020-ik-fitriani.pdf
[fitriani-url]: http://jurnal.atmaluhur.ac.id/index.php/sisfokom/article/view/986
[purwita-pdf]: doc/2020-ni-purwita.pdf
[purwita-url]: http://journal.unipdu.ac.id/index.php/register/article/view/2065
[cendekia-pdf]: doc/2020-p-cendekia.pdf
[cendekia-url]: https://openlibrarypublications.telkomuniversity.ac.id/index.php/engineering/article/view/13104
[elder-archive]: doc/2017-r-elder.htm
[elder-url]: https://blog.robertelder.org/diff-algorithm/
[mit]: http://choosealicense.com/licenses/mit/
