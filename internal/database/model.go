package database

type Document struct {
	ID      int    `db:"id"`
	Content string `db:"content"`
}

type TokenLocation struct {
	DocID int `db:"document_id"`
	Count int `db:"n"`
}
