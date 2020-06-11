import React from 'react';
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'

class Navigation extends React.Component {
  state = {
    key: 'challenges'
  }

  render() {
    const onClick = (e, name) => {
      e.preventDefault();
      this.setState({ key: name });
      this.props.onViewChange(name);
    }

    return (
      <Navbar expand="lg" variant="dark" bg="primary">
        <Nav className="mr-auto" activeKey={this.state.key}>
          <Nav.Link onClick={(e) => onClick(e, "challenges")} eventKey="challenges">Challenges</Nav.Link>
          <Nav.Link onClick={(e) => onClick(e, "runtimes")} eventKey="runtimes">Runtimes</Nav.Link>
        </Nav>
      </Navbar>
    );
  }
}

export default Navigation;
