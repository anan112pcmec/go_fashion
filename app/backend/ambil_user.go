package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

)

// Struktur data login
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struktur respons login
type LoginResponse struct {
	Berhasil bool   `json:"berhasil"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	Nama     string `json:"nama"`
	Alamat   string `json:"alamat"`
	Id       string `json:"id"`
}

type HakAkses struct {
	TipeEntity string
}

type Sesi struct {
	Ctx  context.Context // Context yang membawa informasi
	Role string          // Role dari user
}

type ResponseBalasan struct {
	LoginResponse
	HakAkses
}

// Struktur User sesuai database
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Nama     string `gorm:"column:Nama"`
	Password string `gorm:"column:Password"`
	Role     string `gorm:"column:role"`
	Alamat   string `gorm:"column:Nama_Alamat"`
}

// TableName untuk menentukan nama tabel yang benar
func (User) TableName() string {
	return "user" // Menggunakan tabel "user" bukan "users"
}

// Fungsi untuk menangani login
func Login(w http.ResponseWriter, db *gorm.DB, data json.RawMessage) {
	var loginData LoginData
	var user User

	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)

	if dbName != "kasir_go" {
		fmt.Println("❌ Error: Database yang digunakan bukan 'kasir_go'. Database saat ini:", dbName)
		http.Error(w, `{"berhasil": false, "message": "Akses ditolak: Harus menggunakan database kasir_go"}`, http.StatusForbidden)
		return
	}

	// Decode JSON dari request
	err := json.Unmarshal(data, &loginData)
	if err != nil {
		fmt.Println("❌ Error decoding JSON:", err)
		http.Error(w, `{"berhasil": false, "message": "Format data login tidak valid"}`, http.StatusBadRequest)
		return
	}

	// Debugging: Cetak data yang diterima
	fmt.Println("🔹 Login request diterima:", loginData)

	// Query menggunakan GORM untuk mencari user berdasarkan Nama dan Password
	result := db.Where(`"Nama" = ? AND "Password" = ?`, loginData.Username, loginData.Password).First(&user)

	// Cek jika user tidak ditemukan
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("❌ Login gagal: Username atau password salah")
			http.Error(w, `{"berhasil": false, "message": "Username atau password salah"}`, http.StatusUnauthorized)
		} else {
			fmt.Println("❌ Error saat query database:", result.Error)
			http.Error(w, `{"berhasil": false, "message": "Terjadi kesalahan pada server"}`, http.StatusInternalServerError)
		}
		return
	}

	// Buat token JWT
	token, err := generateJWT(user.Nama, user.Role, user.Alamat, user.ID)
	if err != nil {
		fmt.Println("❌ Error saat membuat token:", err)
		http.Error(w, `{"berhasil": false, "message": "Terjadi kesalahan saat membuat token"}`, http.StatusInternalServerError)
		return
	}

	// Response login
	response := ResponseBalasan{
		LoginResponse{
			Berhasil: true,
			Message:  `Login sukses`,
			Token:    token,
			Nama:     user.Nama,
			Alamat:   user.Alamat,
			Id:       strconv.FormatUint(uint64(user.ID), 10),
		},
		HakAkses{
			TipeEntity: user.Role,
		},
	}

	// Menambahkan token JWT ke dalam response
	response.Message = fmt.Sprintf("✅ Login sukses. Token: %s", token)

	// Set response header dan kirim JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	// Debugging: Informasi user yang berhasil login
	fmt.Println("✅ User berhasil login:", user.Nama, user.Role)
}

// Fungsi untuk membuat token JWT
func generateJWT(username, role, alamat string, id uint) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"Id":       id,
		"alamat":   alamat,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Membuat token dengan signing method HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Secret key untuk signing
	secretKey := os.Getenv("SECRET_KEY")

	fmt.Println(secretKey)
	fmt.Println("ini alamat nya", alamat)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
