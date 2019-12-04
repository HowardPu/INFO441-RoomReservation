import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';

const host = "api.html-summary.me/";

class User extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            client: new WebSocket("wss://" + host + "v1/ws?auth=" + localStorage.getItem('auth'))
        }
    }

    // componentDidMount() {
    //     const client = new WebSocket("wss://" + host + "v1/ws?auth=" + localStorage.getItem('auth'));
    //     client.onopen = () => {
    //       console.log('WebSocket Client Connected');
    //     };

    //     client.onmessage = (message) => {
    //     };

    //     client.onerror = (err) => {
    //         console.log(err);
    //     };

    //     client.onclose = (event) => {
    //         console.log("WebsocketStatus: Closed")
    //     };
    // }
    
    render() {
        return (
            <div>
                <h1>User Board</h1>
                <h2>View Your Reservations</h2>
                <ReservationList ws={this.state.client}/>
                <br />
                <h2>Search Rooms</h2>
                <RoomList ws={this.state.client}/>                
            </div>
        );
    }
}

export default User;