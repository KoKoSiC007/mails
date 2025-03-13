import React, { useEffect, useState } from 'react';
import { APIClient } from '../internal/http';
import { Link } from 'react-router-dom';

export const Mails = () => {
    const [emails, setEmails] = useState([])
    const client = new APIClient()
    useEffect(() => {
        client.get('/mails').then(data => {
            setEmails(data.messages)
        }).catch(error => console.error(error))
    }, [])

    const mails = emails.map((item, i) => (
        <tr>
            <td>{item.id}</td>
            <td>{item.to}</td>
            <td>{item.body}</td>
        </tr>
    ))

    return (
        <div>
            <Link to='/mails/new'>New Mail</Link>
            <table>
                <thead>
                    <tr>
                        <th>Id</th>
                        <th>To</th>
                        <th>Body</th>
                    </tr>
                </thead>
                <tbody>
                    {mails}
                </tbody>
            </table>
        </div>
    )
}