package shell

import (
	"bytes"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"os"
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

type runDependencies struct {
	Process process.Dependencies
}

func systemRunDependencies() runDependencies {
	return runDependencies{Process: process.SystemRunner()}
}

func Run(name string, args ...string) Output {
	return runWithDependencies(systemRunDependencies(), name, args...)
}

func runWithDependencies(dependencies runDependencies, name string, args ...string) (output Output) {
	command := process.Command{Name: name, Args: args}
	log.Debugf("running: %s", strings.Join(append([]string{command.Name}, command.Args...), " "))
	command.Stdout = &output.StdOut
	command.Stderr = &output.StdErr

	if err := process.Execute(dependencies.Process, command); err != nil {
		return output
	}

	return output
}

//func Unzip(file string, outputDir string) (string, error) {
//	return run(exec.Command("unzip", file, "-d", outputDir))
//}

type unzipDependencies struct {
	Files filesystem.Dependencies
}

func systemUnzipDependencies() unzipDependencies {
	return unzipDependencies{Files: filesystem.System()}
}

func Unzip(src string, dest string) (filenames []string, err error) {
	return unzipWithDependencies(systemUnzipDependencies(), src, dest)
}

func unzipWithDependencies(dependencies unzipDependencies, src string, dest string) (filenames []string, err error) {
	r, err := filesystem.OpenZipReader(dependencies.Files, src)
	if err != nil {
		return filenames, err
	}
	defer func() { _ = filesystem.CloseReader(dependencies.Files, r) }()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		log.Debugf("Unzipping dest: %s file: %s => %s", dest, f.Name, fpath)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) && dest != "." {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			err = filesystem.MkdirAll(dependencies.Files, fpath, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		if err = filesystem.MkdirAll(dependencies.Files, filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := filesystem.OpenFile(dependencies.Files, fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := filesystem.OpenZipEntry(dependencies.Files, f)
		if err != nil {
			return filenames, err
		}

		_, _ = filesystem.Copy(dependencies.Files, outFile, rc)

		err = filesystem.Close(dependencies.Files, outFile)
		if err != nil {
			return filenames, err
		}

		err = filesystem.CloseReader(dependencies.Files, rc)
		if err != nil {
			return filenames, err
		}
	}
	return
}
