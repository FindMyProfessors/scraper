package main

import "strings"

type School struct {
	Name       string      `json:"name,omitempty"`
	Professors []Professor `json:"professors"`
	Courses    []Course    `json:"courses"`
}

func (s School) RDFId() string {
	return strings.Replace(strings.ToLower(s.Name), " ", "_", -1)
}

type Professor struct {
	Name         string   `json:"name,omitempty"`
	School       School   `json:"school,omitempty"`
	Teaches      []Course `json:"teaches,omitempty"`
	TotalRatings int      `json:"totalRatings,omitempty"`
	Rating       float64  `json:"rating,omitempty"`
}

func (p Professor) RDFId() string {
	return strings.Replace(strings.ToLower(p.Name), " ", "_", -1)
}

type Course struct {
	Code       string      `json:"code,omitempty"`
	Name       string      `json:"name,omitempty"`
	School     School      `json:"school,omitempty"`
	Professors []Professor `json:"professors,omitempty"`
}
