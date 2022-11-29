package db

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type File struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Path      string
	Size      int64
	Hash      string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var dgdb *gorm.DB

func Setup(dbpath string) error {
	var err error
	dglogger := logger.New(
		log.New(os.Stderr, "", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
		},
	)

	dgdb, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{
		Logger:          dglogger,
		CreateBatchSize: 1000,
	})

	if err != nil {
		return err
	}

	return dgdb.AutoMigrate(&File{})
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
	dgdb.Where("path = ?", fp.Path).Delete(&File{})
	dgdb.Create(fp)
}

func BatchAdd(files []*File) {
	dgdb.Create(files)
}

// DupesSize looks files of the same size and Hash unset
func DupesSize() []File {
	sql := `select * from files where
				size in (select size from files group by size having count(size) > 1)
				and hash = ''`
	var files []File
	dgdb.Raw(sql).Scan(&files)
	return files
}

// Dupes looks for duplicated files
func Dupes() map[string][]File {
	sql := `select * from files where
				hash in (select hash from files where hash != '' group by hash having count(hash) > 1)`
	var files []File
	dgdb.Raw(sql).Scan(&files)

	filesMap := make(map[string][]File)
	for _, f := range files {
		if val, ok := filesMap[f.Hash]; ok {
			filesMap[f.Hash] = append(val, f)
		} else {
			filesMap[f.Hash] = []File{f}
		}
	}

	return filesMap
}

func CheckSize(size int64) bool {
	var file File
	result := dgdb.Where("size = ?", size).First(&file)
	return result.RowsAffected == 1
}

func CheckHash(hash string) bool {
	var file File
	result := dgdb.Where("hash = ?", hash).First(&file)
	return result.RowsAffected == 1
}
