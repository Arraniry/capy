# Capy CLI Tool

Capy adalah sebuah CLI tool yang membantu Anda membuat proyek Go dengan struktur Clean Architecture secara cepat dan mudah. Dengan Capy, Anda dapat membuat proyek baru dan generate berbagai komponen seperti controller, repository, dan usecase.

## Fitur

- Membuat proyek Go baru dengan struktur Clean Architecture.
- Generate komponen seperti controller, repository, dan usecase.
- Pilihan untuk menggunakan berbagai jenis database (MySQL, PostgreSQL, dll).

## Instalasi

Untuk menginstal Capy, pastikan Anda telah menginstal Go di sistem Anda. Kemudian, jalankan perintah berikut:

```bash
go install github.com/arraniry/capy@latest
```

## Penggunaan

### Membuat Proyek Baru

Untuk membuat proyek baru, gunakan perintah berikut:

```bash
capy new [nama-proyek] [database]
```

Contoh:

```bash
capy new my-app mysql
```

### Generate Komponen

Anda juga dapat mengenerate komponen tertentu setelah proyek dibuat. Gunakan perintah berikut:

```bash
capy generate [tipe] [nama]
```

Contoh:

```bash
capy generate controller user
```

### Generate Modul

Untuk mengenerate modul lengkap (model, controller, repository, dan usecase), gunakan perintah berikut:

```bash
capy module [nama-modul]
```

Contoh:

```bash
capy module user
```

## Kontribusi

Jika Anda ingin berkontribusi pada proyek ini, silakan fork repository ini dan buat pull request. Semua kontribusi sangat dihargai!

## Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).
