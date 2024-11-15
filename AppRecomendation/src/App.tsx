import { MoviesProvider } from "@/context";
import { MoviesPage } from "@/pages";

export function App() {
  return (
    <MoviesProvider>
      <MoviesPage />
    </MoviesProvider>
  );
}
