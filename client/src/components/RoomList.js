import React from 'react';
import Form from 'react-bootstrap/Form'
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button'
import Table from 'react-bootstrap/Table'

const host = ""
const jsonHeader =  {
    'Authorization': localStorage.getItem('auth')
}
const getRoomURL = host + "v1/room"

const TestData = [{
    roomName: "Test",
    capacity: 3,
    roomType:"study"
},
{
    roomName: "Test",
    capacity: 3,
    floor:2,
    roomType:"study"
}]

class RoomList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            errMes: '',
            data: TestData,
            showRooms: false,
            name: '',
            floor: '',
            capacity: '',
            type: '',
        }
    }

    onSubmit(e){
        e.preventDefault();
        this.setState({
            showRooms: true
        })
        let userInput = {
            roomName: this.state.name,
            floor: this.state.floor,
            capacity: this.state.capacity,
            roomType: this.state.type
        }
        var data = this.getData(getRoomURL, userInput, jsonHeader);
        this.setState({data: data});
    }

    onChange(e) {
        this.setState(
            {
                notification: '',
                errMes: ''
            }
        )
        switch (e.target.id) {
            case "searchName":
                this.setState({name: e.target.value});
                break;
            case "searchFloor":
                var curFloor;
                if (isNaN(e.target.value) === true) {
                    this.setState({errMes: "Floor must be a number"})
                    curFloor = '';
                } else {
                    curFloor = e.target.value;
                }
                this.setState({floor: curFloor});
                break;
            case "searchCapacity":
                var curCapacity;
                if (isNaN(e.target.value) === true) {
                    this.setState({errMes: "Capacity must be a number"})    
                    curCapacity = '';
                } else {
                    curFloor = e.target.value;
                }
                this.setState({capacity: curCapacity});
                break;
            case "searchType":
                this.setState({type: e.target.value});
                break;
            default:
                break;
        }
    }

    renderData() {
        console.log(this.state.data)
        return this.state.data.map((item, i) => {
            let name = item.roomName;
            let capacity = item.capacity;
            let floor = item.floor;
            let roomType = item.roomType;
            return(
                <tr key={i}>
                    <td key={"name "+name}>{name}</td>
                    <td key={"floor "+name}>{floor ? floor : "n/a"}</td>
                    <td key={"capa "+name}>{capacity ? capacity : "n/a"}</td>
                    <td key={"type "+name}>{roomType}</td>
                    <td key={"btn "+name}><Button onClick={(e, data) => this.onReserve(e, item)}>Reserve</Button></td>
                </tr>
            )
        })
    }
    
    onReserve(e, data) {
        console.log(e);
        console.log(data)
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
            <section>
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
               <Form>
                    <Form.Group controlId="searchName">
                        <Form.Label>Room Name</Form.Label>
                        <Form.Control 
                            value={this.state.name}
                            onChange={(e) => {this.onChange(e)}}
                            placeholder="Enter Room Name" />
                    </Form.Group>

                    <Form.Row>
                        <Form.Group as={Col} controlId="searchFloor">
                            <Form.Label>Floor</Form.Label>
                            <Form.Control 
                                value={this.state.floor}
                                onChange={(e) => {this.onChange(e)}}
                                placeholder="Enter Room Floor" />
                            <Form.Text>(Optional)</Form.Text>
                        </Form.Group>

                        <Form.Group as={Col} controlId="searchCapacity">
                            <Form.Label>Capacity</Form.Label>
                            <Form.Control 
                                value={this.state.capacity}
                                onChange={(e) => {this.onChange(e)}}
                                placeholder="Enter Room Capacity" />
                            <Form.Text>(Optional)</Form.Text>
                        </Form.Group>
                    </Form.Row>

                    <Form.Group controlId="searchType">
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
                            View Rooms
                    </Button>
                </Form>
                <br />
                {this.state.showRooms && this.state.data && 
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
            </section>
        );
    }
}

export default RoomList;