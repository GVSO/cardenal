import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';

const styles = theme => ({
  toolbar: {
    minHeight: 55,
    backgroundColor: theme.palette.error.indigo
  },
  title: {
    fontWeight: 500
  }
});

const Navbar = (props) => {
  const { classes } = props;

  return (
    <AppBar position="static">
      <Toolbar className={classes.toolbar}>
        <Typography variant="title" color="inherit">
          Cardenal
        </Typography>
      </Toolbar>
    </AppBar>
  );
};

Navbar.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Navbar);