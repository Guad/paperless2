import pika
import os
import json
from pathlib import Path
import threading
import base64
import functools

import queues
from parsers.raster import RasterisedDocumentParser

tmpdir = '/tmp/paperless2'
rejections = {}
MAX_REJECTIONS = os.getenv("MAX_REJECTIONS")

if MAX_REJECTIONS == None:
    MAX_REJECTIONS = 5
else:
    MAX_REJECTIONS = int(MAX_REJECTIONS)

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

def post_process(doc, thumbnail, ch, delivery_tag):
    print("Publishing")

    docid = doc['id']

    packet = json.dumps({
        'document': doc,
        'thumbnail': thumbnail,
    })

    ch.basic_publish(
        exchange=queues.DocumentThumbnailComplete,
        routing_key='',
        body=packet,
        properties=pika.BasicProperties(
            delivery_mode = 2, # persistant
        ),
    )

    if docid in rejections:
        del rejections[docid]

    ch.basic_ack(delivery_tag=delivery_tag)

def reject_threadsafe(docid, ch, delivery_tag):
    if not docid in rejections:
        rejections[docid] = 0
    rejections[docid] += 1
    delete = rejections[docid] > MAX_REJECTIONS

    if delete:
        print('Permanently rejecting document ' + docid)
        del rejections[docid]

    ch.basic_reject(delivery_tag=delivery_tag, requeue=(not delete))

def process_file(connection, ch, method, body):
    docid = None
    try:
        packet = json.loads(body)
        doc = packet['document']
        docid = doc['id']

        file = tmpdir + '/' + doc['filename']

        write_b64_data(packet['data'], file)

        parser = RasterisedDocumentParser(file)

        print("Thumbnailing")
        thumbnail_path = parser.get_optimised_thumbnail()

        thumbnail = fget_b64_data(thumbnail_path)

        os.remove(file)
        os.remove(thumbnail_path)

        cb = functools.partial(post_process, doc, thumbnail, ch, method.delivery_tag)
        connection.add_callback_threadsafe(cb)

    except Exception as ex:
        print('ERROR CONVERTING:')
        print(ex)

        cb = functools.partial(reject_threadsafe, docid, ch, method.delivery_tag)
        connection.add_callback_threadsafe(cb)
        

def callback(ch, method, properties, body, args):
    print("Dequeuing")

    (connection) = args

    t = threading.Thread(target=process_file, args=(connection, ch, method, body))
    t.start()
    

def main():
    print('Starting!')
    Path(tmpdir).mkdir(parents=True, exist_ok=True)

    config = readConfig('/config/rabbitmq.json', 'RABBITMQ_SECRETS')

    creds = pika.PlainCredentials(config['username'], config['password'])
    host, port = config['host'].split(':')

    conn = pika.BlockingConnection(pika.ConnectionParameters(host, int(port), '/', creds))

    ch = conn.channel()

    queues.declare(ch)

    print('Setting QoS')
    ch.basic_qos(prefetch_count=1)

    print('Starting consuming')

    on_message_callback = functools.partial(callback, args=(conn))

    ch.basic_consume(
        queue='thumbnail_queue',
        on_message_callback=on_message_callback,
        auto_ack=False,        
        )

    print("Starting reading messages.")
    ch.start_consuming()

if __name__ == "__main__":
    main()