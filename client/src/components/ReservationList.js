import React from 'react';
import Form from 'react-bootstrap/Form'
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button'
import Table from 'react-bootstrap/Table'
import {Redirect} from 'react-router-dom';

const host = "http://localhost"
const jsonHeader =  {
    'Authorization': localStorage.getItem('auth')
}
const getResURL = host + "v1/reserve"

class ReservationList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showRes: false,
            data: {}
        }
    }

    onView(e) {
        e.preventDefault();
    }

    getData(url, userInput, headerInput) {
        fetch(url, {
            method: 'GET',
            mode: "cors",
            headers: headerInput, 
            body: JSON.stringify(userInput)
        }).then(resp => {
            if (resp.ok) {
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            console.log(data);
            return data;
        }).catch(err => {
            var errMes = "Oops something might be wrong! We will fix it soon!"
            console.log(err)
            this.setState({errMes});
            return null;
        })
    }



    render() {
        return (
            <div>
                <Button 
                    variant="primary" 
                    onClick={(e)=>{this.onView(e)}}>
                        View My Reservations
                </Button>
                {this.state.showRes && this.state.data && 
                    <Table variant="dark">
                        <thead>
                            <tr>
                                <th scope="col">Name</th>
                                <th scope="col">Floor</th>
                                <th scope="col">Capacity</th>
                                <th scope="col">Type</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.renderData()}
                        </tbody>
                    </Table>
                }
            </div>
        );
    }
}

export default ReservationList