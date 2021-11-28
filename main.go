package main

import (
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"os"
	"strings"
)

func main() {

	databaseUrl := os.Getenv("DATABASE_URL")
	log.Println("databaseUrl=", databaseUrl)
	conn, err := dgo.DialCloud(databaseUrl, os.Getenv("API_KEY"))
	if err != nil {
		fmt.Println("Error connecting to Dgraph")
		log.Fatal(err)
	}
	ctx := context.Background()

	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	txn := dgraphClient.NewTxn()
	rdf := SchoolToRDF(GetValencia())
	log.Println("rdf=", rdf)
	response, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(rdf)})
	if err != nil {
		log.Fatal(err)
	}
	err = txn.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response=", response)
}

func GetCourse(code string, name string) Course {
	return Course{
		Code:       code,
		Name:       name,
		Professors: []Professor{},
	}
}

func GetProfessor(name string, courses []Course) Professor {
	return Professor{
		Name:         name,
		Teaches:      courses,
		TotalRatings: 10,
		Rating:       4.0,
	}
}

func GetValencia() School {
	return School{
		Name: "Valencia",
		Professors: []Professor{GetProfessor("natalie_angelis", []Course{
			GetCourse("MAC1105H", "College Algebra Honors"),
			GetCourse("MAC1114", "College Trigonometry"),
			GetCourse("MAC2311", "Calculus 1"),
		})},
	}
}

func SchoolToRDF(school School) string {

	var rdf = ""

	schoolRDFId := school.RDFId()

	rdf += fmt.Sprintf("_:%s <School.name> \"%s\" .\n", schoolRDFId, school.Name)
	rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", schoolRDFId, "School")
	courses := CollectCourses(&school)

	for code := range courses {
		rdf += fmt.Sprintf("_:%s <School.courses> _:%s .\n", schoolRDFId, code)
	}
	for _, professor := range school.Professors {
		formattedProfessorName := strings.Replace(strings.ToLower(professor.Name), " ", "_", -1)
		rdf += fmt.Sprintf("_:%s <School.professors> _:%s .\n", schoolRDFId, formattedProfessorName)
	}

	rdf += "\n"

	for _, professor := range school.Professors {
		rdfId := professor.RDFId()
		rdf += fmt.Sprintf("_:%s <Professor.name> \"%s\" .\n", rdfId, professor.Name)
		rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", rdfId, "Professor")
		rdf += fmt.Sprintf("_:%s <Professor.totalRatings> \"%d\" .\n", rdfId, professor.TotalRatings)
		rdf += fmt.Sprintf("_:%s <Professor.rating> \"%f\" .\n", rdfId, professor.Rating)
		rdf += fmt.Sprintf("_:%s <Professor.school> _:%s .\n", rdfId, schoolRDFId)
		for _, courseTaught := range professor.Teaches {
			rdf += fmt.Sprintf("_:%s <Professor.teaches> _:%s .\n", rdfId, courseTaught.Code)
		}
		rdf += "\n"
	}

	for _, course := range courses {
		rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", course.Code, "Course")
		rdf += fmt.Sprintf("_:%s <Course.code> \"%s\" .\n", course.Code, course.Code)
		rdf += fmt.Sprintf("_:%s <Course.name> \"%s\" .\n", course.Code, course.Name)
		rdf += fmt.Sprintf("_:%s <Course.school> _:%s .\n", course.Code, schoolRDFId)

		for _, professor := range course.Professors {
			rdfId := professor.RDFId()
			rdf += fmt.Sprintf("_:%s <Course.professors> _:%s .\n", course.Code, rdfId)
		}
		rdf += "\n"
	}

	return rdf
}

func CollectCourses(school *School) map[string]Course {
	courses := make(map[string]Course)
	for _, professor := range school.Professors {
		for _, course := range professor.Teaches {
			course.Professors = append(course.Professors, professor)
			courses[course.Code] = course
		}
	}
	return courses
}
