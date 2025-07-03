package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type Aksi struct {
	Tujuan string `json:"tujuan"`
	Data   string `json:"data"`
}

func TarikDataPencarianPria(w http.ResponseWriter, r *http.Request, db *gorm.DB, Data string) {
	fmt.Println(Data, "ini sianying")

	var items []BarangCustom

	clean := strings.Trim(Data, `"`)
	// Query ke database
	if err := db.Table("barang_custom").Where("nama = ?", clean).Find(&items).Error; err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Ketemu")
	}

	// Format data ke dalam slice of map[string]interface{}, tanpa kolom "file"
	var dataProduk []map[string]interface{}
	for _, item := range items {
		dataProduk = append(dataProduk, map[string]interface{}{
			"nama":          item.Nama,
			"jenis_pakaian": item.JenisPakaian,
			"deskripsi":     item.Deskripsi,
			"kategori":      item.Kategori,
			"harga":         item.Harga,
			"warna":         item.Warna,
			"stok":          item.Stok,
			"ukuran":        item.Ukuran,
		})
	}

	// Set response header dan encode ke JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataProduk)

	fmt.Println("Data barang tanpa file berhasil dikirim ke frontend")

}

func FotoProdukPencarian(w http.ResponseWriter, data json.RawMessage, db *gorm.DB) {
	// Decode JSON request ke struct RequestData
	var requestData RequestData
	if err := json.Unmarshal(data, &requestData); err != nil {
		http.Error(w, "Gagal membaca data", http.StatusBadRequest)

	}

	// Debugging: Cetak nilai dari request
	fmt.Println("Nama:", requestData.Nama)
	fmt.Println("Jenis Pakaian:", requestData.JenisPakaian)

	// Query langsung dari tabel barang_custom
	var barang struct {
		File []byte `gorm:"column:file"`
	}
	if err := db.Table("barang_custom").Select("file").Where("nama = ? AND jenis_pakaian = ?", requestData.Nama, requestData.JenisPakaian).First(&barang).Error; err != nil {
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)

	}

	// Kirim file gambar dalam bentuk biner
	w.Header().Set("Content-Disposition", "attachment; filename=foto_produk.jpg")
	w.Write(barang.File) // Mengirimkan data gambar dari database
}
