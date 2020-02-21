
DocumentUploadQueue = "document_created"
DocumentOCRComplete = "document_ocr_complete"
DocumentThumbnailComplete = "document_thumbnail_complete"

def declare(chan):

    # 1. Consume queue from DocumentUploadQueue
    # 2. Exchange for DOCUMENT_OCR_COMPLETE

    chan.queue_declare('ocr_queue', durable=True)
    chan.queue_bind(exchange=DocumentUploadQueue, queue='ocr_queue')

    chan.exchange_declare(exchange=DocumentOCRComplete, exchange_type='fanout', durable=True)
    chan.exchange_declare(exchange=DocumentThumbnailComplete, exchange_type='fanout', durable=True)