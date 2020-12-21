package pkg

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Sets up the initial database connection and creates tables if they don't exist
func Initialize() *pg.DB {
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		fmt.Println(err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     configuration.Host,
		User:     configuration.User,
		Password: configuration.Password,
		Database: configuration.Database,
	})

	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic("Could not connect to database")
	}

	orm.RegisterTable((*MoviesGenres)(nil))

	models := []interface{}{
		(*Movie)(nil),
		(*Genre)(nil),
		(*MoviesGenres)(nil),
		(*RawMovie)(nil),
	}

	for _, model := range models {
		err = db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			panic(err)
		}
	}

	//Prepopulating Genres to get past some go-pg
	//annoyances with many2many associations
	genreFile, err := os.Open("./genres.csv")
	if err != nil {
		panic("I could not find the genres csv")
	}

	gen, err := ioutil.ReadAll(genreFile)
	if err != nil {
		panic("Issue with readall on genre file")
	}

	genCsv := string(gen)

	genreString := strings.Split(genCsv, ",")

	for _, element := range genreString {
		var genre Genre

		genre.Genre = element

		db.Model(&genre).Insert()
	}

	//I wish this didn't take forever
	//CopyFrom makes it not take forever!
	movieFile, err := os.Open("./data.tsv")
	if err != nil {
		panic("I could not find the movies tsv")
	}

	_, err = db.CopyFrom(movieFile, "COPY raw_movies FROM STDIN DELIMITER '\t'")
	if err != nil {
		panic(err)
	}

	return db
}

//This will convert the initial seeded data to a format we can use
//Runs as goroutine
func MigrateDate(db *pg.DB) {
	var rawmovie []RawMovie

	db.Model(&rawmovie).Where("tconst <> 'tconst'").Select()

	for _, iterator := range rawmovie {
		var movie Movie
		var genres []Genre
		var genre Genre
		var moviesgenres MoviesGenres

		movie.Tconst = iterator.Tconst
		movie.TitleType = iterator.TitleType
		movie.PrimaryTitle = iterator.PrimaryTitle
		movie.OriginalTitle = iterator.OriginalTitle
		movie.IsAdult, _ = strconv.ParseBool(iterator.IsAdult)
		movie.StartYear, _ = strconv.Atoi(iterator.StartYear)
		movie.EndYear, _ = strconv.Atoi(iterator.EndYear)
		movie.RuntimeMinutes, _ = strconv.Atoi(iterator.RuntimeMinutes)

		db.Model(&movie).Insert()

		genres = genres[:0]
		genreString := strings.Split(iterator.Genres, ",")

		for _, stringName := range genreString {
			db.Model(&genre).Where("genre = ?", stringName).Select()
			moviesgenres.GenreId = genre.Id
			moviesgenres.MovieId = movie.Id
			db.Model(&moviesgenres).Insert()
		}
	}
}
