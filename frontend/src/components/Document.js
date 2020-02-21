import React from 'react';
import { List, Datagrid, ImageField, DeleteButton, TextField, ReferenceInput, SelectInput, EditButton, Edit, SimpleForm, TextInput, NumberInput, DateInput, Create, Filter, DateField, NumberField, ReferenceField } from 'react-admin';
import DownloadField from './DownloadField';
import TagField from './TagField';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import CardMedia from '@material-ui/core/CardMedia';

import Avatar from '@material-ui/core/Avatar';
import DescriptionIcon from '@material-ui/icons/Description';
import IconButton from '@material-ui/core/IconButton';

import Box from '@material-ui/core/Box';

import ArrowDownwardIcon from '@material-ui/icons/ArrowDownward';

const DocumentTitle = ({record}) => {
    return <span>Document {record ? `"${record.title}"` : ''}</span>
}


const DocumentFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <TextInput label="Tags" source="tags" />
    </Filter>
);

const TagParser = (v) => {
    return v?.split(',') ?? [];
}

const TagFormatter = (v) => {
    return v?.join(',') ?? '';
}

const cardStyle = {
    width: 300,
    margin: '0.5em',
    display: 'inline-block',
    verticalAlign: 'top',
    overflow: 'show',
};

const tagsStyle = {
    zIndex: 99,
    position: 'absolute',
    bottom: '6px',
    left: '6px',
}

const singleLine = {
    noWrap: true,
    textOverflow: 'ellipsis',
    style: {
        width: '170px',
        display: 'block',
    },
}

const DocumentGrid = ({ ids, data, basePath }) => (
    <div style={{ margin: '1em' }}>
    {ids.map(id =>
        <Card key={id} style={cardStyle} variant="outlined">             
            <CardHeader                
                title={data[id].title ? <TextField record={data[id]} source="title" {...singleLine} /> : <TextField record={data[id]} {...singleLine} source="filename" />}
                subheader={<DateField record={data[id]} source="timestamp" />}
                action={<IconButton component="a" href={`/api/fetch/${data[id].id}/${data[id].filename}`}  download target="_blank" ><ArrowDownwardIcon/></IconButton>}
                avatar={<Avatar><DescriptionIcon /></Avatar>}>
                </CardHeader>
            
            {
                data[id].thumbnail_path ? 
                <CardMedia image={`https://paperless2.kolhos.chichasov.es/api/thumb/${data[id].id}`} title="Thumbnail" style={{height: '400px', position: "relative"}}>
                    <div style={tagsStyle}>
                        <TagField record={data[id]} source="tags"/>
                    </div>
                </CardMedia> : null
            }
            {
                data[id].thumbnail_path ? null :
                    <CardContent style={{padding: 0}}>
                        <Box display="flex" justifyContent="center" alignItems="center" style={{height: 400}}>
                            <DescriptionIcon color="disabled" />                            
                        </Box>
                    </CardContent>
            }

            <CardActions style={{ textAlign: 'right' }}>
                <EditButton resource="document" basePath={basePath} record={data[id]} />
                <DeleteButton resource="document" record={data[id]}/>
            </CardActions>
        </Card>
    )}
    </div>
);

DocumentGrid.defaultProps = {
    data: {},
    ids: [],
};


export const DocumentList = props => (
    <List {...props} filters={<DocumentFilter/>} perPage={9}>
        <DocumentGrid/>
    </List>
);


export const DocumentEdit = props => (
    <Edit {...props} title={<DocumentTitle/>}>
        <SimpleForm>
            <TextInput source="id" disabled />
            <TextInput source="title" />
            <TextInput source="content" />
            <TextInput source="filename" />
            <DateInput source="timestamp" disabled />
            <TextInput source="tags"  parse={TagParser} format={TagFormatter} />
            <DownloadField source="id"/>
        </SimpleForm>
    </Edit>
);