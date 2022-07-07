package main

import (
	"flag"
	"log"
	"os"
	"practice/three/database"
	"practice/three/database/entities"
	"practice/three/server"
	"sync"

	"github.com/joho/godotenv"
)

var migrate = flag.Bool("migrate", false, "Migrate table to DB")
var seed = flag.Uint("seed", 0, "Populate DB with fake data x values where -seed=x")

func main() {
	s := server.BookServer{}
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Problem with .env file %v\n", err)
	}

	s.ConnectDatabase()

	if *migrate {
		err = database.MigrateEntities()
	}

	if *seed >  0 {
		entities.SeedBooksForAuthor(*seed)
	}
	
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		log.Printf("starting gRPC API server... \n")
		s.Serve(os.Getenv("Port"))
		wg.Done()
	}()

	wg.Wait()
}
