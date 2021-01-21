package lafzi

// StorageType is the type of data storage that will be used to store the indexes.
type StorageType uint8

const (
	// Bolt uses Bolt database, a simple and fast pure Go key/value database based on LMDB project.
	Bolt StorageType = iota

	// SQLite uses SQLite database, a small, fast, self-contained, high-reliability, full-featured,
	// relational SQL database engine.
	SQLite
)

type dataStorage interface {
	close()
	saveDocuments(documents ...Document) error
	findTokens(tokens ...string) ([]documentTokens, error)
}
