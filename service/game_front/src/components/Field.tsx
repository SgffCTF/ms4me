import '../styles/Minefield.css';


export const Field = () => {
    return (
        <div className="container-fluid">
            <div className="minefield d-grid">
                {[...Array(64)].map((_, idx) => (
                <button key={idx} className="cell btn btn-white">
                    &nbsp;
                </button>
                ))}
            </div>
        </div>
    )
}