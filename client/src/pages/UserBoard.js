import React from 'react';
import RoomList from '../components/RoomList';


class User extends React.Component {
    
    render() {
        return (
            <div>
                <h1>User Board</h1>
                <br />
                <h2>Search Rooms</h2>
                <RoomList />
            </div>
        );
    }
}

export default User;