import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import {Redirect, Link} from 'react-router-dom';

const host = "https://api.html-summary.me/"
const signupURL = host + "/v1/users"
const jsonHeader =  {'Content-Type': 'application/json'}

export class Signup extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            email: '',
            password: '',
            passwordReenter: '',
            username: '',
            disabled: true
        };
    }

    // set the state with corresponding event field and value
    // and then detect whether disable the `confirm` button or not
    handleChange(event) {
        let value = event.target.value;
        let fieid = event.target.name;
        let change = this.state;
        change[fieid] = value;
        change.errorMessage = null;
        change.disabled = this.getEmailStatus() !== " alert alert-success" || 
                            this.getPasswordStatus() !== " alert alert-success" ||
                            this.getReenterPasswordStatus() !== " alert alert-success" ||
                            this.getUserNameStatus() !== " alert alert-success";
        this.setState(change);
    }

    // when the email is not well formatted, it will give email input box warning color
    // when the email is well formatted, it will give emaol input box success color
    // when there is no input, it will give email input box no notification color 
    getEmailStatus() {
        if(this.state.email.length > 0) {
            let split = this.state.email.split("@");
            if(split.length === 2) {
                let last = split[1].split(".");
                if(last.length >= 2) {
                    let suffix = last[last.length - 1];
                    if(suffix === "com" || suffix === "edu" || suffix === "gov" || suffix === "org") {
                        return(" alert alert-success");
                    }
                }
            }
            return(" alert alert-warning");
        } else {
            return("");
        } 
    }

    // give the password input box danger color if the length is less than 6
    // give the password input box success color if the length is greater than 6
    // and give the password input box no extra color if there is no imput
    getPasswordStatus() {
        if(this.state.password.length > 0) {
            if(this.state.password.length < 6) {
                return(" alert alert-warning");
            } else {
                return(" alert alert-success");
            }
        } else {
            return("");
        }
    }

    // this method will check whether the reenter password is the same as the input password
    // if the same, give input box success color
    // if not, give input box danger color
    // if no input, give it no extra color
    getReenterPasswordStatus() {
        if(this.state.passwordReenter.length > 0) {
            if(this.state.password === this.state.passwordReenter && this.state.password.length >= 6) {
                return(" alert alert-success");
            } else {
                return(" alert alert-warning");
            }
        } else {
            return("");
        }
    }

    // if the user enters something, give username input box success color
    // or no extra color otherwise
    getUserNameStatus() {
        if(this.state.username.length > 0) {
            return(" alert alert-success");
        } else {
            return("");
        }
    }

    render() {
        let userType = this.props.appState.userType
        if (userType === "Admin") {
            return <Redirect to='/admin'/>;
        }
        if (userType === "Normal") {
            return <Redirect to='/user'/>;
        }
        if(this.state.redirect) {
            return <Redirect to='/signin'/>
        }
        return(
            <div className="sign-up-container">
                <div className="hypnotize sign-up-content" >
                    <div className="sign-up-item">
                        <h1>Your reservation begins here.</h1>
                    </div>
                    {this.state.errorMessage &&
                        <div><p className="alert alert-warning">{this.state.errorMessage}</p></div>
                    }
                    <div className="sign-up-item"> 
                        <h2>Please enter your email</h2>
                        <div className={"form-group" + this.getEmailStatus()}>
                            <input className="form-control"
                                name="email"
                                value={this.state.email}
                                onChange={(event) => this.handleChange(event)}
                                placeholder="Example: JohnSmith@gmail.com"/>
                        </div>
                    </div>

                    <div className="sign-up-item"> 
                        <h2>Please enter your password</h2>
                        <div className={"form-group" + this.getPasswordStatus()}>
                            <input type="password" className="form-control"
                                name="password"
                                value={this.state.password}
                                onChange={(event) => this.handleChange(event)}
                                placeholder="Longer than 6 characters"/>
                        </div>
                    </div>
                    <div className="sign-up-item"> 
                        <h2>Please reenter your password</h2>
                        <div className={"form-group" + this.getReenterPasswordStatus()}>
                            <input type="password" className="form-control"
                                name="passwordReenter"
                                value={this.state.passwordReenter}
                                onChange={(event) => this.handleChange(event)}/>
                        </div>
                    </div>

                    <div className="sign-up-item"> 
                        <h2>Please enter your username</h2>
                        <div className={"form-group" + this.getUserNameStatus()}>
                            <input className="form-control"
                                name="username"
                                value={this.state.username}
                                onChange={(event) => this.handleChange(event)}
                                placeholder="Example: uwlaziestperson1"/>
                        </div>
                    </div>
                    <div className="sign-up-item button-container">
                        <div>
                            <button className="btn btn-danger mr-2" 
                                    onClick={() => this.setState({redirect: true})}>
                                    Cancel
                            </button>
                        </div>

                        <div> 
                            <button className="btn btn-primary mr-2" disabled={this.state.disabled} onClick={() => {
                                this.props.signUpHandler(this.state.email, this.state.password, this.state.username)}}>
                                    Sign Up
                            </button>
                        </div>
                    </div>

                    <div className="sign-up-item"> 
                        <Link to="/sign-in">Have an account? One step to your expedition!</Link>
                    </div> 
                </div>
            </div>
        );
    }
}
export default Signup;