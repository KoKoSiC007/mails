import { useEffect, useState } from "react"
import { APIClient } from "../internal/http"
import { Form, Button } from "react-bootstrap"
import Select from 'react-select';

export const Currencies = () => {
    const [startValue, setStart] = useState()
    const [endValue, setEnd] = useState()
    const [cData, setCData] = useState([])
    const client = new APIClient()
    const [currencyOptions, setCurs] = useState([])
    const [selectedCurs, setSelectedCurs] = useState([])

    useEffect(() => {
        client.get('/currencies').then(data => {
            const allowedCurs = []
            data.forEach(element => {
                if (element.enable) {
                    allowedCurs.push({value: element.name, label: element.name})
                }
            });
            setCurs(allowedCurs)
        }).catch(error => console.error(error))
    }, [])

    const submit = (e) => {
        e.preventDefault()
        console.warn(startValue, endValue, selectedCurs)
        client.post(`/currencies/report`, {startDate: startValue, endDate: endValue, currencies: selectedCurs.map((item, i) => item.value)}).then(data => {
            if (!!data) { 
                setCData(data)
            } else {
                setCData([])
            }
        }).catch(error => console.error(error))
    }

    const currencies = cData.map((item, i) => (
        <tr>
            <td>{item.name}</td>
            <td>{item.maxRate}</td>
            <td>{item.minRate}</td>
            <td>{item.avgRate}</td>
        </tr>
    ))

    return (
        <div>
            <Form onSubmit={submit}>
                <Form.Group>
                    <Form.Label>Start Date</Form.Label>
                    <Form.Control
                        type="date"
                        placeholder="2025-02-01"
                        onChange={e => {
                            setStart(e.target.value)
                        }}
                    />
                </Form.Group>

                <Form.Group>
                    <Form.Label>End Date</Form.Label>
                    <Form.Control
                        type="date"
                        placeholder="2025-02-01"
                        onChange={e => {
                            setEnd(e.target.value)
                        }}
                    />
                </Form.Group>
                <Select 
                    options={currencyOptions}
                    isMulti
                    onChange={e => {
                        setSelectedCurs(e)
                    }}
                    value={selectedCurs}
                />
                <Button
                variant="primary"
                type="submit"
                className="w-100 mt-3"
                > Show </Button>
            </Form>

            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Max</th>
                        <th>Min</th>
                        <th>Avg</th>
                    </tr>
                </thead>
                <tbody>
                    {currencies}
                </tbody>
            </table>
        </div>
    )
}