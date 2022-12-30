module poc-backend

go 1.12

require (
	github.com/babygoat/logging-system/backend v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.7.0
	github.com/golang/protobuf v1.3.3
	github.com/jinzhu/gorm v1.9.10
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.8.0
	github.com/sirupsen/logrus v1.4.2
	google.golang.org/genproto v0.0.0-20191230161307-f3c370f40bfb
)

replace github.com/babygoat/logging-system/backend => ./
