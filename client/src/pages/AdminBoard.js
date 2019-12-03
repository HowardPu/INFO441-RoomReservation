import React from 'react';
import AddRoomForm from '../components/AddRoomForm'
import DeleteRoomForm from '../components/DeleteRoomForm'

class Admin extends React.Component {

    render() {
      return (
          <div>
                <h1>Administrator Board</h1>        
                <AddRoomForm />
                <hr />
                <DeleteRoomForm />
          </div>
      );
  }

}

export default Admin;