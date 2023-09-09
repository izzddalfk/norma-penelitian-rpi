# Report Benchmark

## Go version: Clustering

Durasi pengujian: 5 menit

Hasil benchmark bisa dibaca [disini](./clustering_go.json).

Spesifikasi Cluster:

- Orchestration: microk8s (Kubernetes v1.27.5)
- Virtual Machine (Control Plane / Master)
  - CPU 2 Core@3.4GHz
  - RAM 2GB
  - Ubuntu Server 20.04.1
- Raspberry Pi CM4 (Worker)
  - CPU 4 Core@1.5GHz
  - RAM 4GB
  - eMMC: 16GB
  - Ubuntu Server 20.04.5
- Raspberry Pi 4B (Worker)
  - CPU 4 Core@1.5GHz
  - RAM: 8GB
  - sdCard: 64GB Class 10
  - Ubuntu Server 22.04.3

Koneksi menggunakan WiFi

Komputer yang dipakai untuk mengirimkan request benchmark, MacBook Pro 2021 M1 Pro.
