package master

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (master *Master) genresHandler(response *Genres) {
	*response = Genres{Genresname: master.movieGenreNames}
}

func (master *Master) MoviesGenresHandler(response *[]MovieGenres) {
	*response = master.getMoviesGenres()
}

func (master *Master) moviesTitlesHandler(response *MoviesTitles) {
	*response = MoviesTitles{Title: master.movieTitles}
}

func (master *Master) getMoviesByGenresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	genre := r.URL.Query().Get("id")
	id, err := strconv.Atoi(genre)
	if err != nil {
		http.Error(w, "Error al convertir el id del género a entero", http.StatusBadRequest)
		return
	}
	movies := master.getMoviesByGenre(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
