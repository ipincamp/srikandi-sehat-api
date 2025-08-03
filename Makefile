.PHONY: debug build seeding

debug:
	@go run .

build:
	@go build .

seeding:
	@go run database/seeders/main.go
