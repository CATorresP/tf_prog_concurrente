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

  constructor(url: string = "ws://localhost:9000") {
    this.url = url;
  }

  private sendRequest<T>(
    path: string,
    data?: Record<string, unknown>
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(this.url + path);

      ws.onopen = () => {
        console.log("Sending data to server", JSON.stringify({ data }));
        ws.send(JSON.stringify({ ...data }));
      };

      ws.onmessage = (event) => {
        const response = JSON.parse(event.data.toString());
        console.log("Received data from server", response);
        if (response.error) {
          reject(response.error);
        } else {
          resolve(response);
        }

        setTimeout(() => {
          ws.close();
        }, 600000);
      };

      ws.onerror = (error) => {
        console.error("WebSocket error", error);
        reject(error);
        ws.close();
      };

      ws.onclose = () => {
        console.log("Disconnected from server");
      };
    });
  }

  async getMovies(): Promise<MoviesTitles> {
    return this.sendRequest<MoviesTitles>("/movies/titles");
  }

  async getGenres(): Promise<Genres> {
    return this.sendRequest<Genres>("/genres");
  }

  async getMoviesByGenre(): Promise<MovieGenres[]> {
    return this.sendRequest<MovieGenres[]>("/movies/genres");
  }

  async getGenresMovies(genreId: number): Promise<MovieTitleId[]> {
    return this.sendRequest<MovieTitleId[]>("/genres/movies", { id: genreId });
  }

  async getRecommendations(
    movieRatings: MovieRatings
  ): Promise<Recommendations> {
    const { moviesRatings, userId, quantity, genreIds } = movieRatings;
    return this.sendRequest<Recommendations>("/recommendations", {
      moviesRatings,
      userId,
      quantity,
      genreIds,
    });
  }
}
