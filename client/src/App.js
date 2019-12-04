import React from 'react';
import { BrowserRouter as Router, Switch, Route, Redirect} from 'react-router-dom'
import User from './pages/UserBoard'
import Admin from './pages/AdminBoard'
import Signin from './pages/SigninPage'
import Reserve from './pages/ReserveRoom'
import Signup from './pages/SignupPage'
import {Room} from './pages/Room'
import {SpecificRoom} from './pages/SpecificRoom'
import {Equipment} from './pages/Equipment'
import {Issues} from './pages/Issues'
import './App.css';


const host = "https://api.html-summary.me" 
const signinURL = host + "/v1/sessions"
const signoutURL = host + "/v1/sessions/mine"
const signupURL = host + "/v1/users"

class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
        authToken: "",
        userType: "",
        userName: "",
        error: null
    }

    this.handleSignIn = this.handleSignIn.bind(this)
    this.handleSignOut = this.handleSignOut.bind(this)
    this.handleSignUp = this.handleSignUp.bind(this)

  }

  componentDidMount() {
    
  }

  handleSignUp(email, password, userName) {
    fetch(signupURL, {
      method: 'POST',
      mode: "cors",
      headers: {'Content-Type': 'application/json'}, 
      body: JSON.stringify({
        "email": email,
        "password": password,
        "passwordConf": password,
        "userName": userName
      })
    }).then(resp => {
        if (resp.ok) {
            this.setState({
              authToken: resp.headers.get("Authorization")
            })
            return resp.json();
        } else {
            throw new Error(resp.status)
        }
    }).then(data => {
        this.setState({
          userType: data.userType,
          userName: data.userName
        })
    }).catch(err => {
        console.log(err)
        this.setState({
          "error": "Log In Failed"
        })
    })
  }

  handleSignOut() {
    fetch(signoutURL, {
        method: 'DELETE',
        mode: "cors",
        headers: {'Authorization': this.state.authToken}
    }).then(() => {
      this.setState({
        authToken: "",
        userType: "",
        userName: "",
        error: null
      })
    }).catch(err => {
        console.log(err)
    })
  }

  handleSignIn(email, password) {
    fetch(signinURL, {
        method: 'POST',
        mode: "cors",
        headers: {'Content-Type': 'application/json'}, 
        body: JSON.stringify({
          "email": email,
          "password": password
        })
    }).then(resp => {
        if (resp.ok) {
            this.setState({
              authToken: resp.headers.get("Authorization")
            })
            return resp.json();
        } else {
            throw new Error(resp.status)
        }
    }).then(data => {
      console.log(data)
      this.setState({
        userType: data.userType,
        userName: data.userName
      })
    }).catch(err => {
        console.log(err)
        this.setState({
          "error": "Log In Failed"
        })
    })
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
              <Route exact path='/signin'  render={(routerProps) => {
                  return <Signin {...routerProps} appState={this.state} signInHandler={this.handleSignIn} />
              }} />
              <Route exact path='/signup' render={(routerProps) => {
                  return <Signup {...routerProps} appState={this.state} signUpHandler={this.handleSignUp} />
              }} />
              <Route exact path='/admin' render={(routerProps) => {
                  return <Admin {...routerProps} appState={this.state} signOutHandler={this.handleSignOut} />
              }}/>

              <Route exact path='/equipment' render={(routerProps) => {
                  return <Equipment {...routerProps} appState={this.state} signOutHandler={this.handleSignOut} />
              }}/>

              <Route exact path='/room' render={(routerProps) => {
                  return <Room {...routerProps} appState={this.state} signOutHandler={this.handleSignOut} />
              }}/>

              <Route exact path='/specificRoom' render={(routerProps) => {
                  return <SpecificRoom {...routerProps} appState={this.state} signOutHandler={this.handleSignOut} />
              }}/>

              <Route exact path='/issues' render={(routerProps) => {
                  return <Issues {...routerProps} appState={this.state} signOutHandler={this.handleSignOut} />
              }}/>
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