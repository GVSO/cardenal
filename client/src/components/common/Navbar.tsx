import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import { AppBar, Toolbar, Typography } from '@material-ui/core/';

const styles = theme => ({
  toolbar: {
    minHeight: 55,
    backgroundColor: theme.palette.error.indigo
  },
  title: {
    fontWeight: 500
  }
});

const Navbar = (props: Props) => {
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

export default withStyles(styles)(Navbar);

interface Props {
  classes: {
    toolbar: 'string'
  }
}
