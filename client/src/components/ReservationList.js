import React from 'react';
import Button from 'react-bootstrap/Button'
import Table from 'react-bootstrap/Table'
import {Redirect} from 'react-router-dom';
import Alert from 'react-bootstrap/Alert'


const host = "http://api.html-summary.me"
const jsonHeader =  {
    'Authorization': localStorage.getItem('auth')
}
const resURL = host + "/v1/reserve"

const TestData = [{
    id: 3,
    tranDate:    "2019-11-13",
    reserveDate: "2019-11-13",
    roomName:    "MGH201",
    beginTime:   17,
    endTime:     20,
    roomType:    "Study",
},
{
    id:2,
    tranDate:    "2019-11-13",
    reserveDate: "2019-11-13",
    roomName:    "MGH201",
    beginTime:   17,
    endTime:     20,
    roomType:    "Study",
}]

class ReservationList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showRes: false,
            data: []
        }
    }

    // componentWillUnmount() {
    //     this.setState(
    //         {
    //             showRes: false,
    //             data: []
    //         }
    //     )
    // }

    onView(e) {
        e.preventDefault();
        // this.setState({showRes: true})

        var data = this.getData(resURL, jsonHeader)
    }

    getData(url, headerInput) {
        fetch(url, {
            method: 'GET',
            mode: "cors",
            headers: headerInput, 
        }).then(resp => {
            if (resp.ok) {
                return resp.json();
            } else {
                throw new Error(resp.status)
            }
        }).then(data => {
            this.setState({data})
            this.setState({showRes: true})
            console.log(data);
            return data;
        }).catch(err => {
            var errMes = "Oops something might be wrong! We will fix it soon!"
            console.log(err)
            this.setState({errMes});
            return null;
        })
    }

    decoder(num) {
        var hour;
        var min = "00";
        if (num%2 == 1) {
            hour = (num - 1)/2;
            min = "30"
        } else {
            hour = num / 2;
        }
        return hour + ":" + min
    }

    renderData() {
        console.log(this.state.data)
        return this.state.data.map((item, i) => {
            let roomName = item.roomName;
            let reserveDate = item.reserveDate;
            let beginTime = this.decoder(item.beginTime);
            let endTime = this.decoder(item.endTime);
            let roomType = item.roomType;
            let id = item.id;
            let date = reserveDate.split("-");
            let year = date[0]
            let month = date[1]
            let day  = date[2]
            let cur = new Date();
            var showBtn = false;
            if (cur.getFullYear() <= year && cur.getMonth() <= month && cur.getDay() < day) {
                showBtn = true;
            }
            return(
                <tr key={i}>
                    <td key={"name "+i}>{roomName}</td>
                    <td key={"type "+i}>{roomType}</td>
                    <td key={"begin "+i}>{beginTime}</td>
                    <td key={"end "+i}>{endTime}</td>
                    <td key={"datae "+i}>{reserveDate}</td>
                    {showBtn === true && 
                        <td key={"btn "+i}>
                            <Button value={id} onClick={(e) => {this.removeRes(e)}}>
                                Cancel Reservation
                            </Button>
                        </td>}
                </tr>
            )
        })
    }


    renderSuccessAlert() {
        return (
            <Alert variant="success" onClose={(e) => this.onCloseMes(e)} dismissible>
                <Alert.Heading>Successfully Reservered!</Alert.Heading>
                    <p>
                        You have successfullly remove your reservation at Room {this.state.data.roomName}
                    </p>
            </Alert>
        )
    }

    removeRes(e) {
        e.preventDefault()
        let id = e.target.value
        fetch(resURL+id, {
            method: 'DELETE',
            mode: "cors",
            headers: jsonHeader
        }).then(resp => {
            if (resp.status == 200) {
                console.log("reservation canceled")

                return;
            } else {
                throw Error(resp.status)
            }
        }).catch(err => {
            var errMes = "Oops something might be wrong! We will fix it soon!"
            console.log(err)
            this.setState({errMes});
        })
    }

    render() {
        return (
            <div>
                {this.state.errMes && <div className="errMes">{this.state.errMes}</div>}
                <Button 
                    variant="primary" 
                    onClick={(e)=>{this.onView(e)}}>
                        View My Reservations
                </Button>
                {this.state.showRes && this.state.data && 
                    <Table variant="dark">
                        <thead>
                            <tr>
                                <th scope="col">Room</th>
                                <th scope="col">Type</th>
                                <th scope="col">Begin Time</th>
                                <th scope="col">End Time</th>
                                <th scope="col">Date</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.renderData()}
                        </tbody>
                    </Table>
                }
            </div>
        );
    }
}

export default ReservationList