package main

import (
	"context"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"os"
	"scraper/database"
)

func main() {

	dgraph := database.Connect(os.Getenv("DATABASE_URL"), os.Getenv("API_KEY"))

	txn := dgraph.NewTxn()
	rdf := database.SchoolToRDF(database.GetValencia())
	log.Println("rdf=", rdf)
	response, err := txn.Mutate(context.Background(), &api.Mutation{SetNquads: []byte(rdf)})
	if err != nil {
		log.Fatal(err)
	}
	err = txn.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response=", response)
}
