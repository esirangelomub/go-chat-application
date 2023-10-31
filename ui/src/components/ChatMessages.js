import React from "react";

const ChatMessages = ({ messages }) => {
    return (
        <div className="chat-messages">
            {messages.map((msg, index) => (
                <div className="message-item" key={index}>
                    <strong>{msg.username}</strong>
                    <span>{new Date(msg.timestamp).toLocaleTimeString()}</span>
                    <p>{msg.content}</p>
                </div>
            ))}
        </div>
    );
}

export default ChatMessages;