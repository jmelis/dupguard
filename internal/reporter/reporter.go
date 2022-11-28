package reporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func Check(paths []string) {
	for _, path := range paths {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err.Error())
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if info.Size() == 0 {
				return nil
			}

			// check size
			if db.CheckSize(info.Size()) {
				fmt.Println(path)
			}
			return nil
		})
	}
}
