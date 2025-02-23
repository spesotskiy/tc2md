package tc2mdc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// One line comment
const OLC string = "//"

type TestMethod struct {
	name     string
	tags     []string
	scenario string
	steps    []string
}
type TOCLine struct {
	index   int
	caption string
	link    string
	gitLink string
}

type MDFile struct {
	title       string
	packageName string
	toc         map[string]TOCLine
	methods     []TestMethod
}

// Converts multi line text with one line comments into a MarkDown text
func Convert(comments []string) ([]string, []string, error) {
	errorMessage := isInputEmpty(&comments)
	if errorMessage != "" {
		return nil, nil, errors.New(errorMessage)
	}
	fmt.Println("Start converting...")

	var packageName string
	var mdText, toc []string
	var isFuncStarted bool

	rePackage, _ := regexp.Compile(`^package\s(?P<name>\w+)`)
	reFunc, _ := regexp.Compile(`^func\s(?P<name>Test\w+)\(t \*testing\.T\)`)
	reMarker, _ := regexp.Compile(`^\s(#|##|>|-|--|---)\s[^\s]`) // all MD markers to search for
	for _, origLine := range comments {
		trimmedLine := strings.TrimSpace(origLine)
		switch {
		case strings.HasPrefix(origLine, "package"):
			{
				packageName = addPackageHeader(origLine, rePackage, &mdText)
			}
		case strings.HasPrefix(origLine, "func"): // start of func
			{
				isFuncStarted = addFuncHeader(origLine, reFunc, &mdText)
			}
		case strings.HasPrefix(origLine, "}"): // end of func
			{
				addTopLinkToFuncEnd(isFuncStarted, packageName, &mdText)
			}
		case strings.HasPrefix(trimmedLine, OLC):
			{
				addOneLineComment(trimmedLine, reMarker, &mdText)
			}
		}
	}
	return mdText, toc, nil
}

func isInputEmpty(comments *[]string) string {
	if len(*comments) == 0 {
		return "nil input"
	}
	if len(*comments) == 1 && (*comments)[0] == "" {
		return "empty input"
	}
	return ""
}

func addPackageHeader(origLine string, rePackage *regexp.Regexp, mdText *[]string) string {
	result := getMatchesMap(rePackage, origLine)
	packageName := result["name"]
	convertedLine := "## `" + packageName + "`"
	if packageName != "" {
		*mdText = append(*mdText, convertedLine)
	}
	return packageName
}

func addFuncHeader(origLine string, reFunc *regexp.Regexp, mdText *[]string) bool {
	result := getMatchesMap(reFunc, origLine)
	convertedLine := "#### `" + result["name"] + "`"
	if result["name"] != "" {
		*mdText = append(*mdText, "---")
		*mdText = append(*mdText, convertedLine)
		return true
	}
	return false
}

func addOneLineComment(trimmedLine string, reMarker *regexp.Regexp, mdText *[]string) {
	convertedLine, isConverted := convertLineComment(trimmedLine, reMarker)
	if isConverted {
		*mdText = append(*mdText, convertedLine)
	}
}

func addTopLinkToFuncEnd(isFuncStarted bool, packageName string, mdText *[]string) []string {
	if isFuncStarted {
		isFuncStarted = false
		*mdText = append(*mdText, "")
		*mdText = append(*mdText, getLinkToTop(packageName))
	}
	return *mdText
}

func getLinkToTop(name string) string {
	if name == "" {
		name = "top"
	}
	return "[top](#" + name + ")"
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
