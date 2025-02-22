package tc2mdc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// One line comment
const OLC string = "//"

// Converts multi line text with one line comments into a MarkDown text
func Convert(comments []string) ([]string, error) {

	if len(comments) == 0 {
		return nil, errors.New("nil input")
	}
	if len(comments) == 1 && comments[0] == "" {
		return nil, errors.New("empty input")
	}

	fmt.Println("Start converting...")
	var mdText []string
	var convertedLine string
	var isConverted bool
	rePackage, _ := regexp.Compile(`^package\s(?P<name>\w+)`)
	reFunc, _ := regexp.Compile(`^func\s(?P<name>Test\w+)\(t \*testing\.T\)`)
	reMarker, _ := regexp.Compile(`^\s(#|##|>|-|--|---)\s[^\s]`) // all MD markers to search for
	for _, line := range comments {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "package"):
			{
				result := getMatchesMap(rePackage, line)
				convertedLine = "## `" + result["name"] + "`"
				isConverted = result["name"] != ""
			}
		case strings.HasPrefix(line, "func"):
			{
				result := getMatchesMap(reFunc, line)
				convertedLine = "#### `" + result["name"] + "`"
				isConverted = result["name"] != ""
			}
		case strings.HasPrefix(line, OLC):
			{
				convertedLine, isConverted = convertLineComment(line, reMarker)
			}
		default:
			isConverted = false
		}
		if isConverted {
			mdText = append(mdText, convertedLine)
		}
	}
	return mdText, nil
}

func getMatchesMap(re *regexp.Regexp, line string) map[string]string {
	header := re.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = header[i]
		}
	}
	return result
}

func convertLineComment(line string, re *regexp.Regexp) (string, bool) {
	line = line[2:] // trim OLC
	if re.MatchString(line) {
		marker := re.FindStringSubmatch(line)
		line = line[1:] // trim a leading space
		switch marker[1] {
		case "#", "##":
			{
				line = "##" + line
			}
		case "-", ">":
			{
				// line as is
			}
		case "--":
			{
				line = "  " + line[1:]
			}
		case "---":
			{
				line = "    " + line[2:]
			}
		}
		return line, true
	}
	return "", false
}
