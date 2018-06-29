package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func useVersion(version string, useProject bool) (string, error) {
	if useProject {
		return getProjectGoVersion()
	}
	return version, nil
}

func getProjectGoVersion() (string, error) {
	ver, err := parseTravisYml(".travis.yml")
	if err != nil {
		return "", fmt.Errorf("failed to detect the project's golang version: %v", err)
	}

	return ver, nil
}

func parseTravisYml(name string) (string, error) {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	var re = regexp.MustCompile(`(?mi)^go:\s*\r?\n\s*-\s+(\S+)\s*$`)
	matches := re.FindAllStringSubmatch(string(file), 1)
	if len(matches) == 0 {
		return "", fmt.Errorf("go not found in %v", name)
	}

	goVersion := matches[0][1]
	return goVersion, nil
}
