import '../styles/Minefield.css';
import { useState } from 'react';

type CellType = "closed" | "empty" | "flag" | "mine" | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;

const getCellClass = (type: CellType) => {
  switch(type) {
    case "closed":
      return "btn btn-secondary";
    case "empty":
      return "btn btn-light";
    case "flag":
      return "btn btn-danger";
    case "mine":
      return "btn btn-dark";
    case 1:
      return "btn btn-light text-primary";
    case 2:
      return "btn btn-light text-success";
    case 3:
      return "btn btn-light text-danger";
    case 4:
      return "btn btn-light text-info";
    case 5:
      return "btn btn-light text-warning";
    case 6:
      return "btn btn-light text-secondary";
    case 7:
      return "btn btn-light text-dark";
    case 8:
      return "btn btn-light text-muted";
    default:
      return "btn btn-light";
  }
};

export const Field = () => {
  const [cells] = useState<CellType[]>(() =>
    Array(64).fill(null).map(() => {
      const options: CellType[] = ["closed", "empty", "flag", "mine", 1, 2, 3, 4, 5, 6, 7, 8];
      //   return options[Math.floor(Math.random() * options.length)];
      return options[0];
    })
  );

  return (
    <div className="container-fluid">
      <div className="minefield d-grid" style={{gridTemplateColumns: "repeat(8, 1fr)", gap: "2px"}}>
        {cells.map((type, idx) => (
          <button key={idx} className={`${getCellClass(type)} cell`}>
            {/* Показываем иконки для флага и мины */}
            {type === "flag" && "🚩"}
            {type === "mine" && "💣"}
            {/* Цифры показываем сами как текст */}
            {(typeof type === "number") && type}
            {/* Пустые и закрытые — пустые кнопки */}
          </button>
        ))}
      </div>
    </div>
  );
};
