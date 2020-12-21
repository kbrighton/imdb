package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kbrighton/imdb/pkg"
	"net/http"
)

func main() {

	//Database is a singleton
	Manager = Initialize()
	r := gin.Default()

	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"body": "OK",
		})
	})

	//Versioning API
	api := r.Group("/api/v1")
	{
		api.GET("/health", GetHealth)
		api.GET("/movies/tconst/:tconst", GetTconst)
		api.GET("/movies/startYear/:startYear", GetStartYear)
		api.GET("/movies/genre/:genre", GetGenre)
	}

	//Lets migrate our data
	go MigrateData(Manager)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
