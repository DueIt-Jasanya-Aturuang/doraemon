# Nama target default
.DEFAULT_GOAL := help

# Variabel untuk parameter koneksi database
DB_URL := "postgres://dueit_user:dueit_password@27.112.78.47:5432/dueit_db?search_path=auth&sslmode=disable"

name := "m_default_payment_method"
# Variabel untuk direktori migrations
MIGRATIONS_DIR := infra/migrations

# Perintah untuk menjalankan migrasi
migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) up

# Perintah untuk membatalkan migrasi
migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) down

# Perintah untuk memaksa migrasi ke versi tertentu
migrate-force:
	migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) force 1

# Perintah untuk membuat file migrasi baru
create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq create_table_$(name)

# Target bantuan (help) untuk menampilkan panduan penggunaan makefile
help:
	@echo "Panduan Penggunaan Makefile:"
	@echo "make migrate-up            : Menjalankan migrasi"
	@echo "make migrate-down          : Membatalkan migrasi"
	@echo "make migrate-force         : Memaksa migrasi ke versi tertentu"
	@echo "make create-migration name : Membuat file migrasi baru dengan nama 'name'"
