package lafzi

type StorageType uint8

const (
	Bolt StorageType = iota
	SQLite
)

type dataStorage interface {
	close()
	saveEntries(entries ...DatabaseEntry) error
	findTokens(tokens ...string) ([]dbEntryTokens, error)
}
