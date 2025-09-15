# SRIKANDI SEHAT
REST API ini dibangun menggunakan Go (GoFiber) dengan arsitektur yang bersih dan modern, menyediakan fondasi yang kuat untuk aplikasi frontend seperti Flutter.

## Daftar Isi

1. [Requirements](#1-requirements)

2. [Instalasi & Setup](#2-instalasi--setup)

3. [Fondasi & Arsitektur](#3-fondasi--arsitektur)

4. [Fitur & Endpoint](#4-fitur--endpoint)

5. [Komponen Pendukung Kualitas Kode](#5-komponen-pendukung-kualitas-kode)

6. [License](#6-license)

## 1. Requirements
Untuk menjalankan proyek ini di lingkungan development, Anda memerlukan perangkat lunak berikut:

- Go: Versi `1.21` atau yang lebih baru.

- MySQL: Versi `8.0` atau yang lebih baru (MariaDB juga didukung).

- Make: (Opsional, tapi sangat direkomendasikan) Untuk menjalankan perintah shortcut dari Makefile.

- Text Editor: Visual Studio Code dengan ekstensi Go sangat direkomendasikan.

## 2. Instalasi & Setup

Ikuti langkah-langkah berikut untuk menjalankan API ini di mesin lokal Anda.

1. Clone Repositori

    ```bash
    git clone https://github.com/ipincamp/go-srikandi-sehat-api.git
    cd go-srikandi-sehat-api
    ```

2. Konfigurasi Environment
Salin file `.env.example` menjadi `.env` dan sesuaikan nilainya dengan konfigurasi database dan kredensial seeder Anda.

    ```bash
    cp .env.example .env
    ```

    Buka file `.env` dan isi semua variabel yang diperlukan.

3. Instal Dependensi
Jalankan perintah berikut untuk mengunduh semua package yang dibutuhkan.

    ```bash
    go mod tidy
    ```

4. Jalankan Seeder Database
Perintah ini akan membuat semua tabel yang diperlukan dan mengisi data awal (roles, permissions, user default, dan data wilayah).

    ```bash
    go run database/seeders/main.go
    ```

    Atau, jika Anda menggunakan *`Makefile`*, Anda bisa membuat shortcut `make db-seed`.

5. Jalankan Aplikasi
Sekarang Anda siap untuk menjalankan server API.

    ```bash
    go run main.go
    ```

    Server akan berjalan di `http://0.0.0.0:3000` (atau port yang Anda tentukan di `.env`).

## 3. Fondasi & Arsitektur
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

## 4. Fitur & Endpoint
API ini memiliki beberapa fitur inti yang sudah siap untuk diintegrasikan dengan aplikasi frontend.

### Fitur 1: Sistem Autentikasi & Manajemen Akun
Sistem ini menangani semua kebutuhan dasar pengguna, mulai dari pendaftaran hingga pengelolaan akun, dengan keamanan sebagai prioritas utama.

Menggunakan JSON Web Tokens (JWT) untuk autentikasi yang aman. Password disimpan menggunakan hashing bcrypt. Sistem logout diperkuat dengan blocklist di database, memastikan token yang sudah di-logout benar-benar tidak bisa digunakan lagi.

- `POST /api/auth/register`

    Menerima permintaan pendaftaran dan menaruhnya di antrian untuk diproses di background. Response `202 Accepted` dengan pesan bahwa akun sedang diproses.

- `POST /api/auth/login`

    Mengautentikasi user dan mengembalikan JWT.

- `POST /api/user/logout` (Terproteksi)

    Membatalkan token JWT saat ini.

### Fitur 2: Manajemen Profil Pengguna
Pengguna yang sudah login dapat mengelola data personal mereka.

- `GET /api/me` (Terproteksi)

    Mengambil detail profil lengkap dari user yang sedang login.

- `PUT /api/me/details` (Terproteksi)

    Membuat atau memperbarui profil pengguna secara keseluruhan. Juga bisa digunakan untuk mengubah nama.

    Body: Semua field UpdateProfileRequest (lihat `dto/user.dto.go`).

- `PATCH /api/me/password` (Terproteksi)

    Mengubah password user.

    Body: `old_password`, `new_password`, `new_password_confirmation`.


### Fitur 2: Role-Based Access Control (RBAC)
Sistem hak akses yang terinspirasi dari `spatie/laravel-permission`, memungkinkan kontrol yang sangat detail terhadap apa yang bisa dilakukan oleh setiap user.

Terdapat dua level hak akses utama: **Admin** dan **User**. Endpoint tertentu hanya bisa diakses oleh user dengan role "Admin".

- `GET /api/admin/users` (Hanya Admin)

    Mengambil daftar semua user (kecuali admin lain) dengan sistem paginasi.

    > Query Params

    - `page` (opsional, default: 1): Nomor halaman yang ingin ditampilkan.

    - `limit` (opsional, default: 10): Jumlah data per halaman.

    > Contoh Response

    ```json
    {
        "status": true,
        "message": "Users fetched successfully",
        "data": {
            "data": [ /* ... daftar user ... */ ],
            "meta": {
                "limit": 10,
                "total_rows": 50,
                "total_pages": 5,
                "current_page": 1,
                "previous_page": null,
                "next_page": 2
            }
        }
    }
    ```

- `GET /api/admin/users/:id` (Hanya Admin)

    Mengambil detail user spesifik berdasarkan UUID.

### Fitur 3: API Data Wilayah Indonesia
Menyediakan data wilayah administrasi Indonesia yang bisa digunakan untuk fitur pemilihan alamat.

Endpoint ini menyajikan data provinsi, kabupaten, kecamatan, dan desa yang sudah di-seed ke dalam database.

- `GET /api/regions/provinces`

- `GET /api/regions/regencies?province_code=...`

- `GET /api/regions/districts?regency_code=...`

- `GET /api/regions/villages?district_code=...`

## 5. Komponen Pendukung Kualitas Kode
- Validasi Lanjutan: Menggunakan `go-playground/validator` dengan aturan validasi kustom (misalnya, `password_strength`) dan pesan error yang informatif untuk memastikan integritas data.

- Response JSON Standar: Semua response dari API mengikuti format yang konsisten (`status`, `message`, `data`) untuk memudahkan parsing di aplikasi client.

- Seeder Transaksional: Terdapat sistem seeder yang idempotent (aman dijalankan berkali-kali) dan transaksional. Seeder ini mengisi data awal untuk `roles`, `permissions`, user default (Admin & User) dari `.env`, dan data wilayah.

- Penggunaan: Cukup jalankan `make db-seed` (jika dikonfigurasi di Makefile) atau `go run database/seeders/main.go` untuk menyiapkan seluruh lingkungan development dalam satu perintah.


## 6. License
Proyek ini dilisensikan di bawah [MIT License](LICENSE).
