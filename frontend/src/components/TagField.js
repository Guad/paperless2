import Chip from '@material-ui/core/Chip';
import React from 'react';
import Typography from '@material-ui/core/Typography'

const TagField = ({ record, source }) => {
  let array = record[source]
 
   if (typeof array === 'string') {
       array = [array];
   }

  if (typeof array === 'undefined' || array === null || array.length === 0) {
    return <div></div>
  } else {
    array = array.map(v => v.trim()).filter(v => !!v);
    return (
      <>
        {array.map(item => <Chip
            size="small"
            color="primary"
            label={item}
            key={item}
            style={{marginRight: 4, marginTop: 4}}/>)}
      </>
    )    
  }
}
TagField.defaultProps = { addLabel: true }

export default TagField;
