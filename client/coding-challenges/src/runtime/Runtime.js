import React from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';

class Runtime extends React.Component {
  state = {
    runtimes: [
      {
        display: '',
        image: '',
        installed: false,
        isInstalling: false
      }
    ]
  }

  componentDidMount() {
    // Get all the runtimes
    fetch('http://localhost:8080/runtimes', {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      }
    })
    .then((res) => {
      return res.json()
    })
    .then((data) => {
      data.sort((a, b) => {
        if (a.display < b.display) {
          return -1
        }

        if (a.display > b.display) {
          return 1
        }

        return 0
      }).map(item => item.isInstalling = false)
      this.setState({ runtimes: data })
    })
    .catch(console.log)
  }

  render() {
    const { runtimes } = this.state;

    const install = (e, runtime) => {
      e.preventDefault();

      let runtimes = this.state.runtimes;

      let ind = runtimes.indexOf(runtime);
      runtimes[ind].isInstalling = true;

      this.setState({ runtimes: runtimes });

      fetch('http://localhost:8080/runtimes/' + runtime.name + '/install', {
        method: 'POST',
        credentials: 'same-origin',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({})
      })
      .then((res) => {
        return res.json()
      })
      .then((data) => {
        let runtimes = this.state.runtimes;

        let ind = runtimes.indexOf(runtime);
        runtimes[ind].installed = data.installed;

        this.setState({ runtimes: runtimes });
      })
      .catch(console.log);
    }

    return <Container>
      <Row>
        <Col>
          {runtimes.map((item, key) =>
            <Card key={key} style={{margin: '10px'}}>
              <Card.Body>
                <Card.Title>{item.display}</Card.Title>
                <Card.Text><i>{item.image}</i></Card.Text>
                <Card.Text>{item.installed ? "" : "Not "}Installed</Card.Text>
                <Button variant="primary" disabled={item.installed || item.isInstalling} onClick={(e) => install(e, item)}>Install</Button>
              </Card.Body>
            </Card>
          )}
        </Col>
      </Row>
    </Container>
  }
}

export default Runtime;
