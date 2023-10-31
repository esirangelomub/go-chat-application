// ChatPrompt.js
import React, { useState } from 'react';

const ChatPrompt = ({ ws }) => {
    const [message, setMessage] = useState('');

    const handleSendMessage = () => {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({
                type: 'message',
                content: message
            }));
            setMessage(''); // Clear the input field after sending
        } else {
            console.error('WebSocket is not open.');
        }
    };

    return (
        <div className="chat-prompt">
            <div className="input-group mb-3 w-100">
                <input type="text"
                       className="form-control"
                       placeholder="Type a Message"
                       value={message}
                       onChange={(e) => setMessage(e.target.value)}
                       onKeyPress={(e) => {
                           if (e.key === 'Enter' && !e.shiftKey) {
                               e.preventDefault();  // Prevents adding a newline
                               handleSendMessage();
                           }
                       }}/>
                    <button className="btn btn-outline-secondary"
                            type="button"
                            id="button-send"
                            onClick={handleSendMessage}>Send</button>
            </div>
        </div>
    );
}

export default ChatPrompt;
