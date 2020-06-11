import React from 'react';

import Navigation from './components/Navigation';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Challenge from './challenge/Challenge';
import Runtime from './runtime/Runtime';

import './App.css';

class App extends React.Component {
  state = {
    currentViewName: "challenges"
  }

  render() {
    const currentView = this.state.currentViewName === "challenges" ? <Challenge /> : <Runtime />;

    const onViewChange = (view) => {
      this.setState({ currentViewName: view });
    }    

    return <Container fluid="true">
      <Row>
        <Col>
          <Navigation onViewChange={onViewChange} />
        </Col>
      </Row>
      <Row>
        <Col>
          {currentView}
        </Col>
      </Row>
    </Container>;
  }
}

export default App;
