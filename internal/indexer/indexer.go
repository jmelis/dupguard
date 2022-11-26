package indexer

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmelis/dupguard/internal/db"
	"github.com/jmelis/dupguard/internal/hasher"
)

func indexSize(paths []string) {
	for _, path := range paths {
		log.Println("Indexing:", path)
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

			file := &db.File{Path: path, Size: info.Size()}
			file.Add()

			return nil
		})
	}
}

// func IndexHash1M() {
// 	files := db.DupesSize()
// 	for _, f := range files {
// 		f.Hash1M = hasher.Hash1M(f.Path)
// 		f.Add()
// 	}
// }

func IndexHash() {
	files := db.DupesSize()
	log.Println("Hashing files:", len(files))
	for _, f := range files {
		f.Hash = hasher.Hash(f.Path)
		f.Add()
	}
}

func Index(paths []string) {
	indexSize(paths)
	IndexHash()
}
