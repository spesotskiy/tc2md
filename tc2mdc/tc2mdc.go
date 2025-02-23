package tc2mdc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// One line comment
const OLC string = "//"
const (
	GWT       = 0
	common    = 1
	indented  = 2
	indented2 = 3
)

type TestStep struct {
	kind    int
	comment string
}

type TestMethod struct {
	name     string
	tags     []string
	scenario string
	steps    []TestStep
}

type TOCLine struct {
	index   int
	caption string
	link    string
	gitLink string
}

type TestData struct {
	title       string
	packageName string
	toc         map[string]TOCLine
	methods     []TestMethod
}

func ConvertDataToMD(data *TestData) []string {

	return nil
}

func Parse(codeLines []string) (*TestData, error) {
	errorMessage := isInputEmpty(&codeLines)
	if errorMessage != "" {
		return nil, errors.New(errorMessage)
	}
	fmt.Println("Start parsing...")

	var testData = new(TestData)
	var isFuncStarted bool

	rePackage, _ := regexp.Compile(`^package\s(?P<name>\w+)`)
	reFunc, _ := regexp.Compile(`^func\s(?P<name>Test\w+)\(t \*testing\.T\)`)
	reMarker, _ := regexp.Compile(`^\s(#|##|>|-|--|---)\s[^\s]`) // all MD markers to search for
	for _, origLine := range codeLines {
		trimmedLine := strings.TrimSpace(origLine)
		switch {
		case strings.HasPrefix(origLine, "package"):
			{
				parsePackageHeader(origLine, rePackage, testData)
			}
		case strings.HasPrefix(origLine, "func"): // start of func
			{
				isFuncStarted = parseFunc(origLine, reFunc, testData)
			}
		case strings.HasPrefix(trimmedLine, OLC):
			{
				if isFuncStarted {
					parseOneLineComment(trimmedLine, reMarker, &(testData.methods[len(testData.methods)-1]))
				}
			}
		case strings.HasPrefix(origLine, "}"): // end of func
			{
				isFuncStarted = false
			}
		}
	}
	if testData != nil {
		fmt.Printf("Parsed package %v with %d methods.", testData.packageName, len(testData.methods))
	}
	return testData, nil
}

func parseOneLineComment(line string, re *regexp.Regexp, testMethod *TestMethod) {
	line = line[len(OLC):] // trim OLC
	if re.MatchString(line) {
		marker := re.FindStringSubmatch(line)
		line = line[1:] // trim a leading space
		switch marker[1] {
		case "#":
			{
				testMethod.scenario = line[2:]
			}
		case ">":
			{
				testMethod.tags = splitTags(line)
			}
		case "##":
			{
				testMethod.steps = append(testMethod.steps, TestStep{GWT, line})
			}
		case "-", "--", "---":
			{
				testMethod.steps = append(testMethod.steps, TestStep{len(marker[1]), line})
			}
		}
	}
}

func splitTags(line string) []string {
	return []string{line}
}

func parseFunc(origLine string, reFunc *regexp.Regexp, testData *TestData) bool {
	result := getMatchesMap(reFunc, origLine)
	funcName := result["name"]
	if funcName != "" {
		method := TestMethod{name: funcName}
		(*testData).methods = append((*testData).methods, method)
		return true
	}
	return false
}

func parsePackageHeader(origLine string, rePackage *regexp.Regexp, testData *TestData) {
	result := getMatchesMap(rePackage, origLine)
	packageName := result["name"]
	if packageName != "" {
		(*testData).packageName = packageName
	}
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
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
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
