package maven

import (
	"encoding/json"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"io/ioutil"
	"strings"
)

type GraphStyles struct {
	NodeStyles map[string]GraphStyle `json:"node-styles"`
}

type GraphStyle struct {
	Type      string `json:"type"`
	Color     string `json:"color"`
	FillColor string `json:"fill-color"`
	Style     string `json:"style"`
}

func GraphDefaultStyles() GraphStyles {
	styles := map[string]GraphStyle{
		"org.springframework*": {
			Type:      "box",
			Color:     "#95ca7d",
			FillColor: "#b4d4a1",
			Style:     "filled,rounded",
		},
		"org.jetbrains*": {
			Type:      "box",
			Color:     "#F9C5D5",
			FillColor: "#FEE3EC",
			Style:     "filled,rounded",
		},
		",,test": {
			Type:      "box",
			Color:     "#d2d2d2",
			FillColor: "#d3d3d3",
			Style:     "filled,rounded",
		},
	}
	return GraphStyles{NodeStyles: styles}
}

func GraphArgs() []string {
	return []string{"com.github.ferstl:depgraph-maven-plugin:graph",
		"-DshowVersions", "-DshowGroupIds", "-DshowConflicts", "-DshowDuplicates",
		"-DcustomStyleConfiguration=target/dependency-graph-styles.json",
	}
}

func Graph(onlySecondParty bool, excludeTestScope bool, includeFilters, excludeFilters []string) func(project config.Project) error {
	return func(project config.Project) error {
		mvnArgs := GraphArgs()
		defaultStyles := GraphDefaultStyles()
		secondParty, err := project.Type.Model().GetSecondPartyGroupId()

		if err == nil {
			if onlySecondParty {
				includeFilters = append(includeFilters, fmt.Sprintf("%s*", secondParty))
			}
			wildCardProjectGroupId := fmt.Sprintf("%s*", secondParty)
			defaultStyles.NodeStyles[wildCardProjectGroupId] = GraphStyle{
				Type:      "box",
				Color:     "#95ca7d",
				FillColor: "#79B4B7",
				Style:     "filled,rounded",
			}
		}

		if len(includeFilters) > 0 {
			mvnArgs = append(mvnArgs, fmt.Sprintf("-Dincludes=%s", strings.Join(includeFilters, ",")))
		}
		if len(excludeFilters) > 0 {
			mvnArgs = append(mvnArgs, fmt.Sprintf("-Dexcludes=%s", strings.Join(excludeFilters, ",")))
		}

		if excludeTestScope {
			mvnArgs = append(mvnArgs, "-Dscope=compile")
		}

		if err := WriteGraphStyles(defaultStyles, project.Path); err != nil {
			return err
		}

		return RunOn("mvn", mvnArgs...)(project)
	}
}

func WriteGraphStyles(styles GraphStyles, projectPath string) error {
	jsonStyles, err := json.MarshalIndent(styles, "", "    ")
	if err != nil {
		return err
	}

	stylesFile := file.Path("%s/target/dependency-graph-styles.json", projectPath)
	if file.Exists(stylesFile) {
		return nil
	}
	return ioutil.WriteFile(stylesFile, jsonStyles, 0644)
}
