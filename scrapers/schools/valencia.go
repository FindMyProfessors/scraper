package schools

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"scraper/models"
	"scraper/scrapers/rate_my_professor"
	"strings"
)

type ValenciaScraper struct {
	school models.School
}

func (U *ValenciaScraper) StartSchoolScraper() {
	log.Println("Started scraping valencia website")
	request, err := http.NewRequest(
		"GET",
		"https://valenciacollege.edu/academics/schedule-search/search.php?term=202210&INSM%5B%5D=N&INSM%5B%5D=M&INSM%5B%5D=X&INSM%5B%5D=S&Campus=&course_desc=&crn=",
		nil,
	)
	if err != nil {
		log.Fatalln("error occurred: ", err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln("error occurred: ", err)
	}
	bytes, err := io.ReadAll(response.Body)
	log.Println("Opening goquery document to html webscrape")
	if err != nil {
		log.Fatal(err)
	}
	s := string(bytes)
	//log.Println("s=", s)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}

	//log.Println("doc.Text()=", doc.Text())

	log.Println("Finding 'tbody' selector")
	table := doc.Find("tbody")
	//log.Println("table.Text()=", table.Text())

	rowCursor := 0

	professorMap := make(map[string]models.Professor)
	courseMap := make(map[string]models.Course)

	log.Println("Iterating over 'tr' selector")
	table.Find("tr").Each(func(i int, rowSelection *goquery.Selection) {
		log.Println("rowCursor=", rowCursor)
		if rowCursor%3 == 0 {
			cellCursor := 0
			var course models.Course
			var professor models.Professor
			rowSelection.Find("td").Each(func(i int, selection *goquery.Selection) {
				switch cellCursor {
				case 0:
					courseCode := selection.Find("strong").Text()
					courseName := strings.Replace(selection.Text()[len(courseCode):], "\\", "/", -1)
					log.Println("courseCode=", courseCode)
					log.Println("courseName=", courseName)
					course = courseMap[courseCode]
					if course.Code == "" {
						course = models.Course{
							Code: courseCode,
							Name: courseName,
						}
					}
					break
				case 3:
					professorText := selection.Text()
					log.Println("professorText=", professorText)
					var firstName string
					var lastName string
					if strings.Contains(professorText, "STAFF") {
						firstName = "Staff"
						lastName = ""
					} else {
						split := strings.Split(professorText, " ")
						firstName = split[0]
						lastName = split[1]
					}

					professor = professorMap[professorText]
					if professor.FirstName == "" {
						professor = models.Professor{
							FirstName: firstName,
							LastName:  lastName,
							Teaches:   []models.Course{},
						}
					}
					doesTeachThisClass := func(courses []models.Course) bool {
						for _, c := range professor.Teaches {
							if c.Code == course.Code {
								return true
							}
						}
						return false
					}(professor.Teaches)
					if !doesTeachThisClass {
						professor.Teaches = append(professor.Teaches, course)
						course.Professors = append(course.Professors, professor.RDFId())
					}
					break
				}
				cellCursor++
			})
			log.Println("adding course and professor to their respective maps")
			log.Println("course=", course)
			log.Println("professor=", professor)
			courseMap[course.Code] = course
			professorMap[professor.Name()] = professor
		}
		rowCursor++
	})

	for _, professor := range professorMap {
		U.school.AddProfessor(professor)
		log.Println("Adding Professor to map, ", professor)
	}

	for _, course := range courseMap {
		U.school.AddCourse(course)
		log.Println("Adding Course to map, ", course)
	}
}

func (U *ValenciaScraper) GetSchool() models.School {
	log.Println("Getting School")
	if U.school.Name != "" {
		return U.school
	}
	U.school = models.School{
		Name:       "Valencia College",
		Professors: []models.Professor{},
		Courses:    []models.Course{},
	}
	U.StartSchoolScraper()

	U.school = rate_my_professor.Scrape(&U.school, 1544, 17862, 16828, 17341, 13626, 13651)

	fmt.Println("courses=", U.school.Courses)
	fmt.Println("professors=", U.school.Professors)

	return U.school
}
