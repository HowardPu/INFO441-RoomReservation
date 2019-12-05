import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';
import { AuthButton } from '../components/AuthButton';

class User extends React.Component {

    render() {
        return (
            <div>
                <h1>User Board</h1>
                <h2>View Your Reservations</h2>
                <ReservationList appState={this.props.appState} updateState={this.props.updateState}/>
                <br />
                <h2>Search Rooms</h2>
                <RoomList appState={this.props.appState} updateState={this.props.updateState}/>  

                <AuthButton signOutHandler={this.props.signOutHandler} />
            </div>
        );
    }
}

export default User;