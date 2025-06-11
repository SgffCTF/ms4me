import { CellType } from "../components/Field/Cell";

interface Cell {
    value: CellType;
    is_open: boolean;
}

export interface Field {
    rows: number;
    cols: number;
    mines: number;
    cells_open: number;
    mine_is_open: number;
    grid: Array<Array<Cell>>;
}