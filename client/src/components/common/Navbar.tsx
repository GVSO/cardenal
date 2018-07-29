import React, { Component } from 'react'
import { NavLink } from 'react-router-dom'

import { withStyles, WithStyles } from '@material-ui/core/styles' // eslint-disable-line
import { AppBar, Button, Toolbar, Typography } from '@material-ui/core/'
import SearchBar from 'material-ui-search-bar'

import '../../../sass/components/navbar.scss'

const styles = theme => ({
  flex: {
    alignItems: 'center',
    display: 'flex',
    flex: 1
  },
  searchBar: {
    flexBasis: '50%',
    marginLeft: 50
  },
  title: {
    fontWeight: 500
  },
  toolbar: {
    backgroundColor: theme.palette.primary.indigo,
    minHeight: 55
  }
})

export class Navbar extends Component<Props & WithStyles, State> {
  componentWillMount (): void {
    this.setState({
      searchKeyword: ''
    })
  }

  render () {
    const { classes } = this.props

    return (
      <AppBar position='static'>
        <Toolbar classes={{root: 'toolbar'}} className={classes.toolbar}>
          <div className={classes.flex}>
            <NavLink exact to='/'>
              <Typography variant='title' color='inherit'>
                Cardenal
              </Typography>
            </NavLink>
            <SearchBar
              placeholder='Buscar'
              classes={{root: 'search-bar'}}
              className={classes.searchBar}
              value={this.state.searchKeyword}
              onChange={(keyword: string) => this.setState({ searchKeyword: keyword })}
              /* onRequestSearch={() => doSomethingWith(this.state.value)} */
            />
          </div>

          <a href='/api/services/login'><Button color='inherit'>Login</Button></a>
        </Toolbar>
      </AppBar>
    )
  }
}

export default withStyles(styles)<Props>(Navbar)

/* eslint-disable */
interface Props {
  classes?: {
    flex: string,
    searchBar: string,
    title: string,
    toolbar: string
  }
}

interface State {
  searchKeyword: string
}
/* eslint-enable */
