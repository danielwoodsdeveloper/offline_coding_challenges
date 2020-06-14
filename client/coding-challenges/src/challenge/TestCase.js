import React from 'react';

import Card from 'react-bootstrap/Card';

class TestCase extends React.Component {
  render() {
    return <Card border={this.props.testCase.pass ? "success" : "danger"} style={{ marginTop: '10px', marginBottom: '10px' }}>
      <Card.Body>
        <Card.Title style={{ paddingBottom: '15px' }}><b>Test Case {this.props.testCase.number}</b></Card.Title>
        <p><b>Input:</b> [{this.props.testCase.inputs.map(str => "\"" + str + "\"").join(", ")}]</p>
        <p><b>Output:</b> [{this.props.testCase.outputs.map(str => "\"" + str + "\"").join(", ")}]</p>
        <p><b>Expected:</b> [{this.props.testCase.expected.map(str => "\"" + str + "\"").join(", ")}]</p>
      </Card.Body>
    </Card>
  }
}

export default TestCase;