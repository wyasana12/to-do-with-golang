# To-Do-List App Berbasis Rest API dan GoLang

Ini adalah RESTful API untuk mengelola daftar tugas, dibangun menggunakan bahasa pemrograman Go. API ini mencakup fitur autentikasi pengguna dan operasi CRUD (Create, Read, Update, Delete) untuk item-item daftar tugas, Melakukan tambah, download, dan delete untuk attachment di setiap tugas, serta fungsionalitas keranjang sampah (trash) untuk pemulihan. Aplikasi ini juga dilengkapi dengan integrasi Swagger UI untuk dokumentasi API interaktif. Password hash yang digunakan dalam project tersebut adalah algoritma bcrypt.

## Fitur
#### Manajemen Pengguna
- Registrasi pengguna.
- Verifikasi email.
- Login pengguna (menghasilkan token JWT).
- Reset kata sandi.
- Pengambilan dan pembaruan profil pengguna.
 
#### Manajemen Item Daftar Tugas (Membutuhkan Autentikasi)
- Membuat item daftar tugas baru.
- Melihat semua item daftar tugas untuk seorang pengguna.
- Melihat detail item daftar tugas tertentu.
- Memperbarui item daftar tugas yang sudah ada.
- Menghapus item daftar tugas (soft delete ke tempat sampah).
- Hapus massal (bulk delete) item daftar tugas.

#### Manajemen Attachment Setiap Tugas (Membutuhkan Autentikasi)
- Membuat attachment di setiap tugas.
- Melakukan download attachment.
- Menghapus attachment untuk setiap tugas.

#### Manajemen Keranjang Sampah (Membutuhkan Autentikasi)
- Melihat item daftar tugas yang telah dihapus (di tempat sampah).
- Memulihkan item daftar tugas dari tempat sampah.
- Memulihkan massal (bulk restore) item daftar tugas dari tempat sampah.
- Menghapus permanen item daftar tugas dari tempat sampah.
- Hapus permanen massal (bulk permanent delete) item daftar tugas dari tempat sampah.

## Teknologi Yang Digunakan
- Go (Golang).
- Gorilla Mux: Router http dan pencocok URL.
- GORM: Pustaka GORM untuk Go.
- Viper: Konfigurasi Go dengan fangs.
- Logrus: Pencatatan terstruktur untuk Go.
- Bcrypt: Algoritma hashing untuk katasandi.
- go-playground/validator: Validasi terstruktur dan field Go.
- go-golang/jwt: Implementasi JSON Web Token.

## Instalasi dan Pengaturan
#### 1. Kloning Repositori
```bash
git clone https://github.com/wyasana12/to-do-with-golang.git
cd to-do-list-go
```

#### 2. Instalasi Despendensi
```bash
go mod tidy
```

#### 3. Konfigurasi Basis Data
Buat file .env di root proyek dengan isi sebagai berikut (ganti dengan kredensial anda):

```bash
APP_URL="http://localhost:8080"

DB_HOST="127.0.0.1"
DB_USER="root"
DB_PASSWORD=""
DB_DATABASE="to-do-list"
DB_PORT="3306"

JWT_SECRET_KEY=your_secret_key # Ganti dengan string acak yang kuat

PORT=8080

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com # Gunakan email yang akan digunakan untuk mengirim verifikasi
SMTP_PASSWORD=your_app_password # Gunakan App Password jika Anda menggunakan Gmail dengan 2FA
```
###### Catatan untuk Gmail SMTP:
Jika Anda menggunakan Gmail dan mengaktifkan Verifikasi 2 Langkah (2FA), Anda perlu membuat "App password" alih-alih menggunakan kata sandi Gmail biasa. Anda bisa menemukannya di pengaturan akun Google Anda di bagian Keamanan.

#### 4. Menjalankan Migrasi Basis Data
Aplikasi akan secara otomatis menjalankan fungsi AutoMigrate GORM untuk membuat tabel-tabel yang diperlukan (todos dan users) saat aplikasi dimulai.

#### 5. Menjalankan Aplikasi
```bash
go run main.go
```
Server akan berjalan di port yang ditentukan dalam file .env Anda (default: 8000).

#### 6. Mengakses Dokumentasi API (Swagger UI)
Setelah aplikasi berjalan, buka browser Anda dan navigasikan ke:
```bash
http://localhost:8080/swagger-ui/
```
Anda juga dapat mengakses spesifikasi OpenAPI mentah di:
```bash
http://localhost:8080/swagger.yaml
```

## End Point API
Semua rute di bawah /api/user, /api/todos, dan /api/todos/trash memerlukan autentikasi menggunakan token Bearer di header Authorization.

#### Autentikasi (/api/auth)
    Registrasi Pengguna:
        POST /api/auth/register
    Verifikasi Email:
        GET /api/auth/verify-email?token={token}
    Login Pengguna:
        POST /api/auth/login
    Reset Password:
        PUT /api/auth/reset-password
        
#### Pengguna (/api/user)
    Profil Pengguna:
        GET /api/user/profile
    Perbarui Profil Pengguna:
        PUT /api/user/update-profile
        
#### Daftar Tugas (/api/todos)
    Daftar Semua Tugas:
        GET /api/todos
    Buat Tugas Baru:
        POST /api/todos
    Detail Tugas:
        GET /api/todos/{id}
    Perbarui Tugas:
        PUT /api/todos/{id}
    Hapus Tugas (soft delete):
        DELETE /api/todos/{id}
    Hapus Tugas Massal (soft delete):
        DELETE /api/todos/bulk-delete

#### Attachment (/api/todos/{id}/attachment)
    Tambah Attachment:
        POST /api/todos/{id}/attachment
    Download Attachment:
        GET /api/todos/{id}/attachment/{attachment_id}/download
    Delete Attachment:
        DELETE /api/todos/{id}/attachment/{attachment_id}

#### Keranjang Sampah (/api/todos/trash)
    Daftar Tugas Di Tempat Sampah:
        GET /api/todos/trash
    Pulih Tugas Dari Tempat Sampah:
        PUT /api/todos/trash/restore/{id}
    Pulih Tugas Massal Dari Tempat Sampah:
        PUT /api/todos/trash/bulk-restore
    Hapus Permanen Tugas Dari Tempah Sampah:
        DELETE /api/todos/trash/permanent-delete/{id}
    Hapus Permanen Tugas Massal Dari Tempat Sampah:
        DELETE /api/todos/trash/bulk-permanent-delete
        
