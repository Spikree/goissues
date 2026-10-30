package helper

import (
	"encoding/json"
	"os"

	"github.com/Spikree/git-issues/types"
)

func SaveDB(db map[string]types.Issue) error {
	// marshalIndent keeps the file human readable in Git diffs
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(DbPath, data, 0o644)
}
