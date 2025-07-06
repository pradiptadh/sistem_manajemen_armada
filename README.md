# 🚚 Sistem Manajemen Armada

Aplikasi ini merupakan sistem backend manajemen armada menggunakan Golang, GORM, MQTT, RabbitMQ, dan PostgreSQL. Aplikasi dapat:

- Menerima data lokasi kendaraan via MQTT.
- Menyimpan lokasi ke PostgreSQL.
- Memeriksa geofence dan mengirim event ke RabbitMQ.
- Menyediakan API untuk mengecek lokasi terakhir dan histori kendaraan.
- Worker service yang membaca event geofence dari RabbitMQ.

---

## 📦 Struktur Folder

```
├── db/                    # Koneksi database
├── handler/               # HTTP handler
├── model/                 # Struct model (VehicleLocation, GeofenceEvent)
├── mqtt/                  # Subscriber MQTT
├── publisher/             # Publisher mock lokasi kendaraan ke MQTT
├── repository/            # Interface + implementasi PostgreSQL (GORM)
├── worker/                # Worker baca event dari RabbitMQ
├── Dockerfile             # Untuk service utama
├── Dockerfile.publisher   # Untuk publisher (mock data)
├── Dockerfile.worker      # Untuk worker
├── docker-compose.yml     # Semua service berjalan di sini
└── main.go                # Entry point utama
```

---

## 🚀 Cara Menjalankan

### 1. Clone Repo

```bash
git clone https://github.com/pradiptadh/sistem_manajemen_armada.git
cd sistem_manajemen_armada
```

### 2. Jalankan dengan Docker Compose

```bash
docker compose up --build
```

Akan menjalankan:

- PostgreSQL (port 5434)
- RabbitMQ (port 5672 / 15672)
- Mosquitto MQTT (port 1883)
- Service utama (`sistem_manajemen_armada`)
- Publisher mock lokasi kendaraan (`publisher`)
- Worker geofence (`worker`)

---

## 🔌 Endpoint API

- `GET /vehicles/:vehicle_id/location` → Ambil lokasi terakhir
- `GET /vehicles/:vehicle_id/history` → Ambil histori lokasi kendaraan

Contoh:
```bash
curl http://localhost:8080/vehicles/B1234XYZ/location
```

---

## 📡 Format Pesan MQTT

```json
{
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715003456
}
```

---

## 📨 Format Event Geofence ke RabbitMQ

- **Exchange**: `fleet.events`
- **Queue**: `geofence_alerts`

```json
{
  "vehicle_id": "B1234XYZ",
  "event": "geofence_entry",
  "location": {
    "latitude": -6.2088,
    "longitude": 106.8456
  },
  "timestamp": 1715003456
}
```

---

## ⚙️ Environment Variable

Tersimpan di `docker-compose.yml`

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`
- `MQTT_BROKER`, `MQTT_PORT`
- `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_USER`, `RABBITMQ_PASS`

---

## 👨‍💻 Author

**Pradipta Dwi Haryadi**  
GitHub: [@pradiptadh](https://github.com/pradiptadh)

---
