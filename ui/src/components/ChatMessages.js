import React, {useContext, useEffect, useRef} from "react";
import UserContext from "../contexts/UserContext";

const ChatMessages = ({ messages }) => {
    const { userData } = useContext(UserContext);
    const messagesEndRef = useRef(null);

    const getClassName = (msg) => {
        if (msg.username === 'Bot') {
            return 'bot-message';
        } else if (msg.user_id === userData.user_id) {
            return 'user-message-logged';
        } else {
            return 'user-message';
        }
    };

    const getUserIdColor = (userId) => {
        // Generate a color based on user ID
        const hash = userId.toString().split('').reduce((acc, char) => {
            return char.charCodeAt(0) + ((acc << 5) - acc);
        }, 0);

        const color = (hash & 0x00FFFFFF)
            .toString(16)
            .toUpperCase();

        return '#' + '00000'.substring(0, 6 - color.length) + color;
    };

    function formatDate(isoDateString) {
        const inputDate = new Date(isoDateString);
        const currentDate = new Date();

        const isSameDay = (date1, date2) =>
            date1.getFullYear() === date2.getFullYear() &&
            date1.getMonth() === date2.getMonth() &&
            date1.getDate() === date2.getDate();

        const getTimeString = (date) =>
            date.toISOString().substr(11, 8);  // Extracts time as HH:MM:SS from ISO date string

        let formattedDate;

        if (isSameDay(inputDate, currentDate)) {
            // Today
            formattedDate = `Today - ${getTimeString(inputDate)}`;
        } else if (isSameDay(inputDate, new Date(currentDate.setDate(currentDate.getDate() - 1)))) {
            // Yesterday
            formattedDate = `Yesterday - ${getTimeString(inputDate)}`;
        } else {
            // Older date
            const dateStr = inputDate.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: '2-digit' });
            formattedDate = `${dateStr} - ${getTimeString(inputDate)}`;
        }

        return formattedDate;
    }

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(scrollToBottom, [messages]);

    return (
        <div className="chat-messages">
            {messages.map((msg, index) => (
                <div className={`message-item ${getClassName(msg)} `} key={index}>
                    <strong style={{ color: getUserIdColor(msg.user_id) }}>{msg.username}</strong>
                    <p>{msg.content}</p>
                    <small>{formatDate(msg.created_at)}</small>
                </div>
            ))}
            <div ref={messagesEndRef}></div>
        </div>
    );
}

export default ChatMessages;