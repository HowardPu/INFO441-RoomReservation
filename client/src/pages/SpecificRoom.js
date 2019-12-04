import React, { Component } from 'react';
import {Redirect, Link} from 'react-router-dom';
import {AuthButton} from '../components/AuthButton'

// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button
export class SpecificRoom extends Component {
    constructor(props) {
        super(props)
        this.state = {
            back: false
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
        return(
            <div>
                <div>
                    <h1>Specific Room Manager</h1>
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