package spring

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/http"
	"github.com/devdimensionlab/plybuild/pkg/shell"
	"net/url"
	"os"
	"strings"
	"time"
)

var baseUrl = "https://start.spring.io"

func UrlValuesFrom(bootVersion string, config config.ProjectConfiguration) url.Values {
	// see https://github.com/spring-io/initializr#generating-a-project
	params := url.Values{}
	params.Add("groupId", config.GroupId)
	params.Add("artifactId", config.ArtifactId)
	params.Add("packageName", config.Package)
	params.Add("dependencies", strings.Join(config.Dependencies, ","))
	params.Add("javaVersion", "11")
	params.Add("language", config.Language)
	params.Add("description", config.Description)
	params.Add("name", config.Name)
	params.Add("type", "maven-project")
	if bootVersion != "" {
		params.Add("bootVersion", bootVersion)
	}
	//params.Add("baseDir", targetDir)

	return params
}

func GetRoot() (IoRootResponse, error) {
	return getRoot(systemDiscoveryDependencies())
}

func getRoot(dependencies discoveryDependencies) (IoRootResponse, error) {
	var deps IoRootResponse
	err := dependencies.GetJSON(httpclient.Request{URL: baseUrl}, &deps)
	return deps, err
}

func GetDependencies() (IoDependenciesResponse, error) {
	return getDependencies(systemDiscoveryDependencies())
}

func getDependencies(dependencies discoveryDependencies) (IoDependenciesResponse, error) {
	var deps IoDependenciesResponse
	err := dependencies.GetJSON(httpclient.Request{URL: baseUrl + "/dependencies"}, &deps)
	return deps, err
}

func Validate(config config.ProjectConfiguration) error {
	return validate(systemDiscoveryDependencies(), config)
}

type discoveryHTTPClient interface {
	GetJSON(httpclient.Request, interface{}) error
}

type discoveryDependencies struct {
	Client discoveryHTTPClient
}

func (dependencies discoveryDependencies) GetJSON(request httpclient.Request, parsed interface{}) error {
	if dependencies.Client == nil {
		return httpclient.ErrNoClient
	}
	return dependencies.Client.GetJSON(request, parsed)
}

func systemDiscoveryDependencies() discoveryDependencies {
	return discoveryDependencies{Client: packageDiscoveryHTTP{}}
}

type packageDiscoveryHTTP struct{}

func (packageDiscoveryHTTP) GetJSON(request httpclient.Request, parsed interface{}) error {
	return http.GetJson(request.URL, parsed)
}

func validate(dependencies discoveryDependencies, config config.ProjectConfiguration) error {
	if len(config.Dependencies) == 0 {
		return nil
	}

	var invalidDependencies []string
	validDependencies, err := getDependencies(dependencies)
	if err != nil {
		return err
	}

	for _, userDefinedDependency := range config.Dependencies {
		valid := false
		for validDependency := range validDependencies.Dependencies {
			if validDependency == userDefinedDependency {
				valid = true
			}
		}
		if !valid {
			invalidDependencies = append(invalidDependencies, userDefinedDependency)
		}
	}

	if len(invalidDependencies) > 0 {
		validKeys := make([]string, 0, len(validDependencies.Dependencies))
		for k := range validDependencies.Dependencies {
			validKeys = append(validKeys, k)
		}
		return fmt.Errorf("%s not found in valid list of dependencies %s", invalidDependencies, validKeys)
	} else {
		return nil
	}
}

func DownloadInitializer(targetDir string, formData url.Values) error {
	targetArchiveFile, err := archivePath()
	if err != nil {
		return err
	}
	return downloadInitializer(initializerDependencies{IO: packageInitializerIO{}}, targetDir, formData, targetArchiveFile)
}

type initializerIO interface {
	Download(string, string) error
	Unzip(string, string) ([]string, error)
	Delete(string) error
}

type initializerDependencies struct {
	IO initializerIO
}

func (dependencies initializerDependencies) Download(downloadURL, archivePath string) error {
	if dependencies.IO == nil {
		return httpclient.ErrNoClient
	}
	return dependencies.IO.Download(downloadURL, archivePath)
}

func (dependencies initializerDependencies) Unzip(archivePath, targetDirectory string) ([]string, error) {
	if dependencies.IO == nil {
		return nil, filesystem.ErrNoFilesystem
	}
	return dependencies.IO.Unzip(archivePath, targetDirectory)
}

func (dependencies initializerDependencies) Delete(archivePath string) error {
	if dependencies.IO == nil {
		return filesystem.ErrNoFilesystem
	}
	return dependencies.IO.Delete(archivePath)
}

type packageInitializerIO struct{}

func (packageInitializerIO) Download(downloadURL, archivePath string) error {
	return http.Wget(downloadURL, archivePath)
}

func (packageInitializerIO) Unzip(archivePath, targetDirectory string) ([]string, error) {
	return shell.Unzip(archivePath, targetDirectory)
}

func (packageInitializerIO) Delete(archivePath string) error {
	return file.DeleteSingleFile(archivePath)
}

func downloadInitializer(dependencies initializerDependencies, targetDir string, formData url.Values, targetArchiveFile string) error {
	downloadUrl := fmt.Sprintf("%s/starter.zip?%s", baseUrl, formData.Encode())
	log.Infof("Downloading from %s to %s", baseUrl, targetArchiveFile)
	err := dependencies.Download(downloadUrl, targetArchiveFile)
	if err != nil {
		return err
	}

	log.Infof("Unzipping %s to %s", targetArchiveFile, targetDir)
	_, err = dependencies.Unzip(targetArchiveFile, targetDir)
	if err != nil {
		return err
	}

	log.Infof("Deleting archive file: %s", targetArchiveFile)
	err = dependencies.Delete(targetArchiveFile)
	return err
}

func archivePath() (path string, err error) {
	curDir, err := os.Getwd()
	if err != nil {
		return
	}

	now := time.Now().Unix()
	path = file.Path("%s/spring-%d.zip", curDir, now)
	return
}

func DeleteDemoFiles(targetDir string, orderConfig config.ProjectConfiguration) {

	var fileSuffix = ".kt"
	if orderConfig.Language == "java" {
		fileSuffix = ".java"
	}

	testFile, err := file.FindFirst(fileSuffix, file.Path("%s/src/test/%s", targetDir, orderConfig.Language))
	if err != nil {
		log.Warnf("Unable to find testfile, fileSuffix=" + fileSuffix)
	} else {
		log.Debugf("Deleting testfile for languge %s: %s", orderConfig.Language, testFile)
		err := file.DeleteSingleFile(testFile)
		if err != nil {
			log.Warnf("Unable to delete testfile: %s", testFile)
		}
	}

	for _, f := range []string{"HELP.md", "mvnw", "mvnw.cmd"} {
		log.Debugf("Deleting demofile: %s", f)
		err := file.DeleteSingleFile(file.Path("%s/%s", targetDir, f))
		if err != nil {
			log.Warnf(err.Error())
		}
	}
}
