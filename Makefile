OUTPUT_DIR=./_output
SOURCE_DIRS=./

bin=webclient
product: fmtcheck vet
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${bin}.linux ./

build: fmtcheck vet
	go build -o ${OUTPUT_DIR}/${bin} ./

test: fmtcheck vet
	@go test -v -race ./pkg/...

.PHONY: fmtcheck
fmtcheck:
	@gofmt -s -l $(SOURCE_DIRS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
	
.PHONY: vet
vet:
	@go vet $(shell $(PACKAGES))

clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
