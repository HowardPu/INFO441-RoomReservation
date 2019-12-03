import React from 'react';
import { BrowserRouter as Router, Switch, Route, Redirect} from 'react-router-dom'
import User from './pages/UserBoard'
import Admin from './pages/AdminBoard'
import Signin from './pages/SigninPage'
import Reserve from './pages/ReserveRoom'
import Signup from './pages/SignupPage'
import './App.css';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
        isloggedin: false
    }

}

  componentDidMount() {
    
  }

  // connect() {
  //   if (localStorage.getItem('auth')) {
  //     const websocket = new WebSocket("wss://" + "api.awesome-summary.me" + "/v1/ws?auth=" + localStorage.getItem('auth'));

  //   }

  // }


  render() {
    return (
      <div className="App">
        <header className="App-header"></header>

        <Router>
          <div>
            <Switch>
              <Route exact path='/signin' component={Signin} />
              <Route exact path='/signup' component={Signup} />
              <Route exact path='/admin' component={Admin} />
              <Route exact path='/user' component={User} />
              <Route exact path='/reserve' component={Reserve} />
              <Redirect to="/signin" />
            </Switch>
          </div>
      </Router>
      </div>
    );
  }
}

export default App;