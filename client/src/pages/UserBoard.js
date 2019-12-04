import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';

const host = "api.html-summary.me/";
const client = new WebSocket("ws://" + host + "v1/ws?auth=" + localStorage.getItem('auth'));

class User extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            newRes: {}
        }
    }

    componentDidMount() {
        client.onopen = () => {
          console.log('WebSocket Client Connected');
        };

        client.onmessage = (message) => {
            console.log(message);
            // let messageObj = message.json();
            if (message.data.type === "reservation-create") {
                this.setState({newRes: message.data})
            }
            
        };

        client.onerror = (err) => {
            console.log(err);
        };

        client.onclose = (event) => {
            console.log("WebsocketStatus: Closed")
        };
    }
    
    render() {
        return (
            <div>
                <h1>User Board</h1>
                <h2>View Your Reservations</h2>
                <ReservationList />
                <br />
                <h2>Search Rooms</h2>
                <RoomList newRes={this.state.newRes}/>                
            </div>
        );
    }
}

export default User;