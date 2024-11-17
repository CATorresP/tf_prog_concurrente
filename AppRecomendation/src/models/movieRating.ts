export interface MovieRating {
  movieId: number;
  rating: number;
}

export interface MovieId {
  id: number;
}

export interface MovieRatings {
  moviesRatings: MovieRating[];
  userId: number;
  quantity: number;
  genreIds: number[];
}

export interface Recommendation {
  id: number;
  title: string;
  genres: string[];
  rating: number;
  comment: string;
}

export interface Recommendations {
  userId: number;
  recommendations: Recommendation[];
}
