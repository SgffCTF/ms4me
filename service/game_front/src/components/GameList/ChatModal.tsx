import { useEffect, useRef } from "react";
import { Modal } from 'bootstrap';
import { Chat } from "../Chat";
import { Message } from "../../models/models";

interface Props {
    show: boolean;
    setShow: (show: boolean) => void;
    id: string | null;
    messages: Message[] | null;
}

export const ChatModal = (props: Props) => {
    const modalRef = useRef<HTMLDivElement | null>(null);
    const modalInstanceRef = useRef<Modal | null>(null);
    const chatRef = useRef<{ scrollToBottom: () => void } | null>(null);

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
                setTimeout(() => {
                    chatRef.current?.scrollToBottom();
                }, 200)
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
                        <h5 className="modal-title">Чат во время игры</h5>
                        <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Закрыть"></button>
                    </div>
                    <div className="modal-body">
                        {props.messages && props.id && props.messages.length != 0 &&
                        <Chat messages={props.messages} id={props.id} withInput={false} ref={chatRef}></Chat> || 
                        <p>Сообщения отсутствуют или их срок хранения истёк</p>}
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
                    </div>
                </div>
            </div>
        </div>
    );
}
