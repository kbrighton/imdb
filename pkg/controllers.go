package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Configuration struct {
	Password string
	User     string
	Host     string
	Port     uint16
	Database string
}

func GetTconst(c *gin.Context) {

	var movie Movie
	genreName := make([]string, 0)

	tconst := c.Param("tconst")

	err := Manager.Model(&movie).Relation("Genres").Where("tconst = ?", tconst).Select()
	if err != nil {
		panic(err)
	}
	for _, iterator := range movie.Genres {
		genreName = append(genreName, iterator.Genre)
	}

	movie.GenreArray = genreName

	c.JSON(http.StatusOK, gin.H{
		"results": movie,
	})

}

func GetStartYear(c *gin.Context) {

	var movie []Movie

	startYear, _ := strconv.Atoi(c.Param("startYear"))

	//The first movie was made in 1888, anything earlier than that is...well, not correct
	if startYear < 1888 {
		c.JSON(http.StatusNotFound, gin.H{
			"body": "There are no movies earlier than 1888",
		})
	} else {
		count, err := Manager.Model(&movie).Relation("Genres").Where("start_year = ?", startYear).SelectAndCount()
		if err != nil {
			panic(err)
		}
		for movIdx, _ := range movie {
			genreName := make([]string, 0)

			for _, iterator := range movie[movIdx].Genres {
				genreName = append(genreName, iterator.Genre)
			}

			movie[movIdx].GenreArray = genreName
		}

		c.JSON(http.StatusOK, gin.H{
			"count":   count,
			"results": movie,
		})
	}

}

//This is terrible, but go-pg is being difficult
func GetGenre(c *gin.Context) {
	var genres Genre
	var movieIds []string
	var movie []Movie

	genre := strings.Title(c.Param("genre"))

	err := Manager.Model(&genres).Relation("Movies").Where("genre = ?", genre).Select()
	if err != nil {
		panic(err)
	}

	for _, iterator := range genres.Movies {
		movieIds = append(movieIds, strconv.Itoa(iterator.Id))
	}

	csvIds := strings.Join(movieIds, ",")

	count, err := Manager.Model(&movie).Relation("Genres").Where("id in (" + csvIds + ")").SelectAndCount()
	if err != nil {
		panic(err)
	}

	for movIdx, _ := range movie {
		genreName := make([]string, 0)

		for _, iterator := range movie[movIdx].Genres {
			genreName = append(genreName, iterator.Genre)
		}

		movie[movIdx].GenreArray = genreName
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   count,
		"results": movie,
	})

}

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"body": "OK",
	})
}
