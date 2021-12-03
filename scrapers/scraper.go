package scrapers

import (
	"scraper/models"
)

type SchoolScraper interface {
	StartSchoolScraper()
	GetSchool() models.School
}
