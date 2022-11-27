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
				log.Println(err.Error())
				return nil
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

func IndexHash() {
	files := db.DupesSize()
	log.Println("Hashing files:", len(files))

	for i, f := range files {
		f.Hash = hasher.Hash(f.Path)
		f.Add()
		if i%10 == 0 {
			log.Printf("Hashing progress: %d/%d", i, len(files))
		}
	}
}

func Index(paths []string) {
	indexSize(paths)
	IndexHash()
}
