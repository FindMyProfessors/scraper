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
	Year     int
	Semester Semester
	ID       string
}

type School struct {
}

type Professor struct {
	FirstName string
	LastName  string
	RMPId     string
	Courses   map[string]Course
}

type Course struct {
	Name string
	Code string
}
