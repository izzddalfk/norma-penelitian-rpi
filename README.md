# Penelitian Efisiensi Raspberry Pi 4 sebagai Server untuk UMKM

## API General

Seluruh respon dari service ini memiliki struktur dasar sebagai berikut:

- `code` (_Integer_): Representasi kode status dari respon HTTP.
- `status` (_String_): Kode status respon khusus dari aplikasi sebagai identifikasi dari respon.
- `data` (_Object_, OPTIONAL): Respon data dari request terkait jika request sukses.
- `errors` (_Array of Object_, OPTIONAL): Pesan khusus jika terjadi kesalahan.

**Contoh respon sukses:**

```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "message": "Berhasil mendaftarkan akun untuk Abdullah Faqih sebagai siswa"
  }
}
```

Ada dua versi API berdasarkan tingkat UMKM-nya

[1. API UMKM Kecil](#api-umkm-kecil)
[2. API UMKM Besar](#api-umkm-besar)

## API UMKM Kecil

Simulasi aplikasi untuk UMKM dengan usaha penjualan on-demand seperti kafe atau restoran.

### 1. Menampilkan stok barang

GET: `/api/small/stocks`

Endpoint ini digunakan untuk menampilkan daftar dan jumlah stok barang yang dimiliki oleh UMKM.

Query parameters:

- `page` (Number): Halaman yang ingin ditampilkan
- `total_goods` (Number): Jumlah barang yang ingin ditampilkan dalam satu halaman
- `sort` (String): Jenis urutan. Nilai yang valid adalah `ASC` dan `DESC`. Default `DESC`.
- `sort_by` (String): Urut daftar barang berdasarkan property tertentu. Default `id`, diurut berdasarkan ID barang.

Contoh request:

```text
GET /api/small/stocks?sort_by=id HTTP/1.1
```

Contoh response:

```json
{
  "code": 200,
  "status": "OK",
  "data": [
    {
      "id": 1,
      "name": "Kopi Gula Aren",
      "stock": 100,
      "price": 5000
    },
    {
      "id": 2,
      "name": "Kopi Susu",
      "stock": 30,
      "price": 6000
    },
    {
      "id": 3,
      "name": "Pisang Goreng",
      "stock": 150,
      "price": 1500
    }
  ]
}
```

### 2. Menambahkan barang ke keranjang

POST: `/api/small/cart`

Endpoint ini digunakan untuk membuat record transaksi belanja untuk seorang user. Di simulasi ini langsung menggunakan record transaksi untuk menyederhanakan proses. Namun transaksi yang statusnya belum dibayarkan disebut sebagai keranjang belanja.

Request body:

- `cart_id` (Number, _Optional_): ID dari keranjang belanja dari seorang user. Apabila keranjang belanja sudah ada, property ini harus terisi.
- `user_id` (Number): ID dari pengguna.
- `goods_id` (Number): ID barang yang ingin ditambahkan ke dalam keranjang belanja.
- `goods_price` (Number): Harga satuan barang.
- `total` (Number): Jumlah barang yang ditambahkan.

Response:

- `cart_id` (Number): ID keranjang belanja.
- `total_goods` (Number): Jumlah barang yang ada di keranjang belanja saat ini
- `total_amount` (Number): Total belanja keseluruhan saat ini

Contoh Request:

```json
POST /api/small/cart HTTP/1.1
Content-Type: application/json

{
  "user_id": 100,
  "goods_id": 1,
  "goods_price": 2000,
  "total": 3
}
```

Contoh Response:

```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "cart_id": 1,
    "total_goods": 3,
    "total_amount": 6000
  }
}
```

### 3. Melakukan pembayaran / pembelian

POST: `/api/small/pay`

Endpoint ini digunakan untuk membayar barang yang ada pada keranjang belanja dari user. Proses yang terjadi sebenarnya hanyalah mengubah status dari keranjang belanja menjadi "transaksi".

Request body:

- `cart_id` (Number): ID keranjang belanja
- `payment_amount` (Number): Jumlah pembayaran yang dilakukan oleh user

Response:

- `transaction_id` (String): ID transaksi
- `total_amount` (Number): Jumlah harga pembelian barang
- `payment_amount` (Number): Jumlah pembayaran yang dilakukan oleh user
- `return_amount` (Number): Jumlah uang yang dikembalikan oleh merchant kepada user

Contoh request:

```json
{
  "cart_id": 1,
  "payment_amount": 4000
}
```

Contoh response:

```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "transaction_id": "b915141c-82a9-48eb-842f-b4c64794dcb9",
    "total_amount": 4000,
    "payment_amount": 4000,
    "return_amount": 0
  }
}
```

## API UMKM Besar

Simulasi yang memiliki fitur dari UMKM Kecil, dengan tambahan berikut

### 4. Request ongkos kirim

GET: `/api/big/delivery-price`

Endpoint ini digunakan untuk mensimulasikan permintaan perhitungan ongkos kirim kepada service logistik. Service tersebut hanyalah dummy, service tambahan sederhana yg khusus melakukan perhitungan ongkos kirim.

Query parameter:

- `location` (Number): ID lokasi alamat pembeli

### 5. Request logistik (akan ada service tambahan sederhana yg khusus memberikan respon berhasil meminta penjemputan barang)

POST: `/api/big/pickup`

### 6. Modifikasi stok barang

#### 6.1 Menambah stok

POST: `/api/big/{stuff_name}/stocks`

Payload:

- `action` (String): Value-nya `INCR`
- `total` (Number): Jumlah barang yang ditambahkan kedalam stok

#### 6.2 Mengurangi stok

POST: `/api/big/{stuff_name}/stocks`

- `action` (String): Value-nya `DECR`
- `total` (Number): Jumlah barang yang dikurangi dari stok
