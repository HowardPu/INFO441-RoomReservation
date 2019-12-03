import React from 'react';
import { BrowserRouter as Router, Switch, Route, Redirect} from 'react-router-dom'

const getRoomURL = host + "v1/room"

class AdminRoomList extends React.Component {
    constructor(props) {
    }

    render() {
        return (
            <Route {...rest} render={props => (
                localStorage.getItem('auth')
                    ? <Component {...props} />
                    : <Redirect to={{ pathname: '/login', state: { from: props.location } }} />
            )} />
        )
    }
}