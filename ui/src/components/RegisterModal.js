import React, {useState} from 'react';
import {Alert, Button, Form, Modal} from 'react-bootstrap';

const RegisterModal = ({show, handleClose, onSuccess}) => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [alertMessage, setAlertMessage] = useState('');
    const [alertVariant, setAlertVariant] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: name,
                    email: email,
                    password: password
                })
            });

            if (response.ok) {
                onSuccess();
                handleClose();
                setAlertMessage('User Register successful!');
                setAlertVariant('success');
            } else {
                console.error('Failed to create chat room')
                setAlertMessage('Register failed!');
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
                <Modal.Title>Register</Modal.Title>
            </Modal.Header>
            <Form onSubmit={handleSubmit}>
                <Modal.Body>
                    {alertMessage && <Alert variant={alertVariant}>{alertMessage}</Alert>}
                    <Form.Group controlId="formName">
                        <Form.Label>Name</Form.Label>
                        <Form.Control type="text" placeholder="Enter name" value={name}
                                      onChange={(e) => setName(e.target.value)}/>
                    </Form.Group>

                    <Form.Group controlId="formEmail">
                        <Form.Label>Email</Form.Label>
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
                    <Button variant="primary" type="submit">Register</Button>
                </Modal.Footer>
            </Form>
        </Modal>
    );
};

export default RegisterModal;
