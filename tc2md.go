package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"tc2mdc"
)

func main() {
	fmt.Println("Convert test comments to a MD file.")

	testFiles := []string{
		"./tc2mdc/tc2mdparser_test.go",
		"./tc2mdc/tc2mdwriter_test.go",
	}

	for i, testFile := range testFiles {
		code := readTestFile(testFile)

		testData, err := tc2mdc.Parse(code)
		if err != nil {
			log.Fatal(err)
		}

		mdText := tc2mdc.Write(testData)
		saveToMDFile("scenario"+strconv.Itoa(i)+".md", mdText)
	}
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
