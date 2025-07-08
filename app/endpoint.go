package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anan112pcmec/go_fashion/app/backend"
	backendadmin "github.com/anan112pcmec/go_fashion/app/backend-admin"
)

type AjaxRequest struct {
	Tujuan        string          `json:"tujuan"`
	Data          json.RawMessage `json:"data"`
	Nama          string          `json:"nama"`
	Jenis         string          `json:"jenis"`
	Ukuran        string          `json:"ukuran"`
	Jumlah        string          `json:"jumlah"`
	Harga         string          `json:"harga"`
	Stok          string          `json:"stok"`
	NamaUser      string          `json:"namauser"`
	IdUser        string          `json:"iduser"`
	Deskripsi     string          `json:"deskripsi"`
	IdBarang      string          `json:"id_barang"`
	JenisPakaian  string          `json:"jenis_pakaian"`
	Gender        string          `json:"gender"`
	KodeBarang    string          `json:"kode_barang"`
	TanggalLahir  string          `json:"tanggal_lahir"`
	NomorHp       string          `json:"nomor_hp"`
	Username      string          `json:"username"`
	Bio           string          `json:"bio"`
	Email         string          `json:"email"`
	AlamatLengkap any             `json:"alamat_lengkap"`
	Provinsi      string          `json:"provinsi"`
	KabupatenKota string          `json:"kabupaten_kota"`
	Kecamatan     string          `json:"kecamatan"`
	KelurahanDesa string          `json:"kelurahan_desa"`
	KodePos       string          `json:"kode_pos"`
	Koordinat     string          `json:"koordinat"`
	NamaAlamat    string          `json:"namaalamat"`
	JenisAlamat   string          `json:"jenisalamat"`
	Kategori      string          `json:"kategori"`
}

type AjaxRequestFile struct {
	Tujuan string      `json:"tujuan"`
	Data   interface{} `json:"data"`
}

func (server *Lingkup) AjaxEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Hanya menerima metode POST", http.StatusMethodNotAllowed)
		return
	}

	var req AjaxRequest

	contentType := r.Header.Get("Content-Type")

	fmt.Println(contentType)

	if strings.HasPrefix(contentType, "multipart/form-data") {
		fmt.Println("Dia mengirim multipart")

		err := r.ParseMultipartForm(10 << 20) // Maksimal 10MB
		if err != nil {
			http.Error(w, "Gagal memproses multipart form", http.StatusBadRequest)
			return
		}

		tujuan := r.FormValue("tujuan")
		if tujuan == "Masukin_barang" {
			backend.MasukanBarang(w, r, server.database)
			return
		} else if tujuan == "Masukin_Info_Pribadi" {
			backend.MasukinInfoPribadi(w, server.database, r)
			return
		}
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)

		cleanStr := string(req.Data)

		fmt.Println(cleanStr, "inilah inputnya")

		if err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		if len(req.Tujuan) != 0 {
			switch req.Tujuan {
			case "Login":
				backend.Login(w, server.database, req.Data)
			case "Ambil Data Barang Pria":
				fmt.Println("data akan dikirim ke tarikdatasemua di backend")
				backend.TarikDatSemuaPria(w, r, server.database)
			case "ambil_foto":
				backend.FotoProduk(w, req.Data, server.database)
			case "Mencari Barang Pria":
				backend.TarikDataPencarianPria(w, r, server.database, cleanStr)
				fmt.Println("mencari data barang sesuai dengan input")
			case "Mencari Barang Wanita":
				backend.TarikDatSemuaWanita(w, r, server.database)
				fmt.Println("mencari data barang sesuai dengan input")
			case "ambil_baju":
				var dataBaju []map[string]interface{}
				dataBaju = backend.AmbilBaju(w, server.database, "Semua")
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(dataBaju); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "ambil_celana":
				var dataCelana []map[string]interface{}
				dataCelana = backend.AmbilCelana(w, server.database, "Semua")
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(dataCelana); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "barang_keranjang":
				var Databaju = make(map[string]string)
				Databaju = backend.CheckBarang(w, server.database, req.Nama, req.Jenis, req.Jumlah, req.Ukuran, req.Kategori)
				// ✅ Kirim hasil dalam format JSON
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Databaju); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "masukan_barang_ke_keranjang":
				status := backend.MasukanKeranjang(w, server.database, req.Nama, req.Harga, req.Jenis, req.Stok, req.Ukuran, req.Jumlah, req.NamaUser, req.IdUser, req.Deskripsi)
				fmt.Println("ada gak nih =", req.IdUser)
				fmt.Println(status)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(status); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "mengambil_barang_dari_keranjang":
				var hasil []map[string]string
				hasil = backend.AmbilKeranjangUser(w, server.database, req.NamaUser, req.IdUser)

				fmt.Println("dibawah ini adalah hasil menarik keranjang")
				fmt.Println(hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Ambil_Data_Admin_Short":
				var Hasil []map[string]interface{}
				Hasil = backendadmin.AmbilDataShort(w, server.database, req.Jenis)
				fmt.Println("Mengambil Data Kontrol Barang")
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Ambil_Data_ui_short":
				var Hasil map[string]int
				fmt.Println("Mengambil Data Kontrol Barang")
				Hasil = backendadmin.AmbilDataUiShort(w, server.database)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}

			case "update_stok_keranjang_pemesanan":
				var Hasil map[string]string
				fmt.Println("update_stok_keranjang_pemesanan")
				fmt.Println("update_stok_keranjang_pemesanan", req.IdUser, req.Nama, req.Jumlah, req.Ukuran)
				Hasil = backend.UpdateKeranjangPengguna(w, server.database, req.IdUser, req.Nama, req.Jumlah, req.Ukuran)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Hapus_Barang_Admin_ACC":
				var Hasil map[string]interface{}
				fmt.Println("Menghapus Barang Keseluruhan")
				Hasil = backendadmin.HapusBarang(w, server.database, req.IdBarang, req.Nama, req.JenisPakaian, req.Gender, req.Harga)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "ambil_suatu_stok_baju":
				var Hasil []map[string]interface{}
				Hasil = backendadmin.AmbilDataStok(w, server.database, req.Nama, req.Gender)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Hapus_Salah_Satu_stok":
				var Hasil map[string]string
				Hasil = backendadmin.HapusBarangChild(w, server.database, req.Nama, req.IdBarang, req.Jenis, req.Ukuran, req.KodeBarang)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Ambil_Data_User_settings":
				var Hasil map[string]interface{}
				Hasil = backend.AmbilDataUser(w, server.database, req.Nama, req.IdUser)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Masukin_Alamat_Baru":
				var Hasil map[string]string
				Hasil = backend.MasukanAlamatBaru(server.database, req.AlamatLengkap, req.IdUser, req.Provinsi, req.KabupatenKota, req.Kecamatan, req.KelurahanDesa, req.KodePos, req.Koordinat, req.NamaAlamat, req.JenisAlamat)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			case "Ambil_Data_alamat_settings":
				var Hasil []map[string]interface{}
				Hasil = backend.KirimAlamatPengguna(server.database, req.IdUser)
				fmt.Println(Hasil)
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(Hasil); err != nil {
					http.Error(w, "Gagal mengencode JSON", http.StatusInternalServerError)
				}
			default:
				http.Error(w, "Tujuan tidak ditemukan", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Tujuan Request tidak jelas", http.StatusBadRequest)
		}

	}

}
