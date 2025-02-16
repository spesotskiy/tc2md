package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"tc2mdc"
)

func main() {
	fmt.Println("Convert test comments to a MD file.")

	code := readTestFile("./tc2mdc/tc2mdc_test.go")

	mdText, err := tc2mdc.Convert(code)
	if err != nil {
		log.Fatal(err)
	}

	saveToMDFile("scenario.md", mdText)
}

func readTestFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var code []string
	for scanner.Scan() {
		code = append(code, scanner.Text())
	}
	return code
}

func saveToMDFile(path string, mdText []string) {
	mdFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer mdFile.Close()

	for _, line := range mdText {
		_, err := mdFile.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
