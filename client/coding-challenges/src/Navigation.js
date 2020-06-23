import React from 'react';

import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';

class Navigation extends React.Component {
  state = {
    key: 'challenges'
  }

  render() {
    // Called when a the view is changed, passed back to App component
    const onClick = (e, name) => {
      e.preventDefault();
      this.setState({ key: name });
      this.props.onViewChange(name);
    }

    return (
      <Navbar expand="lg" variant="light" bg="warning">
        <Nav className="mr-auto" activeKey={this.state.key} style={{flexDirection: 'row'}}>
          <Navbar.Brand href="#"><img src="/logo.svg" width="30" height="30" alt="Coding Challenges logo" /></Navbar.Brand>
          <Nav.Link onClick={(e) => onClick(e, "challenges")} eventKey="challenges" style={{paddingLeft: '10px', paddingRight: '10px'}}>Challenges</Nav.Link>
          <Nav.Link onClick={(e) => onClick(e, "runtimes")} eventKey="runtimes" style={{paddingLeft: '10px', paddingRight: '10px'}}>Runtimes</Nav.Link>
        </Nav>
      </Navbar>
    );
  }
}

export default Navigation;
