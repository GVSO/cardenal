import React from 'react';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import { Avatar, Grid, Paper, Typography } from '@material-ui/core';

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

const Result = (props: Props) => {
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

export default withStyles(styles)(Result);

interface Props {
  classes: {
    paper: string,
    avatar: string,
    bigAvatar: string
  }
}
