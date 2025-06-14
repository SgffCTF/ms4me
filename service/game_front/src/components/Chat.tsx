import { useState, useRef, useEffect, useImperativeHandle, forwardRef } from "react";
import "../styles/Chat.css";
import { Message } from "../models/models";
import { useAuth } from "../context/AuthProvider";
import { sendMessage } from "../api/ingame";
import { toast } from "react-toastify";

interface Props {
  id: string;
  messages: Message[];
  withInput: boolean;
}

export const Chat = forwardRef((props: Props, ref) => {
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { user } = useAuth();

  useImperativeHandle(ref, () => ({
    scrollToBottom: () => {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }));

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [props.messages])

  const handleKeyDown = async (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (!input.trim()) return;
      
      try {
        await sendMessage(props.id, input);
        setInput("");
      } catch (e: any) {
        toast.error(e.message);
      }
    }
  };

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
        <div className="messages flex-grow-1 d-flex flex-column overflow-auto">
          <div style={{ flexGrow: 1 }} />
          {props.messages.map((msg, i) => (
            <div
              key={i}
              className={`d-flex flex-column mb-2 ${
                msg.creator_id === user?.id ? "align-items-end" : "align-items-start"
              }`}
            >
              <small className="text-muted mb-1">{msg.creator_username}</small>
              <div
                className={`p-2 rounded shadow-sm ${
                  msg.creator_id === user?.id ? "bg-green text-white" : "bg-light"
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

      {props.withInput &&
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
    }
    </div>
  );
});
