FROM alpine:3.7

RUN \
    # Install required dependencies
    apk --no-cache --update add \
        python3 gnupg libmagic bash \
        sudo poppler tesseract-ocr imagemagick ghostscript unpaper optipng \
    # Install temporary build dependencies
    && apk --no-cache add --virtual .build-dependencies \
        git python3-dev poppler-dev gcc g++ musl-dev zlib-dev jpeg-dev postgresql-dev \
    # Bootstrap pip
    && python3 -m ensurepip

WORKDIR /usr/src/paperless2

ADD requirements.txt .

RUN pip3 install -r requirements.txt

ADD . . 

ENTRYPOINT [ "python3", "-u", "main.py" ]