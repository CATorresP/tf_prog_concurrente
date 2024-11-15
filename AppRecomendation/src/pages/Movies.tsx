import { useState, useEffect } from "react";
import { useMovies } from "@/context";
import { useHandleButtonsClick } from "@/hooks";
import {
  MovieRating,
  MovieRatings,
  Recommendations,
  Recommendation,
} from "@/models";
import {
  DynamicBackground,
  MoviesInit,
  GenresButtons,
  MoviesPosters,
} from "@/components";

enum Stage {
  INIT,
  GENRES,
  MOVIES,
  RECOMMENDATIONS,
}

export function MoviesPage() {
  const [stage, setStage] = useState<Stage>(Stage.INIT);
  const {
    genres,
    fetchGenresMovies,
    genresMovies,
    fetchRecommendations,
    recommendations,
  } = useMovies();
  const { clickedButtons, handleButtonClick } = useHandleButtonsClick();
  const [movieRating, setMovieRating] = useState<MovieRating[]>([]);
  const [quantity, setQuantity] = useState<number>(0);

  const handleGenres = () => {
    setStage(Stage.GENRES);
  };

  const handleMovies = () => {
    setStage(Stage.MOVIES);
  };

  const hadleRecommendations = () => {
    setStage(Stage.RECOMMENDATIONS);
  };

  const handleMovieRating = (rating: number, movieId: number) => {
    setMovieRating((prevRatings) => {
      const existingRating = prevRatings.find((r) => r.movieId === movieId);
      if (existingRating) {
        return prevRatings.map((r) =>
          r.movieId === movieId ? { ...r, rating } : r
        );
      } else {
        return [...prevRatings, { movieId, rating }];
      }
    });
  };

  const handleFetchRecommendations = () => {
    const movieRatingsData: MovieRatings = {
      moviesRatings: movieRating,
      userId: 0,
      quantity: quantity,
      genreIds: clickedButtons,
    };
    fetchRecommendations(movieRatingsData);
  };

  const handleNumRecommendations = (numRecommendations: number) => {
    setQuantity(numRecommendations);
  };

  const handleReset = () => {
    setStage(Stage.INIT);
  };

  useEffect(() => {
    if (clickedButtons.length > 0) {
      fetchGenresMovies(clickedButtons);
    }
  }, [clickedButtons]);

  useEffect(() => {
    if (
      movieRating.length > 0 &&
      quantity > 0 &&
      stage === Stage.RECOMMENDATIONS
    ) {
      handleFetchRecommendations();
    }
  }, [movieRating, quantity, stage]);

  const renderComponent = () => {
    switch (stage) {
      case Stage.INIT:
        return <MoviesInit onStart={handleGenres} />;
      case Stage.GENRES:
        return (
          genres && (
            <GenresButtons
              genres={genres?.movieGenreNames}
              clickedButtons={clickedButtons}
              HandleButtonClick={handleButtonClick}
              onStart={handleMovies}
            />
          )
        );
      case Stage.MOVIES:
        return (
          <MoviesPosters
            movies={genresMovies || []}
            handleMovieRating={handleMovieRating}
            movieRatings={movieRating}
            hadleRecommendations={hadleRecommendations}
            handleNumRecommendations={handleNumRecommendations}
          />
        );
      case Stage.RECOMMENDATIONS:
        return (
          recommendations && (
            <RecommendationComponent
              recommendations={recommendations}
              onReset={handleReset}
            />
          )
        );
      default:
        return null;
    }
  };

  return (
    genres && (
      <div>
        <DynamicBackground />
        {renderComponent()}
      </div>
    )
  );
}

export interface RecommendationComponentProps {
  recommendations: Recommendations;
  onReset: () => void;
}

function RecommendationComponent({
  recommendations,
  onReset,
}: RecommendationComponentProps) {
  return (
    <div className="flex flex-col items-center justify-center h-full">
      <h1 className="text-3xl font-bold mb-4">Recomendaciones</h1>
      {recommendations.recommendations.map((rec: Recommendation) => (
        <div
          key={rec.id}
          className="bg-slate-500 p-4 rounded-lg shadow-md mb-4 w-1/2"
        >
          <h2 className="text-2xl font-semibold">{rec.title}</h2>
          <p className="text-gray-300">Géneros: {rec.genres.join(", ")}</p>
          <p className="text-gray-400">Rating: {rec.rating}</p>
          <p className="text-gray-400">Comentario: {rec.comment}</p>
        </div>
      ))}
      <button
        onClick={onReset}
        className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Volver al inicio
      </button>
    </div>
  );
}
