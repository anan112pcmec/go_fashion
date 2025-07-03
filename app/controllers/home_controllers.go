package controllers

import (
	"fmt"
	"net/http"

	"github.com/unrolled/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// Inisialisasi unrolled/render
	fmt.Println("Handler: Home() berjalan")
	renderer := render.New(render.Options{
		Layout: "layout",
	})

	// Data yang akan dirender
	data := map[string]interface{}{
		"title": "Home Title",
		"body":  "Home Description",
	}

	// Render halaman menggunakan renderer (unrolled/render)
	_ = renderer.HTML(w, http.StatusOK, "layout", data)

}

func Produk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ini bagian produk")
}
