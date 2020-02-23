import React from 'react';
import { Admin, Resource } from 'react-admin';

import authProvider from './api/authProvider';

import Dashboard from './components/Dashboard';
import apiProvider from './api/apiProvider';

import DescriptionIcon from '@material-ui/icons/Description';
import LocalOfferIcon from '@material-ui/icons/LocalOffer';

import { DocumentList, DocumentEdit } from './components/Document';
import { TagList, TagEdit, TagCreate } from './components/Tag';
import { createMuiTheme } from '@material-ui/core/styles';
import LoginPage from './components/login/LoginPage';

function App() {
  const theme = createMuiTheme({
    palette: {
      type: 'dark',
      primary: {
        main: '#7986cb',
      },
      secondary: {
        main: '#3f51b5',
      }
    }
  })

  return (
    <div>
      <Admin loginPage={LoginPage} dataProvider={apiProvider} dashboard={Dashboard} authProvider={authProvider} theme={theme}>
        <Resource name="document" list={DocumentList} edit={DocumentEdit} icon={DescriptionIcon} />
        <Resource name="tag" list={TagList} edit={TagEdit} create={TagCreate} icon={LocalOfferIcon}/>
      </Admin>
    </div>
  );
}

export default App;
