package backendadmin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	MasukPada     string `gorm:"column:Masuk_Pada"`
	TerjualS      string `gorm:"TerjualS"`
	TerjualM      string `gorm:"TerjualM"`
	TerjualL      string `gorm:"TerjualL"`
	TerjualX      string `gorm:"TerjualX"`
	TerjualXL     string `gorm:"TerjualXL"`
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

func AmbilDataStok(w http.ResponseWriter, db *gorm.DB, nama, gender string) []map[string]interface{} {
	var hasil []BarangCustom

	// Struct hasil akhir, disesuaikan dengan alias SQL
	type Barang struct {
		ID      int32  `json:"id"`
		Kode    int32  `json:"kode"`
		Nama    string `json:"nama"`
		Harga   int32  `json:"harga"`
		Jenis   string `json:"jenis"`
		Warna   string `json:"Warna"`
		Untuk   string `json:"untuk"`
		IdInduk int32  `json:"id_induk"`
		Status  string `json:"status"`
		Ukuran  string `json:"Ukuran"`
	}

	err := db.Table("barang_custom").Where(`"nama" = ? AND "Gender" = ?`, nama, gender).Find(&hasil).Error
	if err != nil {
		log.Printf("Gagal mengambil data dari tabel barang_custom: %v\n", err)
		http.Error(w, "Gagal mengambil data stok dari database", http.StatusInternalServerError)
		return nil
	}

	if len(hasil) == 0 {
		log.Printf("Tidak ditemukan data barang_custom untuk nama: %s dan gender: %s\n", nama, gender)
		return nil
	}

	var finalData []map[string]interface{}

	for _, item := range hasil {
		var data []Barang

		var table, kodeAlias, namaAlias, hargaAlias, jenisAlias string

		switch strings.ToLower(item.Kategori) {
		case "baju":
			table = `"Baju"`
			kodeAlias = `"Kode_baju"`
			namaAlias = `"Nama_baju"`
			hargaAlias = `"Harga_baju"`
			jenisAlias = `"Jenis_baju"`
		case "celana":
			table = `"Celana"`
			kodeAlias = `"Kode_Celana"`
			namaAlias = `"Nama_Celana"`
			hargaAlias = `"Harga_Celana"`
			jenisAlias = `"Jenis_Celana"`
		case "sepatu":
			table = `"Sepatu"`
			kodeAlias = `"kode_sepatu"`
			namaAlias = `"nama_sepatu"`
			hargaAlias = `"harga_sepatu"`
			jenisAlias = `"jenis_sepatu"`
		default:
			log.Printf("Kategori tidak dikenal: %s, kode: %q\n", item.Kategori, item.ID)
			continue
		}

		query := fmt.Sprintf(`
			SELECT 
				id, 
				%s AS kode, 
				%s AS nama, 
				%s AS harga, 
				%s AS jenis, 
				"Warna", 
				untuk, 
				id_induk, 
				status, 
				"Ukuran" 
			FROM %s 
			WHERE id_induk = ?`, kodeAlias, namaAlias, hargaAlias, jenisAlias, table)

		if err := db.Raw(query, item.ID).Scan(&data).Error; err != nil {
			log.Printf("Gagal mengambil detail untuk id_induk %d dari tabel %s, error: %v\n", item.ID, table, err)
			continue
		}

		// Ubah ke map
		for _, b := range data {
			row := map[string]interface{}{
				"id":       b.ID,
				"kode":     b.Kode,
				"nama":     b.Nama,
				"harga":    b.Harga,
				"jenis":    b.Jenis,
				"warna":    b.Warna,
				"untuk":    b.Untuk,
				"id_induk": b.IdInduk,
				"status":   b.Status,
				"ukuran":   b.Ukuran,
			}
			finalData = append(finalData, row)
		}
	}

	fmt.Println(finalData, "==> Final Data Gabungan")
	return finalData
}

func HapusBarangChild(w http.ResponseWriter, db *gorm.DB, nama, kode, jenis, ukuran, kode_barang string) map[string]string {
	fmt.Println("===> [HapusBarangChild] Proses dimulai")

	// Step 1: Ambil kategori dari tabel barang_custom
	var kategori string
	err := db.Table("barang_custom").
		Select("kategori").
		Where(`"nama" = ? AND "id" = ?`, nama, kode).
		Scan(&kategori).Error

	if err != nil {
		fmt.Println(" Gagal mengambil kategori dari barang_custom:", err)
		return map[string]string{
			"status":  "false",
			"message": "Gagal mengambil kategori dari barang_custom",
		}
	}

	Table := kategori

	// Jadikan lowercaseo
	fmt.Println(" Kategori ditemukan dan dikonversi ke lowercase:", kategori)

	// Step 2: Coba nama kolom dengan huruf besar
	namaKolom := fmt.Sprintf(`"Nama_%s"`, kategori)
	kodeKolom := fmt.Sprintf(`"Kode_%s"`, kategori)
	query := fmt.Sprintf("%s = ? AND %s = ?", namaKolom, kodeKolom)

	fmt.Println(" Mencoba menghapus data dengan kolom:", namaKolom, "&", kodeKolom)
	result := db.Table(kategori).Where(query, nama, kode_barang).Delete(nil)

	if result.Error != nil || result.RowsAffected == 0 {
		// Step 3: Fallback ke nama kolom huruf kecil
		fmt.Println(" Gagal atau tidak ada baris terhapus, mencoba nama kolom huruf kecil...")

		namaKolomLower := fmt.Sprintf(`"Nama_%s"`, strings.ToLower(kategori))
		kodeKolomLower := fmt.Sprintf(`"Kode_%s"`, strings.ToLower(kategori))

		queryLower := fmt.Sprintf("%s = ? AND %s = ?", namaKolomLower, kodeKolomLower)

		resultFallback := db.Table(Table).Where(queryLower, nama, kode_barang).Delete(nil)
		if resultFallback.Error != nil {
			fmt.Println(" Gagal menghapus data dari tabel kategori (fallback):", resultFallback.Error)
			return map[string]string{
				"status":  "false",
				"message": "Gagal menghapus data dari tabel kategori (fallback)",
			}
		}

		if resultFallback.RowsAffected == 0 {
			fmt.Println(" Tidak ada baris yang terhapus dalam fallback.")
			return map[string]string{
				"status":  "false",
				"message": "Data tidak ditemukan atau sudah dihapus",
			}
		}

		fmt.Println(" Data berhasil dihapus (fallback) dari tabel", kategori)
		return map[string]string{
			"status":  "true",
			"message": "Data berhasil dihapus (fallback)",
		}
	}

	// Jika berhasil langsung dari percobaan pertama
	fmt.Println(" Data berhasil dihapus dari tabel", kategori)
	return map[string]string{
		"status":  "true",
		"message": "Data berhasil dihapus",
	}
}
