import React from 'react';
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Alert from 'react-bootstrap/Alert'
import DatePicker from 'react-datepicker';
import "react-datepicker/dist/react-datepicker.css";

const host = ""
const jsonHeader =  {
    'Content-Type': 'application/json',
    'Authorization': localStorage.getItem('auth')
}
const addRoomURL = host + "v1/room"

class ReservationForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            startDate: new Date()
        }
    }

    setStartDate(date) {
        console.log(date)
    }
    renderDatePicker() {
        var curDate = new Date()
        console.log(curDate)
        let minCode = curDate.setMinutes(0);
        console.log(minCode)
        let minTime = curDate.setHours(8, 0);
        let maxTime = curDate.setHours(20, 0);

        return (
            <DatePicker
                selected={this.state.startDate}
                onChange={date => this.setStartDate(date)}
                showTimeSelect
                minDate={curDate}
                minTime={minTime}
                maxTime={maxTime}
                dateFormat="MMMM d, yyyy h:mm aa"
            />
        );
    }

    render() {
        return(
            <Form>
                <Form.Group>
                    <Form.Text>Reservation Starting Time</Form.Text>
                    {this.renderDatePicker()}
                </Form.Group>
                <Button variant="primary" type="submit" onClick={(e) => {this.onSubmit(e)}}>
                    Reserve
                </Button>
            </Form>
        );
    }
}

export default ReservationForm;