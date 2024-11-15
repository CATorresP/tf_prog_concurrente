package master

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (master *Master) genresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	genres := Genres{Genresname: master.movieGenreNames}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (master *Master) MoviesGenresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	genres := master.getMoviesGenres()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (master *Master) moviesTitlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	MoviesTitles := MoviesTitles{Title: master.movieTitles}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoviesTitles)
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
