package pkg

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tkanos/gonfig"
)

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
	db.AddQueryHook(pgdebug.DebugHook{
		Verbose: true,
	})

	//db, err := gorm.Open("postgres", url)
	//db.LogMode(true)
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic("Could not connect to database")
	}

	orm.RegisterTable((*MoviesGenres)(nil))

	//db.AutoMigrate(&Movies{}, &Genres{})
	/*
		//Prepopulating Genres to get past some GORM
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
			var genre Genres

			genre.Genre = element

			db.Create(&genre)
		}

		//I wish this didn't take forever
		movieFile, err := os.Open("./data.tsv")
		if err != nil {
			panic("I could not find the movies tsv")
		}

		reader := csv.NewReader(movieFile)
		reader.Comma = '\t'
		reader.LazyQuotes = true

		//Skip the header
		_, readerErr := reader.Read()
		if readerErr != nil {
			panic("Could not read the header")
		}

		//var movies []Movies

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}

			var movie Movie
			var genres []Genre
			var genre Genre
			var moviesgenres MoviesGenres

			movie.Tconst = row[0]
			movie.TitleType = row[1]
			movie.PrimaryTitle = row[2]
			movie.OriginalTitle = row[3]
			movie.IsAdult, _ = strconv.ParseBool(row[4])

			if row[5] != "\\N" {
				movie.StartYear, _ = strconv.Atoi(row[5])
			}

			if row[6] != "\\N" {
				movie.EndYear, _ = strconv.Atoi(row[6])
			}

			if row[7] != "\\N" {
				movie.RuntimeMinutes, _ = strconv.Atoi(row[7])
			}

			db.Model(&movie).Insert()

			genres = genres[:0]
			genreString := strings.Split(row[8], ",")

			for _, stringName := range genreString {
				if stringName == "\\N" {
					break
				}
				db.Model(&genre).Where("genre = ?", stringName).Select()
				moviesgenres.GenreId = genre.Id
				moviesgenres.MovieId = movie.Id
				db.Model(&moviesgenres).Insert()
			}

		}

	*/

	return db
}
