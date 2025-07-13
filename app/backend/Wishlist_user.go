package backend

import (
	"fmt"

	"gorm.io/gorm"
)

type WishlistData struct {
	Id           int    `gorm:"column:id"`
	IdUser       int    `gorm:"column:id_user"`
	NamaUser     string `gorm:"column:NamaUser"`
	NamaPakaian  string `gorm:"column:NamaPakaian"`
	JenisPakaian string `gorm:"column:Jenis_pakaian"`
}

func WishlistUser(db *gorm.DB, nama, id string) []map[string]interface{} {
	fmt.Println("===> [WishlistUser] Memulai proses pengambilan wishlist")
	fmt.Printf("[INFO] Data masuk: nama = %s, id = %s\n", nama, id)

	var exists bool

	// Step 1: Validasi user
	fmt.Println("[Step 1] Validasi keberadaan user berdasarkan id dan nama")
	err := db.Table("user").
		Select("count(*) > 0").
		Where(`id = ? AND "Nama" = ?`, id, nama).
		Find(&exists).Error

	if err != nil {
		fmt.Println("[ERROR] Terjadi kesalahan saat validasi user:", err)
		return []map[string]interface{}{
			{
				"Status":     "Gagal",
				"Keterangan": "Terjadi kesalahan saat validasi user",
			},
		}
	}

	if !exists {
		fmt.Println("[INFO] User tidak ditemukan dengan id dan nama tersebut")
		return []map[string]interface{}{
			{
				"Status":     "Gagal",
				"Keterangan": "User tidak ditemukan atau tidak valid",
			},
		}
	}
	fmt.Println("[INFO] User valid, melanjutkan proses pengambilan wishlist")

	// Step 2: Ambil semua wishlist user
	var wishlists []WishlistData
	fmt.Println("[Step 2] Mengambil semua data wishlist dari tabel Wishlist")
	err = db.Table("Wishlist").
		Where(`id_user = ? AND "NamaUser" = ?`, id, nama).
		Find(&wishlists).Error

	if err != nil {
		fmt.Println("[ERROR] Gagal mengambil data wishlist:", err)
		return []map[string]interface{}{
			{
				"Status":     "Gagal",
				"Keterangan": "Gagal mengambil data wishlist: " + err.Error(),
			},
		}
	}

	fmt.Printf("[INFO] Jumlah data wishlist ditemukan: %d\n", len(wishlists))

	// Step 3: Konversi hasil ke []map[string]interface{}
	fmt.Println("[Step 3] Mengonversi data struct ke map[string]interface{}")
	var result []map[string]interface{}
	for i, w := range wishlists {
		fmt.Printf("  -> [Item %d] %s (%s)\n", i+1, w.NamaPakaian, w.JenisPakaian)
		result = append(result, map[string]interface{}{
			"id":            w.Id,
			"id_user":       w.IdUser,
			"NamaUser":      w.NamaUser,
			"NamaPakaian":   w.NamaPakaian,
			"Jenis_pakaian": w.JenisPakaian,
		})
	}

	fmt.Println("[SUKSES] Data wishlist berhasil dikembalikan")
	return result
}

func Wishlist(db *gorm.DB, nama, namapakaian, jenis_pakaian string) map[string]string {
	var userID int

	// Cek apakah user dengan nama tersebut ada
	if err := db.Table("user").Select("id").Where(`"Nama" = ?`, nama).Scan(&userID).Error; err != nil {
		return map[string]string{
			"Status":     "Gagal",
			"Keterangan": "User dengan nama '" + nama + "' tidak ditemukan",
		}
	}

	// Cek apakah item sudah ada di wishlist
	var count int64
	db.Table("Wishlist").Where(`id_user = ? AND "NamaPakaian" = ? AND "Jenis_pakaian" = ?`, userID, namapakaian, jenis_pakaian).Count(&count)

	if count > 0 {
		// Jika sudah ada, hapus dari wishlist
		err := db.Table("Wishlist").Where(`id_user = ? AND "NamaPakaian" = ? AND "Jenis_pakaian" = ?`, userID, namapakaian, jenis_pakaian).Delete(nil).Error
		if err != nil {
			return map[string]string{
				"Status":     "Gagal",
				"Keterangan": "Gagal menghapus dari wishlist: " + err.Error(),
			}
		}
		return map[string]string{
			"Status":     "Berhasil",
			"Keterangan": "'" + namapakaian + "' berhasil dihapus dari wishlist",
		}
	} else {
		// Jika belum ada, masukkan ke wishlist
		err := db.Table("Wishlist").Create(map[string]interface{}{
			"id_user":       userID,
			"NamaUser":      nama,
			"NamaPakaian":   namapakaian,
			"Jenis_pakaian": jenis_pakaian,
		}).Error

		if err != nil {
			return map[string]string{
				"Status":     "Gagal",
				"Keterangan": "Gagal menambahkan ke wishlist: " + err.Error(),
			}
		}

		return map[string]string{
			"Status":     "Berhasil",
			"Keterangan": "Berhasil menambahkan '" + namapakaian + "' ke dalam wishlist",
		}
	}
}
