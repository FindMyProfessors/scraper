// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package fmp

import (
	"context"

	"github.com/FindMyProfessors/scraper/model"
	"github.com/Khan/genqlient/graphql"
)

// CreateCourseResponse is returned by CreateCourse on success.
type CreateCourseResponse struct {
	CreateCourse *model.Course `json:"createCourse"`
}

// GetCreateCourse returns CreateCourseResponse.CreateCourse, and is useful for accessing the field via an interface.
func (v *CreateCourseResponse) GetCreateCourse() *model.Course { return v.CreateCourse }

// CreateProfessorResponse is returned by CreateProfessor on success.
type CreateProfessorResponse struct {
	CreateProfessor *model.Professor `json:"createProfessor"`
}

// GetCreateProfessor returns CreateProfessorResponse.CreateProfessor, and is useful for accessing the field via an interface.
func (v *CreateProfessorResponse) GetCreateProfessor() *model.Professor { return v.CreateProfessor }

// CreateReviewResponse is returned by CreateReview on success.
type CreateReviewResponse struct {
	CreateReview *model.Review `json:"createReview"`
}

// GetCreateReview returns CreateReviewResponse.CreateReview, and is useful for accessing the field via an interface.
func (v *CreateReviewResponse) GetCreateReview() *model.Review { return v.CreateReview }

// CreateSchoolResponse is returned by CreateSchool on success.
type CreateSchoolResponse struct {
	CreateSchool *model.School `json:"createSchool"`
}

// GetCreateSchool returns CreateSchoolResponse.CreateSchool, and is useful for accessing the field via an interface.
func (v *CreateSchoolResponse) GetCreateSchool() *model.School { return v.CreateSchool }

// GetCoursesResponse is returned by GetCourses on success.
type GetCoursesResponse struct {
	School *GetCoursesSchool `json:"school"`
}

// GetSchool returns GetCoursesResponse.School, and is useful for accessing the field via an interface.
func (v *GetCoursesResponse) GetSchool() *GetCoursesSchool { return v.School }

// GetCoursesSchool includes the requested fields of the GraphQL type School.
type GetCoursesSchool struct {
	Courses GetCoursesSchoolCoursesCourseConnection `json:"courses"`
}

// GetCourses returns GetCoursesSchool.Courses, and is useful for accessing the field via an interface.
func (v *GetCoursesSchool) GetCourses() GetCoursesSchoolCoursesCourseConnection { return v.Courses }

// GetCoursesSchoolCoursesCourseConnection includes the requested fields of the GraphQL type CourseConnection.
type GetCoursesSchoolCoursesCourseConnection struct {
	Courses    []model.Course `json:"courses"`
	PageInfo   model.PageInfo `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

// GetCourses returns GetCoursesSchoolCoursesCourseConnection.Courses, and is useful for accessing the field via an interface.
func (v *GetCoursesSchoolCoursesCourseConnection) GetCourses() []model.Course { return v.Courses }

// GetPageInfo returns GetCoursesSchoolCoursesCourseConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *GetCoursesSchoolCoursesCourseConnection) GetPageInfo() model.PageInfo { return v.PageInfo }

// GetTotalCount returns GetCoursesSchoolCoursesCourseConnection.TotalCount, and is useful for accessing the field via an interface.
func (v *GetCoursesSchoolCoursesCourseConnection) GetTotalCount() int { return v.TotalCount }

// GetProfessorByRMPIdResponse is returned by GetProfessorByRMPId on success.
type GetProfessorByRMPIdResponse struct {
	ProfessorByRMPId *model.Professor `json:"professorByRMPId"`
}

// GetProfessorByRMPId returns GetProfessorByRMPIdResponse.ProfessorByRMPId, and is useful for accessing the field via an interface.
func (v *GetProfessorByRMPIdResponse) GetProfessorByRMPId() *model.Professor {
	return v.ProfessorByRMPId
}

// GetProfessorsProfessorsProfessorConnection includes the requested fields of the GraphQL type ProfessorConnection.
type GetProfessorsProfessorsProfessorConnection struct {
	Professors []GetProfessorsProfessorsProfessorConnectionProfessorsProfessor `json:"professors"`
	PageInfo   model.PageInfo                                                  `json:"pageInfo"`
	TotalCount int                                                             `json:"totalCount"`
}

// GetProfessors returns GetProfessorsProfessorsProfessorConnection.Professors, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnection) GetProfessors() []GetProfessorsProfessorsProfessorConnectionProfessorsProfessor {
	return v.Professors
}

// GetPageInfo returns GetProfessorsProfessorsProfessorConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnection) GetPageInfo() model.PageInfo { return v.PageInfo }

// GetTotalCount returns GetProfessorsProfessorsProfessorConnection.TotalCount, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnection) GetTotalCount() int { return v.TotalCount }

// GetProfessorsProfessorsProfessorConnectionProfessorsProfessor includes the requested fields of the GraphQL type Professor.
type GetProfessorsProfessorsProfessorConnectionProfessorsProfessor struct {
	Id        string                                                                               `json:"id"`
	FirstName string                                                                               `json:"firstName"`
	LastName  string                                                                               `json:"lastName"`
	Reviews   GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection `json:"reviews"`
	Teaches   GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection `json:"teaches"`
}

// GetId returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessor.Id, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessor) GetId() string { return v.Id }

// GetFirstName returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessor.FirstName, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessor) GetFirstName() string {
	return v.FirstName
}

// GetLastName returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessor.LastName, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessor) GetLastName() string {
	return v.LastName
}

// GetReviews returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessor.Reviews, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessor) GetReviews() GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection {
	return v.Reviews
}

// GetTeaches returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessor.Teaches, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessor) GetTeaches() GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection {
	return v.Teaches
}

// GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection includes the requested fields of the GraphQL type ReviewConnection.
type GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection struct {
	Reviews []model.Review `json:"reviews"`
}

// GetReviews returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection.Reviews, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessorReviewsReviewConnection) GetReviews() []model.Review {
	return v.Reviews
}

// GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection includes the requested fields of the GraphQL type CourseConnection.
type GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection struct {
	Courses    []model.Course `json:"courses"`
	PageInfo   model.PageInfo `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

// GetCourses returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection.Courses, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection) GetCourses() []model.Course {
	return v.Courses
}

// GetPageInfo returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection.PageInfo, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection) GetPageInfo() model.PageInfo {
	return v.PageInfo
}

// GetTotalCount returns GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection.TotalCount, and is useful for accessing the field via an interface.
func (v *GetProfessorsProfessorsProfessorConnectionProfessorsProfessorTeachesCourseConnection) GetTotalCount() int {
	return v.TotalCount
}

// GetProfessorsResponse is returned by GetProfessors on success.
type GetProfessorsResponse struct {
	Professors GetProfessorsProfessorsProfessorConnection `json:"professors"`
}

// GetProfessors returns GetProfessorsResponse.Professors, and is useful for accessing the field via an interface.
func (v *GetProfessorsResponse) GetProfessors() GetProfessorsProfessorsProfessorConnection {
	return v.Professors
}

// GetSchoolsResponse is returned by GetSchools on success.
type GetSchoolsResponse struct {
	Schools GetSchoolsSchoolsSchoolConnection `json:"schools"`
}

// GetSchools returns GetSchoolsResponse.Schools, and is useful for accessing the field via an interface.
func (v *GetSchoolsResponse) GetSchools() GetSchoolsSchoolsSchoolConnection { return v.Schools }

// GetSchoolsSchoolsSchoolConnection includes the requested fields of the GraphQL type SchoolConnection.
type GetSchoolsSchoolsSchoolConnection struct {
	Schools    []model.School `json:"schools"`
	TotalCount int            `json:"totalCount"`
}

// GetSchools returns GetSchoolsSchoolsSchoolConnection.Schools, and is useful for accessing the field via an interface.
func (v *GetSchoolsSchoolsSchoolConnection) GetSchools() []model.School { return v.Schools }

// GetTotalCount returns GetSchoolsSchoolsSchoolConnection.TotalCount, and is useful for accessing the field via an interface.
func (v *GetSchoolsSchoolsSchoolConnection) GetTotalCount() int { return v.TotalCount }

type NewCourse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// GetName returns NewCourse.Name, and is useful for accessing the field via an interface.
func (v *NewCourse) GetName() string { return v.Name }

// GetCode returns NewCourse.Code, and is useful for accessing the field via an interface.
func (v *NewCourse) GetCode() string { return v.Code }

type NewProfessor struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	RmpId     *string `json:"rmpId"`
}

// GetFirstName returns NewProfessor.FirstName, and is useful for accessing the field via an interface.
func (v *NewProfessor) GetFirstName() string { return v.FirstName }

// GetLastName returns NewProfessor.LastName, and is useful for accessing the field via an interface.
func (v *NewProfessor) GetLastName() string { return v.LastName }

// GetRmpId returns NewProfessor.RmpId, and is useful for accessing the field via an interface.
func (v *NewProfessor) GetRmpId() *string { return v.RmpId }

type NewReview struct {
	Quality    float64     `json:"quality"`
	Difficulty float64     `json:"difficulty"`
	Time       string      `json:"time"`
	Tags       []model.Tag `json:"tags"`
	Grade      model.Grade `json:"grade"`
}

// GetQuality returns NewReview.Quality, and is useful for accessing the field via an interface.
func (v *NewReview) GetQuality() float64 { return v.Quality }

// GetDifficulty returns NewReview.Difficulty, and is useful for accessing the field via an interface.
func (v *NewReview) GetDifficulty() float64 { return v.Difficulty }

// GetTime returns NewReview.Time, and is useful for accessing the field via an interface.
func (v *NewReview) GetTime() string { return v.Time }

// GetTags returns NewReview.Tags, and is useful for accessing the field via an interface.
func (v *NewReview) GetTags() []model.Tag { return v.Tags }

// GetGrade returns NewReview.Grade, and is useful for accessing the field via an interface.
func (v *NewReview) GetGrade() model.Grade { return v.Grade }

type NewSchool struct {
	Name string `json:"name"`
}

// GetName returns NewSchool.Name, and is useful for accessing the field via an interface.
func (v *NewSchool) GetName() string { return v.Name }

// RegisterCourseResponse is returned by RegisterCourse on success.
type RegisterCourseResponse struct {
	RegisterProfessorForCourse bool `json:"registerProfessorForCourse"`
}

// GetRegisterProfessorForCourse returns RegisterCourseResponse.RegisterProfessorForCourse, and is useful for accessing the field via an interface.
func (v *RegisterCourseResponse) GetRegisterProfessorForCourse() bool {
	return v.RegisterProfessorForCourse
}

type Semester string

const (
	SemesterFall   Semester = "FALL"
	SemesterSpring Semester = "SPRING"
	SemesterSummer Semester = "SUMMER"
)

type TermInput struct {
	Year     int      `json:"year"`
	Semester Semester `json:"semester"`
}

// GetYear returns TermInput.Year, and is useful for accessing the field via an interface.
func (v *TermInput) GetYear() int { return v.Year }

// GetSemester returns TermInput.Semester, and is useful for accessing the field via an interface.
func (v *TermInput) GetSemester() Semester { return v.Semester }

// __CreateCourseInput is used internally by genqlient
type __CreateCourseInput struct {
	SchoolId string    `json:"schoolId"`
	Input    NewCourse `json:"input"`
}

// GetSchoolId returns __CreateCourseInput.SchoolId, and is useful for accessing the field via an interface.
func (v *__CreateCourseInput) GetSchoolId() string { return v.SchoolId }

// GetInput returns __CreateCourseInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateCourseInput) GetInput() NewCourse { return v.Input }

// __CreateProfessorInput is used internally by genqlient
type __CreateProfessorInput struct {
	SchoolId string       `json:"schoolId"`
	Input    NewProfessor `json:"input"`
}

// GetSchoolId returns __CreateProfessorInput.SchoolId, and is useful for accessing the field via an interface.
func (v *__CreateProfessorInput) GetSchoolId() string { return v.SchoolId }

// GetInput returns __CreateProfessorInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateProfessorInput) GetInput() NewProfessor { return v.Input }

// __CreateReviewInput is used internally by genqlient
type __CreateReviewInput struct {
	ProfessorId string    `json:"professorId"`
	Input       NewReview `json:"input"`
}

// GetProfessorId returns __CreateReviewInput.ProfessorId, and is useful for accessing the field via an interface.
func (v *__CreateReviewInput) GetProfessorId() string { return v.ProfessorId }

// GetInput returns __CreateReviewInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateReviewInput) GetInput() NewReview { return v.Input }

// __CreateSchoolInput is used internally by genqlient
type __CreateSchoolInput struct {
	Input NewSchool `json:"input"`
}

// GetInput returns __CreateSchoolInput.Input, and is useful for accessing the field via an interface.
func (v *__CreateSchoolInput) GetInput() NewSchool { return v.Input }

// __GetCoursesInput is used internally by genqlient
type __GetCoursesInput struct {
	SchoolId string    `json:"schoolId"`
	Input    TermInput `json:"input"`
	After    *string   `json:"after"`
}

// GetSchoolId returns __GetCoursesInput.SchoolId, and is useful for accessing the field via an interface.
func (v *__GetCoursesInput) GetSchoolId() string { return v.SchoolId }

// GetInput returns __GetCoursesInput.Input, and is useful for accessing the field via an interface.
func (v *__GetCoursesInput) GetInput() TermInput { return v.Input }

// GetAfter returns __GetCoursesInput.After, and is useful for accessing the field via an interface.
func (v *__GetCoursesInput) GetAfter() *string { return v.After }

// __GetProfessorByRMPIdInput is used internally by genqlient
type __GetProfessorByRMPIdInput struct {
	RmpId         string `json:"rmpId"`
	IncludeSchool bool   `json:"includeSchool"`
}

// GetRmpId returns __GetProfessorByRMPIdInput.RmpId, and is useful for accessing the field via an interface.
func (v *__GetProfessorByRMPIdInput) GetRmpId() string { return v.RmpId }

// GetIncludeSchool returns __GetProfessorByRMPIdInput.IncludeSchool, and is useful for accessing the field via an interface.
func (v *__GetProfessorByRMPIdInput) GetIncludeSchool() bool { return v.IncludeSchool }

// __GetProfessorsInput is used internally by genqlient
type __GetProfessorsInput struct {
	SchoolId string    `json:"schoolId"`
	After    *string   `json:"after"`
	Term     TermInput `json:"term"`
}

// GetSchoolId returns __GetProfessorsInput.SchoolId, and is useful for accessing the field via an interface.
func (v *__GetProfessorsInput) GetSchoolId() string { return v.SchoolId }

// GetAfter returns __GetProfessorsInput.After, and is useful for accessing the field via an interface.
func (v *__GetProfessorsInput) GetAfter() *string { return v.After }

// GetTerm returns __GetProfessorsInput.Term, and is useful for accessing the field via an interface.
func (v *__GetProfessorsInput) GetTerm() TermInput { return v.Term }

// __RegisterCourseInput is used internally by genqlient
type __RegisterCourseInput struct {
	ProfessorId string    `json:"professorId"`
	CourseId    string    `json:"courseId"`
	Term        TermInput `json:"term"`
}

// GetProfessorId returns __RegisterCourseInput.ProfessorId, and is useful for accessing the field via an interface.
func (v *__RegisterCourseInput) GetProfessorId() string { return v.ProfessorId }

// GetCourseId returns __RegisterCourseInput.CourseId, and is useful for accessing the field via an interface.
func (v *__RegisterCourseInput) GetCourseId() string { return v.CourseId }

// GetTerm returns __RegisterCourseInput.Term, and is useful for accessing the field via an interface.
func (v *__RegisterCourseInput) GetTerm() TermInput { return v.Term }

func CreateCourse(
	ctx context.Context,
	client graphql.Client,
	schoolId string,
	input NewCourse,
) (*CreateCourseResponse, error) {
	req := &graphql.Request{
		OpName: "CreateCourse",
		Query: `
mutation CreateCourse ($schoolId: ID!, $input: NewCourse!) {
	createCourse(schoolId: $schoolId, input: $input) {
		id
	}
}
`,
		Variables: &__CreateCourseInput{
			SchoolId: schoolId,
			Input:    input,
		},
	}
	var err error

	var data CreateCourseResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func CreateProfessor(
	ctx context.Context,
	client graphql.Client,
	schoolId string,
	input NewProfessor,
) (*CreateProfessorResponse, error) {
	req := &graphql.Request{
		OpName: "CreateProfessor",
		Query: `
mutation CreateProfessor ($schoolId: ID!, $input: NewProfessor!) {
	createProfessor(schoolId: $schoolId, input: $input) {
		id
	}
}
`,
		Variables: &__CreateProfessorInput{
			SchoolId: schoolId,
			Input:    input,
		},
	}
	var err error

	var data CreateProfessorResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func CreateReview(
	ctx context.Context,
	client graphql.Client,
	professorId string,
	input NewReview,
) (*CreateReviewResponse, error) {
	req := &graphql.Request{
		OpName: "CreateReview",
		Query: `
mutation CreateReview ($professorId: ID!, $input: NewReview!) {
	createReview(professorId: $professorId, input: $input) {
		id
	}
}
`,
		Variables: &__CreateReviewInput{
			ProfessorId: professorId,
			Input:       input,
		},
	}
	var err error

	var data CreateReviewResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func CreateSchool(
	ctx context.Context,
	client graphql.Client,
	input NewSchool,
) (*CreateSchoolResponse, error) {
	req := &graphql.Request{
		OpName: "CreateSchool",
		Query: `
mutation CreateSchool ($input: NewSchool!) {
	createSchool(input: $input) {
		id
	}
}
`,
		Variables: &__CreateSchoolInput{
			Input: input,
		},
	}
	var err error

	var data CreateSchoolResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func GetCourses(
	ctx context.Context,
	client graphql.Client,
	schoolId string,
	input TermInput,
	after *string,
) (*GetCoursesResponse, error) {
	req := &graphql.Request{
		OpName: "GetCourses",
		Query: `
query GetCourses ($schoolId: ID!, $input: TermInput!, $after: String) {
	school(id: $schoolId) {
		courses(term: $input, first: 50, after: $after) {
			courses {
				id
				name
				code
			}
			pageInfo {
				hasNextPage
			}
			totalCount
		}
	}
}
`,
		Variables: &__GetCoursesInput{
			SchoolId: schoolId,
			Input:    input,
			After:    after,
		},
	}
	var err error

	var data GetCoursesResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func GetProfessorByRMPId(
	ctx context.Context,
	client graphql.Client,
	rmpId string,
	includeSchool bool,
) (*GetProfessorByRMPIdResponse, error) {
	req := &graphql.Request{
		OpName: "GetProfessorByRMPId",
		Query: `
query GetProfessorByRMPId ($rmpId: String!, $includeSchool: Boolean!) {
	professorByRMPId(rmpId: $rmpId) {
		id
		firstName
		lastName
		linked
		school @include(if: $includeSchool) {
			id
			name
		}
	}
}
`,
		Variables: &__GetProfessorByRMPIdInput{
			RmpId:         rmpId,
			IncludeSchool: includeSchool,
		},
	}
	var err error

	var data GetProfessorByRMPIdResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func GetProfessors(
	ctx context.Context,
	client graphql.Client,
	schoolId string,
	after *string,
	term TermInput,
) (*GetProfessorsResponse, error) {
	req := &graphql.Request{
		OpName: "GetProfessors",
		Query: `
query GetProfessors ($schoolId: ID!, $after: String, $term: TermInput!) {
	professors(schoolId: $schoolId, first: 50, after: $after) {
		professors {
			id
			firstName
			lastName
			reviews(first: 1) {
				reviews {
					id
				}
			}
			teaches(term: $term, first: 50) {
				courses {
					id
					code
					name
				}
				pageInfo {
					hasNextPage
				}
				totalCount
			}
		}
		pageInfo {
			hasNextPage
			endCursor
		}
		totalCount
	}
}
`,
		Variables: &__GetProfessorsInput{
			SchoolId: schoolId,
			After:    after,
			Term:     term,
		},
	}
	var err error

	var data GetProfessorsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func GetSchools(
	ctx context.Context,
	client graphql.Client,
) (*GetSchoolsResponse, error) {
	req := &graphql.Request{
		OpName: "GetSchools",
		Query: `
query GetSchools {
	schools(first: 50) {
		schools {
			id
			name
		}
		totalCount
	}
}
`,
	}
	var err error

	var data GetSchoolsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func RegisterCourse(
	ctx context.Context,
	client graphql.Client,
	professorId string,
	courseId string,
	term TermInput,
) (*RegisterCourseResponse, error) {
	req := &graphql.Request{
		OpName: "RegisterCourse",
		Query: `
mutation RegisterCourse ($professorId: ID!, $courseId: ID!, $term: TermInput!) {
	registerProfessorForCourse(professorId: $professorId, courseId: $courseId, term: $term)
}
`,
		Variables: &__RegisterCourseInput{
			ProfessorId: professorId,
			CourseId:    courseId,
			Term:        term,
		},
	}
	var err error

	var data RegisterCourseResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
