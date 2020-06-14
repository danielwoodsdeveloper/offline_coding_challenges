import React from 'react';

import AceEditor from 'react-ace';

import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/mode-python';
import 'ace-builds/src-noconflict/theme-twilight';

class CodeEditor extends React.Component {
  render() {
    // Called when code is edited, pass back to Challenge component
    const handleCodeInput = (value) => {
      this.props.handleCodeInput(value);
    }

    return <AceEditor
      mode={this.props.mode}
      theme="twilight"
      value={this.props.submission}
      onChange={handleCodeInput}
      height='350px'
      width='100%'
      style={{marginBottom: '10px'}}
      highlightActiveLine={true}
      editorProps={{ $blockScrolling: true }}
      setOptions={{
        showLineNumbers: true,
        tabSize: 2,
      }}
    />
  }
}

export default CodeEditor;