package schools

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"scraper/models"
	"strings"
)

type ValenciaScraper struct {
	school *models.School
}

func (U ValenciaScraper) StartRMSScraper() {

}

func (U ValenciaScraper) StartSchoolScraper() {
	log.Println("Started scraping valencia website")
	request, err := http.NewRequest(
		"GET",
		"https://valenciacollege.edu/academics/schedule-search/search.php?term=202210&INSM%5B%5D=N&INSM%5B%5D=M&INSM%5B%5D=X&INSM%5B%5D=S&Campus=&course_desc=&crn=",
		nil,
	)
	if err != nil {
		log.Fatalln("error occurred: ", err)
	}

	//	request.Header = http.Header{
	//		"Connection":                []string{"keep-alive"},
	//		"Upgrade-Insecure-Requests": []string{"1"},
	//		"Origin":                    []string{"https://valenciacollege.edu"},
	//		"Content-Type":              []string{"application/x-www-form-urlencoded"},
	//		"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML like Gecko) Chrome/93.0.4577.63 Safari/537.36"},
	//		"Accept":                    []string{"text/htmlapplication/xhtml+xmlapplication/xml;q=0.9image/avifimage/webpimage/apng*/*;q=0.8application/signed-exchange;v=b3;q=0.9"},
	//		"Sec-GPC":                   []string{"1"},
	//		"Sec-Fetch-Site":            []string{"same-origin"},
	//		"Sec-Fetch-Mode":            []string{"navigate"},
	//		"Sec-Fetch-User":            []string{"?1"},
	//		"Sec-Fetch-Dest":            []string{"document"},
	//		"Referer":                   []string{"https://valenciacollege.edu/academics/schedule-search/"},
	//		"Accept-Language":           []string{"en-USen;q=0.9"},
	//	}

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

	professorMap := make(map[string]*models.Professor)
	courseMap := make(map[string]*models.Course)

	log.Println("Iterating over 'tr' selector")
	table.Find("tr").Each(func(i int, rowSelection *goquery.Selection) {
		log.Println("rowCursor=", rowCursor)
		//log.Println("rowSelection.Text()=", rowSelection.Text())
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
					c := courseMap[courseCode]
					if c == nil {
						course = models.Course{
							Code: courseCode,
							Name: courseName,
						}
					} else {
						course = *c
					}
					break
				case 3:
					professorText := selection.Text()
					log.Println("professorText=", professorText)

					p := professorMap[professorText]
					if p == nil {
						professor = models.Professor{
							Name:    professorText,
							Teaches: []*models.Course{},
						}
					} else {
						professor = *p
					}
					doesTeachThisClass := func(courses []*models.Course) bool {
						for _, c := range professor.Teaches {
							if c.Code == course.Code {
								return true
							}
						}
						return false
					}(professor.Teaches)
					if !doesTeachThisClass {
						professor.Teaches = append(professor.Teaches, &course)
					}
					break
				}

				cellCursor++
			})
			log.Println("adding course and professor to their respective maps")
			log.Println("course=", course)
			log.Println("professor=", professor)
			courseMap[course.Code] = &course
			professorMap[professor.Name] = &professor
		}
		rowCursor++
	})

	for _, professor := range professorMap {
		U.school.Professors = append(U.school.Professors, professor)
	}

	for _, course := range courseMap {
		U.school.Courses = append(U.school.Courses, course)
	}
}

func (U ValenciaScraper) GetSchool() models.School {
	if U.school != nil {
		return *U.school
	}
	U.school = &models.School{
		Name:       "Valencia College",
		Professors: []*models.Professor{},
		Courses:    []*models.Course{},
	}
	U.StartRMSScraper()
	U.StartSchoolScraper()

	return *U.school
}
