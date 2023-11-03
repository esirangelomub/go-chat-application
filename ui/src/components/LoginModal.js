import React, {useState} from 'react';
import {Button, Form, Modal, Alert} from 'react-bootstrap';

const LoginModal = ({show, handleClose, onSuccess}) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [alertMessage, setAlertMessage] = useState('');
    const [alertVariant, setAlertVariant] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/users/generate_token`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('token', data.access_token);
                onSuccess();
                handleClose();
                setAlertMessage('Login successful!');
                setAlertVariant('success');
            } else {
                console.error('Failed to generate token');
                setAlertMessage('Login failed. Please check your credentials.');
                setAlertVariant('danger');
            }
        } catch (error) {
            console.error(error);
            setAlertMessage('An error occurred. Please try again.');
            setAlertVariant('danger');
        }
    };

    return (
        <Modal show={show} onHide={handleClose}>
            <Modal.Header closeButton>
                <Modal.Title>Login</Modal.Title>
            </Modal.Header>
            <Form onSubmit={handleSubmit}>
                <Modal.Body>
                    {alertMessage && <Alert variant={alertVariant}>{alertMessage}</Alert>}
                    <Form.Group controlId="formEmail">
                        <Form.Label>Username</Form.Label>
                        <Form.Control type="text" placeholder="Enter email" value={email}
                                      onChange={(e) => setEmail(e.target.value)}/>
                    </Form.Group>

                    <Form.Group controlId="formPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control type="password" placeholder="Password" value={password}
                                      onChange={(e) => setPassword(e.target.value)}/>
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={handleClose}>Cancel</Button>
                    <Button variant="primary" type="submit">Sign In</Button>
                </Modal.Footer>
            </Form>
        </Modal>
    );
};

export default LoginModal;
