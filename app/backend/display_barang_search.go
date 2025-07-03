package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type BarangCustom struct {
	Nama         string `json:"nama"`
	File         string `json:"file"`
	JenisPakaian string `json:"jenis_pakaian"`
	Deskripsi    string `json:"deskripsi"`
	Kategori     string `json:"kategori"`
	Harga        int    `json:"harga"`
	Warna        string `json:"warna"`
	Stok         int    `json:"stok"`
	Ukuran       string `json:"ukuran"`
}

type BarangDituju struct {
	Nama      string `json:"Nama_baju"`
	Id        int    `json:"id"`
	KodeBaju  int8   `json:"Kode_baju"`
	HargaBaju string `json:"Harga_baju"`
	JenisBaju string `json:"Jenis_baju"`
	Warna     string `json:"warna"`
	Untuk     string `json:"untuk"`
	status    string `json:"status"`
}
type RequestData struct {
	Nama         string `json:"nama"`
	JenisPakaian string `json:"jenis_pakaian"`
}

type FotoProdukM struct {
	ID           uint   `gorm:"primaryKey"`
	Nama         string `gorm:"column:nama"`
	JenisPakaian string `gorm:"column:jenis_pakaian"`
	File         []byte `gorm:"column:file"`
}

func AmbilBaju(w http.ResponseWriter, db *gorm.DB, untuk string) []map[string]interface{} {
	var baju []BarangCustom
	var dataBaju []map[string]interface{}

	if untuk == "Pria" || untuk == "Wanita" {
		if err := db.Table("barang_custom").Find(&baju).Error; err != nil {
			http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		}

		err := db.Table("barang_custom").
			Where(`EXISTS (
			SELECT 1 FROM "Baju"
			WHERE "Baju"."Nama_baju" = "barang_custom"."nama"
			AND "Baju"."Jenis_baju" = "barang_custom"."jenis_pakaian"
			AND "Baju"."untuk" = ?
		)`, untuk).
			Find(&baju).Error

		if err != nil {
			log.Println("Error mengambil data barang_custom:", err)
		} else {
			for _, item := range baju {
				dataBaju = append(dataBaju, map[string]interface{}{
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
		}

	} else if untuk == "Semua" {
		if err := db.Table("barang_custom").Find(&baju).Error; err != nil {
			http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		}

		err := db.Table("barang_custom").
			Where(`EXISTS (
				SELECT 1 FROM "Baju"
				WHERE "Baju"."Nama_baju" = "barang_custom"."nama"
				AND "Baju"."Jenis_baju" = "barang_custom"."jenis_pakaian"
			)`).
			Find(&baju).Error

		if err != nil {
			log.Println("Error mengambil data barang_custom:", err)
		} else {
			for _, item := range baju {
				dataBaju = append(dataBaju, map[string]interface{}{
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
		}
	}

	return dataBaju

}

func AmbilCelana(w http.ResponseWriter, db *gorm.DB, untuk string) []map[string]interface{} {
	var Celana []BarangCustom

	if err := db.Table("barang_custom").Find(&Celana).Error; err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
	}

	var dataCelana []map[string]interface{}

	if untuk == "Pria" || untuk == "Wanita" {

		errcelana := db.Table("barang_custom").
			Where(`EXISTS (
		SELECT 1 FROM "Celana" 
		WHERE "Celana"."Nama_Celana" = "barang_custom"."nama" 
		AND "Celana"."Jenis_Celana" = "barang_custom"."jenis_pakaian"
          AND "Celana"."untuk" = ?
    	)`, untuk).Find(&Celana).Error

		if errcelana != nil {
			log.Println("Error mengambil data barang_custom:", errcelana)
		} else {
			for _, item := range Celana {
				dataCelana = append(dataCelana, map[string]interface{}{
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
		}

	} else if untuk == "Semua" {
		errcelana := db.Table("barang_custom").
			Where(`EXISTS (
			SELECT 1 FROM "Celana" 
			WHERE "Celana"."Nama_Celana" = "barang_custom"."nama" 
			AND "Celana"."Jenis_Celana" = "barang_custom"."jenis_pakaian"
		)`).Find(&Celana).Error

		if errcelana != nil {
			log.Println("Error mengambil data barang_custom:", errcelana)
		} else {
			for _, item := range Celana {
				dataCelana = append(dataCelana, map[string]interface{}{
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
		}
	}

	return dataCelana
}

func ambilKacamata(w http.ResponseWriter, db *gorm.DB) []map[string]interface{} {
	var kacamata []BarangCustom

	if err := db.Table("barang_custom").Find(&kacamata).Error; err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
	}

	var dataKacamata []map[string]interface{}

	errkacamata := db.Table("barang_custom").
		Where(`EXISTS (
		SELECT 1 FROM "Kacamata" 
		WHERE "Kacamata"."Nama_Kacamata" = "barang_custom"."nama" 
		AND "Kacamata"."Jenis_Kacamata" = "barang_custom"."jenis_pakaian"
    )`).Find(&kacamata).Error

	if errkacamata != nil {
		log.Println("Error mengambil data barang_custom:", errkacamata)
	} else {
		for _, item := range kacamata {
			dataKacamata = append(dataKacamata, map[string]interface{}{
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
	}

	return dataKacamata
}

// Fungsi untuk mengambil semua data barang_custom dan mengirimnya sebagai JSON tanpa kolom file
func TarikDatSemuaPria(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var items []BarangCustom

	// Query ke database
	if err := db.Table("barang_custom").Find(&items).Error; err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		return
	}

	var dataProduk []map[string]interface{}

	dataProduk = append(dataProduk, AmbilBaju(w, db, "Pria")...)

	dataProduk = append(dataProduk, AmbilCelana(w, db, "Pria")...)

	dataProduk = append(dataProduk, ambilKacamata(w, db)...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataProduk)

	fmt.Println("Data barang tanpa file berhasil dikirim ke frontend")
}

func TarikDatSemuaWanita(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var items []BarangCustom

	// Query ke database
	if err := db.Table("barang_custom").Find(&items).Error; err != nil {
		http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		return
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

	fmt.Println(dataProduk)

	var didapat []BarangDituju // struct Baju harus didefinisikan sesuai kolom di DB
	for _, data := range dataProduk {
		var baju BarangDituju
		if err := db.Table("Baju").
			Where("Nama_baju = ? AND Kode_baju = ? AND untuk = ?", data["nama"], data["kode_baju"], "Pria").
			First(&baju).Error; err == nil {
			didapat = append(didapat, baju)
			// Jika tidak ditemukan, abaikan errornya (atau log jika perlu)
		}
	}

	// Set response header dan encode ke JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataProduk)

	fmt.Println("Data barang tanpa file berhasil dikirim ke frontend")
}

func FotoProduk(w http.ResponseWriter, data json.RawMessage, db *gorm.DB) {
	// Decode JSON request ke struct RequestData
	var requestData RequestData
	if err := json.Unmarshal(data, &requestData); err != nil {
		errMsg := fmt.Sprintf("❌ Gagal membaca data JSON: %v", err)
		fmt.Println(errMsg)
		http.Error(w, "Gagal membaca data", http.StatusBadRequest)
		return
	}
	fmt.Printf("✅ JSON berhasil dibaca: Nama=%s, JenisPakaian=%s\n", requestData.Nama, requestData.JenisPakaian)

	// Query langsung dari tabel barang_custom
	var barang struct {
		File []byte `gorm:"column:file"`
	}
	if err := db.Table("barang_custom").Select("file").
		Where("nama = ? AND jenis_pakaian = ?", requestData.Nama, requestData.JenisPakaian).
		First(&barang).Error; err != nil {
		errMsg := fmt.Sprintf("❌ Data tidak ditemukan untuk nama '%s' dan jenis_pakaian '%s': %v",
			requestData.Nama, requestData.JenisPakaian, err)
		fmt.Println(errMsg)
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}
	fmt.Printf("✅ Data berhasil ditemukan dan diambil dari database untuk %s / %s\n",
		requestData.Nama, requestData.JenisPakaian)

	// Kirim file gambar dalam bentuk biner
	w.Header().Set("Content-Disposition", "attachment; filename=foto_produk.jpg")
	if _, err := w.Write(barang.File); err != nil {
		errMsg := fmt.Sprintf("❌ Gagal mengirim file ke response: %v", err)
		fmt.Println(errMsg)
	} else {
		fmt.Println("✅ File gambar berhasil dikirim ke client.")
	}
}
