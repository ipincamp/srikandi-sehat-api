# ==============================================================================
# Makefile untuk Proyek Go (Srikandi Sehat API)
# ==============================================================================

# Variabel Proyek
BINARY_NAME=srikandisehat
MAIN_GO=./cmd/api/main.go
MIGRATE_GO=./cmd/migrate/main.go
SEED_GO=./cmd/seed/main.go

# Variabel Lingkungan
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
TIMEZONE=Asia/Jakarta

# Perintah default yang dijalankan jika 'make' dipanggil tanpa target
.DEFAULT_GOAL := help

# ==============================================================================
# DEFINISI PERINTAH
# ==============================================================================

help: ## ℹ️ Tampilkan semua perintah yang tersedia
	@echo "Perintah yang tersedia:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'


# --------------------------------------
# Perintah Build & Run
# --------------------------------------

build-run: ## --- Build & Run ---
	@# Target palsu ini hanya untuk pengelompokan di 'make help'

clean: ## 🧹 Bersihkan artefak build (direktori ./bin, ./tmp)
	@echo "Membersihkan artefak build..."
	@rm -rf ./bin/* ./tmp

build: ## 🏗️ Kompilasi aplikasi Go ke binary di ./bin
	@echo "Mem-build binary..."
	@mkdir -p ./bin
	@go build -o ./bin/$(BINARY_NAME) $(MAIN_GO)

run: build ## 🚀 Jalankan aplikasi (mode production)
	@echo "Menjalankan (mode production)..."
	@ENV=production TZ=$(TIMEZONE) ./bin/$(BINARY_NAME)

dev: air-install ## 🔄 Jalankan aplikasi (dev) dengan auto-reload (membutuhkan 'Air')
	@echo "Menjalankan (mode development) dengan auto-reload..."
	@$(GOBIN)/air

debug: ## 🐞 Jalankan aplikasi dengan debugger (Delve)
	@echo "Memulai debugger (Delve)..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@$(GOBIN)/dlv debug $(MAIN_GO)

# Target internal untuk 'dev', tidak perlu ditampilkan di help
air-install:
	@if ! command -v $(GOBIN)/air &> /dev/null; then \
		echo "Menginstall 'Air' untuk auto-reload..."; \
		go install github.com/air-verse/air@latest; \
	fi

# --------------------------------------
# Perintah Migrasi Database
# --------------------------------------

database-migrations: ## --- Database Migrations ---
	@# Target palsu ini hanya untuk pengelompokan di 'make help'

create-migration: ## 📝 Buat file migrasi baru. Cth: make create-migration name=create_users_table
	@echo "Membuat file migrasi..."
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-migration name=<nama_migrasi>"; \
		exit 1; \
	fi
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	func_name=$$(echo "$(name)" | sed -e 's/_\([a-z]\)/\u\1/g' -e 's/^\([a-z]\)/\u\1/g'); \
	filepath=database/migrations/$${timestamp}_$(name).go; \
	printf 'package migrations\n\nimport (\n\t"github.com/go-gormigrate/gormigrate/v2"\n\t"gorm.io/gorm"\n)\n\nfunc %s() *gormigrate.Migration {\n\t// TODO: Tentukan struct Anda di sini\n\t// type YourStruct struct {}\n\treturn &gormigrate.Migration{\n\t\tID: "%s",\n\t\tMigrate: func(tx *gorm.DB) error {\n\t\t\t// TODO: Implementasi migrasi (buat tabel/kolom)\n\t\t\t// Cth: return tx.AutoMigrate(&YourStruct{})\n\t\t\treturn nil\n\t\t},\n\t\tRollback: func(tx *gorm.DB) error {\n\t\t\t// TODO: Implementasi rollback (hapus tabel/kolom)\n\t\t\t// Cth: return tx.Migrator().DropTable("your_structs")\n\t\t\treturn nil\n\t\t},\n\t}\n}\n' "$$func_name" "$$timestamp" > $$filepath; \
	echo "Berhasil membuat: $$filepath"

migrate: ## ⬆️ Jalankan semua migrasi yang tertunda (up)
	@echo "Menjalankan migrasi (up)..."
	@go run $(MIGRATE_GO) up

migrate-down: ## ⬇️ Batalkan (rollback) migrasi terakhir (down)
	@echo "Me-rollback migrasi terakhir..."
	@go run $(MIGRATE_GO) down

migrate-reset: ## 🔄 HAPUS semua tabel lalu jalankan ulang SEMUA migrasi (ideal untuk dev)
	@echo "Mer-reset database (drop semua tabel & migrasi ulang)..."
	@go run $(MIGRATE_GO) reset

db-drop: ## ⚠️ DANGER! HAPUS semua tabel & JANGAN migrasi ulang (mengosongkan DB)
	@echo "PERHATIAN! Menghapus semua tabel (tanpa migrasi ulang)..."
	@go run $(MIGRATE_GO) drop-all

# --------------------------------------
# Perintah Seeder Database
# --------------------------------------

database-seeders: ## --- Database Seeders ---
	@# Target palsu ini hanya untuk pengelompokan di 'make help'

create-seeder: ## 🌱 Buat file seeder baru. Cth: make create-seeder name=admin_user
	@echo "Membuat file seeder..."
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-seeder name=<nama_seeder>"; \
		exit 1; \
	fi
	@func_name=$$(echo "$(name)" | sed -e 's/_\([a-z]\)/\u\1/g' -e 's/^\([a-z]\)/\u\1/g')Seeder; \
	filepath=database/seeders/$(name).go; \
	printf 'package seeders\n\nimport (\n\t"log"\n\n\t"gorm.io/gorm"\n)\n\nfunc %s(db *gorm.DB) error {\n\t// TODO: Implementasi logika seeder Anda di sini\n\t// Gunakan db.FirstOrCreate() untuk menghindari duplikat\n\tlog.Println("%s berjalan sukses")\n\treturn nil\n}\n' "$$func_name" "$$func_name" > $$filepath; \
	echo "Berhasil membuat: $$filepath"

seed: ## 💾 Jalankan semua seeder untuk mengisi data awal (roles, admin, dll)
	@echo "Menjalankan database seeders..."
	@go run $(SEED_GO)


# ==============================================================================
# PENGATURAN MAKEFILE
# ==============================================================================

# Mendefinisikan target mana yang bukan file
# Ini mencegah 'make' bingung jika ada file/folder dengan nama yang sama
.PHONY: help \
	build-run clean build run dev debug air-install \
	database-migrations create-migration migrate migrate-down migrate-reset db-drop \
	database-seeders create-seeder seed