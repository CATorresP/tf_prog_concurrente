package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"recommendation-service/syncutils"
	"strconv"
)

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

func LoadGenres(genres *Genres) error {
	err := syncutils.LoadJsonFile("config/master.json", &genres)
	if err != nil {
		return fmt.Errorf("Error al cargar el archivo de configuración de géneros: %v", err)
	}
	return nil
}

func LoadMoviesTitles(movies *MoviesTitles) error {
	err := syncutils.LoadJsonFile("config/master.json", &movies)
	if err != nil {
		return fmt.Errorf("Error al cargar el archivo de configuración de películas: %v", err)
	}
	return nil
}

func LoadMoviesGenreIds(moviesGenreIds *MoviesGenreIds) error {
	err := syncutils.LoadJsonFile("config/master.json", &moviesGenreIds)
	if err != nil {
		return fmt.Errorf("Error al cargar el archivo de configuración de películas: %v", err)
	}
	return nil
}

func MappperMovieGenres(idTitle int, moviesTitles *MoviesTitles, moviesGenreIds *MoviesGenreIds, genres *Genres) MovieGenres {
	movieGenres := MovieGenres{}
	movieGenres.Name = moviesTitles.Title[idTitle]
	for _, genreId := range moviesGenreIds.MoviesGenreIds[idTitle] {
		movieGenres.Genres = append(movieGenres.Genres, genres.Genresname[genreId])
	}
	return movieGenres
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

func createClientRecRequest(userId int, quantity int, genreIds []int) ClientRecRequest {
	return ClientRecRequest{
		UserId:   userId,
		Quantity: quantity,
		GenreIds: genreIds,
	}
}

func getMoviesTitles() []string {
	moviesTitles := MoviesTitles{}
	LoadMoviesTitles(&moviesTitles)
	return moviesTitles.Title
}

func moviesTitlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	titles := getMoviesTitles()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(titles)
}

func genresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	genres := Genres{}
	LoadGenres(&genres)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func getMoviesGenres() []MovieGenres {
	moviesTitles := MoviesTitles{}
	LoadMoviesTitles(&moviesTitles)
	moviesGenreIds := MoviesGenreIds{}
	LoadMoviesGenreIds(&moviesGenreIds)
	genres := Genres{}
	LoadGenres(&genres)
	var moviesGenres []MovieGenres
	for i := 0; i < len(moviesTitles.Title); i++ {
		movieGenres := MappperMovieGenres(i, &moviesTitles, &moviesGenreIds, &genres)
		moviesGenres = append(moviesGenres, movieGenres)
	}
	return moviesGenres
}

func MoviesGenresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	genres := getMoviesGenres()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
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

func MappRatingsClient(ratings []MovieRatingsClient, moviesTitles *MoviesTitles) []float64 {
	numMovies := len(moviesTitles.Title)
	arr := make([]float64, numMovies)
	for _, rating := range ratings {
		arr[rating.MovieId] = float64(rating.Rating)
	}
	return arr
}

func sendRequestToRecommendationService(clientRecRequest ClientRecRequest) {
	masterIp := "172.0.20.3"
	masterPort := 9000
	url := fmt.Sprintf("http://%s:%d/recommendation", masterIp, masterPort)
	jsonValue, _ := json.Marshal(clientRecRequest)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalf("Error al enviar la petición al servicio de recomendaciones: %v", err)
	}

	defer resp.Body.Close()
}

func sendRequestToRecommendationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var clientToSend ClientRecToSend
	err := json.NewDecoder(r.Body).Decode(&clientToSend)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición", http.StatusBadRequest)
		return
	}
	moviesTitles := MoviesTitles{}
	LoadMoviesTitles(&moviesTitles)
	clientRecRequest := createClientRecRequest(clientToSend.UserId, clientToSend.Quantity, clientToSend.GenreIds)
	clientRecRequest.Ratings = MappRatingsClient(clientToSend.MoviesRatings, &moviesTitles)

	sendRequestToRecommendationService(clientRecRequest)

	w.Header().Set("Content-Type", "application/json")
}

// obtner peliculas de un género
func getMoviesByGenre(genre int) []MovieTitleWithID {
	moviesTitles := MoviesTitles{}
	LoadMoviesTitles(&moviesTitles)
	moviesGenreIds := MoviesGenreIds{}
	LoadMoviesGenreIds(&moviesGenreIds)
	genres := Genres{}
	LoadGenres(&genres)
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

func getMoviesByGenresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	//getID of genre from request
	genre := r.URL.Query().Get("id")
	id, err := strconv.Atoi(genre)
	if err != nil {
		http.Error(w, "Error al convertir el id del género a entero", http.StatusBadRequest)
		return
	}
	movies := getMoviesByGenre(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func setupRoutes() {
	http.HandleFunc("/movies/titles", moviesTitlesHandler)
	http.HandleFunc("/genres", genresHandler)
	http.HandleFunc("/movies/genres", MoviesGenresHandler)
	http.HandleFunc("/recommendations", sendRequestToRecommendationHandler)
	http.HandleFunc("/genres/movies", getMoviesByGenresHandler)
	log.Println("Server running on port 9000")
	log.Fatal(http.ListenAndServe(":9000", nil))

}

func main() {
	Banner()
	setupRoutes()
}
