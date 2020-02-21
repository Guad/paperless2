import React from 'react';
import { Labeled, ReferenceInput } from 'react-admin';
import { useField } from 'react-final-form';
import Chip from '@material-ui/core/Chip';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core/styles';
import { AutocompleteInput } from 'react-admin';

const useStyles = makeStyles(theme => ({
    root: {
      display: 'flex',
      justifyContent: 'center',
      flexWrap: 'wrap',
      padding: theme.spacing(0.5),
    },
    chip: {
      margin: theme.spacing(0.5),
    },
}));

const TagInput = ({name, label}) => {
    const classes = useStyles();

    const {
        input: { onChange, value },
        meta: { toucher, error }
    } = useField(name);

    const handleDelete = tag => () => {
        onChange(value.filter(t => t !== tag));
    }

    const newTagChange = (tag) => {
        let v = value;

        if (!v || !v.length) {
            v = [];
        }

        onChange(v.filter(t => t !== tag).concat([tag]));
    }

    return (
        <div>
        <Labeled label={label ?? "Tags"}>
            <Paper className={classes.root}>
                {
                    value && value.length ?
                        value.map(t => <Chip
                            key={t}
                            label={t}
                            onDelete={handleDelete(t)}
                            className={classes.chip}
                        />) : <span>No tags</span>
                }                
            </Paper>
        </Labeled>
        <ReferenceInput label="Add Tag" 
            reference="tag"
            name="tag"
            onChange={newTagChange}
            error={false}
            helperText={false}>
            <AutocompleteInput optionText="name" optionValue="name" />
        </ReferenceInput>
        </div>
    );
}

export default TagInput;