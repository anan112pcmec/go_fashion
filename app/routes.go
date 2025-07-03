package app

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"

	"github.com/anan112pcmec/go_fashion/app/controllers"

)

// Render template sesuai request
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	// Path yang akan digunakan untuk template
	var userTemplatePath string
	var dynamicTemplatePath string

	if tmpl == "Login.tmpl" {
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/user.tmpl"
	} else {
		dynamicTemplatePath = "templates/" + tmpl
		userTemplatePath = "templates/user.tmpl"
	}

	fmt.Println("Path userTemplate:", userTemplatePath)
	fmt.Println("Path dynamicTemplate:", dynamicTemplatePath)

	// Gabungan template yang akan diproses
	templates := template.Must(template.ParseFiles(userTemplatePath, dynamicTemplatePath))

	// Cetak hasil parsing template
	fmt.Println("Parsed templates:", templates.DefinedTemplates())

	// Eksekusi template
	err := templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		fmt.Println("Error executing template:", err)
	}
}

func userTemplate(w http.ResponseWriter, tmpl string, tipeEntity string) {
	// Path template
	var userTemplatePath string
	var dynamicTemplatePath string

	fmt.Println("Ini Role dia:", tipeEntity, "| Diakses melalui userTemplate", tmpl)

	switch tmpl {
	case "Login.tmpl":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/user.tmpl"
	case "Pengguna.tmpl":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/pengguna.tmpl"
	case "Pria_Section.tmpl", "Wanita_Section.tmpl", "Anak_Section.tmpl", "BarangDikirim.tmpl":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/pengguna.tmpl"
		fmt.Println("ini section")
	case "keranjang.tmpl":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/pengguna.tmpl"
		fmt.Println("ini keranjang")
	case "BarangDikirim":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/pengguna.tmpl"
		fmt.Println("ini keranjang")
	default:
		dynamicTemplatePath = "templates/" + tmpl
		userTemplatePath = "templates/user.tmpl"
	}

	fmt.Println("Path userTemplate:", userTemplatePath)
	fmt.Println("Path dynamicTemplate:", dynamicTemplatePath)

	templates, err := template.ParseFiles(userTemplatePath, dynamicTemplatePath)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		fmt.Println("Error parsing templates:", err)
		return
	}

	err = templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		fmt.Println("Error executing template:", err)
		return
	}
}

func adminTemplate(w http.ResponseWriter, tmpl string, tipeEntity string) {
	// Path template
	var userTemplatePath string
	var dynamicTemplatePath string

	fmt.Println("Ini Role dia:", tipeEntity, "| Diakses melalui userTemplate")

	switch tmpl {
	case "Login.tmpl":
		dynamicTemplatePath = "templates/user/" + tmpl
		userTemplatePath = "templates/admin.tmpl"
	case "dashboard.tmpl":
		dynamicTemplatePath = "templates/admin/" + tmpl
		userTemplatePath = "templates/admin.tmpl"
	case "produk.tmpl":
		dynamicTemplatePath = "templates/admin/" + tmpl
		userTemplatePath = "templates/admin.tmpl"
		fmt.Println("ini section")
	case "Kontrol_produk.tmpl":
		dynamicTemplatePath = "templates/admin/" + tmpl
		userTemplatePath = "templates/admin.tmpl"
		fmt.Println("ini Kontrolgo")
	default:
		dynamicTemplatePath = "templates/" + tmpl
		userTemplatePath = "templates/admin.tmpl"
	}

	fmt.Println("Path userTemplate:", userTemplatePath)
	fmt.Println("Path dynamicTemplate:", dynamicTemplatePath)

	templates, err := template.ParseFiles(userTemplatePath, dynamicTemplatePath)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		fmt.Println("Error parsing templates:", err)
		return
	}

	err = templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, "Page not found admin", http.StatusNotFound)
		fmt.Println("Error executing template:", err)
		return
	}
}

// Handler untuk menangani request dinamis
func RouterHandler(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Path                    // Ambil path mentah dari request
	path := strings.TrimPrefix(rawPath, "/") // Trim hanya jika tidak mengandung "user"

	// Jika path mengandung kata "user" di mana pun, gunakan rawPath tanpa trim
	if strings.Contains(rawPath, "user") {
		path = rawPath[1:] // Hilangkan hanya karakter '/' di awal
	}

	fmt.Println("Path yang diterima:", rawPath) // Cetak path mentah dari request
	fmt.Println("Path setelah diproses:", path) // Cetak path setelah diproses

	fmt.Println("ini path raw:", rawPath)

	if path == "" {
		path = "home"
	} else if path == "Pria" {
		path = "Pria_Section.tmpl"
	} else if path == "Wanita" {
		path = "Wanita_Section.tmpl"
	} else if path == "Anak" {
		path = "Anak_Section.tmpl"
	} else if path == "Home" {
		path = "Pengguna.tmpl"
	} else if path == "Pengguna.tmpl" {
		path = "Invalid"
	} else if path == "makan" {
		path = "home"
	} else if path == "BarangDikirim" {
		path = "BarangDikirim.tmpl"
		fmt.Println("DiaMengakses Barang")
	}

	RenderTemplate(w, path) // Render template sesuai path
}

func EntityHandler(w http.ResponseWriter, alamat, role, token string) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		http.Error(w, "⛔ Akses ditolak: Token tidak valid", http.StatusUnauthorized)
		return
	}

	// Ambil role dari token
	tokenRole, ok := claims["role"].(string)
	if !ok {
		http.Error(w, "⛔ Token tidak memiliki role yang valid", http.StatusUnauthorized)
		return
	}

	username, err1 := claims["username"].(string)

	if !err1 {
		http.Error(w, "Token tidak memiliki username dan kredensial yang valid", http.StatusBadGateway)
	}

	// Ambil path mentah dari request
	rawPath := alamat
	path := strings.TrimPrefix(rawPath, "/")

	// Jika path mengandung kata "user", gunakan rawPath tanpa trim
	if strings.Contains(rawPath, "user") {
		path = rawPath[1:] // Hilangkan hanya karakter pertama '/'
	} else if strings.Contains(rawPath, "admin") {
		path = rawPath[1:]
	}

	fmt.Println("Path yang diterima:", rawPath)
	fmt.Println("Path setelah diproses:", path)
	fmt.Println("Role dari token:", tokenRole)

	// Tentukan template berdasarkan path
	switch path {
	case "user/Pengguna.tmpl", "Home", "/user/rumah":
		path = "Pengguna.tmpl"
	case "admin/dashboard.tmpl":
		path = "dashboard.tmpl"
	case "user/Pria":
		path = "Pria_Section.tmpl"
	case "user/Wanita":
		path = "Wanita_Section.tmpl"
	case "user/Anak":
		path = "Anak_Section.tmpl"
	case "admin/produk":
		path = "produk.tmpl"
	case "user/makan":
		path = "Pengguna.tmpl"
	case "admin/kontrol":
		path = "Kontrol_produk.tmpl"
	case "user/keranjang":
		path = "keranjang.tmpl"
	case "user/BarangDikirim":
		path = "BarangDikirim.tmpl"

	default:
		path = "Pengguna.tmpl"
	}

	// Jika role dalam token sesuai dengan role yang dikirim dalam parameter, render template
	if role == "user" || role == "User" {
		fmt.Println(username, "Dia mengakses user")
		if tokenRole == role {
			userTemplate(w, path, role)
			fmt.Println("dirender sebagai user")
		} else {
			http.Error(w, "⛔ Akses ditolak: Role tidak sesuai", http.StatusForbidden)
		}
	} else if role == "admin" || role == "Admin" {
		fmt.Println(username, "Dia mengakses admin")
		if tokenRole == role {
			adminTemplate(w, path, role)
			fmt.Println("dirender sebagai admin")
		} else {
			http.Error(w, "⛔ Akses ditolak: Role tidak sesuai", http.StatusForbidden)
		}
	}
}

// Inisialisasi route
func (server *Lingkup) InitializeRoutes() {
	server.ruter = mux.NewRouter()

	server.ruter.HandleFunc("/app/endpoint.go", func(w http.ResponseWriter, r *http.Request) {
		server.AjaxEndpoint(w, r)
	}).Methods("POST")
	// Route untuk Home (/)
	server.ruter.HandleFunc("/", controllers.Home).Methods("GET")

	// Route untuk halaman dinamis (/{page})
	server.ruter.HandleFunc("/{page}", RouterHandler).Methods("GET")

	server.ruter.HandleFunc("/role", RouterHandler).Methods("POST")

	server.ruter.HandleFunc("/{role}/{page}", func(w http.ResponseWriter, r *http.Request) {
		alamat := r.URL.Path // Contohnya "/user/page"
		vars := mux.Vars(r)

		// Ambil query parameter yang dikirim
		role_alamat := vars["role"]
		role := r.URL.Query().Get("role")
		token := r.URL.Query().Get("token")

		// Cetak hasil ke log untuk debugging
		fmt.Println("Alamat:", alamat)
		fmt.Println("Role:", role)
		fmt.Println("Token:", token)
		fmt.Println("alamat dituju:", role_alamat)

		if role_alamat != "user" && role_alamat != "admin" {
			http.Error(w, "Kamu mengakses halaman tidak valid", http.StatusInternalServerError)
		} else {
			EntityHandler(w, alamat, role, token)
		}
	}).Methods("GET")

	// Serve Static Files (CSS, JS, Images)
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	server.ruter.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
}
