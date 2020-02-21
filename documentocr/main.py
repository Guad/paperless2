import pika
import os
import json
from pathlib import Path
import minio
import threading

import queues
from parsers.raster import RasterisedDocumentParser

tmpdir = '/tmp/paperless2'
S3Client = None

def readConfig(defpath, env):
    path = defpath
    if os.getenv(env):
       path = os.getenv(env)

    with open(path, 'r') as f:
        return json.load(f)

def download_file(doc):
    S3Client.fget_object('documents', doc['s3_path'], tmpdir + '/' + doc['filename'])

def upload_thumbnail(doc, path):
    target_path = os.path.join('thumbnails', doc['id'], 'thumbnail.png')

    S3Client.fput_object('documents', target_path, path, content_type='image/png')

    return target_path

def process_file(ch, method, body):
    try:
        doc = json.loads(body)

        file = tmpdir + '/' + doc['filename']

        print("Downloading file")
        download_file(doc)    

        parser = RasterisedDocumentParser(file)

        print("Parsing")        
        content = parser.get_text()
        print("Thumbnailing")
        thumbnail_path = parser.get_optimised_thumbnail()

        print("Uploading")
        s3path = upload_thumbnail(doc, thumbnail_path)

        print("Publishing")

        packet = json.dumps({
            'document': doc,
            'content': content
        })

        ch.basic_publish(
            exchange=queues.DocumentOCRComplete,
            routing_key='',
            body=packet,
            properties=pika.BasicProperties(
                delivery_mode = 2, # persistant
            ),
        )

        packet = json.dumps({
            'document': doc,
            'thumbnail': s3path
        })

        ch.basic_publish(
            exchange=queues.DocumentThumbnailComplete,
            routing_key='',
            body=packet,
            properties=pika.BasicProperties(
                delivery_mode = 2, # persistant
            ),
        )

        os.remove(file)
        os.remove(thumbnail_path)

        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as ex:
        print('ERROR CONVERTING:')
        print(ex)

        ch.basic_reject(delivery_tag=method.delivery_tag)

def callback(ch, method, properties, body):
    print("Dequeuing")

    t = threading.Thread(target=process_file, args=(ch, method, body))
    t.start()
    



def setups3():
    global S3Client 

    config = readConfig('s3.json', 'S3_SECRETS')

    S3Client = minio.Minio(config['endpoint'],
        access_key=config['access_key'],
        secret_key=config['secret_key'],
        region=config['region'],
        secure=True)

def main():
    print('Starting!')
    Path(tmpdir).mkdir(parents=True, exist_ok=True)

    setups3()

    config = readConfig('rabbitmq.json', 'RABBITMQ_SECRETS')

    creds = pika.PlainCredentials(config['username'], config['password'])
    host, port = config['host'].split(':')

    conn = pika.BlockingConnection(pika.ConnectionParameters(host, int(port), '/', creds))

    ch = conn.channel()

    queues.declare(ch)

    print('Setting QoS')
    ch.basic_qos(prefetch_count=1)

    print('Starting consuming')
    ch.basic_consume(
        queue='ocr_queue',
        on_message_callback=callback,
        auto_ack=False,        
        )

    print("Starting reading messages.")
    ch.start_consuming()

if __name__ == "__main__":
    main()