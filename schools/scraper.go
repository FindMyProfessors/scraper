package schools

import "github.com/FindMyProfessors/scraper/model"

type SchoolScraper interface {
	Scrape() (*model.School, error)
	Name() string
	Terms() []model.Term
	SetTerm(term model.Term)
}
