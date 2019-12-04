import React from 'react';
import ReservationForm from '../components/ReservationForm';
import { Link } from 'react-router-dom';

const host = "https://api.html-summary.me"
const getEquipURL = host + "/v1/specificRoom"
const jsonHeader =  {
    'Authorization': localStorage.getItem('auth')
}

class RoomDetailPage extends React.Component {

    constructor(props) {
        super(props);
        console.log(this.props.updateState)
        this.state = {
            data: {},
            equip: null
        }
    }

    componentDidMount() {
        let url = getEquipURL + '?roomname=' + this.props.location.state.roomInfo.roomName;
        this.getData(url, jsonHeader);    
    }

    renderEquip() {
        var equips = []
        this.state.equip.map((item, i) => {
            equips.push(
                <li key={i}>{item.Name}</li>
            )
        })
        return (
            <ul>
                {equips}
            </ul>
        )
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
            console.log(data)
            this.setState({equip:data});
        }).catch(err => {
            var errMes = "Oops something might be wrong! We will fix it soon!"
            console.log(err)
            this.setState({errMes});
            return null;
        })
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
                {this.state.equip && this.state.equip.length !== 0 ? this.renderEquip() : <div>No Equipments Info</div>}
                <br />
                <h2>Reserve the Room</h2>
                <ReservationForm appState={this.props.appState} updateState={this.props.updateState} roomName={this.props.location.state.roomInfo.roomName}></ReservationForm>
                <br />
                <Link to="/user">Back to User Board</Link>
            </div>
        );
    }
}

export default RoomDetailPage;