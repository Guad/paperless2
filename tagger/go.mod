module github.com/guad/paperless2/tagger

go 1.13

require github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71

require (
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/guad/paperless2/backend v0.0.0
)

replace github.com/guad/paperless2/backend => ../backend
