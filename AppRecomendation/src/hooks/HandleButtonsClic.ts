import { useState } from "react";
export function useHandleButtonsClick() {
  const [clickedButtons, setClickedButtons] = useState<number[]>([]);

  const handleButtonClick = (idGenre: number) => {
    if (clickedButtons.includes(idGenre)) {
      setClickedButtons((prev) => prev.filter((id) => id !== idGenre));
    } else {
      setClickedButtons((prev) => [...prev, idGenre]);
    }
  };

  return { clickedButtons, handleButtonClick };
}
