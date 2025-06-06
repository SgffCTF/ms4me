import { useEffect, useRef } from "react";
import { Modal } from 'bootstrap';

interface Props {
    show: boolean;
    setShow: (show: boolean) => void;
}

export const CreateGameModal = (props: Props) => {
    const modalRef = useRef<HTMLDivElement | null>(null);
    const modalInstanceRef = useRef<Modal | null>(null);

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

    return (
        <div ref={modalRef} className="modal fade" tabIndex={-1} aria-hidden="true">
            <div className="modal-dialog">
                <div className="modal-content">
                    <div className="modal-header">
                        <h5 className="modal-title">Заголовок</h5>
                        <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Закрыть"></button>
                    </div>
                    <div className="modal-body">Содержимое окна</div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
                    </div>
                </div>
            </div>
        </div>
    );
}
