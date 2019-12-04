import React from 'react';
import Form from 'react-bootstrap/Form'
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button'
import Select from 'react-select';

const host = "https://api.html-summary.me/"
const addRoomURL = host + "v1/room"

class AddRoomForm extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            errMes: null,
            name: '',
            floor: '',
            capacity: '',
            type: 'Study',
            notification: '',
            valid: false
        }

        let rmTypes = ["Study", "Teamwork", "Demonstration", "Lounge", "Computer Lab", "Other"]
        this.options = []

        rmTypes.forEach((rmType) => {
            this.options.push({value: rmType, label: rmType});
        })

        this.checkValidity = this.checkValidity.bind(this)
        this.setAttribute = this.setAttribute.bind(this)
        this.postData = this.postData.bind(this)
    }

    setAttribute(field, value) {
        let currentState = this.state;
        currentState[field] = value;
        this.setState(currentState);
    }

    checkValidity() {
        let cap = parseInt(this.state.capacity, 10)
        let floor = parseInt(this.state.floor, 10)
        this.setState({
            valid: isNaN(cap) || isNaN(floor) || this.state.name.length == 0 
        })
    }

    onSubmit(){
        let floor = parseInt(this.state.floor, 10)
        let capacity = parseInt(this.state.capacity, 10)

        let userInput = {
            roomName: this.state.name,
            capacity: capacity,
            floor: floor,
            roomType: this.state.type
        }
        this.postData(addRoomURL, userInput);
    }

    postData(url, userInput) {
        fetch(url, {
            method: 'POST',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify(userInput)
        }).then(() => {
            
            this.props.setResearch()
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }


    render() {
        console.log(this.state.valid)
        return (
            <section>
                <h2>Add Room</h2>
                {!this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                <div className="formContainer">
                    <Form>
                        <Form.Group controlId="adminAddName">
                            <Form.Label>Room Name</Form.Label>
                            <Form.Control 
                                value={this.state.name}
                                onChange={(e) => {
                                    this.setAttribute("name", e.target.value)
                                    this.checkValidity()
                                }}
                                placeholder="Enter Room Name" />
                        </Form.Group>

                        <Form.Row>
                            <Form.Group as={Col} controlId="adminAddFloor">
                                <Form.Label>Floor</Form.Label>
                                <Form.Control 
                                    value={this.state.floor}
                                    onChange={(e) => {
                                        this.setAttribute("floor", e.target.value)
                                        this.checkValidity()
                                    }}
                                    placeholder="Enter Room Floor" />
                            </Form.Group>

                            <Form.Group as={Col} controlId="adminAddCapacity">
                                <Form.Label>Capacity</Form.Label>
                                <Form.Control 
                                    value={this.state.capacity}
                                    onChange={(e) => {
                                        this.setAttribute("capacity", e.target.value)
                                        this.checkValidity()
                                    }}
                                    placeholder="Enter Room Capacity" />
                            </Form.Group>
                        </Form.Row>

                        <Select name="roomType" 
                            options={this.options}
                            value={{value: this.state.type, label: this.state.type}}
                            onChange={(event) => {
                                this.setAttribute("type", event.value);
                        }}/>


                        <Button 
                            disabled={this.state.valid}
                            onClick={()=>{this.onSubmit()}}>
                            Add Room
                        </Button>
                    </Form>
                </div>
            </section>
        );
    }
}

export default AddRoomForm;