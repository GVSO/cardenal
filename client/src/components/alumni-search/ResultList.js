import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';

import Result from './Result';


const styles = {
  resultList: {
    marginTop: '1em',
  },
}

const ResultList = (props) => {
  const { classes } = props;

  return (
    <Grid className={classes.resultList} container justify="center">
      {[0, 1, 2].map(value => (
        <Grid item xs={12} md={8} key={value}>
          <Grid container>
            <Grid item xs={12}>
              <Result />
            </Grid>
          </Grid>
        </Grid>
      ))}
    </Grid>
  );
}

ResultList.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ResultList);