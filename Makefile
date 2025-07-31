include .env
export


DB_URL=mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?multiStatements=true

# Perintah untuk membuat file migrasi baru
# Contoh: make migrate-create name=create_users_table
migrate-create:
	@migrate create -ext sql -dir db/migration -seq $(name)

# Perintah untuk menjalankan migrasi (up)
migrate-up:
	@migrate -database "$(DB_URL)" -path db/migration up

# Perintah untuk rollback migrasi (down)
migrate-down:
	@migrate -database "$(DB_URL)" -path db/migration down

# Perintah untuk rollback semua migrasi
migrate-down-all:
	@migrate -database "$(DB_URL)" -path db/migration down -all
