# UTS Pemrograman Jaringan - Barang Tertinggal di Kampus
## Yohanes Bramanta Adita Saputra - 233409009

## 1. Tema & Masalah
Mengimplementasikan layanan pengelolaan **Barang Tertinggal di Kampus**

## 2. Cara Menjalankan 

### Run server-go
```Bash
go run ./server-go
```

### Run client-go
```Bash
# Membuat data barang baru
go run ./klien-go http://localhost:8080 buat "Botol minum biru" Kantin

# Melihat detail barang berdasarkan ID
go run ./klien-go http://localhost:8080 lihat 1
```

### Run HTTP
```Bash
go run ./bukti-tcp-go localhost:8080
```
