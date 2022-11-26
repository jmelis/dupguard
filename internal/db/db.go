package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"

	"github.com/jmelis/dupguard/internal/hasher"
)

type File struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Path      string
	Size      int64
	Hash1M    string
	Hash      string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var dgdb *gorm.DB

func init() {
	var err error
	dgdb, err = gorm.Open(sqlite.Open("dupguard.db"))
	check(err)
	dgdb.AutoMigrate(&File{})
}
func Prune(path string) {
	fileInfo, err := os.Stat(path)
	check(err)

	if fileInfo.IsDir() {
		pathArg := strings.TrimSuffix(path, "/") + "/%"
		dgdb.Where("path LIKE ?", pathArg).Delete(&File{})
	} else {
		dgdb.Where("path = ?", path).Delete(&File{})
	}
}

func (fp *File) Add() {
	var filesSameSize []File
	dgdb.Where("path = ?", fp.Path).Delete(&File{})

	// look for files of the same size
	dgdb.Where("size = ?", fp.Size).Find(&filesSameSize)

	if len(filesSameSize) > 0 {
		fmt.Println("found matches for", fp.Path, fp.Size)
		fp.Hash1M = hasher.Hash1M(fp.Path)
		for _, f := range filesSameSize {
			fmt.Println("-", f.Path)
		}
		fmt.Println()
	}

	dgdb.Create(fp)
}
