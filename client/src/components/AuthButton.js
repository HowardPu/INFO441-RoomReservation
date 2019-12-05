import React, { Component } from 'react';
import {Button} from 'react-bootstrap'
import {Link} from 'react-router-dom';

// This is the authentication button
// when a user is authenticated, it will display the log-off button and
// when a user is not authentciated, it will display the log-in button
export class AuthButton extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <Button variant="danger" className="btn btn-primary mr-2" onClick={() => {
                this.props.signOutHandler()
            }}>Logout</Button>
        )
    }
}