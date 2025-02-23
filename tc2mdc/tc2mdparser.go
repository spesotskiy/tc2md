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
	fmt.Printf("Parsed package %v with %d methods.", testData.packageName, len(testData.methods))
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
				for _, tag := range strings.Split(line[2:], ",") {
					testMethod.tags = append(testMethod.tags, strings.TrimSpace(tag))
				}
			}
		case "##":
			{
				testMethod.steps = append(testMethod.steps, TestStep{GWT, strings.TrimSpace(line[3:])})
			}
		case "-", "--", "---":
			{
				testMethod.steps = append(testMethod.steps, TestStep{len(marker[1]), strings.TrimSpace(line[len(marker[1]):])})
			}
		}
	}
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
