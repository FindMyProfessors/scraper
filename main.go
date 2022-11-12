package main

import (
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"github.com/FindMyProfessors/scraper/schools"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

var (
	SchoolScraperMap = []schools.SchoolScraper{&schools.UCFScraper{}, &schools.ValenciaScraper{}}
)

func main() {
	fmt.Println("Welcome to the FindMyProfessor's Scraper")

	scraper := GetSchoolScraper()

	fmt.Printf("Starting scraping of %s now\n", scraper.Name())
	scraper.SetTerm(GetTerm(scraper))
	scraper.Scrape()
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
