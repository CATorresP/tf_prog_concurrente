export interface GenresButtonsProps {
  genres: string[];
  HandleButtonClick: (idGenre: number) => void;
  clickedButtons: number[];
  onStart: () => void;
}

export function GenresButtons({
  genres,
  HandleButtonClick,
  clickedButtons,
  onStart,
}: GenresButtonsProps) {
  return (
    <div
      className="
          flex flex-col items-center justify-center h-screen bg-s
      "
    >
      <h1
        className="
          text-4xl text-center text-white font-bold p-4"
      >
        Seleccione un género o generos de película para buscar recomendaciones
      </h1>
      <div className="flex flex-wrap justify-center w-4/5">
        {genres.map((genre, index) => (
          <button
            key={genre}
            className={`${
              clickedButtons.includes(index) ? "bg-green-700" : "bg-blue-700"
            } text-white font-bold py-2 px-4 rounded-sm m-2 transition-colors duration-300`}
            onClick={() => HandleButtonClick(index)}
          >
            {genre}
          </button>
        ))}
      </div>
      <button
        className="
            bg-red-700 text-white font-bold py-2 px-4 rounded-sm m-2 transition-colors duration-300"
        onClick={onStart}
      >
        Iniciar
      </button>
    </div>
  );
}
