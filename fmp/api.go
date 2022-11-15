package fmp

import (
	"context"
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"github.com/Khan/genqlient/graphql"
	"net/http"
	"time"
)

type Api struct {
	Client graphql.Client
}

func NewApi(endpoint string) *Api {
	return &Api{Client: graphql.NewClient(endpoint, http.DefaultClient)}
}

func (a *Api) UpsertSchool(ctx context.Context, school *model.School, term *model.Term) error {
	foundSchool, err := GetSchoolByName(ctx, a.Client, school.Name)
	if err != nil {
		return err
	}
	if foundSchool == nil {
		schoolId, err := a.CreateSchool(ctx, school)
		if err != nil {
			return err
		}
		school.ID = &schoolId
	} else {
		school.ID = foundSchool.ID
	}
	termInput := TermInput{
		Year:     term.Year,
		Semester: Semester(term.Semester),
	}

	professorMap, err := a.GetAllProfessors(ctx, school, termInput)
	if err != nil {
		return err
	}

	courseIds, err := a.GetAllCourseIds(ctx, school, termInput)

	if err != nil {
		return err
	}

	for _, course := range school.Courses {
		_, exists := courseIds[course.Code]
		if !exists {
			response, err := CreateCourse(ctx, a.Client, *school.ID, NewCourse{
				Name: course.Name,
				Code: course.Code,
			})
			if err != nil {
				return err
			}
			courseIds[course.Code] = response.CreateCourse.ID
		}
	}

	for _, elem := range school.Professors {
		err = func(professor *model.Professor) error {
			match, ok := professorMap[professor.FirstName+professor.LastName]
			if ok {
				// TODO: set rmpId via FMP API if elem.RMPId is set
				// TODO: update course registrations if already set
				if len(match.Reviews) > 0 {
					newReviews, err := a.InsertNewReviews(ctx, professor, match.Reviews[0])
					if err != nil {
						return err
					}
					if newReviews > 0 {
						fmt.Printf("%s %s has %d new reviews!\n", professor.FirstName, professor.LastName, newReviews)
					}
				}
				for _, course := range professor.Courses {
					_, ok := match.Courses[course.Code]
					if !ok {
						_, err = RegisterCourse(ctx, a.Client, match.ID, courseIds[course.Code], termInput)
						if err != nil {
							return err
						}
					}

				}
			} else {
				response, err := CreateProfessor(ctx, a.Client, *school.ID, NewProfessor{
					FirstName: professor.FirstName,
					LastName:  professor.LastName,
					RmpId:     &professor.RMPId,
				})
				if err != nil {
					return err
				}
				// CREATE ALL REVIEWS. NEW PROFESSOR.

				for _, review := range professor.Reviews {
					_, err = CreateReview(ctx, a.Client, response.CreateProfessor.ID, NewReview{
						Quality:    review.Quality,
						Difficulty: review.Difficulty,
						Time:       review.Date.Format(time.RFC3339),
						Tags:       review.Tags,
						Grade:      review.Grade,
					})
					if err != nil {
						return err
					}
				}
				for _, course := range professor.Courses {
					courseId, ok := courseIds[course.Code]
					if !ok {
						return fmt.Errorf("course with cod %s not found in courses retrieved from FMP", course.Code)
					}
					_, err = RegisterCourse(ctx, a.Client, response.CreateProfessor.ID, courseId, termInput)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}(elem)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	fmt.Printf("Finished uploading school to FMP\n")
	return nil
}

func (a *Api) InsertNewReviews(ctx context.Context, scrapedProfessor *model.Professor, retrievedReview *model.Review) (int, error) {
	var indexUntil int

	// TODO: Optimize using a binary search

	for i, review := range scrapedProfessor.Reviews {
		if review.Date.Equal(retrievedReview.Date) {
			indexUntil = i
		}
	}

	for i := 0; i < indexUntil; i++ {
		review := scrapedProfessor.Reviews[i]
		_, err := CreateReview(ctx, a.Client, scrapedProfessor.ID, NewReview{
			Quality:    review.Quality,
			Difficulty: review.Difficulty,
			Time:       review.Date.Format(time.RFC3339),
			Tags:       review.Tags,
			Grade:      review.Grade,
		})
		if err != nil {
			return -1, err
		}
	}

	return 0, nil

}

func (a *Api) GetAllProfessors(ctx context.Context, school *model.School, term TermInput) (map[string]*model.Professor, error) {
	professorMap := make(map[string]*model.Professor)

	var after *string
	hasNextPage := true
	for hasNextPage {
		response, err := GetProfessors(ctx, a.Client, *school.ID, after, term)
		if err != nil {
			return nil, err
		}

		query := response.Professors

		professors := query.Professors

		for _, elem := range professors {

			professor := model.Professor{
				ID:        elem.Id,
				FirstName: elem.FirstName,
				LastName:  elem.LastName,
				Courses:   map[string]*model.Course{},
				Reviews:   []*model.Review{},
			}

			for _, review := range elem.Reviews.Reviews {
				professor.Reviews = append(professor.Reviews, &review)
			}

			if elem.Teaches.PageInfo.HasNextPage {
				return nil, fmt.Errorf("implement getting next page for professors, %s is teaching %d classes", elem.FirstName+" "+elem.LastName, elem.Teaches.TotalCount)
			}
			for _, course := range elem.Teaches.Courses {
				professor.Courses[course.Code] = &course
			}

			professorMap[professor.FirstName+professor.LastName] = &professor
		}

		after = &query.PageInfo.EndCursor
		hasNextPage = query.PageInfo.HasNextPage
	}
	return professorMap, nil
}

func (a *Api) GetAllCourseIds(ctx context.Context, school *model.School, input TermInput) (map[string]string, error) {
	courseMap := make(map[string]string)

	var after *string
	hasNextPage := true
	for hasNextPage {
		response, err := GetCourses(ctx, a.Client, *school.ID, input, after)
		if err != nil {
			return nil, err
		}
		query := response.School.Courses

		for _, elem := range query.Courses {
			courseMap[elem.Code] = elem.ID
		}

		after = &query.PageInfo.EndCursor
		hasNextPage = query.PageInfo.HasNextPage
	}
	return courseMap, nil
}

func (a *Api) CreateSchool(ctx context.Context, school *model.School) (id string, err error) {
	response, err := CreateSchool(ctx, a.Client, NewSchool{Name: school.Name})
	if err != nil {
		return "", err
	}
	return *response.GetCreateSchool().ID, nil
}

func GetSchoolByName(ctx context.Context, client graphql.Client, schoolName string) (*model.School, error) {
	response, err := GetSchools(ctx, client)
	if err != nil {
		return nil, err
	}
	query := response.Schools
	for _, elem := range query.GetSchools() {
		if schoolName == elem.Name {
			return &elem, nil
		}
	}
	return nil, nil
}
