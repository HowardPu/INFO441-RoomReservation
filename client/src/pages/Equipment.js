import React, { Component } from 'react';
import {Redirect, Link} from 'react-router-dom';
import {AuthButton} from '../components/AuthButton'


const equip = "https://api.html-summary.me/v1/equip"
// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button
export class Equipment extends Component {
    constructor(props) {
        super(props)
        this.state = {
            back: false,
            equips: null,
            newEquipName: "",
            disabled: true,
            research: false,
            deleteEquip: "",
            delDisabled: true,
            oldName: "",
            newName: "",
            updateDisabled: true
        }
        this.searchEquip = this.searchEquip.bind(this)
    }

    componentDidMount() {
        if(this.props.appState.userType == "Admin") {
            this.searchEquip()
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (this.state.research) {
            this.setState({
                research: false
            })
            this.searchEquip()
        }
    }

    addEquip() {
        fetch(equip, {
            method: 'POST',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify({
                "equipName": this.state.newEquipName
            })
        }).then(() => {
            this.setState({
                newEquipName: "",
                research: true,
                disabled: true
            })
        }).catch(err => {
            console.log(err)
        })
    }

    deleteEquip() {
        fetch(equip, {
            method: 'DELETE',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify({
                "equipName": this.state.deleteEquip
            })
        }).then(() => {
            this.setState({
                deleteEquip: "",
                delDisabled: true,
                research: true
            })
        }).catch(err => {
            console.log(err)
        })
    }


    updateEquip() {
        fetch(equip, {
            method: 'PATCH',
            mode: "cors",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.props.appState.authToken
            }, 
            body: JSON.stringify({
                "equipName": this.state.oldName,
                "newName": this.state.newName
            })
        }).then(() => {
            this.setState({
                oldName: "",
                newName: "",
                updateDisabled: true,
                research: true
            })
        }).catch(err => {
            console.log(err)
        })
    }

    searchEquip() {
        fetch(equip, {
            method: 'GET',
            mode: "cors",
            headers: {'Authorization': this.props.appState.authToken}
          }).then(resp => {
            return resp.json();
          }).then(data => {
              this.setState({
                  equips: data.result
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

        if (!this.state.equips) {
            return(
                <div>Loading Equipments</div>
            )
        }
        return(
            <div>
                <div>
                    <h1>Equipment Manager</h1>
                    <AuthButton signOutHandler={this.props.signOutHandler} />
                    <button className="btn btn-primary mr-2"  onClick={() => {
                        this.setState({
                            back: true
                        })
                    }}>Back</button>
                </div>   

                <div>
                    <h2>All Existing Equipments</h2>

                    {this.state.equips.map((equip, index) => {
                        return <p key={"equip-" + index}>{equip}</p>
                    })}
                </div>

                <hr /> 

                <div className="log-in-item"> 
                    <h2>Add New Equipment</h2>
                    <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.newEquipName}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = value.length == 0
                                this.setState({
                                    newEquipName: value,
                                    disabled: disabled
                                })
                            }}/>
                    </div>

                    <button className="btn btn-primary mr-2"
                            onClick={() => {
                                this.addEquip()
                            }} disabled={this.state.disabled}>
                                Add!
                    </button>
                </div>
                <hr />

                <div className="log-in-item"> 
                    <h2>Delete Equipment</h2>
                    <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.deleteEquip}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = value.length == 0
                                this.setState({
                                    deleteEquip: value,
                                    delDisabled: disabled
                                })
                            }}/>
                    </div>

                    <button className="btn btn-primary mr-2"
                            onClick={() => {
                                this.deleteEquip()
                            }} disabled={this.state.delDisabled}>
                                Delete!
                    </button>
                </div>
                
                <hr />

                <div className="log-in-item"> 
                    <h2>Change Equipment Name</h2>

                    <h3>Current Equipment Name</h3>
                    <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.oldName}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = value.length == 0 || this.state.newName.length == 0
                                this.setState({
                                    oldName: value,
                                    updateDisabled: disabled
                                })
                            }}/>
                    </div>

                    <h3>New Equipment Name</h3>
                    <div className={"form-group"}>
                        <input className="form-control"
                            value={this.state.newName}
                            onChange={(event) => {
                                let value = event.target.value;
                                let disabled = value.length == 0 || this.state.oldName.length == 0
                                this.setState({
                                    newName: value,
                                    updateDisabled: disabled
                                })
                            }}/>
                    </div>

                    <button className="btn btn-primary mr-2"
                            onClick={() => {
                                this.updateEquip()
                            }} disabled={this.state.updateDisabled}>
                                Update!
                    </button>
                </div>
            </div>
        )
    }
}