import React from 'react';
import { BrowserRouter as Router, Switch, Route, Redirect} from 'react-router-dom'
// import Home from './pages/Home'
import User from './pages/UserBoard'
import Admin from './pages/AdminBoard'
import Signin from './pages/SigninPage'
import Signup from './pages/SignupPage'
import './App.css';

class App extends React.Component {


  render() {
    return (
      <div className="App">
        <header className="App-header"></header>

        <Router>
          <div>
            <Switch>
              {/* if currentUrl == '/home', render <Signin> */}
              {/* <Route path='/' component={Home} /> */}
              {/* if currentUrl == '/home', render <Signin> */}
              <Route exact path='/signin' component={Signin} />

              {/* if currentUrl == '/about', render <Signup> */}
              <Route exact path='/signup' component={Signup} />
              <Route exact path='/admin' component={Admin} />
              <Route exact path='/user' component={User} />
              <Redirect to="/signin" />
            </Switch>
          </div>
      </Router>
      </div>
    );
  }
}

export default App;