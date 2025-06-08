import { useEffect, useRef, useState } from "react";
import { Modal } from 'bootstrap';
import { updateGame } from "../api/games";
import { toast } from "react-toastify";
import { GameDetails } from "../models/models";

interface Props {
    id: string;
    show: boolean;
    setShow: (show: boolean) => void;
    gameInfo: GameDetails;
}

export const UpdateGameModal = (props: Props) => {
    const modalRef = useRef<HTMLDivElement | null>(null);
    const modalInstanceRef = useRef<Modal | null>(null);
    const nameInput = useRef<HTMLInputElement>(null);
    const [isPublic, setIsPublic] = useState(props.gameInfo.is_public);

    useEffect(() => {
        if (modalRef.current && !modalInstanceRef.current) {
            modalInstanceRef.current = new Modal(modalRef.current);
            modalRef.current.addEventListener('hidden.bs.modal', () => {
                props.setShow(false);
            });
        }
    }, []);

    useEffect(() => {
        if (modalInstanceRef.current) {
            if (props.show) {
                modalInstanceRef.current.show();
            } else {
                modalInstanceRef.current.hide();
            }
        }
    }, [props.show]);

    const handleIsPublic = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const checked = event.target.checked;
        setIsPublic(checked);
    };

    const handleUpdate = async (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        event.preventDefault();
        if (nameInput.current) {
            try {
                await updateGame(props.id, nameInput.current.value, isPublic);
                toast("Игра обновлена");
            } catch (err: any) {
                toast.error(err.message);
            }
        }
        if (modalInstanceRef.current) {
            modalInstanceRef.current.hide();
        }
    }

    return (
        <div ref={modalRef} className="modal fade" tabIndex={-1} aria-hidden="true">
            <div className="modal-dialog">
                <div className="modal-content">
                    <div className="modal-header">
                        <h5 className="modal-title">Создание игры</h5>
                        <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Закрыть"></button>
                    </div>
                    <div className="modal-body">
                        <div className="form-floating mb-3">
                            <input type="name" className="form-control" id="create-game-name" placeholder="Название" ref={nameInput} defaultValue={props.gameInfo.title}/>
                            <label htmlFor="create-game-name">Название</label>
                        </div>
                        <div className="form-check">
                            <input className="form-check-input" type="checkbox" id="create-game-is-public" checked={isPublic} onChange={handleIsPublic}/>
                            <label className="form-check-label" htmlFor="create-game-is-public">
                                Публичная
                            </label>
                        </div>
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-primary" data-bs-dismiss="modal" onClick={handleUpdate}>Создать</button>
                        <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
                    </div>
                </div>
            </div>
        </div>
    );
}
