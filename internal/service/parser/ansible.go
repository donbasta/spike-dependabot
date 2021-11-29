package parser

// func (d *dependencyParser) Parse(fileContent string) ([]Dependency, error) {
// 	lines := strings.Split(fileContent, "\n")
// 	deps := []Dependency{}
// 	buffer := Dependency{}
// 	for i := 1; i < len(lines); i++ {
// 		if len(lines[i]) == 0 {
// 			continue
// 		}
// 		tmpLine := lines[i]
// 		if lines[i][0] == '-' {
// 			buffer = Dependency{}
// 			tmpLine = tmpLine[1:]
// 		}
// 		tmpLine = strings.Trim(tmpLine, " ")
// 		tokens := strings.Split(tmpLine, ":")
// 		attr := tokens[0]
// 		if attr == "src" {
// 			buffer.Url = tokens[1]
// 		} else if attr == "version" {
// 			buffer.Version = MakeVersion(tokens[1])
// 			deps = append(deps, buffer)
// 		}
// 	}
// 	return deps, nil
// }
