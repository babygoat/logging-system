module poc-backend

go 1.12

require (
	github.com/babygoat/logging-system/backend v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/pkg/errors v0.8.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/babygoat/logging-system/backend => ./
