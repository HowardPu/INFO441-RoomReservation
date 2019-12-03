import React from 'react';
import Form from 'react-bootstrap/Form'
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button'
import Alert from 'react-bootstrap/Alert'

const host = ""
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}
const addRoomURL = host + "v1/room"

class AddRoomForm extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            errMes: '',
            name: '',
            floor: '',
            capacity: '',
            type: '',
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
        switch (e.target.id) {
            case "adminAddName":
                this.setState({name: e.target.value});
                break;
            case "adminAddFloor":
                var curFloor;
                if (isNaN(e.target.value) === true) {
                    this.setState({errMes: "Floor must be a number"})
                    curFloor = '';
                } else {
                    curFloor = e.target.value;
                }
                this.setState({floor: curFloor});
                break;
            case "adminAddCapacity":
                var curCapacity;
                if (isNaN(e.target.value) === true) {
                    this.setState({errMes: "Capacity must be a number"})    
                    curCapacity = '';
                } else {
                    curFloor = e.target.value;
                }
                this.setState({capacity: curCapacity});
                break;
            case "adminAddType":
                this.setState({type: e.target.value});
                break;
            default:
                break;
        }
    }

    onSubmit(e){
        e.preventDefault();
        if (!this.state.name) {
            this.setState({errMes: "Please input room name"})
        } else if (!this.state.type) {
            this.setState({errMes: "Please input room type"})
        } else {
            var floor = this.state.floor === '' ? null : this.state.floor;
            var capacity = this.state.capacity === '' ? null : this.state.capacity;

            let userInput = {
                roomName: this.state.name,
                capacity: capacity,
                floor: floor,
                roomType: this.state.roomType
            }
            console.log(userInput)
            this.postData(addRoomURL, userInput, jsonHeader);
        }
    }

    postData(url, userInput, headerInput) {
        fetch(url, {
            method: 'POST',
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
                <h2>Add Room</h2>
                {this.state.notification && <Alert variant="success">{this.state.notification}</Alert>}
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                <div className="formContainer">
                    <Form>
                        <Form.Group controlId="adminAddName">
                            <Form.Label>Room Name</Form.Label>
                            <Form.Control 
                                value={this.state.name}
                                onChange={(e) => {this.onChange(e)}}
                                placeholder="Enter Room Name" />
                        </Form.Group>

                        <Form.Row>
                            <Form.Group as={Col} controlId="adminAddFloor">
                                <Form.Label>Floor</Form.Label>
                                <Form.Control 
                                    value={this.state.floor}
                                    onChange={(e) => {this.onChange(e)}}
                                    placeholder="Enter Room Floor" />
                                <Form.Text>(Optional)</Form.Text>
                            </Form.Group>

                            <Form.Group as={Col} controlId="adminAddCapacity">
                                <Form.Label>Capacity</Form.Label>
                                <Form.Control 
                                    value={this.state.capacity}
                                    onChange={(e) => {this.onChange(e)}}
                                    placeholder="Enter Room Capacity" />
                                <Form.Text>(Optional)</Form.Text>
                            </Form.Group>
                        </Form.Row>

                        <Form.Group controlId="adminAddType">
                            <Form.Label>Type</Form.Label>
                            <Form.Control 
                                value={this.state.type}
                                onChange={(e) => {this.onChange(e)}}                            
                                placeholder="Enter Room Type" />
                        </Form.Group>

                        <Button 
                            variant="primary" 
                            type="submit"
                            onClick={(e)=>{this.onSubmit(e)}}>
                            Add Room
                        </Button>
                    </Form>
                </div>
            </section>
        );
    }
}

export default AddRoomForm;