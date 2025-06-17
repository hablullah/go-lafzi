package database

import (
	"github.com/guregu/null/v6"
	"github.com/hablullah/go-lafzi/internal/phonetic"
)

type Document struct {
	ID         int         `db:"id"`
	Identifier null.String `db:"identifier"`
	Arabic     string      `db:"arabic"`
	Tokens     []phonetic.NGram
}

type DocumentToken struct {
	DocumentID int    `db:"document_id"`
	Token      string `db:"token"`
	Start      int    `db:"start"`
	End        int    `db:"end"`
}
