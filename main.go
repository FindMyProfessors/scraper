package main

import (
	"fmt"
	"log"
	"os"
	"scraper/database"
	"scraper/models"
	"scraper/scrapers"
	"scraper/scrapers/schools"
	"strings"
	"sync"
)

func main() {
	_ = database.Connect(os.Getenv("DATABASE_URL"), os.Getenv("API_KEY"))

	split := strings.Split(strings.ToLower(os.Getenv("SCHOOLS")), ",")
	for _, schoolName := range split {
		var scraper scrapers.SchoolScraper
		log.Println("Scraping ", schoolName)
		switch strings.TrimSpace(schoolName) {
		case "ucf":
			scraper = schools.UCFScraper{}
			break
		case "valencia":
			scraper = &schools.ValenciaScraper{}
			break
		}
		var wg sync.WaitGroup
		wg.Add(len(split))
		go func(school models.School) {
			_, err := database.MutateDatabase(school)
			if err != nil {
				log.Println("Unable to mutate database. ", err)
			}
			fmt.Println("Mutated database for " + school.Name)
			wg.Done()
		}(scraper.GetSchool())
		wg.Wait()
	}
}
