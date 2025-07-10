package migration

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Nama            string `gorm:"type:varchar(255);not null"`
	Wishlist        string `gorm:"type:varchar(255)"`
	Role            string `gorm:"type:varchar(50);not null"`
	Password        string `gorm:"type:varchar(50)"`
	NamaAlamat      string `gorm:"type:varchar(250)"`
	AlamatKoordinat string `gorm:"type:varchar(250)"`
	Email           string `gorm:"type:varchar(50)"`
	NomorHp         int64
	Gender          string `gorm:"type:varchar(50)"`
	Username        string `gorm:"type:varchar(50)"`
	Bio             string `gorm:"type:text"`
	Foto            []byte `gorm:"type:bytea"`
}

type BarangCustom struct {
	ID           uint           `gorm:"primaryKey"`
	File         []byte         `gorm:"type:bytea"`
	JenisPakaian string         `gorm:"size:255"`
	Deskripsi    string         `gorm:"type:text"`
	Kategori     string         `gorm:"size:100"`
	Harga        int64          `gorm:"type:int8"`
	Warna        string         `gorm:"size:50"`
	Stok         int32          `gorm:"type:int4"`
	Ukuran       string         `gorm:"size:50"`
	Nama         string         `gorm:"size:250"`
	Gender       string         `gorm:"size:250"`
	X            int64          `gorm:"type:int8"`
	L            int64          `gorm:"type:int8"`
	M            int64          `gorm:"type:int8"`
	XL           int64          `gorm:"type:int8"`
	S            int64          `gorm:"type:int8"`
	MasukPada    time.Time      `gorm:"type:timestamp"`
	TerjualS     int64          `gorm:"type:int8"`
	TerjualM     int64          `gorm:"type:int8"`
	TerjualL     int64          `gorm:"type:int8"`
	TerjualX     int64          `gorm:"type:int8"`
	TerjualXL    int64          `gorm:"type:int8"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
