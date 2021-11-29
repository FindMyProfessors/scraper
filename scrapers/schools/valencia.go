package schools

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"scraper/models"
	"strconv"
	"strings"
)

var SchoolIDs = []int{1544, 17862, 16828, 17341, 13626, 13651}

type ValenciaScraper struct {
	school models.School
}

func (U *ValenciaScraper) StartRMSScraper() {
	fmt.Println("")
	var scrapeFunction func(page int) bool
	type Model struct {
		Professors []struct {
			FirstName     string `json:"tFname" bson:"firstName"`
			MiddleName    string `json:"tMiddlename" bson:"middleName"`
			LastName      string `json:"tLname" bson:"lastName"`
			RatingsCount  int    `json:"tNumRatings" bson:"ratingCount"`
			RatingClass   string `json:"rating_class" bson:"ratingClass"`
			OverallRating string `json:"overall_rating" bson:"overallRating"`
		} `json:"professors"`
		Total     int    `json:"searchResultsTotal"`
		Remaining int    `json:"remaining"`
		Type      string `json:"type"`
	}
	var completeModel Model

	for _, id := range SchoolIDs {
		func(id int) {
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
		}(id)
	}

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
	for _, anonProfessor := range completeModel.Professors {
		var professor models.Professor
		var index int
		for index, professor = range U.school.Professors {
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
				U.school.Professors[index] = professor
				break
			}
		}
	}
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

func (U ValenciaScraper) GetSchool() models.School {
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
	U.StartRMSScraper()

	println("courses=", U.school.Courses)
	println("professors=", U.school.Professors)

	return U.school
}
