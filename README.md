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