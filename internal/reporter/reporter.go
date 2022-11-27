package reporter

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmelis/dupguard/internal/db"
)

func Dupes() {
	var dupFilesAll [][]string
	for _, files := range db.Dupes() {
		var dupFiles []string
		for _, f := range files {
			dupFiles = append(dupFiles, f.Path)
		}
		dupFilesAll = append(dupFilesAll, dupFiles)
	}

	data, err := json.Marshal(dupFilesAll)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
