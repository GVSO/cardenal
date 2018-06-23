import React, { Component } from 'react';

import Navbar from './components/common/Navbar';
import ResultList from './components/alumni-search/ResultList';

import '../sass/normalize.scss';

class App extends Component {
  render() {
    return (
      <div className="App">
        <Navbar />
        <ResultList />
      </div>
    );
  }
}

export default App;
