package rate_my_professor

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"scraper/models"
	"strconv"
	"strings"
	"sync"
)

type RMPProfessor struct {
	FirstName     string `json:"tFname" bson:"firstName"`
	MiddleName    string `json:"tMiddlename" bson:"middleName"`
	LastName      string `json:"tLname" bson:"lastName"`
	RatingsCount  int    `json:"tNumRatings" bson:"ratingCount"`
	RatingClass   string `json:"rating_class" bson:"ratingClass"`
	OverallRating string `json:"overall_rating" bson:"overallRating"`
}

func Scrape(school *models.School, schoolIDs ...int) models.School {
	var scrapeFunction func(page int) bool
	type Model struct {
		Professors []RMPProfessor `json:"professors"`
		Total      int            `json:"searchResultsTotal"`
		Remaining  int            `json:"remaining"`
		Type       string         `json:"type"`
	}
	completeModel := Model{}

	var wg sync.WaitGroup
	wg.Add(len(schoolIDs))
	for _, id := range schoolIDs {
		go func(id int) {
			scrapeFunction = func(page int) bool {
				url := fmt.Sprintf("https://www.ratemyprofessors.com/filter/professor/?&page=%d&filter=teacherlastname_sort_s+asc&query=**&queryoption=TEACHER&queryBy=schoolId&sid=%d", page, id)
				response, err := http.Get(url)
				if err != nil {
					log.Fatalf("unable to make request: %s\n", err)
				}
				data, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("unable to read request: %s\n", err)
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						fmt.Println("error occurred: ", err)
					}
				}(response.Body)
				var model Model
				err = json.Unmarshal(data, &model)
				if err != nil {
					log.Fatalf("unable unmarshal data: %s\n", err)
				}
				for _, professor := range model.Professors {
					completeModel.Professors = append(completeModel.Professors, professor)
				}
				if model.Remaining == 0 {
					return false
				}
				return scrapeFunction(page + 1)
			}
			scrapeFunction(1)
			wg.Done()
		}(id)
	}
	wg.Wait()

	isSameProfessor := func(professor models.Professor, rmsProfessor struct {
		FirstName     string `json:"tFname" bson:"firstName"`
		MiddleName    string `json:"tMiddlename" bson:"middleName"`
		LastName      string `json:"tLname" bson:"lastName"`
		RatingsCount  int    `json:"tNumRatings" bson:"ratingCount"`
		RatingClass   string `json:"rating_class" bson:"ratingClass"`
		OverallRating string `json:"overall_rating" bson:"overallRating"`
	}) bool {
		if professor.FirstName == rmsProfessor.FirstName {
			if professor.LastName == rmsProfessor.LastName {
				return true
			}
		}

		tryMultiname := func(a string, b string, deliminator string) bool {
			aNames := strings.Split(a, deliminator)
			bNames := strings.Split(b, deliminator)
			if len(aNames) == 1 && len(bNames) == 1 {
				return a == b
			}
			for _, aa := range aNames {
				for _, bb := range bNames {
					if aa == bb {
						return true
					}
				}
			}
			return false
		}

		if tryMultiname(professor.FirstName, rmsProfessor.FirstName, "-") && tryMultiname(professor.LastName, rmsProfessor.LastName, "-") {
			return true
		}
		if len(rmsProfessor.MiddleName) > 0 {
			if tryMultiname(professor.FirstName, rmsProfessor.MiddleName, " ") && tryMultiname(professor.LastName, rmsProfessor.MiddleName, " ") {
				return true
			}
		}
		return false
	}

	var wg2 sync.WaitGroup
	wg2.Add(len(completeModel.Professors))
	for _, anonProfessor := range completeModel.Professors {
		complete := false
		go func(anonProfessor RMPProfessor) {
			for index, professor := range school.Professors {
				if complete {
					break
				}
				if isSameProfessor(professor, anonProfessor) {
					var rating float64
					if anonProfessor.OverallRating == "N/A" {
						rating = 0.0
					} else {
						float, err := strconv.ParseFloat(anonProfessor.OverallRating, 64)
						if err != nil {
							log.Fatalln("error occurred ", err)
						} else {
							rating = float * float64(anonProfessor.RatingsCount)
						}
					}
					professor.Rating += rating
					professor.TotalRatings += anonProfessor.RatingsCount
					school.SetProfessor(index, professor)
					complete = true
				}
			}
			wg2.Done()
		}(anonProfessor)
	}
	wg2.Wait()
	return *school
}
