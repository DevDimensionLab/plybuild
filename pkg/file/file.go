package file

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type findFirstDependencies struct {
	Files filesystem.Dependencies
}

func systemFindFirstDependencies() findFirstDependencies {
	return findFirstDependencies{Files: filesystem.System()}
}

func FindFirst(fileSuffix string, dir string) (result string, err error) {
	return findFirst(systemFindFirstDependencies(), fileSuffix, dir)
}

func findFirst(dependencies findFirstDependencies, fileSuffix string, dir string) (result string, err error) {
	err = filesystem.Walk(dependencies.Files, dir,
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
		return fmt.Errorf("Unable to read %s, %v", file, err)
	}

	err = json.Unmarshal(byteValue, &parsed)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal %s, %v", file, err)
	}

	return nil
}

func ReadXml(file string, parsed interface{}) error {
	byteValue, err := Open(file)
	if err != nil {
		return fmt.Errorf("Unable to open %s, %v", file, err)
	}

	err = xml.Unmarshal(byteValue, &parsed)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal %s, %v", file, err)
	}

	return nil
}

type existsDependencies struct {
	Files filesystem.Dependencies
}

func systemExistsDependencies() existsDependencies {
	return existsDependencies{Files: filesystem.System()}
}

func Exists(filename string) bool {
	return exists(systemExistsDependencies(), filename)
}

func exists(dependencies existsDependencies, filename string) bool {
	return filesystem.Exists(dependencies.Files, filename)
}

type openDependencies struct {
	Files filesystem.Dependencies
}

func systemOpenDependencies() openDependencies {
	return openDependencies{Files: filesystem.System()}
}

func Open(filePath string) ([]byte, error) {
	return open(systemOpenDependencies(), filePath)
}

func open(dependencies openDependencies, filePath string) ([]byte, error) {
	byteValue, err := filesystem.ReadFile(dependencies.Files, filePath)
	if err != nil {
		return []byte{}, err
	}

	return byteValue, nil
}

func OpenLinesStrict(filePath string) ([]string, error) {
	b, err := Open(filePath)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(b), "\n"), nil
}

func OpenLines(filePath string) ([]string, error) {
	b, err := Open(filePath)
	if err != nil {
		return []string{}, nil
	}

	return strings.Split(string(b), "\n"), nil
}

type overwriteDependencies struct {
	Files filesystem.Dependencies
}

func systemOverwriteDependencies() overwriteDependencies {
	return overwriteDependencies{Files: filesystem.System()}
}

func Overwrite(lines []string, filePath string) error {
	return overwrite(systemOverwriteDependencies(), lines, filePath)
}

func overwrite(dependencies overwriteDependencies, lines []string, filePath string) error {
	return filesystem.WriteFile(dependencies.Files, filePath, []byte(strings.Join(lines, "\n")), 0644)
}

func CopyOrMerge(sourceFile string, destinationFile string) error {
	return copyOrMerge(filesystem.System(), sourceFile, destinationFile)
}

func copyOrMerge(dependencies filesystem.Dependencies, sourceFile string, destinationFile string) error {
	if filesystem.Exists(dependencies, destinationFile) {
		return mergeFile(sourceFile, destinationFile)
	}

	return copyFile(dependencies, sourceFile, destinationFile)
}

func mergeFile(sourceFile string, destinationFile string) error {
	if strings.HasSuffix(sourceFile, ".java") || strings.HasSuffix(sourceFile, ".kt") {
		log.Infof(fmt.Sprintf("ignoring merge of java or kt files: %s", sourceFile))
		return nil
	}

	if strings.HasSuffix(sourceFile, ".properties") {
		log.Infof("merging key=val property file %s with %s", sourceFile, destinationFile)
		return MergeKeyValFile(sourceFile, destinationFile, "=")
	}

	log.Infof("merging text file %s with %s", sourceFile, destinationFile)
	return MergeTextFiles(sourceFile, destinationFile)
}

func CopyFile(sourceFile string, destinationFile string) error {
	return copyFile(filesystem.System(), sourceFile, destinationFile)
}

func copyFile(dependencies filesystem.Dependencies, sourceFile string, destinationFile string) error {
	input, err := filesystem.ReadFile(dependencies, sourceFile)
	if err != nil {
		return err
	}

	pathSeparator := string(os.PathSeparator)
	destinationParts := strings.Split(destinationFile, pathSeparator)
	destinationDir := strings.Join(destinationParts[:len(destinationParts)-1], pathSeparator)
	if !filesystem.Exists(dependencies, destinationDir) {
		err = createDirectory(dependencies, destinationDir)
		if err != nil {
			return err
		}
	}

	fileInfo, err := filesystem.Stat(dependencies, sourceFile)
	if err != nil {
		return err
	}

	log.Debugf("copying FROM\t <= %s", sourceFile)
	log.Debugf("copying TO\t => %s", destinationFile)
	return filesystem.WriteFile(dependencies, destinationFile, input, fileInfo.Mode())
}

func createDirectory(dependencies filesystem.Dependencies, path string) error {
	_, err := filesystem.Stat(dependencies, path)
	if os.IsNotExist(err) {
		errDir := filesystem.MkdirAll(dependencies, path, 0755)
		if errDir != nil {
			return err
		}
	}

	return nil
}

func RelPath(sourceDirectory string, filePath string) (string, error) {
	pathSeparator := string(os.PathSeparator)
	directoryParts := strings.Split(sourceDirectory, pathSeparator)
	fileParts := strings.Split(filePath, pathSeparator)

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

	return strings.Join(fileParts[cut:], pathSeparator), nil
}

type createDirectoryDependencies struct {
	Files filesystem.Dependencies
}

func systemCreateDirectoryDependencies() createDirectoryDependencies {
	return createDirectoryDependencies{Files: filesystem.System()}
}

func CreateDirectory(path string) error {
	return createDirectoryWithDependencies(systemCreateDirectoryDependencies(), path)
}

func createDirectoryWithDependencies(dependencies createDirectoryDependencies, path string) error {
	_, err := filesystem.Stat(dependencies.Files, path)
	if os.IsNotExist(err) {
		errDir := filesystem.MkdirAll(dependencies.Files, path, 0755)
		if errDir != nil {
			return err
		}
	}

	return nil
}

type createFileDependencies struct {
	Files filesystem.Dependencies
}

func systemCreateFileDependencies() createFileDependencies {
	return createFileDependencies{Files: filesystem.System()}
}

func CreateFile(path, content string) error {
	return createFile(systemCreateFileDependencies(), path, content)
}

func createFile(dependencies createFileDependencies, path, content string) error {
	return filesystem.WriteFile(dependencies.Files, path, []byte(content), 0644)
}

type openFileDependencies struct {
	Files filesystem.Dependencies
}

func systemOpenFileDependencies() openFileDependencies {
	return openFileDependencies{Files: filesystem.System()}
}

func OpenFile(fileName string) (*os.File, error) {
	return openFile(systemOpenFileDependencies(), fileName)
}

func openFile(dependencies openFileDependencies, fileName string) (*os.File, error) {
	if !filesystem.Exists(dependencies.Files, fileName) {
		if err := filesystem.WriteFile(dependencies.Files, fileName, []byte{}, 0644); err != nil {
			return nil, err
		}
	}

	return filesystem.OpenFile(dependencies.Files, fileName, os.O_APPEND|os.O_WRONLY, 0644)
}

func SearchReplace(filePath string, from string, to string) error {
	if from == "" {
		return nil
	}

	b, err := Open(filePath)
	if err != nil {
		return err
	}

	log.Debugf("string-replacing: %s [%s => %s]", filePath, from, to)
	replaced := strings.ReplaceAll(string(b), from, to)
	return Overwrite(strings.Split(replaced, "\n"), filePath)
}

func MergeKeyValFile(fromFile string, toFile string, separator string) error {
	fromLines, err := OpenLines(fromFile)
	if err != nil {
		return err
	}

	toLines, err := OpenLines(toFile)
	if err != nil {
		return err
	}

	var newLines []string
	for _, fromLine := range fromLines {
		if fromLine == "" || !strings.Contains(fromLine, separator) {
			continue
		}

		var hasLine = false
		fromParts := strings.Split(fromLine, separator)
		fromKey := fromParts[0]

		for _, toLine := range toLines {
			if toLine == "" && !strings.Contains(toLine, separator) {
				continue
			}
			toParts := strings.Split(toLine, separator)
			toKey := toParts[0]
			if fromKey == toKey {
				log.Debugf("ignoring line due to key duplicate found in source %s: '%s' and '%s' in target:%s", fromFile, fromLine, toLine, toFile)
				hasLine = true
			}
		}
		if !hasLine {
			newLines = append(newLines, fromLine)
			log.Debugf("appending line: '%s', to:%s", fromLine, toFile)
		}
	}

	toLines = append(toLines, newLines...)

	return Overwrite(toLines, toFile)
}

func MergeTextFiles(fromFile string, toFile string) error {
	fromLines, err := OpenLines(fromFile)
	if err != nil {
		return err
	}

	toLines, err := OpenLines(toFile)
	if err != nil {
		return err
	}

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
			log.Debugf("appending line: '%s', to:%s", fromLine, toFile)
		}
	}

	toLines = append(toLines, newLines...)

	return Overwrite(toLines, toFile)
}

func Equal(fileA string, fileB string) (bool, error) {
	fileALines, err := OpenLines(fileA)
	if err != nil {
		return false, err
	}
	fileBLines, err := OpenLines(fileB)
	if err != nil {
		return false, err
	}

	if len(fileALines) != len(fileBLines) {
		return false, nil
	}

	for i := range fileALines {
		if fileALines[i] != fileBLines[i] {
			return false, nil
		}
	}

	return true, nil
}

type deleteSingleFileDependencies struct {
	Files filesystem.Dependencies
}

func systemDeleteSingleFileDependencies() deleteSingleFileDependencies {
	return deleteSingleFileDependencies{Files: filesystem.System()}
}

func DeleteSingleFile(filePath string) error {
	return deleteSingleFile(systemDeleteSingleFileDependencies(), filePath)
}

func deleteSingleFile(dependencies deleteSingleFileDependencies, filePath string) error {
	return filesystem.Remove(dependencies.Files, filePath)
}

type deleteAllDependencies struct {
	Files filesystem.Dependencies
}

func systemDeleteAllDependencies() deleteAllDependencies {
	return deleteAllDependencies{Files: filesystem.System()}
}

func DeleteAll(dirPath string) error {
	return deleteAll(systemDeleteAllDependencies(), dirPath)
}

func deleteAll(dependencies deleteAllDependencies, dirPath string) error {
	return filesystem.RemoveAll(dependencies.Files, dirPath)
}

type clearDirDependencies struct {
	Files filesystem.Dependencies
}

func systemClearDirDependencies() clearDirDependencies {
	return clearDirDependencies{Files: filesystem.System()}
}

func ClearDir(dirPath string, excludes []string) error {
	return clearDir(systemClearDirDependencies(), dirPath, excludes)
}

func clearDir(dependencies clearDirDependencies, dirPath string, excludes []string) error {
	files, err := filesystem.Glob(dependencies.Files, filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		var skip = false
		for _, exclude := range excludes {
			if strings.Contains(file, exclude) {
				log.Debugf("Skipping removal of: %s", file)
				skip = true
			}
		}
		if skip {
			continue
		}
		log.Debugf("Removing: %s", file)
		err = filesystem.RemoveAll(dependencies.Files, file)
		if err != nil {
			return err
		}
	}
	return nil
}

type moveDependencies struct {
	Files filesystem.Dependencies
}

func systemMoveDependencies() moveDependencies {
	return moveDependencies{Files: filesystem.System()}
}

func Move(source, destination string) error {
	return move(systemMoveDependencies(), source, destination)
}

func move(dependencies moveDependencies, source, destination string) error {
	return filesystem.Rename(dependencies.Files, source, destination)
}
