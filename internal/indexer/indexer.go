package indexer

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmelis/dupguard/internal/db"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func indexSize(paths []string) {
	for _, path := range paths {
		db.Prune(path)
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf(err.Error())
			}

			if info.IsDir() {
				return nil
			}

			if info.Size() == 0 {
				return nil
			}

			file := db.File{Path: path, Size: info.Size()}
			file.Add()

			return nil
		})
	}
}

func Index(paths []string) {
	indexSize(paths)
	// index with size

	// find size dups, index with hash1m

	// find hash1m dups, index with hash
}
