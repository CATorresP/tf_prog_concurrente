package master

func (master *Master) genresHandler(response *Genres) {
	*response = Genres{Genresname: master.movieGenreNames}
}

func (master *Master) MoviesGenresHandler(response *[]MovieGenres) {
	*response = master.getMoviesGenres()
}

func (master *Master) moviesTitlesHandler(response *MoviesTitles) {
	*response = MoviesTitles{Title: master.movieTitles}
}

type MoviesByGenderRequest struct {
	Id int `json:"id"`
}

func (master *Master) getMoviesByGenresHandler(request *MoviesByGenderRequest, response *[]MovieTitleWithID) error {
	*response = master.getMoviesByGenre(request.Id)
	return nil
}
