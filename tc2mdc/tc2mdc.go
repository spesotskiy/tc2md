package tc2mdc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// One line comment
const OLC string = "//"

// Link to the top of MD file
const TopLink string = "[top](#top)"

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
	var trimmedLine, convertedLine string
	var isConverted, isFuncStarted bool
	rePackage, _ := regexp.Compile(`^package\s(?P<name>\w+)`)
	reFunc, _ := regexp.Compile(`^func\s(?P<name>Test\w+)\(t \*testing\.T\)`)
	reMarker, _ := regexp.Compile(`^\s(#|##|>|-|--|---)\s[^\s]`) // all MD markers to search for
	for _, origLine := range comments {
		trimmedLine = strings.TrimSpace(origLine)
		switch {
		case strings.HasPrefix(origLine, "package"):
			{
				result := getMatchesMap(rePackage, origLine)
				convertedLine = "## `" + result["name"] + "`"
				isConverted = result["name"] != ""
				if isConverted {
					mdText = append(mdText, convertedLine)
				}
			}
		case strings.HasPrefix(origLine, "func"): // start of func
			{
				result := getMatchesMap(reFunc, origLine)
				convertedLine = "#### `" + result["name"] + "`"
				isConverted = result["name"] != ""
				if isConverted {
					mdText = append(mdText, "---")
					mdText = append(mdText, convertedLine)
					isFuncStarted = true
				}
			}
		case strings.HasPrefix(origLine, "}"): // end of func
			{
				if isFuncStarted {
					isFuncStarted = false
					mdText = append(mdText, "")
					mdText = append(mdText, TopLink)
				}
			}
		case strings.HasPrefix(trimmedLine, OLC):
			{
				convertedLine, isConverted = convertLineComment(trimmedLine, reMarker)
				if isConverted {
					mdText = append(mdText, convertedLine)
				}
			}
		default:
			isConverted = false
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
