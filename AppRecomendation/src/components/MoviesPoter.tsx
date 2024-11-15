import { useState } from "react";
import { MovieTitleId, MovieRating } from "@/models";
import { StarRating } from "@/components";

export interface MoviePosterProps {
  movieTitleId: MovieTitleId;
  rating: number;
  setRating: (rating: number) => void;
}

export function MoviePoster({
  movieTitleId,
  rating,
  setRating,
}: MoviePosterProps) {
  const [imageError, setImageError] = useState(false);
  const handleRatingChange = (newRating: number) => {
    setRating(newRating);
  };

  return (
    <>
      {!imageError && (
        <div className="flex flex-col items-center justify-center bg-gray-800 p-2 rounded-lg">
          <img
            className="w-40 h-60 rounded-lg"
            src={`./src/assets/${movieTitleId.id}.jpg`}
            alt={movieTitleId.title}
            onError={() => setImageError(true)}
          />
          <div>
            <h1 className="text-white text-center">{movieTitleId.title}</h1>
          </div>
          <div>
            <h1 className="text-white text-center">Puntuación:</h1>
            <StarRating rating={rating} onRatingChange={handleRatingChange} />
          </div>
        </div>
      )}
    </>
  );
}

export interface MoviesPostersProps {
  movies: MovieTitleId[];
  handleMovieRating: (rating: number, movieId: number) => void;
  movieRatings: MovieRating[];
  hadleRecommendations: () => void;
  handleNumRecommendations: (numRecommendations: number) => void;
}

export function MoviesPosters({
  movies,
  handleMovieRating,
  movieRatings,
  hadleRecommendations,
  handleNumRecommendations,
}: MoviesPostersProps) {
  return (
    <>
      <div className="text-4xl text-center text-white font-bold p-4">
        Puntua las películas que has visto
      </div>
      <div className="text-center text-white font-bold p-4">
        <h1>¿Cuántas recomendaciones quieres?</h1>
        <input
          type="number"
          onChange={(e) => handleNumRecommendations(Number(e.target.value))}
          className="bg-gray-800 text-white font-bold p-2 rounded"
        ></input>
      </div>

      <div className="flex justify-center">
        {movieRatings.length > 0 && (
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            onClick={hadleRecommendations}
          >
            Obtener recomendaciones
          </button>
        )}
      </div>
      <div className="flex flex-wrap justify-center gap-2">
        {movies.map((movie) => {
          const movieRating = movieRatings.find((r) => r.movieId === movie.id);
          return (
            <MoviePoster
              key={movie.id}
              movieTitleId={movie}
              rating={movieRating ? movieRating.rating : 0}
              setRating={(rating) => handleMovieRating(rating, movie.id)}
            />
          );
        })}
      </div>
    </>
  );
}
