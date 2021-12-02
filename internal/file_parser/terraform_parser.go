package parser

import "strings"

type TerraformParser struct {
}

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
	tokens := strings.Split(url, "?")
	u := tokens[0]
	if len(u) > 4 && u[len(u)-4:] == ".git" {
		u = u[:len(u)-4]
	}
	ret.Url = u
	if len(tokens) == 1 {
		ret.Version = MakeVersion("v0.0.0")
	} else {
		ret.Version = MakeVersion(tokens[1][4:])
	}
	return ret, nil
}

func (tf *TerraformParser) Parse(fileContent string) ([]Dependency, error) {
	lines := strings.Split(fileContent, "\n")
	deps := []Dependency{}
	isModule := false
	for i := 0; i < len(lines); i++ {
		tmpLine := strings.Trim(lines[i], " ")
		if len(tmpLine) == 0 {
			continue
		}
		tokens := strings.Fields(tmpLine)
		if tokens[0] == "module" {
			isModule = true
		}
		attr := tokens[0]
		if isModule && (attr == "source") {
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
