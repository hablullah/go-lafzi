package lafzi

type StorageType uint8

const (
	Bolt StorageType = iota
	SQLite
)

type dataStorage interface {
	close()
	saveEntries(entries ...DictionaryEntry) error
	findTokens(tokens ...string) ([]dictionaryEntryTokens, error)
}
