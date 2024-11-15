export interface MoviesInitProps {
  onStart: () => void;
}

export function MoviesInit({ onStart }: MoviesInitProps) {
  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <div className="bg-slate-700 w-1/2 rounded-sm">
        <h1 className="text-4xl text-center text-white font-bold p-4">
          Bienvenido al Sistema de Recomendación de Películas
        </h1>
      </div>

      <button
        className="bg-blue-700 text-white font-bold py-2 px-4 rounded-sm mt-4"
        onClick={onStart}
      >
        Clic para Comenzar
      </button>
    </div>
  );
}
