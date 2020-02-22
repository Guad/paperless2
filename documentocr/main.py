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

def post_process(doc, content, ch, delivery_tag):
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

    ch.basic_ack(delivery_tag=delivery_tag)

def reject_threadsafe(ch, delivery_tag):
    ch.basic_reject(delivery_tag=delivery_tag)

def process_file(connection, ch, method, body):
    try:
        packet = json.loads(body)
        doc = packet['document']

        file = tmpdir + '/' + doc['filename']

        write_b64_data(packet['data'], file)

        parser = RasterisedDocumentParser(file)

        print("Parsing")        
        content = parser.get_text()

        os.remove(file)

        cb = functools.partial(post_process, doc, content, ch, method.delivery_tag)
        connection.add_callback_threadsafe(cb)

    except Exception as ex:
        print('ERROR CONVERTING:')
        print(ex)

        cb = functools.partial(reject_threadsafe, ch, method.delivery_tag)
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
        queue='ocr_queue',
        on_message_callback=on_message_callback,
        auto_ack=False,        
        )

    print("Starting reading messages.")
    ch.start_consuming()

if __name__ == "__main__":
    main()