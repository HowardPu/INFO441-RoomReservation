import React, { Component } from 'react';
import {Redirect, Link} from 'react-router-dom';
import {AuthButton} from '../components/AuthButton'
import {RoomCard} from './Room'
import Select from 'react-select';

// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button
const equip = "https://api.html-summary.me/v1/equip"
const getAllRoom = "https://api.html-summary.me/v1/room"
const specificRoom = "https://api.html-summary.me/v1/specificRoom"

export class SpecificRoom extends Component {
    constructor(props) {
        super(props)
        this.state = {
            back: false,
            rooms: null,
            research: false,
            equips: null,
            curEquip: null,
            curRoom: "",
            roomEquipID: "",
            addDisabled: true,
            research: false,
            deleteID: "",
            deleteDisable: true
        }

        this.options = []
    }

    componentDidMount() {
        if(this.props.appState.userType == "Admin") {
            this.searchAllRooms()
            this.searchEquip()
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (this.state.research) {
            this.setState({
                research: false
            })
            this.searchAllRooms()
        }
    }

    searchEquip() {
        fetch(equip, {
            method: 'GET',
            mode: "cors",
            headers: {'Authorization': this.props.appState.authToken}
          }).then(resp => {
            return resp.json();
          }).then(data => {
              let result = []
              data.result.forEach((equip) => {
                result.push({value: equip, label: equip})
              })
              this.setState({
                  equips: result,
                  curEquip: result[0]
              })
          }).catch(err => {
              console.log(err)
          })
    }

    addEquipInRoom() {
        fetch(specificRoom, {
            method: 'POST',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify({
                "equipName": this.state.curEquip.value,
                "roomName": this.state.curRoom
            })
        }).then(() => {
            this.setState({
                curRoom: "",
                addDisabled: true,
                research: true,
                rooms: null
            })
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }


    deleteEquipInRoom() {
        fetch(specificRoom, {
            method: 'DELETE',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify({
                "roomEquipID": parseInt(this.state.deleteID)
            })
        }).then(() => {
            this.setState({
                deleteID: "",
                deleteDisable: true,
                research: true,
                rooms: null
            })
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
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

        if(!this.state.rooms || !this.state.equips) {
            return(<div>
                Loading rooms/equips
            </div>)
        }
        return(
            <div>
                <div>
                    <h1>Specific Room Manager</h1>
                    <AuthButton signOutHandler={this.props.signOutHandler} />
                    <button className="btn btn-primary mr-2" disabled={this.state.disabled} onClick={() => {
                        this.setState({
                            back: true
                        })
                    }}>Back</button>

                    <hr />  

                    <div>
                        <h1>All Rooms</h1>
                        {this.state.rooms.map((room) => {
                            return <RoomCard key={"room-" + room.id} room={room} appState={this.props.appState} detail={true}/>
                        })}

                    </div>

                    <hr />  

                    <div>
                        <h1>Add Equipment to A Room</h1>

                        <h2>Select A Equipment</h2>
                        <Select name="equip" 
                            options={this.state.equips}
                            value={this.state.curEquip}
                            onChange={(event) => {
                                this.setState({
                                    curEquip: event
                                })
                        }}/>

                        <h2>Room Name</h2>

                        <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.curRoom}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = value.length == 0 
                                this.setState({
                                    curRoom: value,
                                    addDisabled: disabled
                                })
                            }}/>
                        </div>

                        <button className="btn btn-primary mr-2"
                            onClick={() => {
                                this.addEquipInRoom()
                            }} disabled={this.state.addDisabled}>
                                Add!
                        </button>
                    </div>
                </div>   

                <div className="log-in-item"> 
                    <h2>Delete Equipment of A Room</h2>

                    <h2>Please Give A ID</h2>
                    <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.deleteID}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = isNaN(parseInt(value, 10))
                                this.setState({
                                    deleteID: value,
                                    deleteDisable: disabled
                                })
                            }}/>
                    </div>

                    <button className="btn btn-primary mr-2"
                            onClick={() => {
                                this.deleteEquipInRoom()
                            }} disabled={this.state.deleteDisable}>
                                Delete!
                    </button>
                </div>

            </div>
        )
    }
}