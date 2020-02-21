import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Button, IconButton } from '@material-ui/core';
import GetAppIcon from '@material-ui/icons/GetApp';
import VisibilityIcon from '@material-ui/icons/Visibility';

const useStyles = makeStyles({
    button: {
        marginRight: 5,
    },
    icon: {

    }
})

const DownloadField = ({ record = {}, source }) => {
    const classes = useStyles();
    return (
        <span>
            <Button className={classes.button} variant="contained" color="primary" component="a" href={`/api/fetch/${record[source]}`} target="_blank" startIcon={<VisibilityIcon className={classes.icon} />}>
                View
            </Button>
            <Button className={classes.button} variant="contained" color="primary" component="a" href={`/api/fetch/${record[source]}`} target="_blank" startIcon={<GetAppIcon className={classes.icon} />} download>
                Download
            </Button>
        </span>
    );
}

export default DownloadField;