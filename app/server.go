package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Lingkup struct {
	database *gorm.DB
	ruter    *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
	AppHost string
}

type DBconfig struct {
	dbHost string
	dbPass string
	dbPort string
	dbName string
	dbUser string
}

func (server *Lingkup) Inisialisasi(appConfig AppConfig, dbConfig DBconfig) {
	fmt.Println("Masuk Ke \n" + appConfig.AppName)

	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=kasir_go port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.dbHost, dbConfig.dbUser, dbConfig.dbPass, dbConfig.dbPort,
	)

	fmt.Println("DSN yang digunakan:", dsn)

	// Membuka koneksi ke database dengan GORM
	server.database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	fmt.Println("GORM Logger aktif!")

	if err != nil {
		panic("Server Gagal Menyambungkan ke database")
	} else {
		fmt.Println("Berhasil terhubung ke database:", "kasir_go")
		fmt.Println("DB_NAME dari .env:", os.Getenv("DB_NAME"))
	}

	// 🔍 Cek database yang benar-benar digunakan oleh PostgreSQL
	var currentDB string
	server.database.Raw("SELECT current_database();").Scan(&currentDB)
	fmt.Println("Database yang sedang digunakan:", currentDB)

	server.ruter = mux.NewRouter()
	server.InitializeRoutes()

	fmt.Println("DB_NAME dari cekEnv:", dbConfig.dbName)
}

func cekEnv(nilai, fallback string) string {
	if nilai, ok := os.LookupEnv(nilai); ok {
		return nilai
	}

	return fallback
}

func Jalan() {
	var server Lingkup
	var appConfig AppConfig

	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Konfigurasi aplikasi
	appConfig.AppHost = cekEnv("APP_HOST", "localhost")
	appConfig.AppPort = cekEnv("APP_PORT", "8080")
	appConfig.AppName = cekEnv("APP_NAME", "Faiz")

	// Konfigurasi database
	dbConfig := DBconfig{
		dbHost: cekEnv("DB_HOST", "localhost"),
		dbPort: cekEnv("DB_PORT", "5432"),
		dbUser: cekEnv("DB_USER", "postgres"),
		dbPass: cekEnv("DB_PASSWORD", "Faiz"),
		dbName: cekEnv("DB_NAME", "kasir_go"),
	}

	// Inisialisasi server dengan konfigurasi aplikasi dan database
	server.Inisialisasi(appConfig, dbConfig)

	// Menjalankan server
	server.Jalan(appConfig.AppHost + ":" + appConfig.AppPort)

}

func (server *Lingkup) Jalan(alamat string) {
	fmt.Printf("Server berjalan di %s\n", alamat)

	log.Fatal(http.ListenAndServe(alamat, server.ruter))
}
