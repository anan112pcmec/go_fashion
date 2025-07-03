package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type DataRequest struct {
	JenisPakaian string `json:"jenis_pakaian"`
	Deskripsi    string `json:"deskripsi"`
	Kategori     string `json:"kategori"`
	Harga        int32  `json:"harga"`
	Warna        string `json:"warna"`
	Stok         int32  `json:"stok"`
	Ukuran       string `json:"ukuran"`
	File         []byte `json:"file"` // Simpan file sebagai byte array
	Nama         string `json:"nama_produk"`
	Untuk        string `json:"untuk"`
}

// Struct untuk database
type Barang struct {
	ID           uint `gorm:"primaryKey"`
	JenisPakaian string
	Deskripsi    string
	Kategori     string
	Harga        int32
	Warna        string
	Stok         int32
	Ukuran       string
	File         []byte `gorm:"type:longblob"`
	Nama         string `gorm:"column:nama"`
	S            int    `gorm:"column:S"`
	M            int    `gorm:"column:M"`
	L            int    `gorm:"column:L"`
	X            int    `gorm:"column:X"`
	XL           int    `gorm:"column:XL"`
}

type TableBaju struct {
	KodeBaju  int16  `gorm:"column:Kode_baju"`
	NamaBaju  string `gorm:"column:Nama_baju"`
	HargaBaju int32  `gorm:"column:Harga_baju"`
	JenisBaju string `gorm:"column:Jenis_baju"`
	Warna     string `gorm:"column:Warna"`
	Untuk     string `gorm:"column:untuk"`
	Id_induk  int16  `gorm:"column:id_induk"`
	Ukuran    string `gorm:"column:Ukuran"`
}

type TableCelana struct {
	KodeBaju  int16  `gorm:"column:Kode_Celana"`
	NamaBaju  string `gorm:"column:Nama_Celana"`
	HargaBaju int32  `gorm:"column:Harga_Celana"`
	JenisBaju string `gorm:"column:Jenis_Celana"`
	Warna     string `gorm:"column:Warna"`
	Untuk     string `gorm:"column:untuk"`
	Id_induk  int    `gorm:"column:id_induk"`
	Ukuran    string `gorm:"column:Ukuran"`
}

type TableSepatu struct {
	KodeBaju  int16  `gorm:"column:kode_sepatu"`
	NamaBaju  string `gorm:"column:nama_sepatu"`
	HargaBaju int32  `gorm:"column:harga_sepatu"`
	JenisBaju string `gorm:"column:jenis_sepatu"`
	Warna     string `gorm:"column:warna"`
	Untuk     string `gorm:"column:untuk"`
	Id_induk  int    `gorm:"column:id_induk"`
	Ukuran    string `gorm:"column:Ukuran"`
}

type TableKacamata struct {
	KodeKacamata  int16  `gorm:"column:Kode_Kacamata"`
	NamaKacamata  string `gorm:"column:Nama_Kacamata"`
	HargaKacamata int32  `gorm:"column:Harga_Kacamata"`
	JenisKacamata string `gorm:"column:Jenis_Kacamata"`
	Warna         string `gorm:"column:Warna"`
	Untuk         string `gorm:"column:untuk"`
	Id_induk      int16  `gorm:"column:id_induk"`
	Ukuran        string `gorm:"column:Ukuran"`
}

type KolomId struct {
	Id int `gorm:"column:id"`
}

type KolomStok struct {
	Stock   int `gorm:"column:stok"`
	StockS  int `gorm:"column:S"`
	StockM  int `gorm:"column:M"`
	StockL  int `gorm:"column:L"`
	StockX  int `gorm:"column:X"`
	StockXL int `gorm:"column:XL"`
}

type CariBarang struct {
	Nama      string `gorm:"column:jenis_pakaian"`
	Deskripsi string `gorm:"column:deskripsi"`
}

func generateRandomNumber() int {
	return rand.Intn(90000) + 10000
}

func MasukanBarang(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Gagal memproses multipart form", http.StatusBadRequest)
		return
	}

	// Ambil nilai dari form-data
	jenisPakaian := r.FormValue("jenis_pakaian")
	deskripsi := r.FormValue("deskripsi")
	kategori := r.FormValue("kategori")
	harga := r.FormValue("harga")
	warna := r.FormValue("warna")
	NamaBarang := r.FormValue("nama_produk")
	untuk := r.FormValue("untuk")
	sStr := r.FormValue("s")
	mStr := r.FormValue("m")
	lStr := r.FormValue("l")
	xStr := r.FormValue("x")
	xlStr := r.FormValue("xl")
	s, _ := strconv.Atoi(sStr)
	m, _ := strconv.Atoi(mStr)
	l, _ := strconv.Atoi(lStr)
	x, _ := strconv.Atoi(xStr)
	xl, _ := strconv.Atoi(xlStr)

	stok := s + m + l + x + xl
	stokStr := strconv.Itoa(stok)

	NamaBarang = string(NamaBarang)

	// Ambil file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Gagal mengambil file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println(handler)

	// Baca isi file sebagai byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Gagal membaca file", http.StatusInternalServerError)
		return
	}

	// Konversi harga dan stok ke tipe data integer
	var hargaInt, stokInt int32
	fmt.Sscanf(harga, "%d", &hargaInt)
	fmt.Sscanf(stokStr, "%d", &stokInt)

	var cek_Barang []CariBarang // Ubah struct menjadi slice

	if err := db.Table("barang_custom").
		Where("nama = ? AND kategori = ?", NamaBarang, kategori).
		First(&cek_Barang).Error; err == nil {

		fmt.Println("Data ditemukan:", cek_Barang)
		fmt.Println("dengan data dicari adalah:", NamaBarang, "dan dengan kategori", kategori)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sudah Ada Buku dengan judul ini Harap Ganti"})
		return
	}

	// Buat objek barang
	barang_custom := Barang{
		JenisPakaian: jenisPakaian,
		Deskripsi:    deskripsi,
		Kategori:     kategori,
		Harga:        hargaInt,
		Warna:        warna,
		Stok:         stokInt,
		Ukuran:       "S,M,L,X,XL",
		File:         fileBytes,
		Nama:         NamaBarang,
		S:            s,
		M:            m,
		L:            l,
		X:            x,
		XL:           xl,
	}

	// Simpan ke database
	if err := db.Table("barang_custom").Create(&barang_custom).Error; err != nil {
		http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
		return
	}

	fmt.Println(kategori)

	if kategori == "baju" || kategori == "Baju" {
		barangbaju := TableBaju{
			NamaBaju:  NamaBarang,
			HargaBaju: int32(hargaInt),
			JenisBaju: jenisPakaian,
			Warna:     warna,
			Untuk:     untuk,
		}

		var aidi KolomId
		var Stock KolomStok

		// Perbaikan: First() hanya menerima satu struct
		err1 := db.Table("barang_custom").
			Select("id", "stok").
			Where("nama = ?", NamaBarang).
			First(&aidi).Error

		if err1 != nil {
			barangbaju.Id_induk = 0
		} else {
			barangbaju.Id_induk = int16(aidi.Id)
		}

		// Perbaikan: Mengisi Stock sebelum digunakan
		err2 := db.Table("barang_custom").
			Select("stok").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errS := db.Table("barang_custom").
			Select("S").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errM := db.Table("barang_custom").
			Select("M").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errL := db.Table("barang_custom").
			Select("L").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errX := db.Table("barang_custom").
			Select("X").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errXL := db.Table("barang_custom").
			Select("XL").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		if err2 != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errS != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errM != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errX != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errXL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}
		// Pastikan Stock.Stock memiliki nilai sebelum digunakan dalam loop
		for i := 0; i < int(Stock.StockS); i++ {

			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "S"

			if err := db.Table("Baju").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockM); i++ {

			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "M"

			if err := db.Table("Baju").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockL); i++ {

			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "L"

			if err := db.Table("Baju").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockX); i++ {

			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "X"

			if err := db.Table("Baju").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockXL); i++ {

			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "XL"

			if err := db.Table("Baju").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		var cari_barang []CariBarang

		if err := db.Table("barang_custom").
			Select("jenis_pakaian, deskripsi").
			Where("jenis_pakaian = ?", "cwcjb").
			Find(&cari_barang).Error; err != nil {

			http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Data ditemukan:", cari_barang)
		}
	} else if kategori == "Celana" || kategori == "celana" {
		barangbaju := TableCelana{
			NamaBaju:  NamaBarang,
			HargaBaju: int32(hargaInt),
			JenisBaju: jenisPakaian,
			Warna:     warna,
			Untuk:     untuk,
		}

		var aidi KolomId
		var Stock KolomStok

		// Perbaikan: First() hanya menerima satu struct
		err1 := db.Table("barang_custom").
			Select("id", "stok").
			Where("nama = ?", NamaBarang).
			First(&aidi).Error

		if err1 != nil {
			barangbaju.Id_induk = 0
		} else {
			barangbaju.Id_induk = int(aidi.Id)
		}

		// Perbaikan: Mengisi Stock sebelum digunakan
		err2 := db.Table("barang_custom").
			Select("stok").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errS := db.Table("barang_custom").
			Select("S").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errM := db.Table("barang_custom").
			Select("M").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errL := db.Table("barang_custom").
			Select("L").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errX := db.Table("barang_custom").
			Select("X").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errXL := db.Table("barang_custom").
			Select("XL").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		if err2 != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errS != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errM != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errX != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errXL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		for i := 0; i < int(Stock.StockS); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "S"

			if err := db.Table("Celana").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockM); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "M"

			if err := db.Table("Celana").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "L"

			if err := db.Table("Celana").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "X"

			if err := db.Table("Celana").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "XL"

			if err := db.Table("Celana").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		var cari_barang []CariBarang

		if err := db.Table("barang_custom").
			Select("jenis_pakaian, deskripsi").
			Where("jenis_pakaian = ?", "cwcjb").
			Find(&cari_barang).Error; err != nil {

			http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Data ditemukan:", cari_barang)
		}
	} else if kategori == "Sepatu" || kategori == "sepatu" {
		barangbaju := TableSepatu{
			NamaBaju:  NamaBarang,
			HargaBaju: int32(hargaInt),
			JenisBaju: jenisPakaian,
			Warna:     warna,
			Untuk:     untuk,
		}

		var aidi KolomId
		var Stock KolomStok

		// Perbaikan: First() hanya menerima satu struct
		err1 := db.Table("barang_custom").
			Select("id", "stok").
			Where("nama = ?", NamaBarang).
			First(&aidi).Error

		if err1 != nil {
			barangbaju.Id_induk = 0
		} else {
			barangbaju.Id_induk = int(aidi.Id)
		}

		// Perbaikan: Mengisi Stock sebelum digunakan
		err2 := db.Table("barang_custom").
			Select("stok").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errS := db.Table("barang_custom").
			Select("S").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errM := db.Table("barang_custom").
			Select("M").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errL := db.Table("barang_custom").
			Select("L").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errX := db.Table("barang_custom").
			Select("X").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errXL := db.Table("barang_custom").
			Select("XL").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		if err2 != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errS != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errM != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errX != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errXL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		// Ukuran S
		for i := 0; i < int(Stock.StockS); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "S"

			if err := db.Table("Sepatu").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		// Ukuran M
		for i := 0; i < int(Stock.StockM); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "M"

			if err := db.Table("Sepatu").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		// Ukuran L
		for i := 0; i < int(Stock.StockL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "L"

			if err := db.Table("Sepatu").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		// Ukuran X
		for i := 0; i < int(Stock.StockX); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "X"

			if err := db.Table("Sepatu").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		// Ukuran XL
		for i := 0; i < int(Stock.StockXL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangbaju.KodeBaju = dapetinangka
			barangbaju.Ukuran = "XL"

			if err := db.Table("Sepatu").Create(&barangbaju).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		var cari_barang []CariBarang

		if err := db.Table("barang_custom").
			Select("jenis_pakaian, deskripsi").
			Where("jenis_pakaian = ?", "cwcjb").
			Find(&cari_barang).Error; err != nil {

			http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Data ditemukan:", cari_barang)
		}
	} else if kategori == "Kacamata" || kategori == "kacamata" {
		barangKacamata := TableKacamata{
			NamaKacamata:  NamaBarang,
			HargaKacamata: int32(hargaInt),
			JenisKacamata: jenisPakaian,
			Warna:         warna,
		}

		var aidi KolomId
		var Stock KolomStok

		err1 := db.Table("barang_custom").
			Select("id", "stok").
			Where("nama = ?", NamaBarang).
			First(&aidi).Error

		if err1 != nil {
			barangKacamata.Id_induk = 0
		} else {
			barangKacamata.Id_induk = int16(aidi.Id)
		}

		// Perbaikan: Mengisi Stock sebelum digunakan
		err2 := db.Table("barang_custom").
			Select("stok").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		if err2 != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		errS := db.Table("barang_custom").
			Select("S").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errM := db.Table("barang_custom").
			Select("M").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errL := db.Table("barang_custom").
			Select("L").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errX := db.Table("barang_custom").
			Select("X").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		errXL := db.Table("barang_custom").
			Select("XL").
			Where("nama = ?", NamaBarang).
			First(&Stock).Error

		if err2 != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errS != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errM != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errX != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		if errXL != nil {
			http.Error(w, "Gagal mendapatkan stok", http.StatusInternalServerError)
			return
		}

		// Pastikan Stock.Stock memiliki nilai sebelum digunakan dalam loop
		for i := 0; i < int(Stock.StockS); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangKacamata.KodeKacamata = dapetinangka
			barangKacamata.Ukuran = "S"

			if err := db.Table("Kacamata").Create(&barangKacamata).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockM); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangKacamata.KodeKacamata = dapetinangka
			barangKacamata.Ukuran = "M"

			if err := db.Table("Kacamata").Create(&barangKacamata).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangKacamata.KodeKacamata = dapetinangka
			barangKacamata.Ukuran = "L"

			if err := db.Table("Kacamata").Create(&barangKacamata).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockX); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangKacamata.KodeKacamata = dapetinangka
			barangKacamata.Ukuran = "X"

			if err := db.Table("Kacamata").Create(&barangKacamata).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		for i := 0; i < int(Stock.StockXL); i++ {
			dapetinangka := int16(generateRandomNumber())

			if dapetinangka < 0 {
				dapetinangka = -dapetinangka
			}

			barangKacamata.KodeKacamata = dapetinangka
			barangKacamata.Ukuran = "XL"

			if err := db.Table("Kacamata").Create(&barangKacamata).Error; err != nil {
				http.Error(w, "Gagal menyimpan data ke database", http.StatusInternalServerError)
				return
			}
		}

		var cari_barang []CariBarang

		if err := db.Table("barang_custom").
			Select("jenis_pakaian, deskripsi").
			Where("jenis_pakaian = ?", "cwcjb").
			Find(&cari_barang).Error; err != nil {

			http.Error(w, "Gagal mengambil data dari database", http.StatusInternalServerError)
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Data ditemukan:", cari_barang)
		}
	}

	// Beri response sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data berhasil dimasukkan"})
}
