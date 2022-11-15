package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Semester string

type PageInfo struct {
	StartCursor string `json:"startCursor,omitempty"`
	EndCursor   string `json:"endCursor,omitempty"`
	HasNextPage bool   `json:"hasNextPage,omitempty"`
}

var (
	SemesterSpring Semester = "SPRING"
	SemesterSummer Semester = "SUMMER"
	SemesterFall   Semester = "FALL"

	AllSemesters = []Semester{
		SemesterSpring, SemesterSummer, SemesterFall,
	}
)

func (s Semester) GetIndex() int {
	switch s {
	case SemesterSpring:
		return 0
	case SemesterSummer:
		return 1
	case SemesterFall:
		return 2
	}
	return 0
}

type Term struct {
	Year     int      `json:"year,omitempty"`
	Semester Semester `json:"semester,omitempty"`
	ID       string   `json:"id,omitempty"`
}

type School struct {
	ID         *string      `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Professors []*Professor `json:"professors,omitempty"`
	Courses    []*Course    `json:"courses,omitempty"`
}

type Professor struct {
	ID        string             `json:"id,omitempty"`
	FirstName string             `json:"firstname,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
	RMPId     string             `json:"rmpId,omitempty"`
	Courses   map[string]*Course `json:"courses,omitempty"`
	Reviews   []*Review          `json:"reviews,omitempty"`
}

func (p *Professor) UnmarshalJSON(bytes []byte) error {
	var data map[string]any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	id, ok := data["id"].(string)
	if ok {
		p.ID = id
	}

	rmpId, ok := data["rmpId"].(string)
	if ok {
		p.RMPId = rmpId
	}

	firstName, ok := data["firstName"].(string)
	if !ok {
		firstName, ok = data["firstname"].(string)
	}
	p.FirstName = firstName

	lastName, ok := data["lastName"].(string)
	if !ok {
		lastName, ok = data["lastname"].(string)
	}
	p.LastName = lastName

	courses, ok := data["courses"].(map[string]any)
	if ok {
		p.Courses = map[string]*Course{}
		for courseCode, value := range courses {
			courseMap := value.(map[string]any)
			p.Courses[courseCode] = &Course{
				Name: courseMap["name"].(string),
				Code: courseCode,
			}
		}

	}

	reviewsAnyArray, ok := data["reviews"].([]any)
	if ok {
		p.Reviews = make([]*Review, len(reviewsAnyArray), len(reviewsAnyArray))

		for i, reviewAny := range reviewsAnyArray {

			reviewMap := reviewAny.(map[string]any)

			var review Review
			dateString := reviewMap["date"].(string)
			t, err := time.Parse(time.RFC3339, dateString)
			if err != nil {
				return err
			}
			review.Date = t
			review.Difficulty = reviewMap["difficulty"].(float64)
			review.Quality = reviewMap["quality"].(float64)
			review.Grade = Grade(reviewMap["grade"].(string))

			tagAnyArray, ok := reviewMap["tags"].([]any)
			if ok {
				tagArray := make([]Tag, len(tagAnyArray), len(tagAnyArray))

				for i, tagAny := range tagAnyArray {
					tagArray[i] = Tag(tagAny.(string))
				}
				review.Tags = tagArray
			}
			p.Reviews[i] = &review
		}
	}

	return nil
}

func (p *Professor) String() string {
	return fmt.Sprintf("{id: %s, rmpId: %s, firstname: %s, lastname: %s, rmpId: %s, courses:%v, reviews: %v}", p.ID, p.RMPId, p.FirstName, p.LastName, p.RMPId, p.Courses, p.Reviews)
}

type Course struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

func (c *Course) String() string {
	return fmt.Sprintf("{name: %s, code: %s}", c.Name, c.Code)
}

type Tag string

const (
	TagToughGrader            Tag = "TOUGH_GRADER"
	TagGetReadyToRead         Tag = "GET_READY_TO_READ"
	TagParticipationMatters   Tag = "PARTICIPATION_MATTERS"
	TagExtraCredit            Tag = "EXTRA_CREDIT"
	TagGroupProjects          Tag = "GROUP_PROJECTS"
	TagAmazingLectures        Tag = "AMAZING_LECTURES"
	TagClearGradingCriteria   Tag = "CLEAR_GRADING_CRITERIA"
	TagGivesGoodFeedback      Tag = "GIVES_GOOD_FEEDBACK"
	TagInspirational          Tag = "INSPIRATIONAL"
	TagLotsOfHomework         Tag = "LOTS_OF_HOMEWORK"
	TagHilarious              Tag = "HILARIOUS"
	TagBewareOfPopQuizzes     Tag = "BEWARE_OF_POP_QUIZZES"
	TagSoManyPapers           Tag = "SO_MANY_PAPERS"
	TagCaring                 Tag = "CARING"
	TagRespected              Tag = "RESPECTED"
	TagLectureHeavy           Tag = "LECTURE_HEAVY"
	TagTestsAreTough          Tag = "TESTS_ARE_TOUGH"
	TagTestHeavy              Tag = "TEST_HEAVY"
	TagWouldTakeAgain         Tag = "WOULD_TAKE_AGAIN"
	TagTestsNotMany           Tag = "TESTS_NOT_MANY"
	TagSkipClassYouWontPass   Tag = "SKIP_CLASS_YOU_WONT_PASS"
	TagCaresAboutStudents     Tag = "CARES_ABOUT_STUDENTS"
	TagRespectedByStudents    Tag = "RESPECTED_BY_STUDENTS"
	TagExtraCreditOffered     Tag = "EXTRA_CREDIT_OFFERED"
	TagGradedByFewThings      Tag = "GRADED_BY_FEW_THINGS"
	TagAccessibleOutsideClass Tag = "ACCESSIBLE_OUTSIDE_CLASS"
	TagOnlineSavvy            Tag = "ONLINE_SAVVY"
)

func (t Tag) IsValid() bool {
	switch t {
	case TagToughGrader, TagExtraCreditOffered, TagRespectedByStudents, TagCaresAboutStudents, TagTestsNotMany, TagTestsAreTough, TagWouldTakeAgain, TagSkipClassYouWontPass, TagTestHeavy, TagGetReadyToRead, TagParticipationMatters, TagExtraCredit, TagGroupProjects, TagAmazingLectures, TagClearGradingCriteria, TagGivesGoodFeedback, TagInspirational, TagLotsOfHomework, TagHilarious, TagBewareOfPopQuizzes, TagSoManyPapers, TagCaring, TagRespected, TagLectureHeavy, TagGradedByFewThings, TagAccessibleOutsideClass, TagOnlineSavvy:
		return true
	}
	return false
}

func (t Tag) String() string {
	return string(t)
}

type Grade string

const (
	GradeAPlus      Grade = "A_PLUS"
	GradeA          Grade = "A"
	GradeAMinus     Grade = "A_MINUS"
	GradeBPlus      Grade = "B_PLUS"
	GradeB          Grade = "B"
	GradeBMinus     Grade = "B_MINUS"
	GradeCPlus      Grade = "C_PLUS"
	GradeC          Grade = "C"
	GradeCMinus     Grade = "C_MINUS"
	GradeDPlus      Grade = "D_PLUS"
	GradeD          Grade = "D"
	GradeDMinus     Grade = "D_MINUS"
	GradeFPlus      Grade = "F_PLUS"
	GradeF          Grade = "F"
	GradeFMinus     Grade = "F_MINUS"
	GradeIncomplete Grade = "INCOMPLETE"
	GradeWithdrawn  Grade = "WITHDRAWN"
	GradeNotSure    Grade = "NOT_SURE"
	GradeOther      Grade = "OTHER"
)

func (g Grade) IsValid() bool {
	switch g {
	case GradeAPlus, GradeA, GradeAMinus, GradeBPlus, GradeB, GradeBMinus, GradeCPlus, GradeC, GradeCMinus, GradeDPlus, GradeD, GradeDMinus, GradeFPlus, GradeF, GradeFMinus, GradeIncomplete, GradeWithdrawn, GradeNotSure, GradeOther:
		return true
	}
	return false
}

func (g Grade) String() string {
	return string(g)
}

func GetGradeByString(gradeString string) Grade {
	switch gradeString {
	case "A+":
		return GradeAPlus
	case "A":
		return GradeA
	case "A-":
		return GradeAMinus

	case "B+":
		return GradeBPlus
	case "B":
		return GradeB
	case "B-":
		return GradeBMinus

	case "C+":
		return GradeCPlus
	case "C":
		return GradeC
	case "C-":
		return GradeCMinus

	case "D+":
		return GradeDPlus
	case "D":
		return GradeD
	case "D-":
		return GradeDMinus

	case "F+":
		return GradeFPlus
	case "F":
		return GradeF
	case "F-":
		return GradeFMinus

	case "Incomplete":
		return GradeIncomplete
	case "Drop/Withdrawal":
		return GradeWithdrawn
	case "Not sure yet":
		return GradeNotSure
	default:
		fmt.Printf("%s is passed in but unable to be parsed as a Grade\n", gradeString)
		return GradeOther
	}
}

type Review struct {
	ID         string    `json:"id,omitempty"`
	Quality    float64   `json:"quality"`
	Difficulty float64   `json:"difficulty"`
	Date       time.Time `json:"date"`
	Tags       []Tag     `json:"tags"`
	Grade      Grade     `json:"grade"`
}

func (r *Review) String() string {
	return fmt.Sprintf("{quality: %f, difficulty: %f, date: %s, tags: %v, grade: %s}", r.Quality, r.Difficulty, r.Date, r.Tags, r.Grade)
}

func GetTagByString(tagString string) (Tag, error) {
	formattedTagString := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.TrimSpace(
						strings.ToUpper(tagString)),
					" ", "_"),
				"?", ""),
			".", ""),
		"'", "",
	)

	tag := Tag(formattedTagString)
	if !tag.IsValid() {
		return "", fmt.Errorf("%s is an invalid Tag", tagString)
	}
	return tag, nil
}

// RMPTimeConstant Refer to the time.Layout type for more info
const RMPTimeConstant = "2006-01-02 15:04:05 -0700 MST"

func (r *Review) UnmarshalJSON(bytes []byte) (err error) {
	var data map[string]any

	if err = json.Unmarshal(bytes, &data); err != nil {
		fmt.Printf("error is also happening here\n")
		return err
	}

	id, ok := data["id"].(string)
	if ok {
		r.ID = id
	}

	f, ok := data["qualityRating"].(float64)
	if !ok {
		qualityAny, ok := data["quality"]
		if ok {
			r.Quality = qualityAny.(float64)
		}
	} else {
		r.Quality = f
	}

	f2, ok := data["difficultyRatingRounded"].(float64)
	if !ok {
		f2, ok = data["difficulty"].(float64)
	}
	r.Difficulty = f2

	dateString, ok := data["date"].(string)
	if ok {
		t, err := time.Parse(time.RFC3339, dateString)
		if err != nil || t.Year() == 1 {
			// 0001-01-01T00:00:00Z incorrectly parsed string
			t, err = time.Parse(RMPTimeConstant, dateString)
			if err != nil {
				return err
			}
		}
		r.Date = t
	}

	tagsString, ok := data["ratingTags"].(string)
	if ok {
		if len(tagsString) > 0 {
			split := strings.Split(tagsString, "--")

			tags := make([]Tag, 0, len(split))

			for _, elem := range split {
				tag, err := GetTagByString(elem)
				if err != nil {
					return err
				}
				tags = append(tags, tag)
			}
			r.Tags = tags
		} else {
			r.Tags = []Tag{}
		}
	} else {
		tagArray, ok := data["tags"].([]any)
		if ok {
			if len(tagArray) == 0 {
				r.Tags = []Tag{}
			} else {
				var tags []Tag
				for _, elem := range tagArray {
					tag, err := GetTagByString(elem.(string))
					if err != nil {
						return err
					}
					tags = append(tags, tag)
				}
				r.Tags = tags
			}
		}
	}
	gradeString, ok := data["grade"].(string)
	if ok {
		grade := GetGradeByString(gradeString)
		if !grade.IsValid() {
			return fmt.Errorf("%s is an invalid grade", gradeString)
		}
		r.Grade = grade
	}
	return nil
}
