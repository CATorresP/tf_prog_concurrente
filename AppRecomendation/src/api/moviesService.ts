import {
  MoviesTitles,
  Genres,
  MovieGenres,
  MovieTitleId,
  MovieRatings,
  Recommendations,
} from "@/models";
export class MoviesService {
  private url: string;

  constructor() {
    this.url = "http://localhost:9000";
  }

  async getMovies(): Promise<MoviesTitles> {
    const response = await fetch(`${this.url}/movies/titles`);
    const movies = await response.json();
    return movies;
  }

  async getGenres(): Promise<Genres> {
    const response = await fetch(`${this.url}/genres`);
    const genres = await response.json();
    return genres;
  }

  async getMoviesByGenre(): Promise<MovieGenres[]> {
    const response = await fetch(`${this.url}/movies/genres`);
    const moviesGenres = await response.json();
    return moviesGenres;
  }

  async getGenresMovies(genreId: number): Promise<MovieTitleId[]> {
    const response = await fetch(`${this.url}/genres/movies?id=${genreId}`);
    const movies = await response.json();
    return movies;
  }

  async getRecommendations(
    movieRatings: MovieRatings
  ): Promise<Recommendations> {
    const response = await fetch(`${this.url}/recommendations`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(movieRatings),
    });

    if (!response.ok) {
      throw new Error(
        `Error en la respuesta del servicio de recomendaciones: ${response.statusText}`
      );
    }

    const recommendations: Recommendations = await response.json();
    return recommendations;
  }
}
