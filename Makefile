help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: dep-imports ## Run go fmt on all go files
	go fmt $$(go list ./...)
	goimports -format-only -w $$(go list -f {{.Dir}} ./...)

dep-imports:
	go get -d golang.org/x/tools/cmd/goimports@latest

dep-linter:
	go get -d github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: dep-linter ## Run all the linters
	golangci-lint \
        		--enable=dupl \
        		--enable=bodyclose \
        		--enable=prealloc \
        		--enable=gofmt \
        		--enable=gomnd \
        		--enable=unconvert \
        		--enable=unparam \
        		--enable=asciicheck \
        		--enable=exhaustive \
        		--enable=exportloopref \
        		--enable=goconst \
        		--enable=goerr113 \
        		--enable=whitespace \
        		run ./...

protobuf: ## Run protobuf compiler
	@go get -d google.golang.org/grpc
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

	@protoc -I=${GOPATH}/src -I internal/transport/grpc/ping-pong/proto \
		--go_out=internal/transport/grpc/ping-pong/proto/pb \
		--go-grpc_out=internal/transport/grpc/ping-pong/proto/pb \
       	internal/transport/grpc/ping-pong/proto/*.proto



