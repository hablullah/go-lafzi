package lafzi

import (
	"os"

	"go.etcd.io/bbolt"
)

type Dictionary struct {
	db *bbolt.DB
}

func NewDictionary(path string) (*Dictionary, error) {
	db, err := bbolt.Open(path, os.ModePerm, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	return &Dictionary{db}, nil
}

func (dict *Dictionary) AddEntry(identifier int, arabicText string) error {
	return nil
}
