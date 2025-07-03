package backendadmin

import (
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type BarangCustom struct {
	ID            uint   `gorm:"primaryKey;column:id"`
	Jenis_pakaian string `gorm:"column:jenis_pakaian;not null"`
	Deskripsi     string `gorm:"column:deskripsi"`
	Kategori      string `gorm:"column:kategori"`
	Harga         int32  `gorm:"column:harga"`
	Warna         string `gorm:"column:warna"`
	Stok          int    `gorm:"column:stok"`
	Ukuran        string `gorm:"column:ukuran"`
	Nama          string `gorm:"column:nama"`
	Gender        string `gorm:"column:Gender"`
	X             int    `gorm:"column:X"`
	L             int    `gorm:"column:L"`
	M             int    `gorm:"column:M"`
	XL            int    `gorm:"column:XL"`
	S             int    `gorm:"column:S"`
}

func AmbilDataShort(w http.ResponseWriter, db *gorm.DB, Jenis string) []map[string]interface{} {
	var barang []BarangCustom

	if Jenis == "Semua" {
		db.Table("barang_custom").Find(&barang)
	} else {
		// Jika Jenis tidak sama dengan "Semua", Anda mungkin ingin menambahkan logika untuk mengambil data berdasarkan jenis tertentu.
		// Contoh:
		fmt.Println(Jenis, "ini jenis yang dicari")
		db.Table("barang_custom").Where("kategori = ?", Jenis).Find(&barang)
	}

	var result []map[string]interface{}
	for _, b := range barang {
		result = append(result, map[string]interface{}{
			"id_barang":     b.ID,
			"jenis_pakaian": b.Jenis_pakaian,
			"deskripsi":     b.Deskripsi,
			"kategori":      b.Kategori,
			"harga":         b.Harga,
			"warna":         b.Warna,
			"stok":          b.Stok,
			"ukuran":        b.Ukuran,
			"nama":          b.Nama,
			"gender":        b.Gender,
			"x":             b.X,
			"l":             b.L,
			"m":             b.M,
			"xl":            b.XL,
			"s":             b.S,
		})
	}

	result = append(result, map[string]interface{}{
		"Jumlah_Barang": len(barang),
	})

	return result
}

func AmbilDataUiShort(w http.ResponseWriter, db *gorm.DB) map[string]int {
	var totalBaju int64
	var totalCelana int64
	var totalKacamata int64
	var totalTas int64
	var totalSepatu int64
	var totalBarangCustom int64

	// Menghitung total barang per kategori
	if err := db.Table("barang_custom").Where("kategori = ?", "Baju").Count(&totalBaju).Error; err != nil {
		fmt.Println("Gagal menghitung total baju:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Baju: %d\n", totalBaju) // Menampilkan hasil

	if err := db.Table("barang_custom").Where("kategori = ?", "Celana").Count(&totalCelana).Error; err != nil {
		fmt.Println("Gagal menghitung total celana:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Celana: %d\n", totalCelana) // Menampilkan hasil

	if err := db.Table("barang_custom").Where("kategori = ?", "Kacamata").Count(&totalKacamata).Error; err != nil {
		fmt.Println("Gagal menghitung total kacamata:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Kacamata: %d\n", totalKacamata) // Menampilkan hasil

	if err := db.Table("barang_custom").Where("kategori = ?", "Tas").Count(&totalTas).Error; err != nil {
		fmt.Println("Gagal menghitung total tas:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Tas: %d\n", totalTas) // Menampilkan hasil

	if err := db.Table("barang_custom").Where("kategori = ?", "Sepatu").Count(&totalSepatu).Error; err != nil {
		fmt.Println("Gagal menghitung total sepatu:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Sepatu: %d\n", totalSepatu) // Menampilkan hasil

	// Menghitung total barang custom
	if err := db.Table("barang_custom").Where("status = ?", "custom").Count(&totalBarangCustom).Error; err != nil {
		fmt.Println("Gagal menghitung total barang custom:", err) // Menambahkan detail error
	}
	fmt.Printf("Total Barang Custom: %d\n", totalBarangCustom) // Menampilkan hasil

	// Membuat map untuk menyimpan hasil
	data := map[string]int{
		"total_baju":          int(totalBaju),
		"total_celana":        int(totalCelana),
		"total_kacamata":      int(totalKacamata),
		"total_tas":           int(totalTas),
		"total_sepatu":        int(totalSepatu),
		"total_barang_custom": int(totalBarangCustom),
	}

	fmt.Printf("Data yang dikembalikan: %+v\n", data) // Menampilkan data yang akan dikembalikan

	return data
}

func HapusBarang(w http.ResponseWriter, db *gorm.DB, id_barang, nama, jenis_pakaian, gender, harga string) map[string]interface{} {
	// Logging awal
	fmt.Println("Menghapus Barang Dijalankan")

	// Konversi harga ke integer
	hargaInt, err := strconv.Atoi(harga)
	if err != nil {
		return map[string]interface{}{
			"status":  false,
			"message": "Harga tidak valid (bukan angka)",
		}
	}

	// Hapus data dari tabel 'barang_custom' jika semua kolom cocok
	result := db.Table("barang_custom").
		Where(`"id" = ? AND "nama" = ? AND "jenis_pakaian" = ? AND "Gender" = ? AND "harga" = ?`,
			id_barang, nama, jenis_pakaian, gender, hargaInt).
		Delete(nil)

	if result.Error != nil {
		return map[string]interface{}{
			"status":  false,
			"message": "Gagal menghapus data: " + result.Error.Error(),
		}
	}

	if result.RowsAffected == 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "Tidak ada data yang cocok untuk dihapus",
		}
	}

	return map[string]interface{}{
		"status":  true,
		"message": fmt.Sprintf("Berhasil menghapus %d data", result.RowsAffected),
	}
}
