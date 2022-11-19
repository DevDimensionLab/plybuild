package config

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"strings"
)

func (projectSettings *ProjectSettings) stealthWriter(entries []map[string]interface{}, outputFile string) error {
	lines, err := file.OpenLines(outputFile)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		typ := fmt.Sprintf("%v", entry["type"])
		versionIsProperty := fmt.Sprintf("%v", entry["versionIsProperty"])
		versionValue := fmt.Sprintf("%v", entry["versionValue"])
		artifactId := fmt.Sprintf("%v", entry["artifactId"])
		groupId := fmt.Sprintf("%v", entry["groupId"])
		oldVersion := fmt.Sprintf("%v", entry["oldVersion"])
		newVersion := fmt.Sprintf("%v", entry["newVersion"])

		if !strings.Contains(typ, "outdated") {
			continue
		}

		log.Debugf("Stealth upgrading: %s:%s (%s -> %s) in %s", artifactId, groupId, oldVersion, newVersion, outputFile)
		if versionIsProperty == "true" && versionValue != "" {
			versionKey := strings.ReplaceAll(versionValue, "${", "")
			versionKey = strings.ReplaceAll(versionKey, "}", "")
			oldVersionProperty := fmt.Sprintf("<%s>%s</%s>", versionKey, oldVersion, versionKey)
			newVersionProperty := fmt.Sprintf("<%s>%s</%s>", versionKey, newVersion, versionKey)

			for i, l := range lines {
				if strings.Contains(l, oldVersionProperty) {
					lines[i] = strings.ReplaceAll(l, oldVersionProperty, newVersionProperty)
				}
			}
		} else {
			for i, l := range lines {
				if strings.Contains(l, oldVersion) {
					surroundingLines := strings.Join(lines[i-2:i+2], " ")
					// check if artifactId is in the surrounding lines, preventing writing to wrong dependency
					if strings.Contains(surroundingLines, artifactId) {
						lines[i] = strings.ReplaceAll(lines[i], oldVersion, newVersion)
					}
				}
			}
		}
	}

	return file.Overwrite(lines, outputFile)
}
