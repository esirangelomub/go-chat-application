import React from 'react';
import { Route, Routes } from 'react-router-dom';
import ChatPage from "./ChatPage";

const ChatArea = () => {
    return (
        <Routes>
            <Route path="/c/:chatName/:chatId" element={<ChatPage />} />
            {/* Add more routes as needed */}
        </Routes>
    );
}

export default ChatArea;
