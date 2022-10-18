package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		handleError(errors.New("backup needs a directory"))
	}

	target, err := filepath.Abs(os.Args[1])
	handleError(err)

	zipper(target)
}

func zipper(target string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	currentTime := time.Now()
	date := currentTime.Format("20060102150405000000000")
	targetName := path.Join(exPath, "backups", fmt.Sprintf("%s_%s.zip", filepath.Base(target), date))

	file, err := os.Create(targetName)
	handleError(err)

	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(os.Args[1], walker)
	handleError(err)
}
