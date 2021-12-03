package schools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scraper/models"
	"scraper/scrapers"
	"scraper/scrapers/rate_my_professor"
	"strings"
	"sync"
)

type UCFScraper struct {
	school models.School
}

type UCFProfessor struct {
	CoursePrefix  string `json:"prefix"`
	CatalogNumber string `json:"catalog_number"`
	CourseTitle   string `json:"title"`
	NameFirst     string `json:"name_first"`
	NameLast      string `json:"name_last"`
}

func (U *UCFScraper) StartSchoolScraper() {
	modalities := []string{"WW", "V", "M", "RS", "P"}

	courseMap := make(map[string]models.Course, 0)
	professorMap := make(map[string]models.Professor, 0)

	var wg sync.WaitGroup
	for _, mode := range modalities {
		wg.Add(1)
		var url = fmt.Sprintf("https://cdl.ucf.edu/wp-content/themes/cdl/lib/course-search-ajax.php?call=classes&term=1730&prefix=&catalog=&title=&instructor=&career=&college=&department=&mode=%s&_=1631687421296", mode)
		fmt.Println("url=" + url)
		go func(url string) {
			response, err := http.Get(url)
			if err != nil {
				fmt.Println("error occurred: ", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println("error occurred: ", err)
				}
			}(response.Body)
			body, err := io.ReadAll(response.Body)

			var data struct {
				Professors []UCFProfessor `json:"classes"`
			}
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("error occurred: ", err)
			}

			for _, ucfProfessor := range data.Professors {
				ucfProfessor.NameFirst = scrapers.WashName(ucfProfessor.NameFirst)
				ucfProfessor.NameLast = scrapers.WashName(ucfProfessor.NameLast)

				course := courseMap[ucfProfessor.CoursePrefix+ucfProfessor.CatalogNumber]
				if course.Code == "" {
					course = models.Course{
						Code: ucfProfessor.CoursePrefix + ucfProfessor.CatalogNumber,
						Name: strings.ReplaceAll(ucfProfessor.CourseTitle, "\"", ""),
					}
				}
				professor := professorMap[ucfProfessor.NameFirst+"_"+ucfProfessor.NameLast]
				if professor.FirstName == "" {
					professor = models.Professor{
						FirstName: strings.ReplaceAll(ucfProfessor.NameFirst, "/", ""),
						LastName:  strings.ReplaceAll(ucfProfessor.NameLast, "/", ""),
						Teaches:   []models.Course{},
					}
				}
				doesTeachThisClass := func() bool {
					for _, c := range professor.Teaches {
						if c.Code == course.Code {
							return true
						}
					}
					return false
				}()
				if !doesTeachThisClass {
					professor.Teaches = append(professor.Teaches, course)
					course.Professors = append(course.Professors, professor.RDFId())
				}

				professorMap[ucfProfessor.NameFirst+"_"+ucfProfessor.NameLast] = professor
				courseMap[ucfProfessor.CoursePrefix+ucfProfessor.CatalogNumber] = course
			}

			wg.Done()
		}(url)
	}
	wg.Wait()

	for _, professor := range professorMap {
		U.school.AddProfessor(professor)
		log.Println("Adding Professor to map, ", professor)
	}

	for _, course := range courseMap {
		U.school.AddCourse(course)
		log.Println("Adding Course to map, ", course)
	}
}

func (U *UCFScraper) GetSchool() models.School {
	if U.school.Name != "" {
		return U.school
	}
	U.school = models.School{
		Name:       "University of Central Florida",
		Professors: []models.Professor{},
		Courses:    []models.Course{},
	}
	U.StartSchoolScraper()

	U.school = rate_my_professor.Scrape(&U.school, 1082, 15516, 5567, 5400)

	return U.school
}
