package ucf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"github.com/FindMyProfessors/scraper/rmp"
	"github.com/FindMyProfessors/scraper/util"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var UCF_COURSE_MODALITIES = [5]string{"WW", "V", "M", "RS", "P"}
var UCF_RMP_IDS = []int{1082, 15516, 5567, 5400}

type Scraper struct {
	term         model.Term
	CourseMap    map[string]*model.Course
	ProfessorMap map[string]*model.Professor
}

func (u *Scraper) Scrape() (*model.School, error) {
	u.CourseMap = map[string]*model.Course{}
	u.ProfessorMap = map[string]*model.Professor{}

	var wg sync.WaitGroup
	wg.Add(len(UCF_COURSE_MODALITIES))

	for _, modality := range UCF_COURSE_MODALITIES {
		url := fmt.Sprintf(
			"https://cdl.ucf.edu/wp-content/themes/cdl/lib/course-search-ajax.php?call=classes&term=%s&prefix=&catalog=&title=&instructor=&career=&college=&department=&mode=%s&_=1631687421296",
			u.term.ID,
			modality,
		)

		mu := &sync.Mutex{}

		go u.scrape(url, mu, func() {
			wg.Done()
		})
	}
	wg.Wait()

	var courseArray []*model.Course

	for _, course := range u.CourseMap {
		courseArray = append(courseArray, course)
	}

	var professorArray []*model.Professor

	for _, professor := range u.ProfessorMap {
		professorArray = append(professorArray, professor)
	}

	school := &model.School{
		Name:       u.Name(),
		Professors: professorArray,
		Courses:    courseArray,
	}

	api := rmp.NewApi("dGVzdDp0ZXN0")
	err := api.StartScrape(context.Background(), school, UCF_RMP_IDS...)
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (u *Scraper) scrape(url string, mu *sync.Mutex, callback func()) {
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
		Professors []Course `json:"classes"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error occurred: ", err)
	}

	for _, ucfProfessor := range data.Professors {
		ucfProfessor.NameFirst = util.WashName(ucfProfessor.NameFirst)
		ucfProfessor.NameLast = util.WashName(ucfProfessor.NameLast)

		course, ok := u.CourseMap[ucfProfessor.CoursePrefix+ucfProfessor.CatalogNumber]
		if !ok {
			course = &model.Course{
				Code: ucfProfessor.CoursePrefix + ucfProfessor.CatalogNumber,
				Name: strings.ReplaceAll(ucfProfessor.CourseTitle, "\"", ""),
			}
		}
		professor, ok := u.ProfessorMap[ucfProfessor.NameFirst+"_"+ucfProfessor.NameLast]
		if !ok {
			professor = &model.Professor{
				FirstName: strings.ReplaceAll(ucfProfessor.NameFirst, "/", ""),
				LastName:  strings.ReplaceAll(ucfProfessor.NameLast, "/", ""),
				Courses:   map[string]*model.Course{},
			}
		}

		doesTeachThisClass := func() bool {
			for _, c := range professor.Courses {
				if c.Code == course.Code {
					return true
				}
			}
			return false
		}()
		if !doesTeachThisClass {
			professor.Courses[course.Code] = course
		}

		mu.Lock()
		u.ProfessorMap[ucfProfessor.NameFirst+"_"+ucfProfessor.NameLast] = professor
		u.CourseMap[ucfProfessor.CoursePrefix+ucfProfessor.CatalogNumber] = course
		mu.Unlock()
	}
	callback()
}

func (u *Scraper) Terms() (terms []model.Term) {
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

func (u *Scraper) Name() string {
	return "University of Central Florida"
}

func (u *Scraper) SetTerm(term model.Term) {
	u.term = term
}
