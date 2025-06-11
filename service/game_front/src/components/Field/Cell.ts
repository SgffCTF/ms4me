export type CellType = "c" | "0" | "f" | "m" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8";

export const getCellClass = (type: CellType) => {
  switch(type) {
    case "c":
      return "btn btn-secondary";
    case "0":
      return "btn btn-light";
    case "f":
      return "btn btn-danger";
    case "m":
      return "btn btn-dark";
    case "1":
      return "btn btn-light text-primary";
    case "2":
      return "btn btn-light text-success";
    case "3":
      return "btn btn-light text-danger";
    case "4":
      return "btn btn-light text-info";
    case "5":
      return "btn btn-light text-warning";
    case "6":
      return "btn btn-light text-secondary";
    case "7":
      return "btn btn-light text-dark";
    case "8":
      return "btn btn-light text-muted";
    default:
      return "btn btn-light";
  }
};