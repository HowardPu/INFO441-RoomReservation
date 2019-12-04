import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import {Redirect, Link} from 'react-router-dom';

const host = "http://api.html-summary.me/"
const signupURL = host + "/v1/users"
const jsonHeader =  {'Content-Type': 'application/json'}

class Signup extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            errMes: '',
            email: '',
            password: '',
            username: '',
            passwordConf: '',
            signedup: false,
        }
    }

    onChange(e) {
        this.setState(
            {
                notification: '',
                errMes: ''
            }
        )
        switch (e.target.id) {
            case "signupEmail":
                this.setState({email: e.target.value});
                break;
            case "signupUsername":
                this.setState({username: e.target.value});
                break;
            case "signupPassword":
                this.setState({password: e.target.value});
                break;
            case "signupPasswordConf":
                    this.setState({passwordConf: e.target.value});
                    break;
            default:
                break;
        }
    }
    
    onSubmit(e){
        e.preventDefault();
        if (!this.state.username) {
            this.setState({errMes: "Please enter username"})
        } else if (!this.state.email) {
            this.setState({errMes: "Please enter email"})
        } else if (!this.state.password) {
            this.setState({errMes: "Please enter password"})
        }  else if (!this.state.passwordConf) {
            this.setState({errMes: "Please type in you password again to confirm password"})
        }  else if (this.state.passwordConf !== this.state.password) {
            this.setState({errMes: "Password doesn't match with the confirmation!"})
        } else {
            let userInput = {
                email: this.state.email,
                password: this.state.password,
                passwordConf: this.state.passwordConf,
                userName: this.state.username
            }
            this.postData(signupURL, userInput, jsonHeader);
        }
    }

    postData(url, userInput, headerInput) {
        fetch(url, {
            method: 'POST',
            mode: "cors",
            headers: headerInput, 
            body: JSON.stringify(userInput)
        }).then(resp => {
            if (resp.ok) {
                if (!headerInput.Authorization && resp.headers.get('Authorization')) {
                    localStorage.setItem('auth', resp.headers.get('Authorization'));
                }
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            console.log(data);
            this.setState({signedup: true})
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }

    render() {
        if (this.state.signedup === true) {
            return (<Redirect to="/user" />)
        }
        return (
            <div>
                <h1>Sign up</h1>
                <br />
                <Form>
                    {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                    <Form.Group controlId="signupEmail">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control 
                            type="email" 
                            value={this.state.email}
                            onChange={(e) => {this.onChange(e)}}
                            placeholder="Enter email" />
                        <Form.Text className="text-muted">
                        We'll never share your email with anyone else.
                        </Form.Text>
                    </Form.Group>

                    <Form.Group controlId="signupUsername">
                        <Form.Label>Username</Form.Label>
                        <Form.Control 
                            value={this.state.username}
                            onChange={(e) => {this.onChange(e)}}
                            placeholder="Username" />
                    </Form.Group>

                    <Form.Group controlId="signupPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control 
                            value={this.state.password}
                            onChange={(e) => {this.onChange(e)}}
                            type="password" 
                            placeholder="Password" />
                    </Form.Group>

                    <Form.Group controlId="signupPasswordConf">
                        <Form.Label>Confirm your Password</Form.Label>
                        <Form.Control 
                            type="password" 
                            value={this.state.passwordConf}
                            onChange={(e) => {this.onChange(e)}}
                            placeholder="Password Confirmation" />
                    </Form.Group>

                    <Button variant="primary" type="submit" onClick={(e) => this.onSubmit(e)}>
                        Submit
                    </Button>
                </Form>
                <div>Already have an account? <Link to="/signin">Sign in</Link></div>
            </div>
        );
    }
}

export default Signup;