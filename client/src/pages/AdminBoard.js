import React from 'react';
import AddRoomForm from '../components/AddRoomForm'
import DeleteRoomForm from '../components/DeleteRoomForm'
import {AuthButton} from '../components/AuthButton'
import {Redirect, Link} from 'react-router-dom';

class Admin extends React.Component {
  // add equipment
  // update equipment
  // del equipment
  constructor(props) {
    super(props)

    this.state = {
      "room": false,
      "specificRoom": false,
      "equipment": false,
      "issues": false
    }
  }
  

  
  render() {
    let userType = this.props.appState.userType
    if(userType != "Admin") {
      return <Redirect to='/signin'/>
    }
    if(this.state.room) {
      return <Redirect to='/room'/>
    }
    if(this.state.specificRoom) {
      return <Redirect to='/specificRoom'/>
    }
    if(this.state.equipment) {
      return <Redirect to='/equipment'/>
    }
    if(this.state.issues) {
      return <Redirect to='/issues'/>
    }
    return (
        <div>
          <div>
            <h1>Administrator Board</h1>
            <AuthButton signOutHandler={this.props.signOutHandler} />
          </div>    

          <hr /> 
          <div> 
              <button className="btn btn-primary mr-2" onClick={() => {
                this.setState({"room": true})
              }} >
                      Manage Room
              </button>
              <button className="btn btn-primary mr-2" onClick={() => {
                this.setState({
                  "specificRoom": true
                })}} >
                      Manage Specific Room
              </button>

              <button className="btn btn-primary mr-2" onClick={() => {
                this.setState({
                    "equipment": true
                })
              }} >
                      Manage Equipment
              </button>

              <button className="btn btn-primary mr-2" onClick={() => {
                this.setState({"issues": true})
              }} >
                      Update Issues
              </button>
          </div>
              {/* <hr /> 
              <AddRoomForm />
              <hr />
              <DeleteRoomForm /> */}
        </div>
    );
  }
}

export default Admin;