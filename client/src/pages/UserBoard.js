import React from 'react';
import RoomList from '../components/RoomList';
import ReservationList from '../components/ReservationList';

class User extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            newRes: {}
        }
    }

    componentDidMount() {
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