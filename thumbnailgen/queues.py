
DocumentUploadQueue = "document_created"
DocumentThumbnailComplete = "document_thumbnail_complete"

def declare(chan):

    # 1. Consume queue from DocumentUploadQueue
    # 2. Exchange for DOCUMENT_OCR_COMPLETE

    chan.queue_declare('thumbnail_queue', durable=True)
    chan.queue_bind(exchange=DocumentUploadQueue, queue='thumbnail_queue')

    chan.exchange_declare(exchange=DocumentThumbnailComplete, exchange_type='fanout', durable=True)