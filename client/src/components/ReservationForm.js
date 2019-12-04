import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Alert from 'react-bootstrap/Alert'
import DatePicker from 'react-datepicker';
import "react-datepicker/dist/react-datepicker.css";

const host = "https://api.html-summary.me/"
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}
const latest = 42
const reserveURL = host + "v1/reserve"
const getUsedTimeURL = host + "v1/roomUsedTime"
const availableTime = [];
for (var i = 16; i <= 42; i++) {
    availableTime.push(i);
}

class ReservationForm extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            errMes: '',
            startedDate: null,
            startedTime: null,
            duration: 0.5,
            showSuccessMes: false,
            requestInfo: {}
            // excludeTimes: []
        }
    }

    // componentDidUpdate(){
    //     this.renderTimePicker();
    // }

    setStartDate(date) {
        this.setState({startedDate: date})
    }

    renderDatePicker() {
        let curDate = new Date()
        return (
            <DatePicker
                selected={this.state.startedDate}
                onChange={date => this.setStartDate(date)}
                minDate={curDate}
                dateFormat="MMMM d, yyyy"
            />
        );
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        console.log(this.props.newRes)
        if (this.props.newRes !== prevProps.newRes) {
            return true;
        }
        // if (this.props.newRes !== prevProps.newRes) {
        //     // this.renderTimePicker();
        //     this.manipulateTimeData()
        // }
    }


    setStartTime(time) {
        if (!this.state.startedDate) {
            this.setState({errMes: "Please Input the date first "})
        } else {
            this.setState({startedTime: time})
        }
    }
    

    renderTimePicker() {
        console.log("render!")
        let curDate = new Date()
        var excludeTimes = []
        if (this.state.startedDate){
            var minTime = new Date()
            minTime.setHours(8);
            minTime.setMinutes(0);
            if (curDate.getDate() === this.state.startedDate.getDate()
                && curDate.getMonth() === this.state.startedDate.getMonth()
                && curDate.getFullYear() === this.state.startedDate.getFullYear()
                && curDate.getTime() > minTime) {
                minTime = curDate.getTime();
            }
            var maxTime = new Date()
            maxTime.setHours(20);
            maxTime.setMinutes(30);
            excludeTimes = this.manipulateTimeData();
        }

        return (
            <DatePicker
                selected={this.state.startedTime}
                onChange={date => this.setStartTime(date)}
                showTimeSelect
                showTimeSelectOnly
                minDate={curDate}
                minTime={minTime}
                maxTime={maxTime}
                excludeTimes={excludeTimes}
                dateFormat="h:mm aa"
            />
        );
    }

    manipulateTimeData() {
        let month =  this.state.startedDate.getMonth() + 1
        var getURL = getUsedTimeURL + 
            "?roomname=" + this.props.roomName + "&" +
            "year=" + this.state.startedDate.getFullYear() + "&" +
            "month=" + month + "&" +
            "day=" + this.state.startedDate.getDate();
        var timeslots;
        var excludeTimes = [];
        let newSlot = this.props.newRes;
        console.log(newSlot)
        this.getData(getURL).then(data => {
            console.log(data)
            timeslots = data.result;
            if (newSlot) {
                for (let t = newSlot.begin; t <= newSlot.begin + newSlot.duration; t++) {
                    timeslots.add(t);
                }
            }
            
            if (timeslots) {
                timeslots.forEach(receivedTime => {
                    var time = receivedTime - 1;
                    availableTime.splice(availableTime.indexOf(time), 1);
                    for (let j = time; j > time - this.state.duration*2; j--) {
                        let exclude = this.dateGenerate(j)
                        excludeTimes.push(exclude);
                    }
                }); 
            }
        
            for (let l = latest; l >= latest - this.state.duration*2; l--) {
                let exclude = this.dateGenerate(l)
                excludeTimes.push(exclude);
            }
            return excludeTimes;
            // this.setState({excludeTimes})
        }).catch(err => {
            console.log(err)
        })            
        return excludeTimes
    }

    dateGenerate(time) {
        var min = 0;
        var hour = time / 2;
        if (time % 2 === 1) {
            min = 30;
            hour = (time - 1)/2;
        }
        var exclude = new Date();
        exclude.setHours(hour, min);
        return exclude;
    }

    getData(url) {
        console.log(url)
        return fetch(url, {
            method: 'GET',
            mode: "cors",
            headers: {'Authorization': localStorage.getItem('auth')}, 
        }).then(resp => {
           return resp.json()
        })
    }


    onSelectDuration(e) {
        e.preventDefault();
        this.setState({duration: e.target.value})
    }

    onReserve(e) {
        e.preventDefault();
        if (!this.state.startedDate) {
            this.setState({errMes: "Please Enter the Reservation Starting Date"})
        } else if (!this.state.startedTime) {
            this.setState({errMes: "Please Enter the Reservation Starting Time"})
        } else {
            let hours = this.state.startedTime.getHours();
            let min = this.state.startedTime.getMinutes() === 30 ? 1 : 0;
            console.log(hours + " "  + min)
            let beginTime = hours * 2 + min
            let request = {
                year: this.state.startedDate.getFullYear(),
                month: this.state.startedDate.getMonth() + 1,
                day: this.state.startedDate.getDate(),
                roomName: this.props.roomName,
                beginTime: beginTime,
                duration: this.state.duration * 2
            }

            this.setState({requestInfo: request})
            console.log(request)
            this.postData(reserveURL, request, jsonHeader)
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
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            console.log(data);
            this.setState({showSuccessMes: true})
        }).catch(err => {
            var errMes = err.message
            console.log(err)
            this.setState({errMes});
        })
    }    

    renderSuccessAlert() {
        let endTime = this.state.requestInfo.beginTime + this.state.requestInfo.duration
        let endDate = this.dateGenerate(endTime);
        console.log(this.state.requestInfo.beginTime)
        let endHour = endDate.getHours();
        var endMin = endDate.getMinutes() === 0 ? "00" : endDate.getMinutes();

        let startHour = this.state.startedTime.getHours();
        let startMin = this.state.startedTime.getMinutes() === 0 ? "00" : this.state.startedTime.getMinutes();
        let dateString = `${this.state.requestInfo.month}/${this.state.requestInfo.day}/${this.state.requestInfo.year}`;
        let timeString = `${startHour}:${startMin} to ${endHour}:${endMin}`;
        return (
            <Alert variant="success" onClose={(e) => this.onCloseMes(e)} dismissible>
                <Alert.Heading>Successfully Reservered!</Alert.Heading>
                    <p>
                        You have successfullly reserved Room {this.props.roomName} on {dateString} from {timeString}
                    
                    </p>
            </Alert>
        )
    }

    onCloseMes() {
        this.setState({
            showSuccessMes: false,
            startedDate: null,
            startedTime: null,
            duration: 0.5,
            requestInfo: {}
        })
    }

    render() {
        return(
            <Form>
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                <Form.Group controlId="durationSelector">
                    <Form.Label>How many hours would you like to reserve?</Form.Label>
                    <Form.Control 
                        as="select"
                        onChange={(e)=>this.onSelectDuration(e)}>
                        <option>0.5</option>
                        <option>1</option>
                        <option>1.5</option>
                        <option>2</option>
                        <option>2.5</option>
                        <option>3</option>
                        <option>3.5</option>
                        <option>4</option>
                        <option>4.5</option>
                        <option>5</option>
                    </Form.Control>
                </Form.Group>
                <Form.Group>
                    <Form.Label>Reservation Starting Date</Form.Label>
                    <div>{this.renderDatePicker()}</div>
                    <br />
                    <Form.Label>Reservation Starting Time</Form.Label>
                    <div>{this.renderTimePicker()}</div>
                </Form.Group>
                <Button 
                    variant="primary" 
                    type="submit" 
                    onClick={(e) => {this.onReserve(e)}}>
                    Reserve
                </Button>
                {this.state.showSuccessMes === true &&
                   this.renderSuccessAlert()
                }
            </Form>
        );
    }
}

export default ReservationForm;