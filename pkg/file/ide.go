package file

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func CleanIntellijFiles(targetDir string, recursive bool, dryRun bool) (string, error) {
	if recursive {
		return cleanIntellijFiles(targetDir, dryRun)
	} else {
		return cleanIntellijFile(targetDir, dryRun)
	}
}

func cleanIntellijFiles(targetDir string, dryRun bool) (string, error) {
	var filesDeleted = 0
	var dirsDeleted = 0

	imlFiles, err := FindAll(".iml", []string{}, targetDir)
	if err != nil {
		return "", err
	}
	for _, f := range imlFiles {
		filesDeleted += 1
		if err = logAndDelete(f, dryRun); err != nil {
			return "", err
		}
	}

	ideaDirs, err := FindAll(".idea", []string{}, targetDir)
	if err != nil {
		return "", err
	}
	for _, d := range ideaDirs {
		dirsDeleted += 1
		if err = logAndDelete(d, dryRun); err != nil {
			return "", nil
		}
	}

	if dryRun {
		return fmt.Sprintf("Found %d files and %d directories to delete", filesDeleted, dirsDeleted), nil
	} else {
		return fmt.Sprintf("Deleted %d files and %d directories", filesDeleted, dirsDeleted), nil
	}
}

func cleanIntellijFile(targetDir string, dryRun bool) (string, error) {
	var imlFile = "not found"
	var ideaDir = "not found"

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".iml") && !f.IsDir() {
			fileName := Path("%s/%s", targetDir, f.Name())
			imlFile = fileName
			if err = logAndDelete(fileName, dryRun); err != nil {
				return "", err
			}
		}
		if strings.Contains(f.Name(), ".idea") && f.IsDir() {
			dirName := Path("%s/%s", targetDir, f.Name())
			ideaDir = dirName
			if err = logAndDelete(dirName, dryRun); err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf("Iml file: %s, .idea dir: %s", imlFile, ideaDir), nil
}

func logAndDelete(fileName string, dryRun bool) error {
	if dryRun {
		log.Infof("Found .iml file: %s", fileName)
	} else {
		log.Infof("Deleting .iml file: %s", fileName)
		err := DeleteDirectory(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}
