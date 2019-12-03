import React from 'react';
import AddRoomForm from '../components/AddRoomForm'
import DeleteRoomForm from '../components/DeleteRoomForm'

class Admin extends React.Component {
  // add equipment
  // update equipment
  // del equipment
  

  
    render() {
      return (
          <div>
                <h1>Administrator Board</h1>     
                <hr /> 
                <AddRoomForm />
                <hr />
                <DeleteRoomForm />


                /
          </div>
      );
    }
}

export default Admin;