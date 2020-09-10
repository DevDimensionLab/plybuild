package file

import (
	"co-pilot/pkg/logger"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var log = logger.Context()

func FindFirst(fileSuffix string, dir string) (result string, err error) {
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

func FindAll(suffix string, excludes []string, dir string) (result []string, err error) {
	err = filepath.Walk(dir,
		func(path string, fi os.FileInfo, errIn error) error {
			if strings.HasSuffix(path, suffix) && !SuffixIn(path, excludes) {
				result = append(result, path)
			}
			return nil
		})

	if err == io.EOF {
		err = nil
	}
	return
}

func SuffixIn(keyword string, list []string) bool {
	for _, w := range list {
		if strings.Contains(keyword, w) {
			return true
		}
	}
	return false
}

func ReadJson(file string, parsed interface{}) error {
	byteValue, err := Open(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read %s, %v", file, err))
	}

	err = json.Unmarshal(byteValue, &parsed)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal %s, %v", file, err))
	}

	return nil
}

func ReadXml(file string, parsed interface{}) error {
	byteValue, err := Open(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open %s, %v", file, err))
	}

	err = xml.Unmarshal(byteValue, &parsed)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to unmarshal %s, %v", file, err))
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

func Overwrite(lines []string, filePath string) error {
	return ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

func CopyOrMerge(sourceFile string, destinationFile string) error {
	for _, f := range FilesToIgnore() {
		if strings.Contains(sourceFile, f) {
			log.Infof("ignoring copying %s", sourceFile)
			return nil
		}
	}

	if Exists(destinationFile) {
		log.Infof("merging %s with %s", sourceFile, destinationFile)
		return MergeFiles(sourceFile, destinationFile)
	}

	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	destinationParts := strings.Split(destinationFile, "/")
	destinationDir := strings.Join(destinationParts[:len(destinationParts)-1], "/")
	if !Exists(destinationDir) {
		err = CreateDirectory(destinationDir)
		if err != nil {
			return err
		}
	}

	log.Infof("copying \n\tFROM => %s\n\tTO => %s", sourceFile, destinationFile)
	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func FilesToIgnore() []string {
	return []string{"pom.xml", "co-pilot.json", "/.mvn/wrapper/", "Application"}
}

func RelPath(sourceDirectory string, filePath string) (string, error) {

	directoryParts := strings.Split(sourceDirectory, "/")
	fileParts := strings.Split(filePath, "/")

	if len(directoryParts) >= len(fileParts) {
		return "", errors.New("directory cannot be deeper than filePath")
	}

	cut := 0

	for i := range directoryParts {
		if directoryParts[i] == fileParts[i] {
			cut += 1
		} else {
			break
		}
	}

	return strings.Join(fileParts[cut:], "/"), nil
}

func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return err
		}
	}

	return nil
}

func SearchReplace(filePath string, from string, to string) error {
	b, err := Open(filePath)
	if err != nil {
		return err
	}

	replaced := strings.ReplaceAll(string(b), from, to)
	return Overwrite(strings.Split(replaced, "\n"), filePath)
}

func MergeFiles(fromFile string, toFile string) error {
	fromBytes, err := Open(fromFile)
	if err != nil {
		return err
	}
	fromLines := strings.Split(string(fromBytes), "\n")

	toBytes, err := Open(toFile)
	if err != nil {
		return err
	}
	toLines := strings.Split(string(toBytes), "\n")

	var newLines []string
	for _, fromLine := range fromLines {
		var hasLine = false
		for _, toLine := range toLines {
			if fromLine == toLine {
				hasLine = true
			}
		}
		if !hasLine {
			newLines = append(newLines, fromLine)
			log.Infof("appending line: %s", fromLine)
		}
	}

	toLines = append(toLines, newLines...)

	return Overwrite(toLines, toFile)
}
