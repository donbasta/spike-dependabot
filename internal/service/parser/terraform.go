package parser

import (
	"log"
	"strings"
)

type InvalidSourceError struct {
	msg string
}

func (e *InvalidSourceError) Error() string { return e.msg }

func findUrlFromSource(source string) (string, error) {
	remotes := []string{"github.com", "source.golabs.io"}
	for i := 0; i < len(remotes); i++ {
		idx := strings.Index(source, remotes[i])
		if idx != -1 {
			return source[idx:], nil
		}
	}
	return "", &InvalidSourceError{
		msg: "invalid source",
	}
}

func parseSource(source string) (Dependency, error) {
	ret := Dependency{}
	url, err := findUrlFromSource(source)
	if err != nil {
		return ret, err
	}
	log.Println(url)
	tokens := strings.Split(url, "?")
	ret.Url = tokens[0]
	ret.Version = MakeVersion(tokens[1][4:])
	return ret, nil
}

func (d *dependencyParser) Parse(fileContent string) ([]Dependency, error) {
	lines := strings.Split(fileContent, "\n")
	deps := []Dependency{}
	isModule := false
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			continue
		}
		tmpLine := lines[i]
		tmpLine = strings.Trim(tmpLine, " ")
		tokens := strings.Fields(tmpLine)
		if tokens[0] == "module" {
			isModule = true
		}
		attr := tokens[0]
		if isModule && (attr == "source") {
			log.Println(tokens[0])
			log.Println(tokens[1])
			log.Println(tokens[2])
			log.Println(strings.Trim(tokens[2], "\""))
			buffer, err := parseSource(strings.Trim(tokens[2], "\""))
			if err != nil {
				return []Dependency{}, err
			}
			deps = append(deps, buffer)
			isModule = false
		}
	}
	return deps, nil
}
