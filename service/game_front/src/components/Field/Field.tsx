import '../../styles/Minefield.css';
import { CellType, getCellClass } from './Cell';
import { ClickGameEvent } from '../../models/events';
import { openCell } from '../../api/game';
import { toast } from 'react-toastify';


interface Props {
  gameID: string;
  clickGameEvent: ClickGameEvent | null;
}

export const Field = (props: Props) => {
  const closedCells = Array(64).fill('c');

  const handleClick = async (row: number, col: number) => {
    try {
      await openCell(props.gameID, row, col);
    } catch (e: any) {
      toast.error(e.message);
    }
  }

  console.log(props.clickGameEvent?.field.grid);

  return (
    <div className="container-fluid">
      <div className="minefield d-grid" style={{gridTemplateColumns: `repeat(${props.clickGameEvent?.field.cols}, 1fr)`, gap: "2px"}}>
        {props.clickGameEvent === null && closedCells.map((type, idx) => (
          <button key={idx} className={`${getCellClass(type)} cell`} onClick={() => handleClick(Math.floor((idx + 1) / 8), idx % 8)}>

          </button>
        ))}
        {props.clickGameEvent && props.clickGameEvent.field.grid.map((row, idx1) => (
          row.map((cell, idx2) => (
            <button key={`${idx1}_${idx2}`} className={`${getCellClass(cell.value as CellType)} cell`} onClick={() => handleClick(idx1, idx2)}>
            {(cell.value === "f" && "🚩") ||
            (cell.value === "m" && "💣") ||
            (cell.value != "c" && cell.value)}
            </button>
          ))
        ))
        }
      </div>
    </div>
  );
};
