import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import ReservationForm from '../components/ReservationForm';

const host = "http://localhost" //!!change it later
const reserveURL = host + "/v1/reserve"
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}

class ReserveRoom extends React.Component {

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
                <h2>{this.props.location.state.roomName}</h2>
                <div className="roomInfoContainer">
                    {this.props.location.state.floor && 
                        <div className="roomInfoItem"> 
                            <small>Floor</small> 
                            <div>
                            {this.props.location.state.floor}
                            </div>
                            <br />
                        </div>
                    }
                    {this.props.location.state.capacity && 
                        <div className="roomInfoItem"> 
                            <small>Capacity</small> 
                            <div>{this.props.location.state.capacity}</div>
                            <br />
                        </div>
                    }
                    <div className="roomInfoItem"> 
                        <small>Type</small> 
                        <div>
                            {this.props.location.state.roomType}
                        </div>
                        <br />
                    </div>
                </div>
                <br />
                <h3>Equipments</h3>

                <h2>Reserve the Room</h2>
                <ReservationForm></ReservationForm>
            </div>
        );
    }
}

export default ReserveRoom;