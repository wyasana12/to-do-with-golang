# To-Do-List App Berbasis Rest API dan GoLang

Ini adalah RESTful API untuk mengelola daftar tugas, dibangun menggunakan bahasa pemrograman Go. API ini mencakup fitur autentikasi pengguna dan operasi CRUD (Create, Read, Update, Delete) untuk item-item daftar tugas. Password hash yang digunakan dalam project tersebut adalah algoritma bcrypt.

## Fitur
#### Manajemen Pengguna
- Registrasi Pengguna.
- Login Pengguna (menghasilkan token JWT).
- Pengambilan Profil Pengguna.
 
#### Manajemen Item Daftar Tugas (Membutuhkan Autentikas)
- Membuat Item Daftar Tugas Baru.
- Melihat semua item daftar tugas untuk seorang pengguna.
- Melihat detail item daftar tugas tertentu.
- Memperbarui item daftar tugas yang sudah ada.
- Menghapus item daftar tugas.

## Teknologi Yang Digunakan
- Go (Golang).
- Gorilla Mux: Router http dan pencocok URL.
- GORM: Pustaka GORM untuk Go.
- Viper: Konfigurasi Go dengan fangs.
- Logrus: Pencatatan terstruktur untuk Go.
- Bcrypt: Algoritma hashing untuk katasandi.
- go-playground/validator: Validasi terstruktur dan field Go.

## Instalasi dan Pengaturan
#### 1. Kloning Repositori
```bash
git clone https://github.com/wyasana12/to-do-with-golang.git
cd to-do-list-go
```

#### 2. Instalasi Despendensi
```bash
git mod tidy
```

#### 3. Konfigurasi Basis Data
```bash
PORT=8000
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_DATABASE=go_todo_app
```
Ganti nilai-nilai tersebut dengan kredensial MySQL anda.

#### 4. Menjalankan Migrasi Basis Data
Aplikasi akan secara otomatis menjalankan fungsi AutoMigrate GORM untuk membuat tabel-tabel yang diperlukan (todos dan users) saat aplikasi dimulai.

#### 5. Menjalankan Aplikasi
```bash
go run main.go
```
Server akan berjalan di port yang ditentukan dalam file .env Anda (default: 8000).

## End Point API
Semua rute di bawah /api/user dan /api/todo memerlukan autentikasi menggunakan token Bearer di header Authorization.

#### Autentikasi (/api/auth)
    Registrasi Pengguna:
        POST /api/auth/register
    Login Pengguna:
        POST /api/auth/login
        
#### Pengguna (/api/user)
    Profil Pengguna:
        GET /api/user/profile (Membutuhkan autentikasi)
        
#### Daftar Tugas (/api/todo)
    Daftar Semua Tugas:
        GET /api/todo (Membutuhkan autentikasi)
    Buat Tugas Baru:
        POST /api/todo/create (Membutuhkan autentikasi)
    Detail Tugas:
        GET /api/todo/{id} (Membutuhkan autentikasi)
    Perbarui Tugas:
        PUT /api/todo/{id} (Membutuhkan autentikasi)
    Hapus Tugas:
        DELETE /api/todo/{id}/delete (Membutuhkan autentikasi)
        
