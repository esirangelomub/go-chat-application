import React, { useEffect, useState } from "react";
import { useParams } from 'react-router-dom';
import ChatMessages from "./ChatMessages";
import ChatPrompt from "./ChatPrompt";

const ChatPage = () => {
    const { chatName, chatId } = useParams();
    const [messages, setMessages] = useState([]);
    const [ws, setWs] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        const wsInstance = new WebSocket(`${process.env.REACT_APP_WS_URL}/ws/${chatId}?jwt=${token}`);

        wsInstance.onopen = () => {
            console.log("Connected to the chat");
        };

        wsInstance.onmessage = (event) => {
            const messageData = JSON.parse(event.data);
            setMessages(prevMessages => [...prevMessages, messageData]);
        };

        wsInstance.onclose = () => {
            setMessages([]);
            console.log("Disconnected from the chat");
        };

        setWs(wsInstance);

        return () => {
            if (wsInstance) {
                wsInstance.close();
            }
        };
    }, [chatId]);

    return (
        <div className="container">
            <div className="row justify-content-md-center">
                <div className="col-md-8">
                    <h1>{chatName}</h1>
                    <ChatMessages messages={messages}/>
                    <div className="message-count">Messages: {messages.length}</div>
                    <ChatPrompt ws={ws}/>
                </div>
            </div>
        </div>
    );
}

export default ChatPage;
