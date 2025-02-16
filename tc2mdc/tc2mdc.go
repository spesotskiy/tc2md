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
	re, _ := regexp.Compile(`^\s(#|##|>|-|--|---)\s[^\s]`) // all MD markers to search for
	for _, line := range comments {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, OLC) {
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
				mdText = append(mdText, line)
			}
		}
	}
	return mdText, nil
}
