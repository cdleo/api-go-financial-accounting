BUILDPATH=$(CURDIR)
API_NAME=api-go-financial-accounting

build: 
	@echo "Creando Binario ..."
	@go build -ldflags '-s -w' -o $(BUILDPATH)/build/bin/${API_NAME} cmd/main.go
	@echo "Binario generado en build/bin/${API_NAME}"

test: 
	@echo "Ejecutando tests..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

coverage:
	@echo "Coverfile..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out
	@go tool cover -func coverfile_out | grep total | awk '{print substr($$3, 1, length($$3)-1)}' > coverage.txt

.PHONY: test build
