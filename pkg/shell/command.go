package shell

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Run(name string, args ...string) (string, error) {
	return run(exec.Command(name, args...))
}

func run(cmd *exec.Cmd) (string, error) {
	log.Debugf("running: %s", cmd.String())
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s", stdOut.String(), stdErr.String()), err
	}

	return stdOut.String(), nil
}

//func Unzip(file string, outputDir string) (string, error) {
//	return run(exec.Command("unzip", file, "-d", outputDir))
//}

func Unzip(src string, dest string) (filenames []string, err error) {

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		err = outFile.Close()
		if err != nil {
			return filenames, err
		}

		err = rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return
}
