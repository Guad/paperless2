docker build -t guad/paperless-backend ./backend/. && docker push guad/paperless-backend
docker build -t guad/paperless-cleaner ./cleaner/. && docker push guad/paperless-cleaner
docker build -t guad/paperless-tagger ./tagger/. && docker push guad/paperless-tagger
docker build -t guad/paperless-ocr ./documentocr/. && docker push guad/paperless-ocr