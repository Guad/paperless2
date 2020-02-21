import React, { useCallback, useState } from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';

import { useDropzone } from 'react-dropzone';

export default () => {
    const [files, setFiles] = useState([]);
    
    const upload = (file) => {
        const formData = new FormData();

        formData.append('title', '');
        formData.append('document', file);

        fetch('/api/push', {
            method: 'POST',
            body: formData,
        }).then((result) => {
            return result.json()
        }).then((result) => {
            setFiles(files.concat([{
                    id: files.length,
                    filename: result.filename,
                }])
            )
        }).catch((error) => {
            setFiles(files.concat([{
                    id: files.length,
                    error: true,
                }])
            )
        })
    }

    const onDrop = useCallback(acceptedFiles => {
        acceptedFiles.map(upload);
    }, [])
    const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop })

    const renderFile = (file) => (
        <div key={file.id}>
            {file.error ? <span>Error uploading file</span> : <span>Uploaded {file.filename}</span>}
        </div>
    );

    return (
        <Card>
            <CardHeader title="Paperless Administration" />
            <CardContent>
                <div {...getRootProps()} style={{
                    border: '3px dotted grey',
                    textAlign: 'center'
                }}>
                    <input {...getInputProps()} />
                    {
                        isDragActive ?
                            <p>Drop the files here ...</p> :
                            <p>Drag 'n' drop some files here, or click to select files</p>
                    }
                </div>
                <div style={{ marginTop: 5 }}>
                    {files.map(renderFile)}
                </div>
            </CardContent>
        </Card>
    )
};