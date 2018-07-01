import React from 'react'
import { withStyles } from '@material-ui/core/styles'
import { Grid } from '@material-ui/core'

import Result from './Result'

const styles = {
  resultList: {
    marginTop: '1em'
  }
}

const ResultList = (props: Props) => {
  const { classes } = props

  return (
    <Grid className={classes.resultList} container justify='center'>
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
  )
}

export default withStyles(styles)(ResultList)

/* eslint-disable */
interface Props {
  classes: {
    resultList: string
  }
}
/* eslint-enable */
