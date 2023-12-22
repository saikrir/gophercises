package main

import (
	"bufio"
	"context"
	"time"

	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func readYorN(scanner *bufio.Scanner) string {

	for {
		fmt.Print("Ready? Y/N ")
		if scanner.Scan() {
			res := strings.ToUpper(scanner.Text())
			switch res {
			case "Y", "N":
				return res
			default:
				continue
			}
		}
	}
}

func evalAnswer(answer string, scanner *bufio.Scanner, evalResultChan chan<- bool) {

	if scanner.Scan() {
		userAnswer := strings.TrimSpace(scanner.Text())
		evalResultChan <- strings.EqualFold(userAnswer, answer)
	}

}

const timeout = 3

func main() {

	var (
		questionsFile *os.File
		err           error
	)

	if questionsFile, err = os.OpenFile("problems.csv", os.O_RDONLY, 0666); err != nil {
		log.Fatalln("failed to open problems file")
	}

	defer questionsFile.Close()

	fmt.Printf("Welcome to Quiz, You have %d seconds for each Answer \n", timeout)

	scanner := bufio.NewScanner(os.Stdin)

	problemsCsvReader := csv.NewReader(questionsFile)

	records, _ := problemsCsvReader.ReadAll()

	var nCorrectAnswers int

	ctx := context.Background()

	answerChan := make(chan bool)

outer:
	for _, questionAndAns := range records {

		question, answer := questionAndAns[0], questionAndAns[1]
		fmt.Printf("Question %s: \t", question)

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		go evalAnswer(answer, scanner, answerChan)

		select {
		case res := <-answerChan:
			if res {
				nCorrectAnswers++
			}
		case <-ctx.Done():
			fmt.Println("\nSorry, Time is up!")
			break outer
		}
	}

	fmt.Printf("You Have completed the quiz with %d correct answers out of %d \n", nCorrectAnswers, len(records))
}
