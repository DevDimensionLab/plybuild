package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var springBootDownloadUrl = "https://repo.spring.io/release/org/springframework/boot/spring-boot-cli/[RELEASE]/spring-boot-cli-[RELEASE]-bin.zip"

func main() {
	targetDir := "./target"
	argsWithoutProg := os.Args[1:]
	if 0 == len(argsWithoutProg) {
		os.RemoveAll(targetDir)
		os.MkdirAll(targetDir, os.ModePerm)
		springBootCliZip := "./target/spring-boot-cli.zip"
		Wget(springBootDownloadUrl, springBootCliZip)
		Unzip( springBootCliZip, targetDir )
	}

	springExec, err := FindFile("bin/spring", targetDir)
	SpringBootCLI(springExec)

	if err != nil {
		log.Println(err)
	}

}


func SpringBootCLI(springExec string ) error {
	cmd := exec.Command(springExec)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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

func Unzip(file string, outputDir string ) error {
	cmd := exec.Command("Unzip", file,"-d", outputDir )
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}