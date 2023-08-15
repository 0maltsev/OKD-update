package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	//excelCreator("test_output.csv")
	regexer("test_output.txt", "test_output.csv")
}

type TestOutput struct {
	testArea        string
	testDescription string
	additionalTags  string
	testStatus      string
	testDestination string
}

func regexer(testOutputFileName string, tableFileName string) {
	file, _ := os.Open(testOutputFileName)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	records := [][]string{
		{"test_area", "description", "additional_tags", "test_destination"},
	}

	f, err := os.Create(tableFileName)
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	counter := 0
	for scanner.Scan() {
		testOutput := TestOutput{}
		testOutput.testStatus = statusTestAggregator(scanner.Text())
		if testOutput.testStatus != "ERROR" {
			testOutput.testDescription = descriptionTestAggregator(scanner.Text())
			allTagsArray := allTagsAggregator(scanner.Text())

			additionalTagsArray := allTagsArray[1:]

			for _, tag := range additionalTagsArray {
				testOutput.additionalTags = testOutput.additionalTags + tag + " "
			}
			testOutput.testArea = allTagsArray[0]
			testOutput.testDestination = allTagsArray[len(allTagsArray)-1]

			testOutput.testDescription = descriptionTestAggregator(scanner.Text())
			for _, tag := range allTagsArray {
				testOutput.testDescription = strings.Replace(testOutput.testDescription, tag, "", -1)
			}

			records := [][]string{
				{testOutput.testArea, testOutput.testDescription, testOutput.additionalTags, testOutput.testDestination},
			}

			for _, record := range records {
				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}
			}
			counter = counter + 1
		}
	}

	println(counter)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func statusTestAggregator(fileLine string) string {

	matchedSkipped, _ := regexp.MatchString(`^skipped:`, fileLine)
	if matchedSkipped == true {
		return "skipped"
	}
	matchedPassed, _ := regexp.MatchString(`^passed:`, fileLine)
	if matchedPassed == true {
		return "passed"
	}
	matchedFailed, _ := regexp.MatchString(`^failed:`, fileLine)
	if matchedFailed == true {
		return "failed"
	}

	return "ERROR"
}

func descriptionTestAggregator(fileLine string) string {
	re, _ := regexp.Compile(`\"(.*?)\"`)
	res := re.FindAllString(fileLine, -1)[0]
	return res
}

func allTagsAggregator(fileLine string) []string {
	re, _ := regexp.Compile(`\[(.*?)\]`)
	res := re.FindAllString(fileLine, -1)
	return res
}

func excelCreator(tableFileName string) {

	records := [][]string{
		{"test_area", "description", "additional_tags", "test_status", "test_destination"},
	}

	f, err := os.Create(tableFileName)
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
