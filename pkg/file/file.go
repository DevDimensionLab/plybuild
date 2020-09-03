package file

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Find(fileSuffix string, dir string) (result string, err error) {
	err = filepath.Walk(dir,
		func(path string, fi os.FileInfo, errIn error) error {
			if strings.HasSuffix(path, fileSuffix) {
				result = path
				return io.EOF
			}
			return nil
		})

	if err == io.EOF {
		err = nil
	}
	return
}

func ReadJson(file string, parsed interface{}) error {
	byteValue, err := Open(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &parsed)
	if err != nil {
		return err
	}

	return nil
}

func ReadXml(file string, parsed interface{}) error {
	byteValue, err := Open(file)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(byteValue, &parsed)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Open(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()

	return byteValue, nil
}
