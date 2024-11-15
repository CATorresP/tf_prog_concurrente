package master

import "fmt"

func MappRatingsClient(ratings []MovieRatingsClient, moviesTitles *MoviesTitles) []float64 {
	numMovies := len(moviesTitles.Title)
	arr := make([]float64, numMovies)
	for _, rating := range ratings {
		arr[rating.MovieId] = float64(rating.Rating)
	}
	return arr
}

func Banner() {
	fmt.Println("  ____                 _                            _             _   _             ")
	fmt.Println(" |  _ \\ ___  __ _  ___| |_ ___  _ __ ___   ___  __| | __ _ _ __ | |_(_) ___  _ __  ")
	fmt.Println(" | |_) / _ \\/ _` |/ __| __/ _ \\| '_ ` _ \\ / _ \\/ _` |/ _` | '_ \\| __| |/ _ \\| '_ \\ ")
	fmt.Println(" |  _ <  __/ (_| | (__| || (_) | | | | | |  __/ (_| | (_| | | | | |_| | (_) | | | |")
	fmt.Println(" |_| \\_\\___|\\__,_|\\___|\\__\\___/|_| |_| |_|\\___|\\__,_|\\__,_|_| |_|\\__|_|\\___/|_| |_|")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(" __          __  _                            _          __  __                                    ")
	fmt.Println(" \\ \\        / / | |                          | |        |  \\/  |                                   ")
	fmt.Println("  \\ \\  /\\  / /__| | ___ ___  _ __ ___   ___  | |_ ___   | \\  / | __ _ _ __   __ _  __ _  ___ _ __  ")
	fmt.Println("   \\ \\/  \\/ / _ \\ |/ __/ _ \\| '_ ` _ \\ / _ \\ | __/ _ \\  | |\\/| |/ _` | '_ \\ / _` |/ _` |/ _ \\ '__| ")
	fmt.Println("    \\  /\\  /  __/ | (_| (_) | | | | | |  __/ | || (_) | | |  | | (_| | | | | (_| | (_| |  __/ |    ")
	fmt.Println("     \\/  \\/ \\___|_|\\___\\___/|_| |_| |_|\\___|  \\__\\___/  |_|  |_|\\__,_|_| |_|\\__,_|\\__, |\\___|_|    ")
	fmt.Println("                                                                                     __/ |         ")
	fmt.Println("                                                                                    |___/          ")
	fmt.Println("--------------------------------------------------------------------------------")
}

func (master *Master) getMoviesByGenre(genre int) []MovieTitleWithID {
	moviesTitles := MoviesTitles{Title: master.movieTitles}
	moviesGenreIds := MoviesGenreIds{MoviesGenreIds: master.movieGenreIds}
	var moviesGenres []MovieTitleWithID
	for i := 0; i < len(moviesTitles.Title); i++ {
		for _, genreId := range moviesGenreIds.MoviesGenreIds[i] {
			if genreId == genre {
				movie := MovieTitleWithID{
					Title: moviesTitles.Title[i],
					Id:    i,
				}
				moviesGenres = append(moviesGenres, movie)
				break
			}
		}
	}
	return moviesGenres
}

func MappperMovieGenres(idTitle int, moviesTitles *MoviesTitles, moviesGenreIds *MoviesGenreIds, genres *Genres) MovieGenres {
	movieGenres := MovieGenres{}
	movieGenres.Name = moviesTitles.Title[idTitle]
	for _, genreId := range moviesGenreIds.MoviesGenreIds[idTitle] {
		movieGenres.Genres = append(movieGenres.Genres, genres.Genresname[genreId])
	}
	return movieGenres
}

func (master *Master) getMoviesGenres() []MovieGenres {
	moviesTitles := MoviesTitles{Title: master.movieTitles}
	moviesGenreIds := MoviesGenreIds{MoviesGenreIds: master.movieGenreIds}
	genres := Genres{Genresname: master.movieGenreNames}
	var moviesGenres []MovieGenres
	for i := 0; i < len(moviesTitles.Title); i++ {
		movieGenres := MappperMovieGenres(i, &moviesTitles, &moviesGenreIds, &genres)
		moviesGenres = append(moviesGenres, movieGenres)
	}
	return moviesGenres
}

func getComment(rating, max, min, mean float64) string {
	highThreshold1 := mean + (max-mean)*0.90
	highThreshold2 := mean + (max-mean)*0.60
	lowThreshold1 := mean - (mean-min)*0.60
	lowThreshold2 := mean - (mean-min)*0.90

	if rating > highThreshold1 {
		return "Altamente Recomendado. Muy por encima de la media"
	} else if rating > highThreshold2 {
		return "Recomendado. Bastante por encima de la media"
	} else if rating > mean {
		return "Ligeramente Recomendado. Por encima de la media"
	} else if rating > lowThreshold1 {
		return "Ligeramente No Recomendado. Justo por debajo de la media"
	} else if rating > lowThreshold2 {
		return "Poco Recomendado. Bastante por debajo de la media"
	} else {
		return "Muy Poco Recomendado. Muy por debajo de la media"
	}
}
