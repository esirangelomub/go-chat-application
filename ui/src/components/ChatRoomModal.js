import React, {useState} from 'react';
import {Alert, Button, Form, Modal} from 'react-bootstrap';

const ChatRoomModal = ({show, handleClose, onSuccess}) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [alertMessage, setAlertMessage] = useState('');
    const [alertVariant, setAlertVariant] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const response = await fetch(`${process.env.REACT_APP_API_BASE_URL}/chats/rooms`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name: name,
                        description: description
                    })
                });

                if (response.ok) {
                    onSuccess();
                    handleClose();
                    setAlertMessage('Chat Room created successful!');
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
        }
    };

    return (
        <Modal show={show} onHide={handleClose}>
            <Modal.Header closeButton>
                <Modal.Title>New Chat Room</Modal.Title>
            </Modal.Header>
            <Form onSubmit={handleSubmit}>
                <Modal.Body>
                    {alertMessage && <Alert variant={alertVariant}>{alertMessage}</Alert>}
                    <Form.Group controlId="formName">
                        <Form.Label>Name</Form.Label>
                        <Form.Control type="text" placeholder="Enter name" value={name}
                                      onChange={(e) => setName(e.target.value)}/>
                    </Form.Group>

                    <Form.Group controlId="formDescription">
                        <Form.Label>Description</Form.Label>
                        <Form.Control type="text" placeholder="Enter description" value={description}
                                      onChange={(e) => setDescription(e.target.value)}/>
                    </Form.Group>

                </Modal.Body>
                <Modal.Footer>
                    <Button variant="secondary" onClick={handleClose}>Cancel</Button>
                    <Button variant="primary" type="submit">Save Chat Room</Button>
                </Modal.Footer>
            </Form>
        </Modal>
    );
};

export default ChatRoomModal;
