package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

type Hand struct {
	Position string
	URL      string
	Win      bool
	Holding  map[string]string
	Pretty   string
}

func main() {
	fileContent, err := ioutil.ReadFile("./cardschat.txt")
	if err != nil {
		fmt.Printf("Error getting file: %s", err.Error())
	}
	re := regexp.MustCompile(`http.*`)
	urls := re.FindAll(fileContent, -1)
	allUrls := []string{}
	for _, url := range urls {
		allUrls = append(allUrls, string(url))
	}

	re = regexp.MustCompile(`Hero \(Hero\) is ([A-Z]*\+?\d?)`)
	outer := re.FindAllSubmatch(fileContent, -1)
	allPositions := []string{}
	for _, positions := range outer {
		position := positions[1]
		allPositions = append(allPositions, string(position))
	}
	heroHands := []Hand{}

	allHands := strings.Split(string(fileContent), "\n\n♦ ♣ ♥ ♠\n\n")
	allWins := []bool{}

	allHoldings := map[int]map[string]string{}
	for index, hand := range allHands {
		re := regexp.MustCompile(`Hero \(([A-Z]*\+?\d?)\) wins`)
		wins := re.FindAll([]byte(hand), -1)
		allWins = append(allWins, len(wins) > 0)

		re = regexp.MustCompile(`Pre-Flop: \((\d|,)*\) Hero \(Hero\) is [A-Z]*\+?\d? with (\w|\d)(♠|♣|♥|♦) (\w|\d)(♠|♣|♥|♦)`)
		holdings := re.FindAllSubmatch([]byte(hand), -1)
		for _, holdings := range holdings {
			firstCard := string(holdings[2])
			firstSuit := string(holdings[3])
			secondCard := string(holdings[4])
			secondSuit := string(holdings[5])
			allHoldings[index] = map[string]string{"firstCard": firstCard, "firstSuit": firstSuit, "secondCard": secondCard, "secondSuit": secondSuit}
		}

	}

	// allMadeFlops:=[]bool{}
	// Pre-Flop: \((\d|,)*\) Hero \(Hero\) is ([A-Z]*\+?\d?) with (\w|\d)(♠|♣|♥|♦)

	wins := 0
	for index, url := range allUrls {
		hand := Hand{}
		hand.URL = url
		hand.Position = allPositions[index]
		if allWins[index] == true {
			wins++
			hand.Win = true
		} else {
			hand.Win = false
		}
		hand.Holding = allHoldings[index]
		hand.Pretty = allHoldings[index]["firstCard"] + allHoldings[index]["firstSuit"] + " " + allHoldings[index]["secondCard"] + allHoldings[index]["secondSuit"]
		heroHands = append(heroHands, hand)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Printf("\n\n\n Hands Won: %v  Hands Played: %v  Win%%: %v\n", wins, len(allWins), float64(wins)/float64(len(allWins))*100)

	fmt.Fprintf(w, "\n\n %s\t%s\t%s\t%s\t%s\t", "Hand Number", "Holding", "Victorious", "Position ", "URL")
	fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t", "-------------", "-------", "----------", "--------", "----")

	for handNumber, hand := range heroHands {
		fmt.Fprintf(w, "\n %v\t%v\t%v\t%v\t%v\t", handNumber, hand.Pretty, hand.Win, hand.Position, hand.URL)
		fmt.Fprintf(w, "\n %s\t%s\t%s\t%s\t%s\t", "-------------", "-------", "----------", "--------", "----")
	}
}
