# SRIKANDI SEHAT
REST API ini dibangun menggunakan Go (GoFiber) dengan arsitektur yang bersih dan modern, menyediakan fondasi yang kuat untuk aplikasi frontend seperti Flutter.

## 1. Requirements
Untuk menjalankan proyek ini di lingkungan development, Anda memerlukan perangkat lunak berikut:

- Go: Versi `1.21` atau yang lebih baru.

- MySQL: Versi 8.0 atau yang lebih baru (MariaDB juga didukung).

- Make: (Opsional, tapi sangat direkomendasikan) Untuk menjalankan perintah shortcut dari Makefile.

- Text Editor: Visual Studio Code dengan ekstensi Go sangat direkomendasikan.

## 2. Fondasi & Arsitektur
Proyek ini dibangun di atas tumpukan teknologi modern dan mengikuti prinsip-prinsip clean architecture untuk memastikan skalabilitas dan kemudahan pengelolaan.

### Teknologi Utama

- Bahasa: Go

- Framework: Fiber v2 (terinspirasi dari Express.js, sangat cepat)

- Database: MySQL

- ORM: GORM (ORM yang matang untuk Go)

- Validasi: `go-playground/validator/v10`

### Struktur Proyek
Proyek ini menggunakan arsitektur berlapis yang terorganisir dengan baik untuk memisahkan tanggung jawab:

- `config`: Mengelola konfigurasi dari file `.env`.

- `database`: Mengelola koneksi database (MySQL) dan seeder.

- `src/constants`: Menyimpan nilai konstan (enum) untuk `roles` dan `classifications`.

- `src/dto`: (Data Transfer Object) Mendefinisikan `struct` untuk data input (request) dan output (response).

- `src/handlers`: Berisi logika bisnis untuk setiap endpoint.

- `src/middleware`: Berisi middleware untuk proteksi rute (otentikasi JWT, pengecekan role).

- `src/models`: Mendefinisikan `struct` GORM yang merepresentasikan tabel di database.

- `src/routes`: Mendefinisikan semua rute API.

- `src/utils`: Berisi fungsi-fungsi pembantu (helper) seperti hashing, JWT, validasi, dan response standar.

## 3. Fitur-Fitur Utama
API ini memiliki beberapa fitur inti yang sudah siap untuk diintegrasikan dengan aplikasi frontend.

### Fitur 1: Sistem Autentikasi & Manajemen Akun
Sistem ini menangani semua kebutuhan dasar pengguna, mulai dari pendaftaran hingga pengelolaan akun, dengan keamanan sebagai prioritas utama.

- Penjelasan: Menggunakan JSON Web Tokens (JWT) untuk autentikasi yang aman. Password disimpan menggunakan hashing bcrypt. Sistem logout diperkuat dengan blocklist di database, memastikan token yang sudah di-logout benar-benar tidak bisa digunakan lagi.

- Endpoint:

    - `POST /api/auth/register`: Mendaftarkan user baru dan secara otomatis memberinya role "User".

    - `POST /api/auth/login`: Mengautentikasi user dan mengembalikan JWT.

    - `GET /api/user/profile`: (Terproteksi) Mengambil detail profil user yang sedang login.

    - `PATCH /api/user/profile/details`: (Terproteksi) Memperbarui nama atau email user.

    - `PATCH /api/user/profile/password`: (Terproteksi) Mengubah password setelah verifikasi password lama.

    - `POST /api/user/logout`: (Terproteksi) Membatalkan token JWT saat ini.

### Fitur 2: Role-Based Access Control (RBAC)
Sistem hak akses yang terinspirasi dari `spatie/laravel-permission`, memungkinkan kontrol yang sangat detail terhadap apa yang bisa dilakukan oleh setiap user.

- Penjelasan: Terdapat dua level hak akses utama: Admin dan User. Endpoint tertentu hanya bisa diakses oleh user dengan role "Admin".

- Endpoint Khusus Admin:

    - `GET /api/admin/users`: (Hanya Admin) Mengambil daftar semua user, kecuali admin lain.

    - `GET /api/admin/users/:id`: (Hanya Admin) Mengambil detail user spesifik berdasarkan UUID.

### Fitur 3: API Data Wilayah Indonesia
Menyediakan data wilayah administrasi Indonesia yang bisa digunakan untuk fitur pemilihan alamat.

- Penjelasan: Endpoint ini menyajikan data provinsi, kabupaten, kecamatan, dan desa yang sudah di-seed ke dalam database.

- Endpoint:

    - `GET /api/regions/provinces`

    - `GET /api/regions/regencies?province_code=...`

    - `GET /api/regions/districts?regency_code=...`

    - `GET /api/regions/villages?district_code=...`

## 4. Komponen Pendukung Kualitas Kode
- Validasi Lanjutan: Menggunakan `go-playground/validator` dengan aturan validasi kustom (misalnya, `password_strength`) dan pesan error yang informatif untuk memastikan integritas data.

- Response JSON Standar: Semua response dari API mengikuti format yang konsisten (`status`, `message`, `data`) untuk memudahkan parsing di aplikasi client.

- Seeder Transaksional: Terdapat sistem seeder yang idempotent (aman dijalankan berkali-kali) dan transaksional. Seeder ini mengisi data awal untuk `roles`, `permissions`, user default (Admin & User) dari `.env`, dan data wilayah.

- Penggunaan: Cukup jalankan `make db-seed` (jika dikonfigurasi di Makefile) atau `go run database/seeders/main.go` untuk menyiapkan seluruh lingkungan development dalam satu perintah.


## 5. License
Proyek ini dilisensikan di bawah [MIT License](LICENSE).
