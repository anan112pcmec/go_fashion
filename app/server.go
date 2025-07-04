package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	if v, ok := os.LookupEnv(nilai); ok {
		return v
	}

	return fallback
}

// Upgrader WebSocket (biar bisa upgrade HTTP ke WS)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Atur kebijakan origin sesuai kebutuhan, ini buka semua origin
		return true
	},
}

// Handler WebSocket
func (server *Lingkup) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection ke WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Gagal upgrade ke WebSocket:", err)
		return
	}
	defer conn.Close()

	log.Println("Client WebSocket connected:", r.RemoteAddr)

	// Loop baca pesan dari client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Printf("Pesan diterima: %s\n", message)

		// Contoh: Echo balik pesan yang diterima
		err = conn.WriteMessage(messageType, []byte("Server balas: "+string(message)))
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}

	log.Println("Client WebSocket disconnected:", r.RemoteAddr)
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
