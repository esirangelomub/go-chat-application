import React from 'react';

const initialUserData = {
    isLoggedIn: false,
    access_token: null,
    user_id: null,
    email: null
};

const UserContext = React.createContext({
    userData: initialUserData,
    setUserData: () => {}
});

export default UserContext;
