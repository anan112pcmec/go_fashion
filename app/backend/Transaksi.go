package backend

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func VerifikasiTransaksi(
	db *gorm.DB,
	nama_Barang string,
	kordinat string,
	AlamatLengkap interface{},
	ukuran string,
	jumlah string,
	harga string,
	id_user string,
) map[string]string {

	var nama string

	// Query langsung nama user dengan kolom Nama (huruf kapital N)

	fmt.Println(nama_Barang, kordinat, AlamatLengkap, ukuran, jumlah, harga, id_user)
	err := db.Raw(`SELECT "Nama" FROM user WHERE id = ?`, id_user).Scan(&nama).Error
	if err != nil {
		log.Println("Error saat query nama user:", err)
		return map[string]string{
			"Status": "Gagal",
			"Pesan":  "Terjadi kesalahan saat mengambil nama user",
		}
	}

	// Kalau nama kosong berarti user tidak ditemukan
	if nama == "" {
		return map[string]string{
			"Status": "Gagal",
			"Pesan":  "User tidak ditemukan",
		}
	}

	// Berhasil verifikasi
	return map[string]string{
		"Status": "Berhasil",
		"Pesan":  "Transaksi diverifikasi",
		"User":   nama,
		"Barang": nama_Barang,
	}
}
