package master

type MovieTitleWithID struct {
	Title string `json:"title"`
	Id    int    `json:"id"`
}

type Genres struct {
	Genresname []string `json:"movieGenreNames"`
}

type MoviesTitles struct {
	Title []string `json:"movieTitles"`
}

type MoviesGenreIds struct {
	MoviesGenreIds [][]int `json:"movieGenreIds"`
}

type MovieGenres struct {
	Name   string   `json:"Name"`
	Genres []string `json:"genres"`
}

type ClientRecRequest struct {
	UserId   int       `json:"userId"`
	Quantity int       `json:"quantity"`
	GenreIds []int     `json:"genreIds"`
	Ratings  []float64 `json:"ratings"`
}

type MovieRatingsClient struct {
	MovieId int `json:"movieId"`
	Rating  int `json:"rating"`
}

type ClientRecToSend struct {
	UserId        int                  `json:"userId"`
	Quantity      int                  `json:"quantity"`
	GenreIds      []int                `json:"genreIds"`
	MoviesRatings []MovieRatingsClient `json:"moviesRatings"`
}
