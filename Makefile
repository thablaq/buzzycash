.PHONY: swagger
swagger:
	swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal  

.PHONY: run
run: swagger
	go run cmd/main.go          