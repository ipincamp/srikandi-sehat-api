.PHONY: migrate-create migrate-up migrate-down db-seed debug build

GOBIN=$(GOPATH)/bin
GOPATH=$(shell go env GOPATH)

migrate-create:
	@$(eval timestamp := $(shell date +%Y%m%d%H%M%S))
	@$(eval name := $(name))
	@if [ -z "$(name)" ]; then \
		echo "\033[31mError: migration name is required. Example: make migrate-create name=create_users_table\033[0m"; \
		exit 1; \
	fi
	@$(eval filename := database/migrations/$(timestamp)_$(name).go)
	@$(eval funcname := $(shell echo $(name) | sed -r 's/(^|_)([a-z])/\U\2/g'))

	@echo "package migrations" > $(filename)
	@echo "" >> $(filename)
	@echo "import (" >> $(filename)
	@echo "    \"github.com/go-gormigrate/gormigrate/v2\"" >> $(filename)
	@echo "    \"gorm.io/gorm\"" >> $(filename)
	@echo ")" >> $(filename)
	@echo "" >> $(filename)
	@echo "func $(funcname)() *gormigrate.Migration {" >> $(filename)
	@echo "    return &gormigrate.Migration{" >> $(filename)
	@echo "        ID: \"$(timestamp)\"," >> $(filename)
	@echo "" >> $(filename)
	@echo "        Migrate: func(tx *gorm.DB) error {" >> $(filename)
	@echo "            // TODO: Write your migration logic here" >> $(filename)
	@echo "            return nil" >> $(filename)
	@echo "        }," >> $(filename)
	@echo "" >> $(filename)
	@echo "        Rollback: func(tx *gorm.DB) error {" >> $(filename)
	@echo "            // TODO: Write your rollback logic here" >> $(filename)
	@echo "            return nil" >> $(filename)
	@echo "        }," >> $(filename)
	@echo "    }" >> $(filename)
	@echo "}" >> $(filename)

	@echo "\033[32mMigration file created successfully:\033[0m $(filename)"

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

db-seed:
	@go run cmd/seed/main.go

air-install:
	@echo "Installing Air..."
	@go install github.com/air-verse/air@latest
	@echo "Air installed successfully. You can now use 'make debug' to run the application with live reload."

dev: air-install
	@echo "Running in development mode with auto-reload..."
	@$(GOBIN)/air
debug:
	@go run .

build:
	@go build .
