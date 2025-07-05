package backend

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type UserInformation struct {
	Id         int    `gorm:"column:id"`
	Nama       string `gorm:"column:Nama"`
	Wishlist   string `gorm:"column:Wishlist"`
	Role       string `gorm:"column:role"`
	NamaAlamat string `gorm:"column:Nama_Alamat"`
	Koordinat  string `gorm:"column:Alamat_koordinat"`
	Email      string `gorm:"column:Email"`
	NomorHp    int    `gorm:"column:Nomor_Hp"`
	Gender     string `gorm:"column:Gender"`
	Username   string `gorm:"column:Username"`
	Bio        string `gorm:"column:Bio"`
	Foto       []byte `gorm:"column:Foto"`
}

func AmbilDataUser(w http.ResponseWriter, db *gorm.DB, nama, id string) map[string]string {
	log.Println("[AMBIL DATA USER] Memulai proses...")

	var userinfo UserInformation

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("[ERROR] ID tidak valid:", err)
		fmt.Fprintln(w, "ID tidak valid")
		return map[string]string{"status": "false", "message": "ID tidak valid"}
	}

	log.Println("[AMBIL DATA USER] Querying database...")
	err = db.Table("user").Where(`"Nama" = ? AND id = ?`, nama, int32(idInt)).First(&userinfo).Error
	if err != nil {
		log.Println("[ERROR] Data tidak ditemukan:", err)
		fmt.Fprintln(w, "Data tidak ditemukan:", err.Error())
		return map[string]string{"status": "false", "message": "Data tidak ditemukan"}
	}

	log.Println("[AMBIL DATA USER] Data ditemukan, menyusun map...")

	// Konversi int ke string jika perlu
	return map[string]string{
		"status":    "true",
		"id":        strconv.Itoa(userinfo.Id),
		"nama":      userinfo.Nama,
		"wishlist":  userinfo.Wishlist,
		"role":      userinfo.Role, // dijelaskan di bawah
		"alamat":    userinfo.NamaAlamat,
		"koordinat": userinfo.Koordinat,
		"email":     userinfo.Email,
		"nomor_hp":  strconv.Itoa(userinfo.NomorHp),
		"gender":    userinfo.Gender,
		"username":  userinfo.Username,
		"bio":       userinfo.Bio,
	}
}

func MasukinInfoPribadi(w http.ResponseWriter, db *gorm.DB) (hasil map[string]string) {

	hasil = map[string]string{
		"Status": "berhasil",
	}

	return
}
