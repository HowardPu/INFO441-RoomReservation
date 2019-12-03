import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import {Redirect} from 'react-router-dom';

const host = "https://api.awesome-summary.me" //!!change it later
const signinURL = host + "/v1/sessions"
const jsonHeader =  {'Content-Type': 'application/json'}

class Signin extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            email: '',
            password: '',
            errMes: '',
            adminSignedIn: false,
            normalSignedIn: false
        }
    }

    emailOnChange(e) {
        console.log(e.target.value)
        this.setState({email: e.target.value})
    }

    passOnChange(e) {
        console.log(e.target.value)
        this.setState({password: e.target.value})
    }

    onSubmit(e){
        e.preventDefault();
        this.setState({errMes: ""})
        if (!this.state.email) {
            this.setState({errMes: "Please input email"})
        } else if (!this.state.password) {
            this.setState({errMes: "Please input password"})
        } else {
            let userInput = {
                email: this.state.email,
                password: this.state.password
            }
            this.checkSignin(userInput, jsonHeader);
        }
    }

    checkSignin(userInput, headerInput) {
        fetch(signinURL, {
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
                if (data.userType === "Normal") {
                    this.setState({normalSignedIn: true})
                } else {
                    this.setState({adminSignedIn: true})
                }
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }

    render() {
        if (this.state.adminSignedIn === true) {
            return (<Redirect to="/admin" />)
        } else if (this.state.normalSignedIn === true) {
            return (<Redirect to="/user" />)
        }
        return (
            <div>
                <h1>Sign In</h1>
                <Form>
                    <Form.Group controlId="formBasicEmail">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control 
                            value={this.state.email}
                            onChange={(e) => {this.emailOnChange(e)}}
                            type="email" 
                            placeholder="Enter email" />
                        <Form.Text className="text-muted">
                        We'll never share your email with anyone else.
                        </Form.Text>
                    </Form.Group>

                    <Form.Group controlId="formBasicPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control 
                            value={this.state.password} 
                            onChange={(e) => {this.passOnChange(e)}}
                            type="password" 
                            placeholder="Password" />
                    </Form.Group>
                    <Button variant="primary" type="submit" onClick={(e) => {this.onSubmit(e)}}>
                        Sign In
                    </Button>
                </Form>
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
            </div>
        );
    }
}

export default Signin