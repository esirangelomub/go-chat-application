import React from 'react';

const initialUserData = {
    isLoggedIn: false,
    access_token: null,
    email: null
};

const UserContext = React.createContext({
    userData: initialUserData,
    setUserData: () => {}
});

export default UserContext;
