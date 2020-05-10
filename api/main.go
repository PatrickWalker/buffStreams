package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	controllers "github.com/buffup/api/controllers"
	helpers "github.com/buffup/api/helpers"
	migrations "github.com/buffup/api/migrations"

	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/julienschmidt/httprouter"
)

func main() {

	conf := helpers.GetConfig()
	err := migrations.Migrate(conf)
	if err != nil {
		log.Println("Migration Failed. Exiting")
		log.Println(err)
		os.Exit(1)
	}
	//create the required Mux Style routing

	StreamsCon := controllers.NewStreamsController(conf)
	QuestionsCon := controllers.NewQuestionsController(conf)

	router := httprouter.New()
	router.GET("/streams", StreamsCon.List)
	router.GET("/questions/:questionID", QuestionsCon.View)
	log.Println("Starting API Server on port 1323")
	log.Fatal(http.ListenAndServe(":1323", handlers.CompressHandler(router)))
}
