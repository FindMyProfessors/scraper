package valencia

import "github.com/FindMyProfessors/scraper/model"

type ValenciaScraper struct {
	term model.Term
}

func (u *ValenciaScraper) Scrape() (*model.School, error) {
	//TODO implement me
	panic("implement me")
}

func (u *ValenciaScraper) Terms() []model.Term {
	//TODO implement me
	panic("implement me")
}

func (u *ValenciaScraper) Name() string {
	return "Valencia College"
}

func (u *ValenciaScraper) SetTerm(term model.Term) {
	u.term = term
}
