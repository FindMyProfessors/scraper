package rate_my_professor

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"scraper/models"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type RMPProfessor struct {
	FirstName             string  `json:"firstName"`
	MiddleName            string  `json:"middleName"`
	LastName              string  `json:"lastName"`
	RatingsCount          int     `json:"numRatings"`
	OverallRating         float64 `json:"avgRatingRounded"`
	WouldTakeAgainPercent float64 `json:"wouldTakeAgainPercentRounded"`
}

func StartScrape(school *models.School, schoolIDs ...int) models.School {
	arr := make([]RMPProfessor, 0)
	for _, id := range schoolIDs {
		arr = append(arr, scrape(make([]RMPProfessor, 0), "", id)...)
	}

	sort.SliceStable(arr[:], func(i, j int) bool {
		return strings.Compare(arr[i].LastName, arr[j].LastName) == -1
	})
	log.Println(arr)

	return crossReference(school, arr)
}

func crossReference(school *models.School, rmpProfessors []RMPProfessor) models.School {
	isSameProfessor := func(professor models.Professor, rmsProfessor RMPProfessor) bool {
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

	var wg sync.WaitGroup
	wg.Add(len(rmpProfessors))
	for _, rmpProfessor := range rmpProfessors {
		complete := false
		go func(rmpProfessor RMPProfessor) {
			for index, professor := range school.Professors {
				if complete {
					break
				}
				if isSameProfessor(professor, rmpProfessor) {
					// TODO: Look into a collision function, might possibly need to scrape all of the reviews for the professor.
					professor.Rating = rmpProfessor.OverallRating
					professor.TotalRatings = rmpProfessor.RatingsCount
					school.SetProfessor(index, professor)
					complete = true
				}
			}
			wg.Done()
		}(rmpProfessor)
	}
	wg.Wait()
	return *school
}

func scrape(professors []RMPProfessor, cursor string, schoolId int) []RMPProfessor {
	base64SchoolIdCursor := base64.StdEncoding.EncodeToString([]byte("School-" + strconv.Itoa(schoolId)))
	log.Println("base64SchoolIdCursor=", base64SchoolIdCursor)

	m := make(map[string]interface{})
	m["query"] = "query NewSearch($schoolId: ID, $first: Int!, $last: Int!, $cursor: String!){newSearch{teachers(query:{text:\"\",schoolID:$schoolId} first: $first last: $last, after: $cursor) {edges {node {id firstName lastName numRatings avgRatingRounded wouldTakeAgainPercentRounded }} pageInfo {hasNextPage endCursor}}}}"
	m["variables"] = struct {
		SchoolId string `json:"schoolId"`
		First    int    `json:"first"`
		Last     int    `json:"last"`
		Cursor   string `json:"cursor"`
	}{base64SchoolIdCursor, 1000, 1000, cursor}

	jsonBytes, err := json.Marshal(m)
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
		log.Fatalln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("body=", string(body))

	model, err := parse(body)
	if err != nil {
		log.Fatalln(err)
	}

	for _, prof := range model.Data.NewSearch.Teachers.Professors {
		prof.Professor.WouldTakeAgainPercent = math.Ceil(prof.Professor.WouldTakeAgainPercent*100) / 100
		prof.Professor.OverallRating = math.Ceil(prof.Professor.OverallRating*100) / 100
		professors = append(professors, prof.Professor)
	}

	pageInfo := model.Data.NewSearch.Teachers.PageInfo

	log.Println("EndCursor=", pageInfo.EndCursor)
	log.Println("HasNextPage=", pageInfo.HasNextPage)

	if pageInfo.HasNextPage {
		return scrape(professors, pageInfo.EndCursor, schoolId)
	}
	return professors
}

type Model struct {
	Data struct {
		NewSearch struct {
			Teachers struct {
				Professors []struct {
					Professor RMPProfessor `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"teachers"`
		} `json:"newSearch"`
	} `json:"data"`
}

func parse(body []byte) (Model, error) {
	var model Model
	err := json.Unmarshal(body, &model)
	if err != nil {
		return Model{}, err
	}
	return model, nil
}
