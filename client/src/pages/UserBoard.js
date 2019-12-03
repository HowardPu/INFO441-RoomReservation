import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';

const host = "";
const client = new WebSocket("wss://" + host + "/v1/ws?auth=" + localStorage.getItem('auth'));

class User extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            newRes: {}
        }
    }

    componentWillMount() {
        client.onopen = () => {
          console.log('WebSocket Client Connected');
        };

        client.onmessage = (message) => {
            console.log(message);
            let messageObj = message.json();
            if (messageObj.type === "reservation-create") {
                this.setState({newRes: messageObj})
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