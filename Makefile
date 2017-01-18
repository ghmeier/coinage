
deps:
	go get -u gopkg.in/gin-gonic/gin.v1
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/pborman/uuid

run:
	go build
	./expresso-billing