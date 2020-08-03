package spring

import (
	"os"
	"os/exec"
	"spring-boot-co-pilot/pkg/util"
)

func SpringBootCLI(springExec string, arg ...string) error {
	cmd := exec.Command(springExec, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitFrom(config InitConfiguration) []string {
	var output []string

	output = append(output, "init")
	output = append(output, "-g="+config.GroupId)
	output = append(output, "-a="+config.ArtifactId)
	output = append(output, "--package="+config.Package)
	//output = append(output, "--d=...") ... hva var dette igjen?
	output = append(output, "-j=11")
	output = append(output, "--language=kotlin")
	output = append(output, "--description="+config.Description)
	output = append(output, "--name="+config.Name)
	output = append(output, "webservice") // outputdir

	return output
}

func DefaultConfiguration() InitConfiguration {
	return InitConfiguration{
		GroupId:     "com.example.demo",
		ArtifactId:  "demo-webservice",
		Package:     "com.example.demo",
		Name:        "webservice",
		Description: "demo webservice",
	}
}

func GetRoot() (IoRootResponse, error) {
	var deps IoRootResponse
	err := util.GetUrl("http://start.spring.io", &deps)
	return deps, err
}

func GetInfo() (IoInfoResponse, error) {
	var deps IoInfoResponse
	err := util.GetUrl("http://start.spring.io/actuator/info", &deps)
	return deps, err
}

func GetDependencies() (IoDependenciesResponse, error) {
	var deps IoDependenciesResponse
	err := util.GetUrl("http://start.spring.io/dependencies", &deps)
	return deps, err
}
