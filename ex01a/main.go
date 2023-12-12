package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

const CSV_FILE_PATH = "./problems.csv"

func loadQuestions() [][]string {
	csvFile, err := os.Open(CSV_FILE_PATH)
	if err != nil {
		log.Fatalln("Failed to open file ", CSV_FILE_PATH, err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalln("Failed to read all rows from CSV File")
	}

	return records
}

func getUserInput(scanner *bufio.Scanner) (string, error) {
	if scanner.Scan() {
		return strings.Trim(scanner.Text(), " "), nil
	}
	return "", fmt.Errorf("failed to read input from user")
}

func main() {
	var records [][]string = loadQuestions()
	fmt.Printf("Welcome to Quiz, You have %d questions to answer\n Lets Start \n", len(records))
	var successCtr, failureCtr int
	scanner := bufio.NewScanner(os.Stdin)
	for i, record := range records {

	repeat:

		fmt.Printf("Question #%d: ", i+1)
		fmt.Printf("What is %s ? ", record[0])

		answer, err := getUserInput(scanner)
		switch {
		case err != nil:
			log.Fatalln(err)
		case answer == "":
			goto repeat
		case strings.EqualFold(answer, record[1]):
			successCtr++
		default:
			failureCtr++
		}
	}

	fmt.Printf("You have completed the quiz, your score is %d/%d \n", successCtr, len(records))
}
