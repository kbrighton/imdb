package pkg

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Sets up the initial database connection and creates tables if they don't exist
func Initialize() *pg.DB {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("json")

	var config Configurations

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Couldn't decode into struct")
	}

	db := pg.Connect(&pg.Options{
		Addr:     config.IMDB_DB_HOST,
		User:     config.IMDB_DB_USER,
		Password: config.IMDB_DB_PASSWORD,
		Database: config.IMDB_DB_NAME,
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
	movieFile, err := os.Open("./title.basics.tsv")
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
func MigrateData(db *pg.DB) {
	err := db.Model((*RawMovie)(nil)).Where("tconst <> 'tconst'").ForEach(func(rm *RawMovie) error {
		var movie Movie
		var genres []Genre
		var genre Genre
		var moviesgenres MoviesGenres

		movie.Tconst = rm.Tconst
		movie.TitleType = rm.TitleType
		movie.PrimaryTitle = rm.PrimaryTitle
		movie.OriginalTitle = rm.OriginalTitle
		movie.IsAdult, _ = strconv.ParseBool(rm.IsAdult)
		movie.StartYear, _ = strconv.Atoi(rm.StartYear)
		movie.EndYear, _ = strconv.Atoi(rm.EndYear)
		movie.RuntimeMinutes, _ = strconv.Atoi(rm.RuntimeMinutes)

		db.Model(&movie).Insert()

		genres = genres[:0]
		genreString := strings.Split(rm.Genres, ",")

		for _, stringName := range genreString {
			db.Model(&genre).Where("genre = ?", stringName).Select()
			moviesgenres.GenreId = genre.Id
			moviesgenres.MovieId = movie.Id
			db.Model(&moviesgenres).Insert()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
