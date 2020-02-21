import React, { useCallback, useState } from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';

import { useDropzone } from 'react-dropzone';
import { Title } from 'react-admin';

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
            if (!result.ok) {
                throw new Error('Upload failed');
            }

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
            {file.error ? <span>Error uploading file.</span> : <span>Uploaded {file.filename}</span>}
        </div>
    );

    return (
        <Card>
            <Title title="Upload Document" />
            <CardHeader title="Upload Document" />
            <CardContent>
                <div {...getRootProps()} style={{
                    border: '3px dotted grey',
                    textAlign: 'center',
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'center',
                    height: 200,
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