import pika
import os
import json
from pathlib import Path
import threading
import base64

import queues
from parsers.raster import RasterisedDocumentParser

tmpdir = '/tmp/paperless2'

def readConfig(defpath, env):
    path = defpath
    if os.getenv(env):
       path = os.getenv(env)

    with open(path, 'r') as f:
        return json.load(f)

def write_b64_data(data, path):
    data_bytes = base64.b64decode(data)
    with open(path, "wb") as f:
        f.write(data_bytes)

def fget_b64_data(path):
    with open(path, "rb") as f:
        data_bytes = f.read()
        return base64.b64encode(data_bytes).decode('utf8')

def process_file(ch, method, body):
    try:
        packet = json.loads(body)
        doc = packet['document']

        file = tmpdir + '/' + doc['filename']

        write_b64_data(packet['data'], file)

        parser = RasterisedDocumentParser(file)

        print("Parsing")        
        content = parser.get_text()
        print("Thumbnailing")
        thumbnail_path = parser.get_optimised_thumbnail()

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
            'thumbnail': fget_b64_data(thumbnail_path),
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
    

def main():
    print('Starting!')
    Path(tmpdir).mkdir(parents=True, exist_ok=True)

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