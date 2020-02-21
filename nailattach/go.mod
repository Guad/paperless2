module github.com/guad/paperless2/nailattach

go 1.13

replace github.com/guad/paperless2/backend => ../backend

require (
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/guad/paperless2/backend v0.0.0-00010101000000-000000000000
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
)
