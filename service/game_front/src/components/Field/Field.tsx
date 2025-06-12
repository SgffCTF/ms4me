import '../../styles/Minefield.css';
import { CellType, getCellClass } from './Cell';
import { ClickGameEvent, RoomParticipant } from '../../models/events';
import { openCell, setFlag } from '../../api/ingame';
import { toast } from 'react-toastify';
import { useAuth } from '../../context/AuthProvider';


interface Props {
  gameID: string;
  clickGameEvent: ClickGameEvent | null;
  fieldOwnerID: number | null;
}

export const Field = (props: Props) => {
  const closedCells = Array(64).fill('c');
  const { user } = useAuth();
  let participant: RoomParticipant | undefined = props.clickGameEvent?.participants.find((val) => {
    if (val.id == props.fieldOwnerID) {
      return val;
    }
  });

  const handleClick = async (e: React.MouseEvent, row: number, col: number) => {
    e.preventDefault();
    if (props.fieldOwnerID !== user?.id) {
      toast.error("Это не твоё поле");
      return;
    }

    try {
      if (e.button === 0) {
        await openCell(props.gameID, row, col);
      } else if (e.button === 2) {
        await setFlag(props.gameID, row, col);
      }
    } catch (err: any) {
      toast.error(err.message);
    }
  };

  return (
    <div className="container-fluid">
      <div className="minefield d-grid" style={{gridTemplateColumns: `repeat(${participant?.field.cols}, 1fr)`, gap: "2px"}}>
        {!participant && closedCells.map((type, idx) => (
          <button
          key={idx}
          className={`${getCellClass(type)} cell`}
          onClick={(e) => handleClick(e, Math.floor((idx + 1) / 8), idx % 8)}
          >
          </button>
        ))}
        {participant && participant.field.grid.map((row, idx1) => (
          row.map((cell, idx2) => (
            <button
            key={`${idx1}_${idx2}`}
            className={`${getCellClass(cell.value as CellType)} cell`}
            onClick={(e) => handleClick(e, idx1, idx2)}
            onContextMenu={(e) => {
              e.preventDefault();
              handleClick(e, idx1, idx2);
            }}>
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
