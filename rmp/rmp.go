package rmp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FindMyProfessors/scraper/model"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func StartScrape(school *model.School, schoolIDs ...int) error {
	var professorArray = []*model.Professor{}
	for _, id := range schoolIDs {
		fmt.Printf("Scraping %s with rmp id %d\n", school.Name, id)
		professors, err := scrape([]*model.Professor{}, "", id)
		if err != nil {
			return err
		}
		professorArray = append(professorArray, professors...)
	}

	sort.SliceStable(professorArray[:], func(i, j int) bool {
		return strings.Compare(professorArray[i].LastName, professorArray[j].LastName) == -1
	})
	log.Println(professorArray)
	crossReference(school, professorArray)
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
	mu := &sync.Mutex{}
	for _, rmpProfessor := range rmpProfessors {
		go func(rmpProfessor *model.Professor) {
			defer wg.Done()
			for i, professor := range school.Professors {
				if isSameProfessor(professor, rmpProfessor) {

					mu.Lock()
					professor.Reviews = rmpProfessor.Reviews
					school.Professors[i] = professor
					mu.Unlock()
					log.Printf("found match for %s %s\n", professor.FirstName, professor.LastName)
					break
				}
			}
		}(rmpProfessor)
	}
	wg.Wait()
}

func scrape(professors []*model.Professor, cursor string, schoolId int) ([]*model.Professor, error) {
	rmpResponseModel, err := makeRequest(schoolId, cursor)
	if err != nil {
		return nil, err
	}

	for _, prof := range rmpResponseModel.Data.NewSearch.Teachers.Professors {
		rmpProfessor := prof.Professor

		reviews := make([]*model.Review, 0, len(rmpProfessor.Reviews.Edges))

		for _, elem := range rmpProfessor.Reviews.Edges {
			reviews = append(reviews, &elem.Rating)
		}

		professors = append(professors, &model.Professor{
			FirstName: rmpProfessor.FirstName,
			LastName:  rmpProfessor.LastName,
			RMPId:     rmpProfessor.RMPId,
			Reviews:   reviews,
		})
	}

	pageInfo := rmpResponseModel.Data.NewSearch.Teachers.PageInfo

	log.Println("EndCursor=", pageInfo.EndCursor)
	log.Println("HasNextPage=", pageInfo.HasNextPage)

	if pageInfo.HasNextPage {
		return scrape(professors, pageInfo.EndCursor, schoolId)
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
