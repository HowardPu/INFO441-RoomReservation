import React from 'react';
import { Link} from 'react-router-dom'


class Home extends React.Component {
    render(){
        return (
            <ul>
                <li>
                    <Link to="/signin">Sign In</Link>
                </li>
                <li>
                    <Link to="/signup">Sign up</Link>
                </li>
            </ul>
        )
    }
}

export default Home;