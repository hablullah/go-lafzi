package database

type Document struct {
	ID       int    `db:"id"`
	Arabic   string `db:"arabic"`
	Phonetic string `db:"phonetic"`
}

type DocumentToken struct {
	DocumentID int    `db:"document_id"`
	Token      string `db:"token"`
}
