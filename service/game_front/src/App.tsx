import { useState, useRef, useEffect } from 'react'
import './App.css'

function App() {
  const [message, setMessage] = useState("");
  const [logs, setLogs] = useState<string[]>([]);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const wsURL = import.meta.env.VITE_WS_URI;
    const socket = new WebSocket(wsURL);
    socketRef.current = socket;

    socket.onopen = (ev) => { console.log("connected") };
    socket.onclose = (ev) => { console.log("closing conn") };
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
      <div style={{ padding: "20px" }}>
        <h1>WebSocket Chat</h1>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Enter message"
        />
        <button onClick={sendMessage}>Send</button>

        <div>
          <h2>Logs:</h2>
          <ul>
            {logs.map((log, index) => (
              <li key={index}>{log}</li>
            ))}
          </ul>
        </div>
      </div>
    </>
  )
}

export default App
