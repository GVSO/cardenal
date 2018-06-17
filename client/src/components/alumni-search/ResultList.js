import React from 'react';
import { withStyles } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';

import Result from './Result';


const styles = {
  resultWrapper: {
    boxShadow: '10px 0 10px -4px rgba(0, 0, 0, 0.2), -10px 0 10px -4px rgba(0, 0, 0, 0.2)'
  }
}

const ResultList = (props) => {
  const { classes } = props;

  return (
    <Grid container justify="center">
      {[0, 1, 2].map(value => (
        <Grid item xs={12} md={8} key={value}>
          <Grid container className={classes.resultWrapper}>
            <Grid item xs={12}>
              <Result />
            </Grid>
          </Grid>
        </Grid>
      ))}
    </Grid>
  );
}

export default withStyles(styles)(ResultList);