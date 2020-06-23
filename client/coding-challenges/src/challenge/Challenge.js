import React from 'react';

import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Dropdown from 'react-bootstrap/Dropdown';
import Button from 'react-bootstrap/Button';
import TestCase from './TestCase';
import CodeEditor from './CodeEditor';

class Challenge extends React.Component {
  state = {
    tests: [],
    runtimes: [],
    selectedTest: {
      title: '',
      description: ''
    },
    selectedRuntime: {
      display: '',
      name: '',
      template: ''
    },
    submission: '',
    isSubmitting: false,
    submissionResults: []
  }

  componentDidMount() {
    // Get all tests
    fetch('http://localhost:8080/tests', {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      }
    }).then((res) => {
      return res.json()
    }).then((data) => {
      data.sort((a, b) => {
        return a.number - b.number
      })
      this.setState({ tests: data, selectedTest: data[0] })
    }).catch(console.log)

    // Get all the runtimes
    fetch('http://localhost:8080/runtimes', {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      }
    }).then((res) => {
      return res.json()
    }).then((data) => {
      data.sort((a, b) => {
        if (a.display < b.display) {
          return -1
        }

        if (a.display > b.display) {
          return 1
        }

        return 0
      })
      this.setState({ runtimes: data, selectedRuntime: data[0], submission: data[0].template })
    }).catch(console.log)
  }

  render() {
    const { tests, runtimes, selectedTest, selectedRuntime, submission, isSubmitting, submissionResults } = this.state;
    
    // Multiple Java runtimes, so map all to same Ace mode
    const mode = selectedRuntime.name.includes("java") ? "java" : selectedRuntime.name;

    // Called when a new test is selected
    const handleNav = (e, test) => {
      e.preventDefault();
      this.setState({ selectedTest: test, submissionResults: [] });
    }

    // Called when language dropdown is mutated
    const handleRuntimeChange = (e, runtime) => {
      e.preventDefault();
      this.setState({ selectedRuntime: runtime, submission: runtime.template });
    }

    // Called when code is edited
    const handleCodeInput = (value) => {
      this.setState({ submission: value })
    }

    // POST code to server and load results into the state
    const submitCode = (e) => {
      e.preventDefault();

      this.setState({ isSubmitting: true });

      fetch('http://localhost:8080/tests/' + selectedTest.number + '/submission', {
        method: 'POST',
        credentials: 'same-origin',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          runtime: selectedRuntime.name,
          code: submission.split("\n")
        })
      }).then((res) => {
        return res.json()
      }).then((data) => {
        let testCases = data.testcases;
        testCases.sort((a, b) => {
          return a.number - b.number
        })
        this.setState({ submissionResults: testCases, isSubmitting: false });
      }).catch(console.log);
    }

    return <Container fluid="true">
      <Row style={{minHeight: '100%'}}>
        <Col>
          {/* Select tests */}
          <Nav style={{paddingTop: '10px'}}>
            {tests.map((item, key) =>
              <Nav.Link key={key} onClick={(e) => handleNav(e, item)}>{item.title}</Nav.Link>
            )}
          </Nav>
        </Col>

        <Col xs={10}>
          <div style={{padding: '10px'}}>
            {/* Title and description */}
            <h1>{selectedTest.title}</h1>
            <div style={{paddingBottom: '10px'}} dangerouslySetInnerHTML={{ __html: selectedTest.description }} />

            {/* Select the language */}
            <Dropdown style={{paddingBottom: '10px', flexDirection: 'column'}}>
              <Dropdown.Toggle variant="warning" id="language-select">
                {selectedRuntime.display}
              </Dropdown.Toggle>
              <Dropdown.Menu>
                {runtimes.filter((item) => item.installed).map((item, key) =>
                  <Dropdown.Item key={key} onClick={(e) => handleRuntimeChange(e, item)}>{item.display}</Dropdown.Item>
                )}
              </Dropdown.Menu>
            </Dropdown>
            
            {/* Coding window */}
            <CodeEditor handleCodeInput={handleCodeInput} submission={submission} mode={mode} />
            <Button style={{marginBottom: '10px'}} variant="warning" onClick={submitCode} disabled={isSubmitting}>Submit Code</Button>
            
            {/* Submission results */}
            {submissionResults.map((testCase, key) =>
              <TestCase key={key} testCase={testCase} />
            )}
          </div>
        </Col>
      </Row>
    </Container>
  }
}

export default Challenge;
