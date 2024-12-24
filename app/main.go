package main

import (
	"SmartRecipe/api/server"
	"SmartRecipe/database"
)

func main() {
	database.NewDatabaseInstance()
	//scrapers.StartScraping()
	s := server.APIServer{Addr: ":8080"}
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
