package rmp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"github.com/Khan/genqlient/graphql"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Api struct {
	Client graphql.Client
}

type AuthenticationTransportWrapper struct {
	Password string
}

func (d AuthenticationTransportWrapper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", "Basic "+d.Password)
	return http.DefaultTransport.RoundTrip(request)
}

func NewApi(password string) *Api {
	httpClient := &http.Client{
		Transport:     AuthenticationTransportWrapper{Password: password},
		Timeout:       time.Minute * 1,
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           http.DefaultClient.Jar,
	}
	return &Api{Client: graphql.NewClient("https://www.ratemyprofessors.com/graphql", httpClient)}
}

func (a *Api) StartScrape(ctx context.Context, school *model.School, schoolIDs ...int) error {
	var professorArray []*model.Professor
	for _, id := range schoolIDs {
		fmt.Printf("Scraping %s with rmp id %d\n", school.Name, id)
		base64SchoolIdCursor := base64.StdEncoding.EncodeToString([]byte("School-" + strconv.Itoa(id)))
		professors, err := a.scrape(ctx, []*model.Professor{}, "", base64SchoolIdCursor)
		if err != nil {
			return fmt.Errorf("unable to scrape: %v", err)
		}
		professorArray = append(professorArray, professors...)
	}

	sort.SliceStable(professorArray[:], func(i, j int) bool {
		return strings.Compare(professorArray[i].LastName, professorArray[j].LastName) == -1
	})
	crossReference(school, professorArray)
	for _, elem := range professorArray {
		fmt.Printf("professor=%v\n", elem)
	}
	return nil
}

func isSameProfessor(professor *model.Professor, rmpProfessor *model.Professor) bool {
	//fmt.Printf("professor=%v\n", professor)
	//fmt.Printf("rmpProfessor=%v\n", rmpProfessor)

	if professor.FirstName == rmpProfessor.FirstName {
		if professor.LastName == rmpProfessor.LastName {
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

	if tryMultiname(professor.FirstName, rmpProfessor.FirstName, "-") && tryMultiname(professor.LastName, rmpProfessor.LastName, "-") {
		return true
	}
	return false
}

func crossReference(school *model.School, rmpProfessors []*model.Professor) {
	var wg sync.WaitGroup
	wg.Add(len(rmpProfessors))
	for _, rmpProfessor := range rmpProfessors {
		go func(rmpProfessor *model.Professor) {
			for _, professor := range school.Professors {
				if isSameProfessor(professor, rmpProfessor) {
					professor.Reviews = rmpProfessor.Reviews
					professor.RMPId = rmpProfessor.RMPId
					log.Printf("found match for %s %s\n", professor.FirstName, professor.LastName)
					break
				}
			}
			wg.Done()
		}(rmpProfessor)
	}
	wg.Wait()

}

func (a *Api) scrape(ctx context.Context, professors []*model.Professor, cursor string, schoolId string) ([]*model.Professor, error) {
	response, err := NewSearch(ctx, a.Client, &schoolId, 1000, cursor)
	if err != nil {
		return nil, err
	}

	for _, prof := range response.NewSearch.Teachers.Edges {
		rmpProfessor := prof.Node

		var reviews = make([]*model.Review, len(rmpProfessor.Ratings.Edges), len(rmpProfessor.Ratings.Edges))

		for i, elem := range rmpProfessor.Ratings.Edges {
			fmt.Printf("elem.Rating=%v\n", elem)
			rmpRating := elem.Node
			t, err := time.Parse(model.RMPTimeConstant, *rmpRating.Date)
			if err != nil {
				return nil, err
			}
			var tags []model.Tag
			tagsString := *rmpRating.RatingTags
			if len(tagsString) > 0 {
				split := strings.Split(tagsString, "--")

				tags = make([]model.Tag, 0, len(split))

				for _, elem := range split {
					tag, err := model.GetTagByString(elem)
					if err != nil {
						return nil, err
					}
					tags = append(tags, tag)
				}
			} else {
				tags = []model.Tag{}
			}

			grade := model.GetGradeByString(*rmpRating.Grade)
			if !grade.IsValid() {
				return nil, fmt.Errorf("%s is an invalid grade", *rmpRating.Grade)
			}

			reviews[i] = &model.Review{
				Quality:    float64(*rmpRating.QualityRating),
				Difficulty: *rmpRating.DifficultyRatingRounded,
				Date:       t,
				Tags:       tags,
				Grade:      grade,
			}
		}

		professors = append(professors, &model.Professor{
			FirstName: *rmpProfessor.FirstName,
			LastName:  *rmpProfessor.LastName,
			RMPId:     *rmpProfessor.Id,
			Reviews:   reviews,
		})

		fmt.Printf("reviews=%v\n", reviews)
	}

	pageInfo := response.NewSearch.Teachers.PageInfo

	log.Println("EndCursor=", pageInfo.EndCursor)
	log.Println("HasNextPage=", pageInfo.HasNextPage)

	if pageInfo.HasNextPage {
		return a.scrape(ctx, professors, *pageInfo.EndCursor, schoolId)
	}
	return professors, nil
}

func parse(body []byte) (*RMPResponseModel, error) {
	var rmpResponseModel RMPResponseModel
	err := json.Unmarshal(body, &rmpResponseModel)
	if err != nil {
		return &RMPResponseModel{}, err
	}
	return &rmpResponseModel, nil
}

func makeRequest(schoolId int, cursor string) (*RMPResponseModel, error) {
	base64SchoolIdCursor := base64.StdEncoding.EncodeToString([]byte("School-" + strconv.Itoa(schoolId)))
	//log.Printf("base64SchoolIdCursor=%s\n", base64SchoolIdCursor)

	variableMap := map[string]any{}
	variableMap["query"] = `query NewSearch($schoolId:ID,$first:Int!,$cursor:String!){newSearch{teachers(query:{text:"",schoolID:$schoolId}first:$first after:$cursor){edges{node{id firstName lastName ratings(first:500){edges{node{qualityRating difficultyRatingRounded date ratingTags grade}}}}}pageInfo{hasNextPage endCursor}}}}`
	variableMap["variables"] = struct {
		SchoolId string `json:"schoolId"`
		First    int    `json:"first"`
		Cursor   string `json:"cursor"`
	}{base64SchoolIdCursor, 1000, cursor}

	jsonBytes, err := json.Marshal(variableMap)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("json:" + string(jsonBytes))

	request, err := http.NewRequest("POST", "https://www.ratemyprofessors.com/graphql", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Authorization", "Basic dGVzdDp0ZXN0")
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Println("body=", string(body))

	return parse(body)
}
