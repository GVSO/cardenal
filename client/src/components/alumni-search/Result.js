import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import Avatar from '@material-ui/core/Avatar';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';

const styles = {
  root: {
    flexGrow: 1,
  },
  paper: {
    height: 100
  },
  avatar: {
    margin: '1em .8em .5em 1em',
  },
  bigAvatar: {
    width: 60,
    height: 60,
  },
};

const Result = (props) => {
  const { classes } = props;

  return (
    <Paper className={classes.paper} elevation={2}>
      <Grid container>
        <Grid item sm={2} lg={1}>
          <Avatar
            alt="Adelle Charles"
            src="https://dummyimage.com/200x200/000000/fff"
            className={classNames(classes.avatar, classes.bigAvatar)}
          />

        </Grid>
        <Grid item>
          <Typography variant="headline">
            Text
          </Typography>
        </Grid>
      </Grid>
    </Paper>
  );
}

Result.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Result);