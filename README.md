# KataGenzi API Backend

## Project Description

This is the API backend for the KataGenzi article platform. Built with Go (Golang) and the Fiber framework, this project provides a set of RESTful endpoints to manage users, articles (posts), tags, and image uploads. The architecture is designed to be *decoupled*, allowing a frontend or any other client application to interact with it independently.

## Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Fiber
- **Database ORM:** GORM
- **Database:** Configurable for PostgreSQL, MySQL, or SQLite (example uses PostgreSQL).
- **Environment Variables:** godotenv
- **Authentication:** JSON Web Tokens (JWT)
- **Image Storage:** Cloudinary

## Project Structure

```
backend/
├── database/         # Konfigurasi dan koneksi database (GORM)
├── handlers/         # Logika untuk menangani permintaan API (Controllers)
├── middleware/       # Middleware untuk request (e.g., autentikasi JWT)
├── models/           # Representasi tabel database (struct GORM)
├── utils/            # Fungsi utilitas (e.g., JWT, Cloudinary, validasi)
├── go.mod            # Manajemen dependensi Go
├── go.sum            # Checksum dependensi
├── main.go           # Titik masuk aplikasi dan definisi rute
└── .env.example      # Contoh file untuk variabel lingkungan
```

## Panduan Instalasi dan Pengaturan

### 1. Prasyarat
Pastikan perangkat lunak berikut terinstal di mesin Anda:
- Go (v1.21+ direkomendasikan)
- Git
- Server Database (seperti PostgreSQL, MySQL, atau lainnya yang didukung GORM)

### 2. Pengaturan Awal
1.  Kloning repositori ini:
    ```bash
    git clone https://github.com/mohamadsolkhannawawi/article-project-backend.git
    ```
2.  Masuk ke direktori proyek:
    ```bash
    cd article-project-backend
    ```

### 3. Konfigurasi Backend (Go)
1.  Instal dependensi Go:
    ```bash
    go mod tidy
    ```
2.  Buat file lingkungan dengan menyalin dari contoh:
    ```bash
    cp .env.example .env
    ```
3.  Buka file `.env` dan konfigurasikan koneksi database serta kredensial lainnya. Contoh di bawah ini menggunakan **PostgreSQL**.

    ```env
    # --- DATABASE ---
    # Sesuaikan dengan driver database yang Anda gunakan (e.g., postgres, mysql)
    DB_URL="host=localhost user=your_user password=your_password dbname=katagenzi_db port=5432 sslmode=disable"

    # --- JWT ---
    JWT_SECRET="your_super_secret_key"
    JWT_EXPIRES_IN=72h # Durasi token (contoh: 72 jam)

    # --- CLOUDINARY ---
    CLOUDINARY_URL="cloudinary://api_key:api_secret@cloud_name"
    ```

4.  **Pengaturan Database:**
    -   Jalankan server database Anda (misalnya PostgreSQL).
    -   Buat database baru dengan nama yang Anda tentukan di file `.env` (contoh: `katagenzi_db`).

5.  Jalankan server pengembangan backend:
    ```bash
    go run main.go
    ```
    Aplikasi akan berjalan dan secara otomatis menjalankan migrasi database saat pertama kali dimulai. API akan tersedia di `http://localhost:3000`.

### 4. Mengakses Aplikasi
-   Setelah backend berjalan, Anda dapat mulai menggunakannya dengan aplikasi frontend atau alat pengujian API seperti Postman.
-   Untuk panduan pengujian menggunakan Postman, silakan merujuk ke file `POSTMAN_QUICK_START.md` di root proyek.
