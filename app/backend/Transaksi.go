package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

)

type BarangJson struct {
	NamaBarang string `json:"nama"`
	Ukuran     string `json:"ukuran"`
	Jumlah     string `json:"jumlah"`
	Harga      string `json:"harga"`
}

type Transaksi struct {
	ID               int32     `gorm:"column:id"`
	IDUser           int32     `gorm:"column:id_user"`
	NamaPembeli      string    `gorm:"column:Nama_pembeli"`
	AlamatLengkap    string    `gorm:"column:Alamat_lengkap"`
	KoordinatAlamat  string    `gorm:"column:Koordinat_Alamat"`
	NamaItem         string    `gorm:"column:Nama_Item"`
	KodeBarang       string    `gorm:"column:kode_barang"` // jsonb tidak ada tipe data langsung di Go, gunakan string
	Harga            float64   `gorm:"column:harga"`
	Jumlah           int32     `gorm:"column:jumlah"`
	KodePemesanan    string    `gorm:"column:kode_pemesanan"`
	MetodePembayaran string    `gorm:"column:Metode_pembayaran"`
	BuktiPembayaran  string    `gorm:"column:Bukti_pembayaran"`
	StatusPengiriman string    `gorm:"column:Status_pengiriman"`
	PromoAtauDiskon  string    `gorm:"column:Promo_atau_diskon"`
	TanggalDanWaktu  time.Time `gorm:"column:Tanggal_dan_waktu"`
}

func VerifikasiTransaksi(
	db *gorm.DB,
	Barangnya json.RawMessage,
	kordinat string,
	AlamatLengkap interface{},
	id_user string,
	MetodePembayaran string,
) map[string]string {

	var user User
	var kategori BarangCustom
	var daftarBarang []BarangJson

	if err := json.Unmarshal(Barangnya, &daftarBarang); err != nil {
		log.Println("Gagal decode daftar barang:", err)
		return map[string]string{
			"Status": "Gagal",
			"Pesan":  "Data daftar barang tidak valid",
		}
	}

	// Ambil nama user dari ID
	err := db.First(&user, "id = ?", id_user).Error
	if err != nil {
		log.Println("Error saat query nama user:", err)
		return map[string]string{
			"Status": "Gagal",
			"Pesan":  "Terjadi kesalahan saat mengambil nama user",
		}
	}

	if user.Nama == "" {
		return map[string]string{
			"Status": "Gagal",
			"Pesan":  "User tidak ditemukan atau nama kosong",
		}
	}

	for _, barang := range daftarBarang {
		nama_Barang := barang.NamaBarang
		ukuran := barang.Ukuran
		jumlahbarang, _ := strconv.Atoi(barang.Jumlah)

		err1 := db.Table("barang_custom").First(&kategori, "nama = ?", nama_Barang).Error
		if err1 != nil {
			log.Println("Error saat query ambil kategori dari tabel 'barang_custom':", err1)
			return map[string]string{
				"Status": "Gagal",
				"Pesan":  "Terjadi kesalahan saat mengambil kategori dari tabel barang_custom",
			}
		}

		if kategori.Kategori == "" {
			return map[string]string{
				"Status": "Gagal",
				"Pesan":  "Kategori Pakaian tidak ditemukan",
			}
		}

		fmt.Println(kategori.Kategori)

		var KodePakaian string
		var NamaPakaianKolom string

		switch kategori.Kategori {
		case "Baju":
			KodePakaian = `"Kode_baju"`
			NamaPakaianKolom = `"Nama_baju"`
		case "Celana":
			KodePakaian = `"Kode_Celana"`
			NamaPakaianKolom = `"Nama_Celana"`
		case "Sepatu":
			KodePakaian = `"kode_sepatu"`
			NamaPakaianKolom = `"nama_sepatu"`
		case "Kacamata":
			KodePakaian = `"Kode_Kacamata"`
			NamaPakaianKolom = `"Nama_Kacamata"`
		}

		var hasilInt []int64
		whereClause := fmt.Sprintf(`status != ? AND status IS NOT NULL AND status != '' AND "Ukuran" = ? AND %s = ?`, NamaPakaianKolom)
		err2 := db.
			Table(kategori.Kategori).
			Select(KodePakaian).
			Where(whereClause, "sudah dibeli", ukuran, nama_Barang).
			Order("RANDOM()").
			Limit(jumlahbarang).
			Scan(&hasilInt).Error

		if err2 != nil {
			log.Println("Gagal ambil data acak dari tabel:", kategori.Kategori, "Error:", err2)
			continue // lanjut ke barang berikutnya jika gagal ambil data
		}

		// Konversi []int64 → []string
		hasilStr := make([]string, len(hasilInt))
		for i, v := range hasilInt {
			hasilStr[i] = strconv.FormatInt(v, 10)
		}

		formatHasil := "{" + strings.Join(hasilStr, ":") + "}"
		fmt.Println("Format hasil:", formatHasil)
	}

	return map[string]string{
		"Status": "Berhasil",
		"Pesan":  "Transaksi diverifikasi",
		"User":   user.Nama,
	}
}
