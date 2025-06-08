import { useEffect, useRef, useState } from "react";
import { Field } from "../components/Field";
import { Chat } from "../components/Chat";

export const Game = () => {
  const [message, setMessage] = useState("");
  const [logs, setLogs] = useState<string[]>([]);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const wsURL = import.meta.env.VITE_WS_URI;
    const socket = new WebSocket(wsURL);
    socketRef.current = socket;

    socket.onopen = () => { console.log("connected") };
    socket.onclose = () => { console.log("closing conn") };
    socket.onmessage = (ev) => {
      setLogs((prevLogs) => [...prevLogs, `Received: ${ev.data}`]);
    };

    return () => {
      socket.close();
    };
  }, []);

  const sendMessage = () => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.send(message);
      setLogs((prevLogs) => [...prevLogs, `Sent: ${message}`]);
    }
  };

  return (
    <>
      <h1 className="text-center mt-5 mb-5">Игра</h1>
      <div className="container-fluid">
        <div className="row d-flex h-all">
          <div className="col-4">
            <Field></Field>
          </div>
          <div className="col-4"> 
            <Field></Field>
          </div>
          <div className="col-4">
            <Chat></Chat>
          </div>
        </div>
      </div>
    </>
  )
}