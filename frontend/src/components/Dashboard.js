import React, { useCallback, useState } from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';

import { Title } from 'react-admin';

import Dropzone from './Dropzone';

export default () => {
    return (
        <Card>
            <Title title="Upload Document" />
            <CardHeader title="Upload Document" />
            <CardContent>
                <Dropzone/>
            </CardContent>
        </Card>
    )
};