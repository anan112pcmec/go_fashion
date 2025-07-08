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
	Nama      string `gorm:"column:Nama"`
	Harga     string `gorm:"column:Harga"`
	Ukuran    string `gorm:"column:Ukuran"`
	Jenis     string `gorm:"column:Jenis"`
	Jumlah    int32  `gorm:"column:Ukuran"`
	Deskripsi string `gorm:"column:deskripsi"`
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

	// Update status baju menjadi "dikeranjangi"
	err1 := db.Exec(`
		UPDATE "Baju" SET "status" = 'dikeranjangi' 
		WHERE "Nama_baju" = ? AND "Ukuran" = ? AND "Jenis_baju" = ?
	`, nama, ukuran, jenis).Error

	if err1 != nil {
		return fmt.Sprintf("Gagal mengubah status baju: %v", err1)
	}

	return "Berhasil menambahkan atau memperbarui keranjang dan mengubah status baju"
}

func UpdateKeranjangPengguna(w http.ResponseWriter, db *gorm.DB, id, nama_barang, jumlah, ukuran string) map[string]string {
	var count int64
	var harga int64

	// 1. Cek apakah barang dengan nama dan ukuran tersedia di barang_custom
	err1 := db.Table(`"barang_custom"`).
		Where(`"nama" = ? `, nama_barang).
		Count(&count).Error

	if err1 != nil || count == 0 {
		hasil := map[string]string{
			"Status": "Gagal: Barang dengan ukuran tidak ditemukan di barang_custom",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 2. Ambil harga berdasarkan nama dan ukuran
	err2 := db.Table(`"barang_custom"`).
		Select(`"harga"`).
		Where(`"nama" = ?`, nama_barang).
		Scan(&harga).Error

	if err2 != nil {
		hasil := map[string]string{
			"Status": "Gagal Mendapatkan Harga Satuan",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 3. Konversi jumlah
	jumlahInt, err3 := strconv.Atoi(jumlah)
	if err3 != nil {
		hasil := map[string]string{
			"Status": "Jumlah tidak valid",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 4. Hitung total harga
	totalHarga := int64(jumlahInt) * harga
	totalHargaStr := strconv.FormatInt(totalHarga, 10)

	// 5. Update ke tabel keranjang
	err4 := db.Table(`"keranjang"`).
		Where(`"id_user" = ? AND "Nama" = ? AND "Ukuran" = ?`, id, nama_barang, ukuran).
		Updates(map[string]interface{}{
			"Jumlah": jumlahInt,
			"Harga":  totalHargaStr,
		}).Error

	if err4 != nil {
		hasil := map[string]string{
			"Status": "Gagal Mengupdate Keranjang",
		}
		fmt.Println(hasil)
		return hasil
	}

	// 6. Sukses
	hasil := map[string]string{
		"Status": "Berhasil",
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
		Select(`"Nama", "Harga", "Ukuran", "Jenis", "Jumlah", "deskripsi"`).
		Where(`"id_user" = ?`, user_id).
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
			"nama":      barang.Nama,
			"harga":     barang.Harga,
			"jumlah":    fmt.Sprintf("%d", barang.Jumlah),
			"jenis":     barang.Jenis,
			"ukuran":    barang.Ukuran,
			"deskripsi": barang.Deskripsi,
		}

		fmt.Printf("   ↳ Data dikonversi: %+v\n", m)
		result = append(result, m)
	}

	fmt.Println("✅ Semua data berhasil diproses dan dikembalikan")
	return result
}
