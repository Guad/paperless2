import React from 'react';
import { Admin, Resource, ListGuesser, EditGuesser } from 'react-admin';

import authProvider from './api/authProvider';

import Dashboard from './components/Dashboard';
import apiProvider from './api/apiProvider';

import DescriptionIcon from '@material-ui/icons/Description';
import LocalOfferIcon from '@material-ui/icons/LocalOffer';

import { DocumentList, DocumentEdit } from './components/Document';
import { TagList, TagEdit, TagCreate } from './components/Tag';

function App() {
  return (
    <div>
      <Admin dataProvider={apiProvider} dashboard={Dashboard} authProvider={authProvider}>
        <Resource name="document" list={DocumentList} edit={DocumentEdit} icon={DescriptionIcon} />
        <Resource name="tag" list={TagList} edit={TagEdit} create={TagCreate} icon={LocalOfferIcon}/>
      </Admin>
    </div>
  );
}

export default App;
