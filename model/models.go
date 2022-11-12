package model

type Semester string

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
	Name       string       `json:"name,omitempty"`
	Professors []*Professor `json:"professors,omitempty"`
	Courses    []*Course    `json:"courses,omitempty"`
}

type Professor struct {
	FirstName string             `json:"firstname,omitempty"`
	LastName  string             `json:"lastName,omitempty"`
	RMPId     string             `json:"id,omitempty"`
	Courses   map[string]*Course `json:"courses,omitempty"`
	Reviews   []*Review          `json:"reviews,omitempty"`
}

type Course struct {
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type Review struct {
	Quality    float64 `json:"qualityRating"`
	Difficulty float64 `json:"difficultyRatingRounded"`
	DateString string  `json:"date"`
	RatingTags string  `json:"ratingTags"`
	Grade      string  `json:"grade"`
}
