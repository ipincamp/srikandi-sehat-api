.PHONY: app-debug app-build db-seed

app-debug:
	@go run .

app-build:
	@go build .

db-seed:
	@go run seeders/seeder.go