# go-boilerplate

Deskripsi singkat:

Proyek ini adalah template boilerplate untuk layanan backend menggunakan Go. Berfungsi sebagai panduan dan standar struktur kode, konvensi penamaan, serta contoh implementasi utilitas, repository, service, dan handler HTTP (Gin).

**Struktur**
Penjelasan singkat tiap direktori/top-level file dan peranannya dalam proyek (tujuan: standarisasi struktur dan mempermudah on-boarding):

- Root files
	- `go.mod`, `go.sum`: dependency dan versi module.
	- `Makefile`: shortcut tugas developer (build, run, test, migrate).
	- `README.md`: dokumentasi proyek.

- `cmd/`
	- Berisi entrypoint aplikasi. Setiap subfolder mewakili satu binary/entrypoint (mis. `main`). Di dalamnya diparsing flag, inisialisasi konfigurasi, dan memanggil wiring aplikasi.

- `internal/` :Kode aplikasi yang tidak diekspor ke package lain. Struktur internal merepresentasikan boundary arsitektur aplikasi.
	- `apps/`: wiring dan bootstrap aplikasi (menghubungkan repo, service, transport, migration, health checks).
	- `configs/`: definisi struct konfigurasi, loader, dan validasi konfigurasi (file/env precedence).
	- `dbs/`: inisialisasi koneksi database, pool, helper transaksi, dan health checks DB.
	- `dtos/`: transport-level data shapes (request/response DTO), dipisah dari domain entities.
	- `entities/`: domain models dan value objects.
	- `repositories/`: definisi interface repository dan implementasi penyimpanan atau data access (SQL/NoSQL/cache/RestApi). Pisahkan interface dan implementasi untuk memudahkan mocking.
	- `services/`: use-cases / business logic yang mengorkestrasi repositori dan external clients.
	- `transports/`: adapter transport (HTTP, gRPC, RMQ); untuk HTTP: `transports/http/router.go`, `handlers/`, `middlewares/`.
	- `utils/`
		- `validation/`: aturan validasi input (dipanggil sebelum masuk ke service).
		- `logs/`: konfigurasi logger dan helper structured logging.

- `internal/repositories/_mock/` dan test files
	- Tempat mock otomatis untuk interfaces (digenerate dalam tests) dan test-related helpers.

- `tmp/`
	- Direktori untuk artifacts lokal atau binary hasil build selama development; tidak untuk file penting.

Pedoman tambahan
	- Handler HTTP dan middleware diletakkan di `transports/http/handlers` dan `transports/http/middlewares`.
	- Tests mengikuti struktur yang sama dengan file yang dites (mis. `services/example_service_test.go`).
	- File konfigurasi contoh (mis. `.env.stage.example`) disimpan di root, jangan commit file dengan secret.

Tujuan struktur ini: menjaga separation of concerns, memudahkan testing, dan membuat dependency wiring eksplisit.


**Arsitektur & Aliran Data (Layer-by-Layer)**

Bagian ini menjelaskan secara rinci tanggung jawab tiap layer dalam proyek dan bagaimana data mengalir di antara layer-layer tersebut. Tidak ada contoh kode di bagian ini — hanya aturan dan penjelasan konseptual untuk standarisasi arsitektur.

- Entrypoint / CLI (`cmd`)
	- Tanggung jawab: parsing flag/argument, memilih `stage`, inisialisasi konfigurasi, logger, dan koneksi infrastruktur (DB, broker). Melakukan bootstrap aplikasi dan menyerahkan kontrol ke layer `apps`.
	- Aliran: menentukan stage → memuat konfigurasi → inisialisasi infrastruktur → memanggil wiring aplikasi.

- Configs (`internal/configs`)
	- Tanggung jawab: memuat, memvalidasi, dan menyediakan konfigurasi read-only ke seluruh aplikasi.
	- Aliran: nilai config diinject ke konstruktor (dependency injection) untuk digunakan oleh apps, repositori, dan services.

- Apps / Wiring (`internal/apps`)
	- Tanggung jawab: menyusun dan menghubungkan semua dependensi (repo, service, handler), mengatur lifecycle (migrations, health checks), dan memulai transport.
	- Aliran: menerima config & infra → instantiate komponen → registrasi route & middleware → mulai server.

- Transports (HTTP) — Router & Server (`transports/http`)
	- Tanggung jawab: menerima koneksi masuk, routing, dan lifecycle server.
	- Aliran request: masuk → middleware → handler.
Env, Stage & Mode

.env files menyimpan variabel lingkungan; jangan commit file yang berisi kredensial. Simpan contoh variabel di `.env.example` atau `.env.stage.example`.

File env per stage

- Gunakan file terpisah per stage untuk memudahkan deployment dan pengujian, misalnya:
	- `.env.stage.local`
	- `.env.stage.development`
	- `.env.stage.staging`
	- `.env.stage.production`

- Pilih konvensi nama yang konsisten di seluruh tim (mis. gunakan `production` bukan `prod`, kecuali jika ada mapping jelas).

`--stage`

- Tujuan: memilih konfigurasi environment (lokal, development, staging, production).
- Pengaruh: menentukan file env yang akan dimuat, level logging default, koneksi ke resource (DB/queue), dan perilaku fitur.

`--mode`

- Tujuan: menentukan subsistem yang dijalankan oleh binary.
- Nilai yang umum: `http`, `grpc`, `worker`, `migrate`, `cli`.
- Perilaku: setiap mode hanya inisialisasi komponen yang diperlukan (mis. `migrate` hanya menyiapkan koneksi DB dan menjalankan migrasi).

Loader & precedence

- Precedence konfigurasi yang disarankan: flag CLI (`--stage`) > environment variables > file konfigurasi.
- Loader biasanya membentuk nama file berdasarkan nilai `--stage` (mis. `.env.stage.production`) dan memuatnya pada startup. 

Rekomendasi operasional

- Pastikan `Makefile` dan dokumentasi memakai nama stage yang sama dengan file `.env.stage.*`.
- Mode dan stage bersifat orthogonal: panggil `--mode` untuk menentukan subsistem dan `--stage` untuk menentukan konfigurasi environment.
- Jangan commit `.env.stage.production` yang berisi secret; gunakan secret manager untuk produksi.
	- Log level harus sesuai: `Info` untuk alur normal, `Error` untuk kegagalan.

7. Testing
	- Nama file test `xxx_test.go`; gunakan table-driven tests dan mock untuk dependencies.
	- Untuk unit tests, mock external resources (DB, network) dan hanya tes logika bisnis.

**Aturan Standarisasi Kode (Ketat)**

1. Tangani semua error secara eksplisit
	- Jangan mengabaikan error dengan menggunakan blank identifier `_` atau mengabaikannya tanpa alasan kuat.
	- Selalu propagasi atau wrap error menggunakan `fmt.Errorf("...: %w", err)` agar konteks tidak hilang.

2. Jangan gunakan `_` untuk mengabaikan hasil penting
	- Gunakan `_` hanya untuk mengabaikan nilai yang benar-benar tidak relevan. Jangan gunakan untuk mengabaikan error.

3. Gunakan `gofmt` / `go fmt` / `gofumpt` untuk formatting
	- Terapkan formatting otomatis sebagai pre-commit hook atau di CI. Gunakan `goimports` untuk merapikan imports.

4. Gunakan static analysis dan linters
	- Jalankan `go vet`, `staticcheck`, dan `golangci-lint` di CI. Anggap temuan linter sebagai aturan, bukan rekomendasi.

5. Hindari `panic` di library; gunakan pada aplikasi cli/main saja
	- Untuk paket yang bisa di-reuse, kembalikan error, jangan panggil `panic`.

6. Gunakan `context.Context` pada boundary transport/IO
	- Terima `context.Context` sebagai parameter pertama pada fungsi yang melakukan I/O, request handling, atau operasi panjang.

7. Batasi penggunaan variabel global yang mutable
	- Favor dependency injection dan konstruktor (`NewX`) untuk membuat instance; global hanya untuk konstanta atau konfigurasi immutable.

8. Timeout dan cancellation eksplisit
	- Untuk request jaringan atau DB, pastikan ada timeout atau gunakan context dengan deadline.

9. Dokumentasikan semua simbol yang diekspor
	- Gunakan komentar di atas tipe, fungsi, dan variabel yang diekspor agar `godoc` dan tim lain paham intent.

10. Hindari magic numbers dan string
	 - Tempatkan konfigurasi dan konstanta di satu file konstanta atau di package `configs`.

11. Desain untuk testabilitas
	 - Program harus mudah diuji: gunakan interface untuk dependency (repo, clients) dan sediakan cara untuk inject mock pada tests.

12. Praktek keamanan dasar
	 - Jangan commit secrets. Validasi input, lakukan escaping output, dan gunakan library yang dipelihara untuk enkripsi/crypto.

13. CI wajib: build, test, lint, vet
	 - Pastikan pipeline CI melakukan `go test -race`, linter, dan build untuk mencegah regresi.

14. Kode harus idiomatik Go
	 - Ikuti Effective Go dan Go Code Review Comments: nama pendek, error-first returns, tidak ada setter/getter yang berlebihan, gunakan slices/maps idiomatik.

15. Konsistensi penamaan
	 - Hindari underscore di identifier publik; gunakan `CamelCase` untuk exported names dan `camelCase` untuk private names.


Penjelasan Makefile

Makefile memudahkan developer menjalankan tugas umum. Target yang biasa ada:

- `make build` : compile binary (`go build ./...`).
- `make run` : jalankan aplikasi (`go run ./cmd/main`).
- `make test` : jalankan unit test (`go test ./...`).

Penjelasan `.env` file

.env menyimpan variabel lingkungan lokal (development). Contoh isian:

```
APP_ENV=development
PORT=8080
DB_DSN=user:pass@tcp(localhost:3306)/dbname
```

Jangan commit file `.env` berisi kredensial. Gunakan `.env.example` untuk contoh.

Ketika menjalankan aplikasi, flag `--stage` digunakan untuk memilih file yang sesuai. Contoh Makefile menjalankan application dengan flag:

```
./cmd/$(BINARY_NAME) --mode http --stage prod
```

Tips:

- Buat file `.env.example` atau `.env.stage.example` yang tidak berisi kredensial tapi menunjukkan variabel yang diperlukan.
- Jangan commit file `.env.stage.production` yang berisi secret; gunakan secret manager untuk produksi.
- Pastikan nilai `--stage` pada Makefile konsisten dengan nama file `.env.stage.*` yang Anda pakai.

Mode dan Stage

`--mode` dan `--stage` adalah flag runtime yang umum dipakai untuk memilih cara aplikasi berjalan dan konfigurasi environment.

1) `--mode`
	- Tujuan: menentukan subsistem yang dijalankan oleh binary.
	- Nilai yang umum:
		- `http` : jalankan HTTP server (biasanya Gin). Ekspos REST/HTTP API.
		- `grpc` : jalankan gRPC server (jika aplikasi mendukung gRPC).
		- `worker` : jalankan worker/consumer untuk background jobs (queue consumer, worker loop).
		- `migrate` : jalankan proses migrasi database dan keluar.
		- `cli` / `task` : jalankan tugas CLI tertentu atau runner task.
	- Perilaku: setiap mode hanya menginisialisasi komponen yang diperlukan untuk fungsinya. Misalnya `migrate` hanya menyiapkan koneksi DB dan menjalankan migrasi, tanpa memulai server.

2) `--stage`
	- Tujuan: memilih konfigurasi environment (lokal, development, staging, production).
	- Contoh nilai: `local`, `development`, `staging`, `production` (atau singkatan jika Anda konsisten, mis. `prod`).
	- Pengaruh: penentuan file env yang akan dimuat (`.env.stage.<stage>`), level logging default, koneksi ke resource berbeda (DB, queue), dan perilaku fitur (mis. feature flags).

3) Kombinasi & precedence
	- Mode menentukan subsistem; stage menentukan konfigurasi untuk subsistem tersebut.
	- Precedence konfigurasi: flag CLI (`--stage`) > environment variables > file konfigurasi. Pastikan loader config Anda menghormati urutan ini.

4) Praktik dan rekomendasi
	- Gunakan nama stage yang konsisten dan deskriptif (`production` lebih disarankan daripada `prod` kecuali Anda mendefinisikan mapping jelas).
	- Pastikan Makefile dan dokumentasi memakai nama stage yang sama dengan file `.env.stage.*` yang ada di repo.
	- Mode harus bersifat orthogonal terhadap stage: Anda bisa menjalankan `--mode http --stage development` maupun `--mode worker --stage production`.

5) Contoh perintah (dokumentasi — bukan contoh kode implementasi)

```bash
# jalankan HTTP server pada stage production
./cmd/main --mode=http --stage=production

# jalankan migrasi pada stage staging
./cmd/main --mode=migrate --stage=staging

# jalankan worker di environment development
./cmd/main --mode=worker --stage=development
```

Catatan: pastikan loader environment Anda mencari file `.env.stage.<stage>` atau memetakan singkatan stage ke nama file yang sesuai.

Penjelasan config

Config didefinisikan sebagai struct (mis. `configs.Config`) dan di-load dari file + environment variables (prioritas env > file).

Contoh singkat:

```go
type Config struct {
	 AppName string
	 Port    int
	 DBDSN   string
}
```


Langkah singkat untuk menjalankan aplikasi lokal:

1. Copy file environment contoh dan sesuaikan:

```bash
cp .env.example .env.stage.local
# edit .env sesuai kebutuhan
```

2. Install dependensi dan build (Go modules):

```bash
go mod tidy
make build
```

3. Jalankan aplikasi:

```bash
make run
# atau
go run ./cmd/main --mode=http --stage=local
```

Menjalankan unit test:

```bash
make test
# atau
go test ./... -v -coverprofile=coverage.out
```


