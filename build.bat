docker build -t guad/paperless-backend ./backend/. && docker push guad/paperless-backend
docker build -t guad/paperless-cleaner ./cleaner/. && docker push guad/paperless-cleaner
docker build -t guad/paperless-tagger ./tagger/. && docker push guad/paperless-tagger
docker build -t guad/paperless-nailattach ./nailattach/. && docker push guad/paperless-nailattach
docker build -t guad/paperless-ocr ./documentocr/. && docker push guad/paperless-ocr
docker build -t guad/paperless-frontend ./frontend/. && docker push guad/paperless-frontend
docker build -t guad/paperless-thumbnail ./thumbnailgen/. && docker push guad/paperless-thumbnail