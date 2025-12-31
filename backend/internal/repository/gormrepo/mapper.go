package gormrepo

import (
	"encoding/json"
)

func marshalAuthors(authors []string) ([]byte, error) {
	if authors == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(authors)
}

func unmarshalAuthors(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var authors []string
	if err := json.Unmarshal(raw, &authors); err != nil {
		return nil
	}
	return authors
}
