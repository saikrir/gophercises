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

func evalAnswer(answer string, userAnswer string, evalResultChan chan<- bool) {
	evalResultChan <- strings.EqualFold(userAnswer, answer)
}

func main() {

	var (
		questionsFile *os.File
		err           error
	)

	if questionsFile, err = os.OpenFile("problems.csv", os.O_RDONLY, 0666); err != nil {
		log.Fatalln("failed to open problems file")
	}

	defer questionsFile.Close()

	fmt.Printf("%10s\n", "Welcome to Quiz ")

	scanner := bufio.NewScanner(os.Stdin)

	problemsCsvReader := csv.NewReader(questionsFile)

	records, _ := problemsCsvReader.ReadAll()

	var nCorrectAnswers int
	ctx := context.Background()
	answerChan := make(chan bool, len(records))

	for _, questionAndAns := range records {

		question, answer := questionAndAns[0], questionAndAns[1]

	readInput:
		fmt.Printf("Question: %s \n", question)
		fmt.Print("Press Any Key to Start the Timer ....\n")
		if scanner.Scan() {
			fmt.Println()
			fmt.Printf("Answer \t:")
		}

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if scanner.Scan() {
			userAnswer := strings.TrimSpace(scanner.Text())
			if userAnswer == "" {
				goto readInput
			}
			go evalAnswer(answer, userAnswer, answerChan)
		}

		select {
		case res := <-answerChan:
			if res {
				nCorrectAnswers++
			}
		case <-ctx.Done():
			log.Fatalln("Timeup!")
		}
	}

	fmt.Printf("You Have completed the quiz with %d correct answers out of %d \n", nCorrectAnswers, len(records))
}
