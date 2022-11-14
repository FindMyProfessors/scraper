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

func (a *Api) UpsertSchool(ctx context.Context, school *model.School) error {
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

	professorMap, err := a.GetAllProfessors(ctx, school)
	if err != nil {
		return err
	}

	for _, elem := range school.Professors {
		match, ok := professorMap[elem.FirstName+elem.LastName]
		if !ok {
			professor, err := CreateProfessor(ctx, a.Client, *school.ID, NewProfessor{
				FirstName: elem.FirstName,
				LastName:  elem.LastName,
				RmpId:     &elem.RMPId,
			})
			if err != nil {
				return err
			}
			match = professor.CreateProfessor
			match.FirstName = elem.FirstName
			match.LastName = elem.LastName
			match.RMPId = elem.RMPId
			// CREATE ALL REVIEWS. NEW PROFESSOR.

			for _, review := range elem.Reviews {
				_, err = CreateReview(ctx, a.Client, match.ID, NewReview{
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
		} else {
			if len(match.Reviews) > 0 {
				newReviews, err := a.InsertNewReviews(ctx, elem, match.Reviews[0])
				if err != nil {
					return err
				}
				if newReviews > 0 {
					fmt.Printf("%s %s has %d new reviews!\n", elem.FirstName, elem.LastName, newReviews)
				}
			}
		}
	}
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

func (a *Api) GetAllProfessors(ctx context.Context, school *model.School) (map[string]*model.Professor, error) {
	var professorMap map[string]*model.Professor

	var after *string
	hasNextPage := true
	for hasNextPage {
		response, err := GetProfessors(ctx, a.Client, *school.ID, after)
		if err != nil {
			return professorMap, err
		}
		query := response.Professors

		professors := query.Professors

		for _, elem := range professors {
			professorMap[elem.FirstName+elem.LastName] = &elem
		}

		after = &query.PageInfo.EndCursor
		hasNextPage = query.PageInfo.HasNextPage
	}
	return professorMap, nil
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
