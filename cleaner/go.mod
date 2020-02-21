module github.com/guad/paperless2/cleaner

go 1.13

require github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71

require (
	github.com/guad/paperless2/backend v0.0.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
)

replace github.com/guad/paperless2/backend => ../backend
