package helper

import (
	"encoding/json"
	"os"

	"github.com/Spikree/git-issues/types"
)

const DbPath = ".issues_db.json"

// Load database state from disk into memory
func loadDB() (map[string]types.Issue, error) {
	db := make(map[string]types.Issue)

	// if the database isnt initilised yet, return initialised map
	if _, err := os.Stat(DbPath); os.IsNotExist(err) {
		return db, nil
	}

	data, err := os.ReadFile(DbPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &db)
	return db, err
}
