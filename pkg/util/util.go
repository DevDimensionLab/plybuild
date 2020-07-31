package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindFile(fileSuffix string, dir string) (result string, err error) {
	err = filepath.Walk(dir,
		func(path string, fi os.FileInfo, errIn error) error {
			if strings.HasSuffix( path, fileSuffix) {
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


func Wget(url, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Unzip(file string, outputDir string) error {
	cmd := exec.Command("unzip", file, "-d", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}


func ReadJson(file string, parsed interface{}) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &parsed)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	return nil
}