package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kbrighton/imdb/pkg"
)

func main() {

	//Database is a singleton
	Manager = Initialize()
	r := gin.Default()

	//Versioning API
	api := r.Group("/api/v1")
	{
		api.GET("/health", GetHealth)
		api.GET("/movies/tconst/:tconst", GetTconst)
		api.GET("/movies/startYear/:startYear", GetStartYear)
		api.GET("/movies/genre/:genre", GetGenre)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
