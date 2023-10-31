import React, { useContext, useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSignOutAlt, faFileText } from '@fortawesome/free-solid-svg-icons';
import UserContext from "../contexts/UserContext";
import LoginModal from './LoginModal';
import RegisterModal from "./RegisterModal";
import ChatRoomModal from "./ChatRoomModal";

const Sidebar = () => {
    const { userData, setUserData } = useContext(UserContext);
    const [showLogin, setShowLogin] = useState(false);
    const [showRegister, setShowRegister] = useState(false);
    const [showChatRoom, setShowChatRoom] = useState(false);
    const [chatRooms, setChatRooms] = useState([]);

    const handleLoginShow = () => setShowLogin(true);
    const handleLoginClose = () => setShowLogin(false);

    const handleRegisterShow = () => setShowRegister(true);
    const handleRegisterClose = () => setShowRegister(false);

    const handleChatRoomShow = () => setShowChatRoom(true);
    const handleChatRoomClose = () => setShowChatRoom(false);

    const handleLoginSuccess = async () => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('No token found');
            }
            const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/users/me`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                const data = await response.json();
                setUserData({
                    isLoggedIn: true,
                    access_token: token,
                    email: data.email
                });
                setShowLogin(false);
            } else {
                console.error('Failed to generate token')
            }
        } catch (error) {
            console.error(error);
            handleSignOut();
        }
    };

    const handleRegisterSuccess = async() => {
        console.log('chegou');
    }

    const handleChatRoomSuccess = async() => {
        console.log('chegou chat room');
    }

    const handleSignOut = () => {
        setUserData({
            isLoggedIn: false,
            access_token: null,
            email: null
        });
    };

    useEffect(() => {
        const fetchChatRooms = async () => {
            if (userData.isLoggedIn) {
                try {
                    const token = localStorage.getItem('token');
                    const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/chats/rooms`, {
                        method: 'GET',
                        headers: {
                            'Authorization': `Bearer ${token}`,
                            'Content-Type': 'application/json'
                        }
                    });
                    if (response.ok) {
                        const data = await response.json();
                        setChatRooms(data);
                    } else {
                        console.error('Failed to fetch chat rooms');
                    }
                } catch (error) {
                    console.error('Error fetching chat rooms:', error);
                }
            }
        };

        fetchChatRooms();
    }, [userData.isLoggedIn, userData.access_token]);

    return (
        <div className="sidebar">
            <div className="sidebar-header">
                <button className="btn btn-outline-primary w-100" disabled={!userData.isLoggedIn} onClick={handleChatRoomShow}>New Chat Room</button>
            </div>

            {/* Body */}
            <div className="sidebar-body">
                <ul>
                    {userData.isLoggedIn ? (
                        chatRooms.map(room => (
                            <li key={room.id}>
                                <FontAwesomeIcon icon={faFileText} /> {room.name}
                                <Link to={`/c/${room.id}`}>
                                    <button className="btn btn-outline-info ml-2">Open</button>
                                </Link>
                            </li>
                        ))
                    ) : (
                        <li className="text-center">Please sign in to see the chat room list</li>
                    )}
                </ul>
            </div>

            {/* Footer */}
            <div className="sidebar-footer">
                {userData.isLoggedIn ? (
                    <>
                        <p className="text-truncate" style={{ width: '85%' }}>{userData.email}</p>
                        <button className="btn btn-outline-secondary" onClick={handleSignOut}>
                            <FontAwesomeIcon icon={faSignOutAlt} />
                        </button>
                    </>
                ) : (
                    <div className="d-flex w-100">
                        <button className="btn btn-outline-success w-50" onClick={handleLoginShow}>Sign In</button>
                        <button className="btn btn-outline-info w-50" onClick={handleRegisterShow}>Register</button>
                    </div>
                )}
            </div>

            <LoginModal show={showLogin} handleClose={handleLoginClose} onSuccess={handleLoginSuccess} />
            <RegisterModal show={showRegister} handleClose={handleRegisterClose} onSuccess={handleRegisterSuccess} />
            <ChatRoomModal show={showChatRoom} handleClose={handleChatRoomClose} onSuccess={handleChatRoomSuccess} />
        </div>
    );
}

export default Sidebar;
