import React, { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';

import LinearProgress from '@material-ui/core/LinearProgress'
import Card from '@material-ui/core/Card'

import CardHeader from '@material-ui/core/CardHeader';
import CardMedia from '@material-ui/core/CardMedia';
import CardContent from '@material-ui/core/CardContent';
import CardActions from '@material-ui/core/CardActions';
import Collapse from '@material-ui/core/Collapse';
import Avatar from '@material-ui/core/Avatar';

import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import BackupIcon from '@material-ui/icons/Backup';
import DoneIcon from '@material-ui/icons/Done';

import { makeStyles } from '@material-ui/core/styles';

import { green, red } from '@material-ui/core/colors';

const styles = {
    cards: {
        marginTop: 5,
    }
}

const useStyles = makeStyles(theme => ({
    error: {
      color: theme.palette.getContrastText(red[500]),
      backgroundColor: red[500],
    },
    green: {
      color: '#fff',
      backgroundColor: green[500],
    },
  }));
  
  

export default (props) => {
    const classes = useStyles();
    const [messages, setMessages] = useState([]);
    
    const upload = (file) => {
        const formData = new FormData();

        formData.append('title', '');
        formData.append('document', file);

        const id = file.path + messages.length;

        setMessages((state) => [{
            id: id,
            status: 'loading',
            filename: file.path,
        }].concat(state));

        fetch('/api/push', {
            method: 'POST',
            body: formData,
        }).then((result) => {
            if (!result.ok) {
                throw new Error('Upload failed');
            }

            return result.json()
        }).then((result) => {
            setMessages(state => {
                const index = state.findIndex(v => v.id == id);
                console.log(state, id, index);

                return state.slice(0, index).concat([{
                    ...state[index],
                    status: 'done',
                }]).concat(state.slice(index+1));
            });

        }).catch((error) => {
            setMessages(state => {
                const index = state.findIndex(v => v.id == id);
                console.log(state, id, index);
                return state.slice(0, index).concat([{
                    ...state[index],
                    status: 'error',
                }]).concat(state.slice(index+1));
            });

        })
    }

    const onDrop = useCallback(acceptedFiles => {
        acceptedFiles.map(upload);
    })

    const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop })

    const renderFile = (file) => (
        <Card key={file.id} variant="outlined" style={styles.cards}>
            <CardHeader
                avatar={
                    <Avatar className={file.status == 'error' ? classes.error : file.status == 'done' ? classes.green : null}>
                        {
                            file.status == 'loading' ? 
                                <BackupIcon/> :
                                file.status == 'done' ? <DoneIcon/> :
                                <ErrorOutlineIcon/>
                        }
                    </Avatar>
                }
                title={file.filename}
                subheader={file.status == 'loading' ? <LinearProgress/> : 
                        file.status == 'done' ? 'File was uploaded successfuly.' : 'Error while uploading file.'}
                />
        </Card>
    );

    return (
        <div>
            <div {...getRootProps()} style={{
                border: '3px dotted grey',
                textAlign: 'center',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                height: props.height ?? 200,
            }}>
                <input {...getInputProps()} />
                {
                    isDragActive ?
                        <p>Drop the files here ...</p> :
                        <p>Drag 'n' drop some files here, or click to select files</p>
                }
            </div>
            <div style={{ marginTop: 5 }}>
                {messages.map(renderFile)}
            </div>
        </div>
    )
};