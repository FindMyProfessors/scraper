package database

import (
	"context"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"scraper/models"
)

var (
	Dgraph *dgo.Dgraph
)

func Connect(url string, apiKey string) *dgo.Dgraph {
	if Dgraph != nil {
		return Dgraph
	}
	var conn *grpc.ClientConn
	var err error
	if apiKey == "" {
		dialOpts := append([]grpc.DialOption{},
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
		conn, err = grpc.Dial(url, dialOpts...)
	} else {
		conn, err = dgo.DialCloud(url, apiKey)
	}
	if err != nil {
		fmt.Println("Error connecting to Dgraph")
		log.Fatal(err)
	}
	Dgraph = dgo.NewDgraphClient(api.NewDgraphClient(conn))
	return Dgraph
}

func MutateDatabase(school models.School) (*api.Response, error) {
	txn := Dgraph.NewTxn()
	rdf := SchoolToRDF(school)
	log.Println("rdf=", rdf)
	ctx := context.Background()
	response, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(rdf)})
	if err != nil {
		return nil, err
	}
	err = txn.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func SchoolToRDF(school models.School) string {

	var rdf = ""

	schoolRDFId := school.RDFId()

	rdf += fmt.Sprintf("_:%s <School.name> \"%s\" .\n", schoolRDFId, school.Name)
	rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", schoolRDFId, "School")
	//school = models.CollectCourses(school)
	for _, course := range school.Courses {
		rdf += fmt.Sprintf("_:%s <School.courses> _:%s .\n", schoolRDFId, course.Code)
	}
	for _, professor := range school.Professors {
		if len(professor.FirstName) == 0 {
			log.Println("Professor name length = 0, ", professor)
		} else {
			rdf += fmt.Sprintf("_:%s <School.professors> _:%s .\n", schoolRDFId, professor.RDFId())
		}
	}

	rdf += "\n"

	for _, professor := range school.Professors {
		rdfId := professor.RDFId()
		rdf += fmt.Sprintf("_:%s <Professor.name> \"%s\" .\n", rdfId, professor.Name())
		rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", rdfId, "Professor")
		rdf += fmt.Sprintf("_:%s <Professor.totalRatings> \"%d\" .\n", rdfId, professor.TotalRatings)
		rdf += fmt.Sprintf("_:%s <Professor.rating> \"%f\" .\n", rdfId, professor.Rating)
		rdf += fmt.Sprintf("_:%s <Professor.school> _:%s .\n", rdfId, schoolRDFId)
		for _, courseTaught := range professor.Teaches {
			rdf += fmt.Sprintf("_:%s <Professor.teaches> _:%s .\n", rdfId, courseTaught.Code)
		}
		rdf += "\n"
	}

	for _, course := range school.Courses {
		rdf += fmt.Sprintf("_:%s <dgraph.type> \"%s\" .\n", course.Code, "Course")
		rdf += fmt.Sprintf("_:%s <Course.code> \"%s\" .\n", course.Code, course.Code)
		rdf += fmt.Sprintf("_:%s <Course.name> \"%s\" .\n", course.Code, course.Name)
		rdf += fmt.Sprintf("_:%s <Course.school> _:%s .\n", course.Code, schoolRDFId)

		for _, professor := range course.Professors {
			rdf += fmt.Sprintf("_:%s <Course.professors> _:%s .\n", course.Code, professor)
		}
		rdf += "\n"
	}

	return rdf
}
