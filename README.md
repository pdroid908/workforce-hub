# WAQI Air Quality Checker Web

Aplikasi web untuk memantau kualitas udara real-time berbasis data **WAQI (World Air Quality Index)**.

## Fitur Utama
- Cek kualitas udara via **GPS** pengguna.
- Cek kualitas udara via **pencarian kota/stasiun**.
- Menampilkan nilai **AQI**, status, polutan dominan, suhu, kelembapan, dan saran kesehatan.
- Data daftar stasiun global disimpan di **Redis cache (TTL 5 jam)** agar loading lebih cepat.
- API backend menggunakan **Go + Gin**, frontend statis menggunakan **HTML + Tailwind CSS**.

## Struktur Singkat
- `/index.html` → frontend
- `/api/index.go` → entrypoint API untuk Vercel
- `/waqi/waqi.go` → handler & logic integrasi WAQI
- `/redis/r.go` → koneksi dan helper Redis cache
- `/vercel.json` → konfigurasi build & route Vercel

## Environment Variables
Buat environment variable berikut:

- `WAQI_TOKEN` : token API dari WAQI
- `REDIS_URL` : URL koneksi Redis (contoh Upstash)

## Menjalankan Proyek (Lokal)
Proyek ini disiapkan untuk Vercel (Go serverless + static frontend).

1. Install Vercel CLI
   ```bash
   npm i -g vercel
   ```
2. Jalankan lokal
   ```bash
   vercel dev
   ```
3. Buka aplikasi di URL lokal yang diberikan Vercel.

## Endpoint API
### `GET /api/cities`
Mengambil daftar stasiun global untuk dropdown (cache Redis 5 jam).

### `GET /api/air-quality`
Mengambil data kualitas udara live.

Parameter query:
- `lat` & `lng` (opsional): mode GPS
- `city` (opsional): mode nama kota / UID stasiun (mis. `@1234`)

Contoh:
```http
GET /api/air-quality?lat=-7.797068&lng=110.370529
GET /api/air-quality?city=yogyakarta
GET /api/air-quality?city=@1234
```

## Sumber Data
- WAQI Official API: https://aqicn.org/api/

## Catatan
Jika `WAQI_TOKEN` atau `REDIS_URL` belum diatur, sebagian fitur API akan gagal diproses.
