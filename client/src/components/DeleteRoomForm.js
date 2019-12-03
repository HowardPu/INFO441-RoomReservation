import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Alert from 'react-bootstrap/Alert'

const host = ""
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}
const delRoomURL = host + "v1/room"

class DeleteRoomForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            errMes: '',
            name: '',
            notification: ''
        }
    }

    onChange(e) {
        this.setState(
            {
                notification: '',
                errMes: ''
            } 
        )
        this.setState({name: e.target.value})
    }

    onSubmit(e){
        e.preventDefault();
        if (!this.state.name) {
            this.setState({errMes: "Please input room name"})
        } else {
            let userInput = {roomName: this.state.name}
            console.log(userInput)
            this.patchData(delRoomURL, userInput, jsonHeader);
        }
    }

    patchData(url, userInput, headerInput) {
        fetch(url, {
            method: 'PATCH',
            mode: "cors",
            headers: headerInput, 
            body: JSON.stringify(userInput)
        }).then(resp => {
            if (resp.ok) {
                if (!headerInput.Authorization && resp.headers.get('Authorization')) {
                    localStorage.setItem('auth', resp.headers.get('Authorization'));
                }
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            console.log(data);
            let mes = "successfully added room" + data.roomName;
            this.setState({notification: mes})
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }


    render() {
        return (
            <section>
                <h2>Delete Room</h2>
                {this.state.notification && <Alert variant="success">{this.state.notification}</Alert>}
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                <div className="formContainer">
                    <Form>
                        <Form.Group controlId="adminDelName">
                            <Form.Label>Room Name</Form.Label>
                            <Form.Control 
                                value={this.state.name}
                                onChange={(e) => {this.onChange(e)}}
                                placeholder="Enter Room Name" />
                        </Form.Group>

                        <Button 
                            variant="danger" 
                            type="submit"
                            onClick={(e)=>{this.onSubmit(e)}}>
                            Delete Room
                        </Button>
                    </Form>
                </div>
            </section>
        );
    }
}

export default DeleteRoomForm;