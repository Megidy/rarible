start:
	@docker-compose up --build -d
install-dependencies:
	@go install github.com/golang/mock/mockgen@latest
test:
	@go test ./... -v
gen-mock:
	@mockgen -source=internal/client/interface.go -destination=internal/client/mock/mock_raribleclient.go -package=client
	@mockgen -source=internal/service/interface.go -destination=internal/service/mock/mock_nft_service.go -package=service