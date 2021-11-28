package scrapers

import (
	"scraper/models"
)

type SchoolScraper interface {
	StartRMSScraper()
	StartSchoolScraper()
	GetSchool() models.School
}
