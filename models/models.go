package models

import (
	"fmt"
	"log"
	"strings"
)

type School struct {
	Name       string      `json:"name,omitempty"`
	Professors []Professor `json:"professors"`
	Courses    []Course    `json:"courses"`
}

func (s *School) AddProfessor(professor Professor) []Professor {
	s.Professors = append(s.Professors, professor)
	return s.Professors
}

func (s *School) AddCourse(course Course) []Course {
	s.Courses = append(s.Courses, course)
	return s.Courses
}

func (s *School) SetProfessor(index int, professor Professor) []Professor {
	s.Professors[index] = professor
	fmt.Println("Set professor: ", professor.String())
	return s.Professors
}

func (s School) String() string {
	return fmt.Sprintf("{name=%s, professors=%s, courses=%s}", s.Name, s.Professors, s.Courses)
}

func (s School) RDFId() string {
	return strings.Replace(strings.Replace(strings.Replace(strings.ToLower(s.Name), " ", "_", -1), "'", "_", -1), "\\", "/", -1)
}

type Professor struct {
	FirstName    string   `json:"first_name,omitempty"`
	LastName     string   `json:"last_name,omitempty"`
	Teaches      []Course `json:"teaches,omitempty"`
	TotalRatings int      `json:"totalRatings,omitempty"`
	Rating       float64  `json:"rating,omitempty"`
}

func (p Professor) String() string {
	return fmt.Sprintf("{first_name=%s, last_name=%s, teaches=%s, totalRatings=%d, rating=%f}", p.FirstName, p.LastName, p.Teaches, p.TotalRatings, p.Rating)
}

func (p Professor) RDFId() string {
	rdf := strings.Replace(strings.Replace(strings.Replace(strings.ToLower(p.FirstName+" "+p.LastName), " ", "_", -1), "'", "_", -1), "\\", "/", -1)
	if rdf == "" {
		log.Fatalln(p.String())
	}
	return rdf
}

func (p Professor) Name() string {
	return p.FirstName + " " + p.LastName
}

type Course struct {
	Code       string   `json:"code,omitempty"`
	Name       string   `json:"name,omitempty"`
	Professors []string `json:"professors,omitempty"`
}

func (c Course) String() string {
	return fmt.Sprintf("{name=%s, code=%s,  professors=%s}", c.Name, c.Code, c.Professors)
}
