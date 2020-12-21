package pkg

import (
	"github.com/go-pg/pg/v10"
)

var Manager *pg.DB

type Movie struct {
	Id             int      `json:"-"`
	Tconst         string   `json:"tconst"`
	TitleType      string   `json:"titleType"`
	PrimaryTitle   string   `json:"primaryTitle"`
	OriginalTitle  string   `json:"originalTitle"`
	IsAdult        bool     `json:"isAdult" pg:",use_zero"`
	StartYear      int      `json:"startYear,omitempty"`
	EndYear        int      `json:"endYear,omitempty"`
	RuntimeMinutes int      `json:"runtimeMinutes,omitempty"`
	Genres         []Genre  `json:"-" pg:"many2many:movies_genres"`
	GenreArray     []string `json:"genres" pg:"-"`
}

type Genre struct {
	Id     int     `json:"-"`
	Genre  string  `json:""`
	Movies []Movie `json:"-" pg:"many2many:movies_genres"`
}

type MoviesGenres struct {
	MovieId int
	GenreId int
}

type RawMovie struct {
	Tconst         string `json:"tconst"`
	TitleType      string `json:"titleType"`
	PrimaryTitle   string `json:"primaryTitle"`
	OriginalTitle  string `json:"originalTitle"`
	IsAdult        string `json:"isAdult" pg:",use_zero"`
	StartYear      string `json:"startYear,omitempty"`
	EndYear        string `json:"endYear,omitempty"`
	RuntimeMinutes string `json:"runtimeMinutes,omitempty"`
	Genres         string `json:"genres"`
}
