import React, {useEffect, useState} from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import UserContext from './contexts/UserContext';
import Sidebar from './components/Sidebar';
import ChatArea from './components/ChatArea';

function App() {
    const [user, setUser] = useState({
        isLoggedIn: false,
        token: null,
        email: null
    });

    const setUserData = (data) => {
        setUser(prevState => ({...prevState, ...data}));
    };

    useEffect(() => {
        const fetchData = async () => {
            const token = localStorage.getItem('token');
            if (token) {
                try {
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
                            token: token,
                            user_id: data.id,
                            email: data.email
                        });
                    } else {
                        console.error('Failed to generate token');
                        setUserData({
                            isLoggedIn: false,
                            access_token: null,
                            user_id: null,
                            email: null
                        });
                    }
                } catch (error) {
                    console.error("There was an error fetching user data:", error);
                    setUserData({
                        isLoggedIn: false,
                        access_token: null,
                        user_id: null,
                        email: null
                    });
                }
            } else {
                setUserData({
                    isLoggedIn: false,
                    access_token: null,
                    user_id: null,
                    email: null
                });
                localStorage.removeItem('token')
            }
        }
        fetchData()
    }, []);

    return (
        <Router>
            <UserContext.Provider value={{userData: user, setUserData}}>
                <div className="d-flex">
                    <Sidebar/>
                    <ChatArea/>
                </div>
            </UserContext.Provider>
        </Router>
    );
}

export default App;
