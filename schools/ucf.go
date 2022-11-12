package schools

import (
	"github.com/FindMyProfessors/scraper/model"
	"strconv"
	"time"
)

type UCFScraper struct {
	term model.Term
}

func (u *UCFScraper) Scrape() *model.School {
	return nil
}

func (u *UCFScraper) Terms() (terms []model.Term) {
	// SPRING 2021  = 1710
	// FALL 2022    = 1760

	// How are UCF Term IDs calculated?
	// UCF was founded June 10th 1963
	// SPRING 2023 = 1770
	// Using that ID we can work in reverse
	// 1770/10 = 177, divide by 10 because the ID is incremented by 10 every semester
	// 3 semesters per year so 177/3 = 59
	// 2022-59 = 1963
	// UCF was founded in 1963
	currentYear := time.Now().Year() + 1
	// TODO: Only add next semester after, determine current semester and next

	for i := 0; i < 2; i++ {
		year := currentYear - i
		yearDifference := year - 1963

		semesterIndex := 2

		for i := 1; i < 4; i++ {
			semesterDifference := (yearDifference * 3) - i
			termId := semesterDifference * 10

			terms = append(terms, model.Term{
				Year:     year,
				Semester: model.AllSemesters[semesterIndex],
				ID:       strconv.Itoa(termId),
			})

			semesterIndex--
		}

	}

	return terms
}

func (u *UCFScraper) Name() string {
	return "University of Central Florida"
}

func (u *UCFScraper) SetTerm(term model.Term) {
	u.term = term
}
