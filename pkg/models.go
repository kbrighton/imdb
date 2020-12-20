package pkg

import (
	"github.com/go-pg/pg/v10"
)

var Manager *pg.DB

type Movie struct {
	Id             int      `json:"-"`
	Tconst         string   `json:"tconst"`
	TitleType      string   `json:"title_type"`
	PrimaryTitle   string   `json:"primary_title"`
	OriginalTitle  string   `json:"original_title"`
	IsAdult        bool     `json:"is_adult" pg:",use_zero"`
	StartYear      int      `json:"start_year,omitempty"`
	EndYear        int      `json:"end_year,omitempty"`
	RuntimeMinutes int      `json:"runtime_minutes,omitempty"`
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
