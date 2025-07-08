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

func ambilSepatu(w http.ResponseWriter, db *gorm.DB, untuk string) []map[string]interface{} {
	var sepatu []BarangCustom
	var dataSepatu []map[string]interface{}

	log.Println("Fungsi ambilSepatu dijalankan dengan parameter untuk =", untuk)

	if untuk == "Pria" || untuk == "Wanita" {
		log.Println("Mencoba mengambil data sepatu untuk:", untuk)

		err := db.Table("barang_custom").
			Where(`EXISTS (
				SELECT 1 FROM "Sepatu"
				WHERE "Sepatu"."nama_sepatu" = "barang_custom"."nama"
				AND "Sepatu"."jenis_sepatu" = "barang_custom"."jenis_pakaian"
				AND "Sepatu"."untuk" = ?
			)`, untuk).
			Find(&sepatu).Error

		if err != nil {
			log.Println("Error mengambil data barang_custom (sepatu pria/wanita):", err)
		} else {
			log.Printf("Berhasil mengambil %d data sepatu untuk %s\n", len(sepatu), untuk)
			for i, item := range sepatu {
				log.Printf("Data ke-%d: %+v\n", i+1, item)
				dataSepatu = append(dataSepatu, map[string]interface{}{
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
		log.Println("Mencoba mengambil semua data sepatu")

		err := db.Table("barang_custom").
			Where(`EXISTS (
				SELECT 1 FROM "Sepatu"
				WHERE "Sepatu"."nama_sepatu" = "barang_custom"."nama"
				AND "Sepatu"."jenis_sepatu" = "barang_custom"."jenis_pakaian"
			)`).
			Find(&sepatu).Error

		if err != nil {
			log.Println("Error mengambil semua data sepatu:", err)
		} else {
			log.Printf("Berhasil mengambil %d data semua sepatu\n", len(sepatu))
			for i, item := range sepatu {
				log.Printf("Data ke-%d: %+v\n", i+1, item)
				dataSepatu = append(dataSepatu, map[string]interface{}{
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

	} else {
		log.Println("Parameter 'untuk' tidak valid:", untuk)
	}

	log.Println("Fungsi ambilSepatu selesai dijalankan")

	return dataSepatu
}

func ambilKacamata(w http.ResponseWriter, db *gorm.DB, untuk string) []map[string]interface{} {
	var kacamata []BarangCustom
	var dataKacamata []map[string]interface{}

	if untuk == "Pria" || untuk == "Wanita" {
		err := db.Table("barang_custom").
			Where(`EXISTS (
				SELECT 1 FROM "Kacamata"
				WHERE "Kacamata"."Nama_Kacamata" = "barang_custom"."nama"
				AND "Kacamata"."Jenis_Kacamata" = "barang_custom"."jenis_pakaian"
				AND "Kacamata"."untuk" = ?
			)`, untuk).
			Find(&kacamata).Error

		if err != nil {
			log.Println("Error mengambil data barang_custom (kacamata pria/wanita):", err)
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

	} else if untuk == "Semua" {
		err := db.Table("barang_custom").
			Where(`EXISTS (
				SELECT 1 FROM "Kacamata"
				WHERE "Kacamata"."Nama_Kacamata" = "barang_custom"."nama"
				AND "Kacamata"."Jenis_Kacamata" = "barang_custom"."jenis_pakaian"
			)`).
			Find(&kacamata).Error

		if err != nil {
			log.Println("Error mengambil semua data kacamata:", err)
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

	dataProduk = append(dataProduk, ambilSepatu(w, db, "Semua")...)

	dataProduk = append(dataProduk, AmbilBaju(w, db, "Semua")...)

	dataProduk = append(dataProduk, AmbilCelana(w, db, "Semua")...)

	dataProduk = append(dataProduk, ambilKacamata(w, db, "Semua")...)

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
	fmt.Println("===> [FotoProduk] Memulai proses pengambilan foto produk")

	// Step 1: Decode JSON request ke struct RequestData
	var requestData RequestData
	fmt.Println("[Step 1] Mencoba decode JSON...")
	if err := json.Unmarshal(data, &requestData); err != nil {
		errMsg := fmt.Sprintf("❌ [Step 1] Gagal membaca data JSON: %v", err)
		fmt.Println(errMsg)
		http.Error(w, "Gagal membaca data", http.StatusBadRequest)
		return
	}
	fmt.Printf("✅ [Step 1] JSON berhasil dibaca: Nama='%s', JenisPakaian='%s'\n", requestData.Nama, requestData.JenisPakaian)

	// Step 2: Query data gambar dari tabel barang_custom
	var barang struct {
		File []byte `gorm:"column:file"`
	}
	fmt.Printf("[Step 2] Mencoba mengambil data gambar dari database untuk: Nama='%s', JenisPakaian='%s'\n", requestData.Nama, requestData.JenisPakaian)
	if err := db.Table("barang_custom").
		Select("file").
		Where("nama = ? AND jenis_pakaian = ?", requestData.Nama, requestData.JenisPakaian).
		First(&barang).Error; err != nil {
		errMsg := fmt.Sprintf("❌ [Step 2] Data tidak ditemukan di database: %v", err)
		fmt.Println(errMsg)
		http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
		return
	}
	fmt.Println("✅ [Step 2] Data gambar berhasil ditemukan di database.")

	// Step 3: Kirim gambar ke client dalam bentuk biner
	fmt.Println("[Step 3] Mengirim gambar ke client...")
	w.Header().Set("Content-Disposition", "attachment; filename=foto_produk.jpg")
	w.Header().Set("Content-Type", "image/jpeg") // Bisa juga "application/octet-stream" jika format file tidak pasti
	if _, err := w.Write(barang.File); err != nil {
		errMsg := fmt.Sprintf("❌ [Step 3] Gagal mengirim file ke response: %v", err)
		fmt.Println(errMsg)
	} else {
		fmt.Println("✅ [Step 3] File gambar berhasil dikirim ke client.")
	}

	fmt.Println("✅ [FotoProduk] Proses selesai tanpa error.")
}
