import { createContext, useContext, useState, useEffect } from "react";
import { MoviesService } from "@/api/moviesService";
import {
  MoviesTitles,
  Genres,
  MovieGenres,
  MovieTitleId,
  Recommendations,
  MovieRatings,
} from "@/models";

const moviesService = new MoviesService();

interface MoviesContextProps {
  movies: MoviesTitles | null;
  genres: Genres | null;
  moviesByGenre: MovieGenres[] | null;
  genresMovies: MovieTitleId[] | null;
  recommendations: Recommendations | null;
  fetchMoviesByGenre: () => void;
  fetchGenresMovies: (genreId: number[]) => void;
  fetchRecommendations: (movieRatings: MovieRatings) => void;
}

interface MoviesContextProviderProps {
  children: React.ReactNode | React.ReactNode[];
}

export const MoviesContext = createContext<MoviesContextProps>(
  {} as MoviesContextProps
);

export const useMovies = () => {
  const context = useContext(MoviesContext);
  if (!context) {
    throw new Error("useMovies must be used within a MoviesProvider");
  }
  return context;
};

export const MoviesProvider = ({ children }: MoviesContextProviderProps) => {
  const [movies, setMovies] = useState<MoviesTitles | null>(null);
  const [genres, setGenres] = useState<Genres | null>(null);
  const [moviesByGenre, setMoviesByGenre] = useState<MovieGenres[] | null>(
    null
  );
  const [genresMovies, setGenresMovies] = useState<MovieTitleId[] | null>(null);
  const [recommendations, setRecommendations] =
    useState<Recommendations | null>(null);

  const fetchMovies = async () => {
    const movies = await moviesService.getMovies();
    setMovies(movies);
  };

  const fetchGenres = async () => {
    const genres = await moviesService.getGenres();
    setGenres(genres);
  };

  const fetchMoviesByGenre = async () => {
    const moviesGenres = await moviesService.getMoviesByGenre();
    setMoviesByGenre(moviesGenres);
  };

  const fetchGenresMovies = async (genreIds: number[]) => {
    let genresMovies: MovieTitleId[] = [];
    for (let i = 0; i < genreIds.length; i++) {
      const movies = await moviesService.getGenresMovies(genreIds[i]);
      genresMovies = genresMovies.concat(movies);
    }
    genresMovies = genresMovies.filter(
      (v, i, a) => a.findIndex((t) => t.id === v.id) === i
    );
    setGenresMovies(genresMovies);
  };

  const fetchRecommendations = async (movieRatings: MovieRatings) => {
    const recommendations = await moviesService.getRecommendations(
      movieRatings
    );
    console.log(recommendations);
    setRecommendations(recommendations);
  };

  useEffect(() => {
    fetchMovies();
    fetchGenres();
  }, []);

  return (
    <MoviesContext.Provider
      value={{
        movies,
        genres,
        moviesByGenre,
        genresMovies,
        recommendations,
        fetchMoviesByGenre,
        fetchGenresMovies,
        fetchRecommendations,
      }}
    >
      {children}
    </MoviesContext.Provider>
  );
};
