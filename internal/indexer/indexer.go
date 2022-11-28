package indexer

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmelis/dupguard/internal/db"
	"github.com/jmelis/dupguard/internal/hasher"
	"github.com/schollz/progressbar/v3"
)

func indexSize(paths []string) {
	var files []*db.File
	for _, path := range paths {
		log.Println("Indexing:", path)
		bar := progressbar.Default(-1)
		db.Prune(path)
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			bar.Add(1)

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
			files = append(files, file)

			return nil
		})
		bar.Exit()
	}

	db.BatchAdd(files)
}

func indexHash() {
	files := db.DupesSize()
	log.Println("Hashing files:", len(files))

	bar := progressbar.Default(int64(len(files)))
	for _, f := range files {
		f.Hash = hasher.Hash(f.Path)
		f.Add()
		bar.Add(1)
	}
}

func Index(paths []string) {
	indexSize(paths)
	indexHash()
}
