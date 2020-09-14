package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Recursive(targetDir string, keyword string) (files []string, err error) {
	err = filepath.Walk(targetDir,
		func(path string, fi os.FileInfo, errIn error) error {
			hit, err := Grep(path, keyword)
			if err != nil {
				return err
			}
			if hit {
				files = append(files, path)
			}
			return nil
		})

	if err == io.EOF {
		err = nil
	}

	return
}

func Grep(filepath string, keyword string) (hit bool, err error) {
	lines, err := OpenLines(filepath)
	if err != nil {
		return
	}

	for _, line := range lines {
		if strings.Contains(line, keyword) {
			hit = true
			return
		}
	}

	return
}
