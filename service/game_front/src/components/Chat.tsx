import { useState, useRef, useEffect } from "react";
import "../styles/Chat.css";

type Message = {
  text: string;
  author: "me" | "other";
  authorName: string;
};

export const Chat = () => {
  const [messages, setMessages] = useState<Message[]>([
    { text: "Добро пожаловать!", author: "other", authorName: "Сервер" },
    { text: "Напишите сообщение...", author: "other", authorName: "Сервер" }
  ]);
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const sendMessage = () => {
    if (input.trim()) {
      setMessages((prev) => [
        ...prev,
        { text: input, author: "me", authorName: "Вы" }
      ]);
      setInput("");
      
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          { text: "Автоматический ответ", author: "other", authorName: "Бот" }
        ]);
      }, 1000);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="chat-card card h-chat w-chat d-flex flex-column" style={{ borderRadius: 15 }}>
      <div
        className="card-header d-flex justify-content-between align-items-center p-3 bg-green text-white border-bottom-0"
        style={{ borderTopLeftRadius: 15, borderTopRightRadius: 15 }}
      >
        <i className="fas fa-angle-left" />
        <p className="mb-0 fw-bold">Чат</p>
        <i className="fas fa-times" />
      </div>

      <div className="card-body d-flex flex-column flex-grow-1 overflow-hidden px-3 py-2">
        <div className="messages flex-grow-1 overflow-auto d-flex flex-column justify-content-end">
          {messages.map((msg, i) => (
            <div
              key={i}
              className={`d-flex flex-column mb-2 ${
                msg.author === "me" ? "align-items-end" : "align-items-start"
              }`}
            >
              <small className="text-muted mb-1">{msg.authorName}</small>
              <div
                className={`p-2 rounded shadow-sm ${
                  msg.author === "me" ? "bg-green text-white" : "bg-light"
                }`}
                style={{ maxWidth: "75%" }}
              >
                {msg.text}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
      </div>

      <div className="card-footer bg-white border-top p-3">
        <textarea
          className="form-control"
          placeholder="Напиши сообщение..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          rows={2}
        />
      </div>
    </div>
  );
};
