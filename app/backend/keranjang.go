package backend

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Keranjang struct {
	ID     uint
	Nama   string
	Ukuran string
	Jenis  string
	Jumlah int
	Harga  int
}

type Ambil_Barang struct {
	Nama        string `gorm:"column:Nama"`
	Harga       string `gorm:"column:Harga"`
	Ukuran      string `gorm:"column:Ukuran"`
	Jenis       string `gorm:"column:Jenis"`
	Jumlah      int32  `gorm:"column:Ukuran"`
	Deskripsi   string `gorm:"column:deskripsi"`
	HargaSatuan int32  `gorm:"column:harga"`
}

// Pastikan GORM pakai tabel "keranjang"
func (Keranjang) TableName() string {
	return "keranjang"
}

// Fungsi untuk memasukkan data ke tabel keranjang
func MasukanKeranjang(w http.ResponseWriter, db *gorm.DB, nama, harga, jenis, stok, ukuran, jumlah, nama_user, id_user, deskripsi string) string {
	hargaInt, err := strconv.Atoi(harga)
	if err != nil {
		return "Gagal mengkonversi harga"
	}

	jumlahInt, err := strconv.Atoi(jumlah)
	if err != nil {
		return "Gagal mengkonversi jumlah"
	}

	var existingJumlah int
	err = db.Raw(`
		SELECT "Jumlah" FROM "keranjang" 
		WHERE "id_user" = ? AND "Nama" = ? AND "Ukuran" = ? AND "Jenis" = ? AND "deskripsi" = ?
	`, id_user, nama, ukuran, jenis, deskripsi).Row().Scan(&existingJumlah)

	newJumlah := existingJumlah + jumlahInt
	newTotalHarga := hargaInt * newJumlah

	// Jika kamu butuh string dari newTotalHarga
	HargaFinal := strconv.Itoa(newTotalHarga)

	if err == nil {
		// Barang sudah ada, update jumlah dan harga

		if err != nil {
			fmt.Println("Gagal Mengkonversi Harga Final")
			return "Gagal Mengkonversi Harga Final"
		}

		errUpdate := db.Exec(`
			UPDATE "keranjang" SET "Jumlah" = ?, "Harga" = ?
			WHERE "id_user" = ? AND "Nama" = ? AND "Ukuran" = ? AND "Jenis" = ? AND "deskripsi" = ?
		`, newJumlah, HargaFinal, id_user, nama, ukuran, jenis, deskripsi).Error

		if errUpdate != nil {
			return fmt.Sprintf("Gagal memperbarui jumlah keranjang: %v", errUpdate)
		}
	} else {
		// Barang belum ada, insert baru
		err = db.Exec(`
			INSERT INTO "keranjang" ("id_user", "Nama", "Ukuran", "Jenis", "Jumlah", "Harga", "deskripsi") 
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, id_user, nama, ukuran, jenis, jumlahInt, HargaFinal, deskripsi).Error

		if err != nil {
			return fmt.Sprintf("Gagal menambahkan ke keranjang: %v", err)
		}
	}

	return "Berhasil menambahkan atau memperbarui keranjang dan mengubah status baju"
}

func UpdateKeranjangPengguna(w http.ResponseWriter, db *gorm.DB, id, nama_barang, jumlah, ukuran string) map[string]string {
	var harga int64
	var kategori string

	// 1. Ambil kategori dari barang_custom berdasarkan nama barang
	err1 := db.Table(`"barang_custom"`).
		Select(`"kategori"`).
		Where(`"nama" = ?`, nama_barang).
		Scan(&kategori).Error

	if err1 != nil || kategori == "" {
		hasil := map[string]string{
			"Status": "Gagal: Kategori tidak ditemukan di barang_custom",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 2. Cek dan ambil harga dari tabel kategori (misalnya: tabel "Baju", "Sepatu", dll)
	var count int64
	var tablekategori string
	var tableUkuran string
	if kategori == "Baju" {
		tablekategori = "baju"
		tableUkuran = "Ukuran"
	} else if kategori == "Celana" {
		tablekategori = "Celana"
		tableUkuran = "Ukuran"
	} else if kategori == "Sepatu" {
		tablekategori = "sepatu"
		tableUkuran = "Ukuran"
	}
	// Bentuk nama kolom dinamis, misalnya: Nama_Baju

	var namaKolom string
	Harga := "Harga"

	if kategori == "Celana" || kategori == "Baju" || kategori == "Kacamata" {
		namaKolom = fmt.Sprintf("Nama_%s", tablekategori)
	} else {
		namaKolom = fmt.Sprintf("nama_%s", tablekategori)
		Harga = "harga"
	}

	// Query dengan nama kolom dinamis
	err2 := db.Table(fmt.Sprintf(`"%s"`, kategori)).
		Where(fmt.Sprintf(`"%s" = ? AND "%s" = ?`, namaKolom, tableUkuran), nama_barang, ukuran).
		Count(&count).Error

	if err2 != nil || count == 0 {
		hasil := map[string]string{
			"Status": fmt.Sprintf("Gagal: Barang dengan nama '%s' dan ukuran '%s' tidak ditemukan di tabel kategori '%s'", nama_barang, ukuran, kategori),
		}
		fmt.Println(hasil)
		return hasil
	}

	// 3. Ambil harga dari tabel kategori
	err3 := db.Table(fmt.Sprintf(`"%s"`, kategori)).
		Select(fmt.Sprintf(`"%s_%s"`, Harga, tablekategori)).
		Where(fmt.Sprintf(`"%s" = ? AND "%s" = ?`, namaKolom, tableUkuran), nama_barang, ukuran).
		Scan(&harga).Error

	if err3 != nil {
		hasil := map[string]string{
			"Status": "Gagal Mendapatkan Harga dari Tabel Kategori",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 4. Konversi jumlah
	jumlahInt, err4 := strconv.Atoi(jumlah)
	if err4 != nil {
		hasil := map[string]string{
			"Status": "Jumlah tidak valid",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 5. Hitung total harga
	totalHarga := int64(jumlahInt) * harga
	totalHargaStr := strconv.FormatInt(totalHarga, 10)

	// 6. Update ke tabel keranjang
	err5 := db.Table(`"keranjang"`).
		Where(`"id_user" = ? AND "Nama" = ? AND "Ukuran" = ?`, id, nama_barang, ukuran).
		Updates(map[string]interface{}{
			"Jumlah": jumlahInt,
			"Harga":  totalHargaStr,
		}).Error

	if err5 != nil {
		hasil := map[string]string{
			"Status": "Gagal Mengupdate Keranjang",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 7. Sukses
	hasil := map[string]string{
		"Status":      "Berhasil",
		"Kategori":    kategori,
		"HargaSatuan": fmt.Sprintf("%d", harga),
		"HargaTotal":  totalHargaStr,
	}
	fmt.Println(hasil)
	return hasil
}

func CheckBarang(w http.ResponseWriter, db *gorm.DB, nama string, jenis string, jumlah string, ukuran, kategori string) map[string]string {
	var count int64
	ukuranGede := strings.ToUpper(ukuran)
	hasil := make(map[string]string)

	// Mapping kategori ke nama tabel & kolom
	var (
		namaTabel   string
		kolomNama   string
		kolomJenis  string
		jenisBarang string
	)

	switch kategori {
	case "Baju":
		namaTabel = `"Baju"`
		kolomNama = `"Nama_baju"`
		kolomJenis = `"Jenis_baju"`
		jenisBarang = "Baju"
	case "Celana":
		namaTabel = `"Celana"`
		kolomNama = `"Nama_Celana"`
		kolomJenis = `"Jenis_Celana"`
		jenisBarang = "Celana"
	case "Kacamata":
		namaTabel = `"Kacamata"`
		kolomNama = `"Nama_Kacamata"`
		kolomJenis = `"Jenis_Kacamata"`
		jenisBarang = "Kacamata"
	case "Sepatu":
		namaTabel = "Sepatu"
		kolomNama = "nama_sepatu"
		kolomJenis = "jenis_sepatu"
		jenisBarang = "Sepatu"
	default:
		hasil["status"] = "Gagal"
		hasil["item"] = "Jenis barang tidak dikenali"
		return hasil
	}

	// Query stok
	err := db.Table(namaTabel).
		Where(fmt.Sprintf(`%s = ? AND %s = ? AND "Ukuran" = ?`, kolomNama, kolomJenis), nama, jenis, ukuranGede).
		Count(&count).Error

	if err != nil {
		log.Println("Gagal hitung stok", jenisBarang, ":", err)
		hasil["status"] = "Gagal"
		hasil["item"] = "Terjadi kesalahan saat menghitung stok"
		return hasil
	}

	jumlahInt, err := strconv.ParseInt(jumlah, 10, 64)
	if err != nil {
		log.Println("Gagal konversi jumlah:", err)
		hasil["status"] = "Gagal"
		hasil["item"] = "Jumlah yang diberikan tidak valid"
		return hasil
	}

	if count < jumlahInt {
		hasil["status"] = "Gagal"
		hasil["item"] = fmt.Sprintf("Kelebihan order, sisa stok hanya %d", count)
	} else {
		sisa := count - jumlahInt
		hasil["status"] = "Berhasil"
		hasil["item"] = fmt.Sprintf("%s tersedia, sisa stok setelah order: %d", jenisBarang, sisa)
	}

	return hasil
}

func AmbilKeranjangUser(w http.ResponseWriter, db *gorm.DB, nama, user_id string) []map[string]string {
	fmt.Println("===> Memulai pengambilan data keranjang user")
	fmt.Println("Parameter diterima: nama =", nama, ", user_id =", user_id)

	var barangs []Ambil_Barang

	// Menjalankan query ke database
	fmt.Println("Menjalankan query untuk mengambil data dari tabel keranjang")
	err := db.Table("keranjang").
		Select(`keranjang."Nama", keranjang."Harga", keranjang."Ukuran", keranjang."Jenis", keranjang."Jumlah", keranjang."deskripsi", barang_custom."harga"`).
		Joins(`JOIN barang_custom ON keranjang."Nama" = barang_custom."nama" AND keranjang."deskripsi" = barang_custom."deskripsi"`).
		Where(`keranjang."id_user" = ?`, user_id).
		Find(&barangs).Error

	// Cek apakah ada error
	if err != nil {
		fmt.Println("❌ Gagal mengambil data keranjang dari database:", err)
		return nil
	}

	fmt.Println("✅ Data berhasil diambil dari database")
	fmt.Printf("Jumlah baris yang ditemukan: %d\n", len(barangs))

	// Proses konversi data ke bentuk map[string]string
	var result []map[string]string
	for i, barang := range barangs {
		fmt.Printf("➤ Memproses baris ke-%d: %+v\n", i+1, barang)

		m := map[string]string{
			"nama":         barang.Nama,
			"harga":        barang.Harga,
			"jumlah":       fmt.Sprintf("%d", barang.Jumlah),
			"jenis":        barang.Jenis,
			"ukuran":       barang.Ukuran,
			"deskripsi":    barang.Deskripsi,
			"harga_satuan": fmt.Sprintf("%d", barang.HargaSatuan),
		}

		fmt.Printf("   Data dikonversi: %+v\n", m)
		result = append(result, m)
	}

	fmt.Println("✅ Semua data berhasil diproses dan dikembalikan")
	return result
}

func HapusBarangDariKeranjang(db *gorm.DB, nama, ukuran, jumlah, namauser, iduser string) map[string]interface{} {
	// Step 1: Konversi iduser ke int
	id, err := strconv.Atoi(iduser)
	if err != nil {
		return map[string]interface{}{
			"status": "gagal",
			"pesan":  "ID user tidak valid",
		}
	}

	// Step 2: Validasi user
	var user struct {
		ID   int
		Nama string
	}

	if err := db.Table(`"user"`).
		Where(`"id" = ?`, id).
		First(&user).Error; err != nil {
		return map[string]interface{}{
			"status": "gagal",
			"pesan":  "User tidak ditemukan",
		}
	}

	// Step 3: Hapus dari keranjang
	if err := db.Table(`"keranjang"`).
		Where(`"id_user" = ? AND "Nama" = ? AND "Ukuran" = ? AND "Jumlah" = ?`, iduser, nama, ukuran, jumlah).
		Delete(nil).Error; err != nil {
		return map[string]interface{}{
			"status": "gagal",
			"pesan":  "Gagal menghapus barang dari keranjang",
		}
	}

	// Step 4: Sukses
	return map[string]interface{}{
		"status": "berhasil",
	}
}
