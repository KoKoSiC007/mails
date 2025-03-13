import { useState } from "react"
import { APIClient } from "../internal/http"
import { useNavigate } from "react-router-dom"
import { Button, Form } from "react-bootstrap"

export const NewMail = () => {
    const [to, setTo] = useState("")
    const [body, setBody] = useState("")
    const client = new APIClient()
    const navigate = useNavigate()

    const submit = (e) => {
        e.preventDefault()
        client.post('/mail', {to: to, body: body})
        .then(_ => navigate('/mails'))
        .catch(error => console.error(error))
    }

    return (
        <div
        style={{height: "100vh"}}
        className="d-flex justify-content-center align-items-center">
            <div style={{width: 300}}>
                <h1 className="text-center">New Mail</h1>
                <Form onSubmit={submit}>
                    <Form.Group>
                        <Form.Label>To</Form.Label>
                        <Form.Control 
                            type="string"
                            placeholder="To"
                            onChange={e => {
                                setTo(e.target.value)
                            }}
                        />
                    </Form.Group>
                    <Form.Group>
                        <Form.Label>Body</Form.Label>
                        <Form.Control 
                            type="text"
                            placeholder="Body"
                            onChange={e => {
                                setBody(e.target.value)
                            }}
                        />
                    </Form.Group>
                    <Button
                    variant="primary"
                    type="submit"
                    className="w-100 mt-3"
                    > Send </Button>
                </Form>
            </div>
        </div>
    )
}