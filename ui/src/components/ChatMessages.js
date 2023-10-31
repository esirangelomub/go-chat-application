import React from "react";

const ChatMessages = ({ messages }) => {
    return (
        <div className="chat-messages">
            {messages.map((msg, index) => (
                <div className={`message-item ${msg.username === 'Bot' && 'bot-message'}`} key={index}>
                    <div className="w-100 d-flex justify-content-between align-content-center">
                        <strong>{msg.username}</strong>
                        <small>{new Date(msg.timestamp).toLocaleTimeString()}</small>
                    </div>
                    <p>{msg.content}</p>
                </div>
            ))}
        </div>
    );
}

export default ChatMessages;