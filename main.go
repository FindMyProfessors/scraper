package main

import (
	"encoding/json"
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"github.com/FindMyProfessors/scraper/schools"
	"github.com/FindMyProfessors/scraper/schools/ucf"
	"github.com/FindMyProfessors/scraper/schools/valencia"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strconv"
	"strings"
)

var (
	SchoolScraperMap = []schools.SchoolScraper{&ucf.Scraper{}, &valencia.ValenciaScraper{}}
)

func main() {
	fmt.Println("Welcome to the FindMyProfessor's Scraper")

	scraper := GetSchoolScraper()

	fmt.Printf("Starting scraping of %s now\n", scraper.Name())
	term := GetTerm(scraper)
	scraper.SetTerm(term)
	school, err := scraper.Scrape()
	if err != nil {
		panic(err)
	}
	if ShouldWriteToFile() {
		fileName := fmt.Sprintf("%s-%d-%s.json", school.Name, term.Year, term.Semester)
		fmt.Printf("Writing the scraped school data to file %s\n", fileName)

		marshal, err := json.Marshal(*school)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(fileName, marshal, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Wrote %s to %s\n", school.Name, fileName)
	}

	if ShouldSendToFMP() {

	}
}

func GetSchoolScraper() schools.SchoolScraper {
	for i, scraper := range SchoolScraperMap {
		fmt.Printf("%d - %s\n", i, scraper.Name())
	}
	for {
		fmt.Printf("Which school would you like to scrape? ")

		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Printf("Unable to scan choice: %v\n", err)
			continue
		}

		if choice > len(SchoolScraperMap) {
			fmt.Printf("%d is an invalid selection, please try again\n\n", choice)
			continue
		}
		return SchoolScraperMap[choice]
	}
}

func GetTerm(scraper schools.SchoolScraper) model.Term {
	for {
		terms := scraper.Terms()

		for i, term := range terms {
			fmt.Printf("%d - %s %s\n", i, strconv.Itoa(term.Year), cases.Title(language.English).String(strings.ToLower(string(term.Semester))))
		}

		fmt.Printf("Which term would you like to scrape? ")

		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Printf("Unable to scan choice: %v\n", err)
			continue
		}

		if choice > len(terms) {
			fmt.Printf("%d is an invalid selection, please try again\n\n", choice)
			continue
		}
		return terms[choice]
	}
}

func ShouldWriteToFile() bool {
	for {
		fmt.Printf("Would you like to write the data to a file (y/n)? ")

		var choice string
		_, err := fmt.Scanf("%s", &choice)
		if err != nil {
			fmt.Printf("Unable to scan choice: %v\n", err)
			continue
		}

		switch strings.ToLower(choice) {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Printf("%s is an invalid option. Choose 'y' or 'n'.\n", choice)
			continue
		}
	}
}

func ShouldSendToFMP() bool {
	for {
		fmt.Printf("Would you like to send the data to FMP (y/n)? ")

		var choice string
		_, err := fmt.Scanf("%s", &choice)
		if err != nil {
			fmt.Printf("Unable to scan choice: %v\n", err)
			continue
		}

		switch strings.ToLower(choice) {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Printf("%s is an invalid option. Choose 'y' or 'n'.\n", choice)
			continue
		}
	}
}
