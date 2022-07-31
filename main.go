package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	// read problems from quiz.csv file
	if fObj, err := os.Open(fileName); err == nil {
		// open file, create reader
		csvReader := csv.NewReader(fObj)
		// read file
		if cLine, err := csvReader.ReadAll(); err == nil {
			//use parseProblem to parse each line
			// return problems and error

			return parseProblem(cLine), nil
		} else {
			return nil, fmt.Errorf("Error in parsing csv  file %s: %v", fileName, err.Error())
		}

	} else {
		return nil, fmt.Errorf("Error in opening file %s: %v", fileName, err.Error())
	}

}

func main() {
	//input name of the file
	fName := flag.String("file", "quiz.csv", "path of the file")
	// set duration of timer
	timer := flag.Int("timer", 30, "timer for the quiz")
	flag.Parse()
	// pull questions from fileName with problem puller
	problems, err := problemPuller(*fName)
	// handle error
	if err != nil {
		exit(fmt.Sprintf("Something Went Wrong: %s", err.Error()))
	}
	// count correct ansers
	correctAnswers := 0
	// initialize timer
	timeObject := time.NewTimer(time.Duration(*timer) * time.Second)
	answersChan := make(chan string)
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = ", i+1, p.question)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answersChan <- answer
		}()
		select {
		case <-timeObject.C:
			fmt.Println("\nTime is up!")
			break problemLoop
		case iAns := <-answersChan:
			if iAns == p.answer {
				correctAnswers++
			}
			if i == len(problems)-1 {
				close(answersChan)
			}
		}
	}
	fmt.Printf("Your results were %d out of %d problems \n", correctAnswers, len(problems))
	fmt.Println("Press Enter to Quit")
	// quit the program
	<-answersChan
}

func parseProblem(lines [][]string) []problem {
	//parser each line of questions and parse into problem
	// return problem and error
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{question: lines[i][0], answer: lines[i][1]}
	}
	return r

}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
