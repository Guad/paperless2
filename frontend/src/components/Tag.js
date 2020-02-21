import React from 'react';
import { List, Datagrid, TextField, ReferenceInput, SelectInput, EditButton, Edit, SimpleForm, TextInput, NumberInput, DateInput, Create, Filter, DateField, NumberField, ReferenceField } from 'react-admin';

const TagTitle = ({record}) => {
    return <span>Tag {record ? `"${record.name}"` : ''}</span>
}


const TagFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Search" source="name" alwaysOn />
    </Filter>
);

export const TagList = props => (
    <List {...props} filters={<TagFilter/>}>
        <Datagrid rowClick="edit">
            <TextField source="name" />
            <TextField source="regex" />
            <EditButton/>
        </Datagrid>
    </List>
);

export const TagEdit = props => (
    <Edit {...props} title={<TagTitle/>}>
        <SimpleForm>
            <TextInput source="id" disabled />
            <TextInput source="name" />
            <TextInput source="regex" />
        </SimpleForm>
    </Edit>
);


export const TagCreate = props => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <TextInput source="regex" />
        </SimpleForm>
    </Create>
);