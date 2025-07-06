package backend

import (
	"encoding/json"
	"fmt"
	"io"
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

type RequestUserSettings struct {
	Tujuan       string `json:"tujuan"`
	IDUser       string `json:"iduser"`
	NamaUser     string `json:"namauser"`
	NomorHP      string `json:"nomor_hp"`
	Username     string `json:"username"`
	Bio          string `json:"bio"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
	TanggalLahir string `json:"tanggal_lahir"`
}

func AmbilDataUser(w http.ResponseWriter, db *gorm.DB, nama, id string) map[string]interface{} {
	log.Println("[AMBIL DATA USER] Memulai proses...")

	var userinfo UserInformation

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println("[ERROR] ID tidak valid:", err)
		fmt.Fprintln(w, "ID tidak valid")
		return map[string]interface{}{"status": "false", "message": "ID tidak valid"}
	}

	log.Println("[AMBIL DATA USER] Querying database...")
	err = db.Table("user").Where(`"Nama" = ? AND id = ?`, nama, int32(idInt)).First(&userinfo).Error
	if err != nil {
		log.Println("[ERROR] Data tidak ditemukan:", err)
		fmt.Fprintln(w, "Data tidak ditemukan:", err.Error())
		return map[string]interface{}{"status": "false", "message": "Data tidak ditemukan"}
	}

	log.Println("[AMBIL DATA USER] Data ditemukan, menyusun map...")

	// Konversi int ke string jika perlu
	return map[string]interface{}{
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
		"foto":      userinfo.Foto,
	}
}

func MasukinInfoPribadi(w http.ResponseWriter, db *gorm.DB, r *http.Request) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Form tidak bisa diparsing", http.StatusBadRequest)
		return
	}

	// Ambil inputan form
	data := RequestUserSettings{
		Tujuan:       r.FormValue("tujuan"),
		IDUser:       r.FormValue("iduser"),
		NamaUser:     r.FormValue("namauser"),
		NomorHP:      r.FormValue("nomor_hp"),
		Username:     r.FormValue("username"),
		Bio:          r.FormValue("bio"),
		Email:        r.FormValue("email"),
		Gender:       r.FormValue("gender"),
		TanggalLahir: r.FormValue("tanggal_lahir"),
	}

	fmt.Println("Data diterima:", data)

	// Validasi user
	var existing UserInformation
	err = db.Table("user").Where(`"id" = ? AND "Nama" = ?`, data.IDUser, data.NamaUser).First(&existing).Error
	if err != nil {
		fmt.Println("User tidak ditemukan:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  "gagal",
			"Message": "User tidak valid atau tidak ditemukan",
		})
		return
	}

	// Ambil file foto jika ada
	var fotoData []byte
	hapusFoto := false
	file, _, err := r.FormFile("foto")
	if err == nil {
		defer file.Close()
		fotoData, _ = io.ReadAll(file)
	} else {
		fmt.Println("Tidak ada file foto dikirim, akan reset foto user.")
		hapusFoto = true
	}

	// Konversi nomor HP ke int
	nomorHpInt, err := strconv.Atoi(data.NomorHP)
	if err != nil {
		fmt.Println("Nomor HP tidak valid:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  "gagal",
			"Message": "Nomor HP tidak valid",
		})
		return
	}

	// Siapkan update map
	updateMap := map[string]interface{}{
		"Username": data.Username,
		"Bio":      data.Bio,
		"Email":    data.Email,
		"Nomor_Hp": nomorHpInt,
		"Gender":   data.Gender,
	}

	if hapusFoto {
		updateMap["Foto"] = nil // Hapus foto dari database
	} else {
		updateMap["Foto"] = fotoData
	}

	// Update data ke DB
	err = db.Table("user").
		Where(`"id" = ? AND "Nama" = ?`, data.IDUser, data.NamaUser).
		Updates(updateMap).Error

	if err != nil {
		fmt.Println("Gagal update data:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"Status":  "gagal",
			"Message": "Gagal update data user",
		})
		return
	}

	// Respon sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Status":  "berhasil",
		"Message": "Data user berhasil diperbarui",
	})
}
