# UTS Pemrograman Jaringan - Barang Tertinggal di Kampus
Yohanes Bramanta Adita Saputra - 233409009

### 1. Tema & Masalah
Mengimplementasikan layanan pengelolaan **Barang Tertinggal di Kampus**

### 2. Cara Menjalankan 

#### Run server-go
```Bash
go run ./server-go
```

#### Run client-go
```Bash
# Membuat data barang baru
go run ./klien-go http://localhost:8080 buat "Botol minum biru" Kantin

# Melihat detail barang berdasarkan ID
go run ./klien-go http://localhost:8080 lihat 1
```

#### Run HTTP
```Bash
go run ./bukti-tcp-go localhost:8080
```

### 3. Endpoint Table

| Method | Path | Success Status | Failed Status | `curl` example command |
| :---: | :---: | :---: | :---: | :--- |
| `POST` | `/barang` | 201 | 400, 409 | `curl -X POST http://localhost:8080/barang -H "Content-Type: application/json" -d '{"nama":"Botol Minum Biru","lokasi":"kantin"}'` |
| `GET` | `/barang` | 200 | - | `curl http://localhost:8080/barang` |
| `GET` | `/barang?status=tersimpan` | 200 | 400 | `curl "http://localhost:8080/barang?status=tersimpan"` |
| `GET` | `/barang/{id}` | 200 | 400,404 | `curl http://localhost:8080/barang/1` |
| `POST` | `/barang/{id}/klaim` | 200 | 400,404,409 | `curl -X POST http://localhost:8080/barang/1/klaim` |
| `DELETE` | `/barang/{id}` | 200 | 400,404, 409 | `curl -X DELETE http://localhost:8080/barang/1` |
| `GET` | `/healthz` | 200 | - | `curl http://localhost:8080/healthz` |

### 4. Retry Policy
- **Jumlah Retry & Jeda**: Klien melakukan retry maks 5 kali. Jika server mengirimkan respon header `Retry-After`, klien akan tahu durasi penundaan
- **Kondisi Retry (`lihat`)**: Retry otomatis dipicu apabila terjadi error pada jaringan, status HTTP `429`, atau status HTTP `5xx`. Status `400` dan `404` tidak akan di -retry karena merupakan kesalahan dari klien
- **Kenapa Perintah `tambah` Tidak di-retry**: Command `tambah` menggunakan metode HTTP `POST` yang bersifat **non-idempotent (operasi yang jika dijalankan berulang kali dengan input yang sama, menghasilkan perubahan status yang berbeda)**. Jika request di-retry secara otomatis saat penundaan jaringan, akan beresiko terjadi duplikasi data di server apabila request sebelumnya berhasil diproses

### 5. AI Usage
AI Usage pada proyek ini:
- Menggunakan Gemini
- Pembuatan Skema basic standard library `net/http` Go 1.22 untuk perbandingan string case-insentive, penyusunan library klien HTTP, dan penyusunan draft Dokumentasi `README.md`