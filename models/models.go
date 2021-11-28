package models

import (
	"fmt"
	"strings"
)

type School struct {
	Name       string       `json:"name,omitempty"`
	Professors []*Professor `json:"professors"`
	Courses    []*Course    `json:"courses"`
}

func (s School) String() string {
	return fmt.Sprintf("{name=%s, professors=%s, courses=%s}", s.Name, s.Professors, s.Courses)
}

func CollectCourses(school *School) {
	for _, professor := range school.Professors {
		for _, course := range professor.Teaches {
			course.Professors = append(course.Professors, professor.Name)
		}
	}
}

func (s School) RDFId() string {
	return strings.Replace(strings.Replace(strings.Replace(strings.ToLower(s.Name), " ", "_", -1), "'", "_", -1), "\\", "/", -1)
}

type Professor struct {
	Name         string    `json:"name,omitempty"`
	Teaches      []*Course `json:"teaches,omitempty"`
	TotalRatings int       `json:"totalRatings,omitempty"`
	Rating       float64   `json:"rating,omitempty"`
}

func (p Professor) String() string {
	return fmt.Sprintf("{name=%s, teaches=%s, totalRatings=%d, rating=%f}", p.Name, p.Teaches, p.TotalRatings, p.Rating)
}

func (p Professor) RDFId() string {
	return strings.Replace(strings.Replace(strings.Replace(strings.ToLower(p.Name), " ", "_", -1), "'", "_", -1), "\\", "/", -1)
}

type Course struct {
	Code       string   `json:"code,omitempty"`
	Name       string   `json:"name,omitempty"`
	Professors []string `json:"professors,omitempty"`
}

func (c Course) String() string {
	return fmt.Sprintf("{name=%s, code=%s,  professors=%s}", c.Name, c.Code, c.Professors)
}
