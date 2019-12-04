import React from 'react';
import Form from 'react-bootstrap/Form'
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button'
import Table from 'react-bootstrap/Table'
import {Redirect} from 'react-router-dom';

const host = "http://api.html-summary.me"
const jsonHeader =  {
    'Authorization': localStorage.getItem('auth')
}
const getRoomURL = host + "/v1/room"

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
            data: [],
            showRooms: false,
            name: '',
            floor: '',
            capacity: '',
            type: '',
            clickReserve: false,
            reserveRoom: {},
        }
    }

    onSubmit(e){
        e.preventDefault();
        let roomName = this.state.name ? this.state.name : "*"
        let roomType = this.state.type ? this.state.type : "*"
        var url = `${getRoomURL}?roomname=${roomName}&roomtype=${roomType}`
        if (this.state.floor) {
            url = url + "&floor=" + this.status.floor
        }

        if (this.state.capacity) {
            url = url + "&floor=" + this.status.capacity
        }
        console.log(url)
        this.getData(url, jsonHeader);
    }


    onChange(e) {
        this.setState(
            {
                notification: '',
                errMes: '',
                reserveRoom: null
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
            let roomName = item.roomName;
            let capacity = item.capacity;
            let floor = item.floor;
            let roomType = item.roomType;
            return(
                <tr key={i}>
                    <td key={"name "+i}>{roomName}</td>
                    <td key={"floor "+i}>{floor ? floor : "n/a"}</td>
                    <td key={"capa "+i}>{capacity ? capacity : "n/a"}</td>
                    <td key={"type "+i}>{roomType}</td>
                    <td key={"btn "+i}><Button value={item} onClick={() => {this.setState({clickReserve: true, reserveRoom: item})}}>Reserve</Button></td>
                </tr>
            )
        })
    }

    getData(url, headerInput) {
        fetch(url, {
            method: 'GET',
            mode: "cors",
            headers: headerInput, 
        }).then(resp => {
            if (resp.ok) {
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            this.setState({data:data});
            this.setState({ showRooms: true})
        }).catch(err => {
            var errMes = "Oops something might be wrong! We will fix it soon!"
            console.log(err)
            this.setState({errMes});
            return null;
        })
    }


    render() {
        if (this.state.clickReserve && this.state.reserveRoom !== null) {
            console.log(this.state.reserveRoom)
            var passState = {
                roomInfo: this.state.reserveRoom
            };

            if (this.props.newRes && 
                this.props.newRes.roomName === this.state.reserveRoom.roomName) {
                passState.newRes = this.props.newRes
            }
            return (<Redirect to={{pathname:'/reserve', state:passState}} />)
        }
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
                {this.state.showRooms === true && this.state.data && 
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