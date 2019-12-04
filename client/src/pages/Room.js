import React, { Component } from 'react';
import {Redirect, Link} from 'react-router-dom';
import {AuthButton} from '../components/AuthButton';
import AddRoomForm from '../components/AddRoomForm'
import DeleteRoomForm from '../components/DeleteRoomForm'

// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button

const getAllRoom = "https://api.html-summary.me/v1/room"

export class Room extends Component {
    constructor(props) {
        super(props)
        this.state = {
            back: false,
            rooms: null,
            rearchRoom: this.props.appState.authToken,
            research: false
        }
        this.searchAllRooms = this.searchAllRooms.bind(this)
        this.setResearch = this.setResearch.bind(this)
    }

    componentDidMount() {
        if(this.props.appState.userType == "Admin") {
            this.searchAllRooms()
        }
    }

    setResearch() {
        this.setState({
            research: true
        })
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (this.state.research) {
            console.log(this.state.research)
            console.log(prevState.research)
            this.setState({
                research: false
            })
            this.searchAllRooms()
        }
    }

    searchAllRooms() {
        fetch(getAllRoom, {
            method: 'GET',
            mode: "cors",
            headers: {'Authorization': this.props.appState.authToken}
          }).then(resp => {
            return resp.json();
          }).then(data => {
              this.setState({
                  rooms: data
              })
          }).catch(err => {
              console.log(err)
          })
    }

    render() {
        let userType = this.props.appState.userType
        if(userType != "Admin") {
            return <Redirect to='/signin'/>
        }
        if(this.state.back) {
            return <Redirect to='/admin'/>
        }

        if(!this.state.rooms) {
            return(<div>
                Loading rooms
            </div>)
        }

        return(
            <div>
                <div>
                    <h1>Room Manager</h1>
                    <AuthButton signOutHandler={this.props.signOutHandler} />
                    <button className="btn btn-primary mr-2" onClick={() => {
                        this.setState({
                            back: true
                        })
                    }}>Back</button>
                </div> 

                <hr />  

                <div>
                    {this.state.rooms.map((room) => {
                        return <RoomCard key={"room-" + room.id} room={room} appState={this.props.appState}/>
                    })}

                </div>

                <hr /> 
                <AddRoomForm appState={this.props.appState} setResearch={this.setResearch} />

                <hr />
                <DeleteRoomForm appState={this.props.appState} setResearch={this.setResearch} />
            </div>
        )
    }
}


export class RoomCard extends Component {
    constructor(props) {
        super(props);

        this.state = {
            "equipments": [],
            "issues": [],
            "detail": false 
        }

        this.searchAllEquipments = this.searchAllEquipments.bind(this)
    }


    // search representatives and its rating
    // while changing the data to display
    // once those data are changed
    // (userrep is handled by its own star algorithm so it does not need to listen)
    componentDidMount(){
        this.searchAllEquipments()
    }

    searchAllEquipments() {
        let searchRoomURL = "https://api.html-summary.me/v1/specificRoom?roomname=" + this.props.room.roomName
        fetch(searchRoomURL, {
            method: 'GET',
            mode: "cors",
            headers: {'Authorization': this.props.appState.authToken}
            }).then(resp => {
                return resp.json();
            }).then(data => {
                console.log(data)
                this.setState({
                    equipments: data
                })
            }).catch(err => {
                console.log(err)
            })
    }

    render() {
        if (!this.state.equipments || !this.state.issues) {
            return(<div>Loading Cards</div>)
        }
    
        return(
            <div>
                {(!this.state.detail) &&
                    <div className="card" style={
                        {"width": "18rem"}
                    }>
                        <div className="card-body">
                            <h5 className="card-title mb-2 text-dark">{this.props.room.roomName}</h5>
                            
                            <h6 className="card-subtitle mb-2 text-muted">{"ID: " + this.props.room.id}</h6>
        
                            <p className="card-text mb-2 text-dark">{"Capacity: " + this.props.room.capacity}</p>
                            <p className="card-text mb-2 text-dark">{"Room Type: " + this.props.room.roomType}</p>
                            <button className="btn btn-primary mr-2"
                                onClick={() => {
                                    this.setState({
                                        detail: true
                                    })
                                }}>
                                    More Detail
                            </button>
                        </div>
                    </div>
                }
                {this.state.detail &&
                    <div className="card" style={
                        {"width": "18rem"}
                    }>
                        <div className="card-body">
                            <h5 className="card-title mb-2 text-dark">{this.props.room.roomName}</h5>
                            {this.state.equipments.length == 0 &&
                                <p className="card-text mb-2 text-dark">
                                        {"No Equipment"}
                                </p>
                            }
                            {this.state.equipments.map((equipment) => {
                                return <p key={"equip-room-" + equipment.roomEquipID} className="card-text mb-2 text-dark">
                                        {"Name: " + equipment.Name + " --- ID: " + equipment.roomEquipID}
                                </p>
                                
                                
                            })}
                            
                            <button className="btn btn-primary mr-2"
                                onClick={() => {
                                    this.setState({
                                        detail: false
                                    })
                                }}>
                                    Back
                            </button>
                        </div>
                    </div>
                }
            </div>
        );
        
    }
}