package shell

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/logger"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Output struct {
	StdOut bytes.Buffer
	StdErr bytes.Buffer
	Err    error
}

func (output Output) String() string {
	return fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s", output.StdOut.String(), output.StdErr.String())
}

func (output Output) FormatError() error {
	return logger.ExternalError(output.Err, output.String())
}

func Run(name string, args ...string) Output {
	return run(exec.Command(name, args...))
}

func run(cmd *exec.Cmd) (output Output) {
	log.Debugf("running: %s", cmd.String())
	cmd.Stdout = &output.StdOut
	cmd.Stderr = &output.StdErr

	if err := cmd.Run(); err != nil {
		return output
	}

	return output
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
		log.Debugf("Unzipping dest: %s file: %s => %s", dest, f.Name, fpath)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) && dest != "." {
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
