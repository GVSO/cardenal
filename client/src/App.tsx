import React, { Component } from 'react'
import { Route, Switch } from 'react-router-dom'

import Navbar from './components/common/Navbar'
import ResultList from './components/alumni-search/ResultList'

import '../sass/normalize.scss'

class App extends Component {
  render () {
    return (
      <React.Fragment>
        <Navbar />

        <Switch>
          <Route exact path='/' component={ResultList} />
          <Route exact path='/login' component={ResultList} />
        </Switch>
      </React.Fragment>
    )
  }
}

export default App
