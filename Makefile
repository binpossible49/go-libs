generate: ## Generate proto
	@protoc \
		-I common_proto/ \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/ \
		-I $(GOPATH)/src/github.com/gogo/protobuf/ \
		--gogo_out=plugins=grpc:common_proto \
		--govalidators_out=gogoimport=true:common_proto \
		--grpc-gateway_out=common_proto \
		common_proto/*.proto

test: ## Run test for whole project
	@go test -v ./...

lint: ## Run linter
	@golangci-lint run ./...