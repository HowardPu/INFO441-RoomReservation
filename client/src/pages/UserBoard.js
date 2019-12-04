import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';

class User extends React.Component {
    constructor(props) {
        super(props);
    }


    render() {
        return (
            <div>
                <h1>User Board</h1>
                <h2>View Your Reservations</h2>
                <ReservationList appState={this.props.appState} updateState={this.props.updateState}/>
                <br />
                <h2>Search Rooms</h2>
                <RoomList appState={this.props.appState} updateState={this.props.updateState}/>                
            </div>
        );
    }
}

export default User;