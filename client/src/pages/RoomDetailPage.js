import React from 'react';
import ReservationForm from '../components/ReservationForm';
import { Link } from 'react-router-dom';

const host = "https://api.html-summary.me/" //!!change it later
const reserveURL = host + "/v1/reserve"
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}

class RoomDetailPage extends React.Component {

    constructor(props) {
        super(props);
        console.log(this.props.location.state)
        this.state = {
            data: {}
        }
    }

    render() {
        return(
            <div>
                <h1>Room Detail</h1>
                <br/>
                <h2>{this.props.location.state.roomInfo.roomName}</h2>
                <div className="roomInfoContainer">
                    {this.props.location.state.roomInfo.floor && 
                        <div className="roomInfoItem"> 
                            <small>Floor</small> 
                            <div>
                            {this.props.location.state.roomInfo.floor}
                            </div>
                            <br />
                        </div>
                    }
                    {this.props.location.state.roomInfo.capacity && 
                        <div className="roomInfoItem"> 
                            <small>Capacity</small> 
                            <div>{this.props.location.state.roomInfo.capacity}</div>
                            <br />
                        </div>
                    }
                    <div className="roomInfoItem"> 
                        <small>Type</small> 
                        <div>
                            {this.props.location.state.roomInfo.roomType}
                        </div>
                        <br />
                    </div>
                </div>
                <br />
                <h3>Equipments</h3>
                    
                <h2>Reserve the Room</h2>
                <ReservationForm newRes={this.props.location.state.newRes} roomName={this.props.location.state.roomInfo.roomName}></ReservationForm>
                <Link to="/user">Back to User Board</Link>
            </div>
        );
    }
}

export default RoomDetailPage;