import { APIClient } from "../internal/http"
import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useNavigate } from "react-router-dom";


export const Auth = () => {
    const [emailValue, setEmail] = useState("")
    const [passValue, setPassword] = useState("")
    let client = new APIClient()
    const navigate = useNavigate()

    const submit = (e) => {
        e.preventDefault()
        client.login(`/user/login?login=${emailValue}&password=${passValue}`)
        .then(() => navigate('/'))
        .catch(() => navigate('/login'))
    }
    return (
        <div
        style={{ height: "100vh" }}
        className="d-flex justify-content-center align-items-center"
        >
        <div style={{ width: 300 }}>
            <h1 className="text-center">Sign in</h1>
            <Form onSubmit={submit}>
            <Form.Group>
                <Form.Label>Email address</Form.Label>
                <Form.Control
                type="email"
                placeholder="Enter email"
                onChange={e => {
                    setEmail(e.target.value);
                }}
                />
            </Form.Group>

            <Form.Group>
                <Form.Label>Password</Form.Label>
                <Form.Control
                type="password"
                placeholder="Password"
                onChange={e => {
                    setPassword(e.target.value);
                }}
                />
            </Form.Group>
            <Button
                variant="primary"
                type="submit"
                className="w-100 mt-3"
            > Sign in  </Button>
            </Form>
        </div>
        </div>
    ) 
}