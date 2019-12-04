import React, { Component } from 'react';
import {Redirect, Link} from 'react-router-dom';
import {AuthButton} from '../components/AuthButton'

// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button

const getAllRoom = "https://api.html-summary.me/v1/room"

export class Room extends Component {
    constructor(props) {
        super(props)
        this.state = {
            back: false,
            rooms: null
        }
    }

    componentDidMount() {
        if(this.props.appState.userType == "Admin") {
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
        console.log(this.state.rooms)
        return(
            <div>
                <div>
                    <h1>Room Manager</h1>
                    <AuthButton signOutHandler={this.props.signOutHandler} />
                </div>   


                <button className="btn btn-primary mr-2" disabled={this.state.disabled} onClick={() => {
                    this.setState({
                        back: true
                    })
                }}>Back</button>
            </div>
        )
    }
}