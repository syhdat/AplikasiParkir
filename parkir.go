package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

// ============================================================
// TIPE BENTUKAN (Struct)
// ============================================================

type Petugas struct {
	ID       int
	Username string
	Password string
	Nama     string
	Aktif    bool
}

type Kendaraan struct {
	NomorPolisi string
	JenisKend   string // "mobil" atau "motor"
}

type Transaksi struct {
	ID          int
	Kendaraan   Kendaraan
	WaktuMasuk  time.Time
	WaktuKeluar time.Time
	Selesai     bool
	BiayaParkir float64
	IDPetugas   int
}

// ============================================================
// KONSTANTA & VARIABEL GLOBAL (array utama)
// ============================================================

const (
	MAX_PETUGAS   = 50
	MAX_TRANSAKSI = 200

	TARIF_MOTOR_JAM1 = 2000.0
	TARIF_MOTOR_NEXT = 1000.0
	TARIF_MOBIL_JAM1 = 5000.0
	TARIF_MOBIL_NEXT = 3000.0
)

var (
	arrPetugas    [MAX_PETUGAS]Petugas
	jumlahPetugas int

	arrTransaksi    [MAX_TRANSAKSI]Transaksi
	jumlahTransaksi int

	petugasLogin *Petugas // nil = belum login
	isAdmin      bool
	reader       = bufio.NewReader(os.Stdin)

	nextIDPetugas   = 1
	nextIDTransaksi = 1
)

// ============================================================
// FUNGSI UTILITAS
// ============================================================

func inputString(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func inputInt(prompt string) int {
	var n int
	fmt.Print(prompt)
	fmt.Scan(&n)
	reader.ReadString('\n')
	return n
}

func inputFloat(prompt string) float64 {
	var f float64
	fmt.Print(prompt)
	fmt.Scan(&f)
	reader.ReadString('\n')
	return f
}

func cetakGaris(karakter string, panjang int) {
	for i := 0; i < panjang; i++ {
		fmt.Print(karakter)
	}
	fmt.Println()
}

func cetakJudul(judul string) {
	cetakGaris("=", 60)
	fmt.Printf("  %s\n", judul)
	cetakGaris("=", 60)
}

func tekanEnterLanjut() {
	fmt.Print("\nTekan Enter untuk melanjutkan...")
	reader.ReadString('\n')
}

// hitungDurasi: fungsi menghitung selisih jam (pembulatan ke atas) — rekursif
func hitungJamRekursif(menit int) int {
	if menit <= 0 {
		return 0
	}
	if menit <= 60 {
		return 1
	}
	return 1 + hitungJamRekursif(menit-60)
}

// hitungBiaya: fungsi hitung biaya parkir
func hitungBiaya(jenis string, waktuMasuk, waktuKeluar time.Time) float64 {
	durasi := waktuKeluar.Sub(waktuMasuk)
	menitTotal := int(math.Ceil(durasi.Minutes()))
	if menitTotal < 1 {
		menitTotal = 1
	}
	jam := hitungJamRekursif(menitTotal)

	var biaya float64
	if jenis == "motor" {
		if jam <= 1 {
			biaya = TARIF_MOTOR_JAM1
		} else {
			biaya = TARIF_MOTOR_JAM1 + float64(jam-1)*TARIF_MOTOR_NEXT
		}
	} else {
		if jam <= 1 {
			biaya = TARIF_MOBIL_JAM1
		} else {
			biaya = TARIF_MOBIL_JAM1 + float64(jam-1)*TARIF_MOBIL_NEXT
		}
	}
	return biaya
}

// formatRupiah: fungsi format angka ke rupiah
func formatRupiah(nominal float64) string {
	return fmt.Sprintf("Rp %.0f", nominal)
}

// formatWaktu: fungsi format waktu
func formatWaktu(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006 15:04")
}

// ============================================================
// SEARCH — Sequential Search
// ============================================================

// sequentialSearchPetugas: cari petugas berdasarkan username (return index, -1 jika tidak ada)
func sequentialSearchPetugas(username string) int {
	found := -1
	i := 0
	for i < jumlahPetugas && found == -1 {
		if arrPetugas[i].Username == username {
			found = i
		}
		i++
	}
	return found
}

// sequentialSearchTransaksiByNopol: cari SEMUA transaksi aktif berdasarkan nomor polisi
func sequentialSearchTransaksiByNopol(nopol string) int {
	found := -1
	i := 0
	for i < jumlahTransaksi && found == -1 {
		if strings.EqualFold(arrTransaksi[i].Kendaraan.NomorPolisi, nopol) && !arrTransaksi[i].Selesai {
			found = i
		}
		i++
	}
	return found
}

// ============================================================
// SEARCH — Binary Search (array harus terurut berdasarkan ID)
// ============================================================

// binarySearchPetugasByID: cari petugas berdasarkan ID (array harus terurut ID asc)
func binarySearchPetugasByID(id int) int {
	low := 0
	high := jumlahPetugas - 1
	found := -1
	for low <= high && found == -1 {
		mid := (low + high) / 2
		if arrPetugas[mid].ID == id {
			found = mid
		} else if arrPetugas[mid].ID < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return found
}

// binarySearchTransaksiByID: cari transaksi berdasarkan ID
func binarySearchTransaksiByID(id int) int {
	low := 0
	high := jumlahTransaksi - 1
	found := -1
	for low <= high && found == -1 {
		mid := (low + high) / 2
		if arrTransaksi[mid].ID == id {
			found = mid
		} else if arrTransaksi[mid].ID < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return found
}

// ============================================================
// SORT — Selection Sort
// ============================================================

// selectionSortTransaksiByWaktu: urutkan salinan transaksi berdasarkan waktu masuk
func selectionSortTransaksiByWaktu(arr []Transaksi, ascending bool) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		idx := i
		for j := i + 1; j < n; j++ {
			if ascending {
				if arr[j].WaktuMasuk.Before(arr[idx].WaktuMasuk) {
					idx = j
				}
			} else {
				if arr[j].WaktuMasuk.After(arr[idx].WaktuMasuk) {
					idx = j
				}
			}
		}
		arr[i], arr[idx] = arr[idx], arr[i]
	}
}

// selectionSortTransaksiByBiaya: urutkan berdasarkan biaya parkir
func selectionSortTransaksiByBiaya(arr []Transaksi, ascending bool) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		idx := i
		for j := i + 1; j < n; j++ {
			if ascending {
				if arr[j].BiayaParkir < arr[idx].BiayaParkir {
					idx = j
				}
			} else {
				if arr[j].BiayaParkir > arr[idx].BiayaParkir {
					idx = j
				}
			}
		}
		arr[i], arr[idx] = arr[idx], arr[i]
	}
}

// ============================================================
// SORT — Insertion Sort
// ============================================================

// insertionSortTransaksiByNopol: urutkan berdasarkan nomor polisi (alfabet)
func insertionSortTransaksiByNopol(arr []Transaksi, ascending bool) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		if ascending {
			for j >= 0 && arr[j].Kendaraan.NomorPolisi > key.Kendaraan.NomorPolisi {
				arr[j+1] = arr[j]
				j--
			}
		} else {
			for j >= 0 && arr[j].Kendaraan.NomorPolisi < key.Kendaraan.NomorPolisi {
				arr[j+1] = arr[j]
				j--
			}
		}
		arr[j+1] = key
	}
}

// insertionSortPetugasByNama: urutkan petugas berdasarkan nama
func insertionSortPetugasByNama(arr []Petugas, ascending bool) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		if ascending {
			for j >= 0 && arr[j].Nama > key.Nama {
				arr[j+1] = arr[j]
				j--
			}
		} else {
			for j >= 0 && arr[j].Nama < key.Nama {
				arr[j+1] = arr[j]
				j--
			}
		}
		arr[j+1] = key
	}
}

// ============================================================
// PROSEDUR — Inisialisasi Data Awal
// ============================================================

func inisialisasiData() {
	// Admin default
	arrPetugas[0] = Petugas{
		ID:       nextIDPetugas,
		Username: "admin",
		Password: "admin123",
		Nama:     "Administrator",
		Aktif:    true,
	}
	nextIDPetugas++
	jumlahPetugas++

	// Petugas contoh
	arrPetugas[1] = Petugas{
		ID:       nextIDPetugas,
		Username: "budi",
		Password: "budi123",
		Nama:     "Budi Santoso",
		Aktif:    true,
	}
	nextIDPetugas++
	jumlahPetugas++
}

// ============================================================
// MODUL LOGIN
// ============================================================

func login() bool {
	cetakJudul("LOGIN APLIKASI PARKIR MALL")
	username := inputString("Username : ")
	password := inputString("Password : ")

	idx := sequentialSearchPetugas(username)
	if idx == -1 {
		fmt.Println("\n[!] Username tidak ditemukan.")
		return false
	}
	if arrPetugas[idx].Password != password {
		fmt.Println("\n[!] Password salah.")
		return false
	}
	if !arrPetugas[idx].Aktif {
		fmt.Println("\n[!] Akun tidak aktif.")
		return false
	}

	petugasLogin = &arrPetugas[idx]
	isAdmin = (username == "admin")
	fmt.Printf("\n[✓] Selamat datang, %s!\n", petugasLogin.Nama)
	tekanEnterLanjut()
	return true
}

func logout() {
	fmt.Printf("\n[✓] %s telah logout.\n", petugasLogin.Nama)
	petugasLogin = nil
	isAdmin = false
	tekanEnterLanjut()
}

// ============================================================
// MODUL ADMIN — Kelola Petugas
// ============================================================

func tampilMenuAdmin() {
	for {
		cetakJudul("MENU ADMIN")
		fmt.Println("  1. Tambah Petugas")
		fmt.Println("  2. Edit Petugas")
		fmt.Println("  3. Hapus Petugas")
		fmt.Println("  4. Daftar Petugas")
		fmt.Println("  5. Cari Petugas (by ID)")
		fmt.Println("  0. Kembali")
		cetakGaris("-", 60)
		pilihan := inputInt("Pilihan : ")
		switch pilihan {
		case 1:
			tambahPetugas()
		case 2:
			editPetugas()
		case 3:
			hapusPetugas()
		case 4:
			daftarPetugas()
		case 5:
			cariPetugasByID()
		case 0:
			return
		default:
			fmt.Println("[!] Pilihan tidak valid.")
			tekanEnterLanjut()
		}
	}
}

func tambahPetugas() {
	cetakJudul("TAMBAH PETUGAS")
	if jumlahPetugas >= MAX_PETUGAS {
		fmt.Println("[!] Data petugas sudah penuh.")
		tekanEnterLanjut()
		return
	}
	nama := inputString("Nama lengkap : ")
	username := inputString("Username     : ")

	// cek username unik
	if sequentialSearchPetugas(username) != -1 {
		fmt.Println("[!] Username sudah digunakan.")
		tekanEnterLanjut()
		return
	}

	password := inputString("Password     : ")

	p := Petugas{
		ID:       nextIDPetugas,
		Username: username,
		Password: password,
		Nama:     nama,
		Aktif:    true,
	}
	arrPetugas[jumlahPetugas] = p
	jumlahPetugas++
	nextIDPetugas++
	fmt.Printf("\n[✓] Petugas '%s' berhasil ditambahkan (ID: %d).\n", nama, p.ID)
	tekanEnterLanjut()
}

func editPetugas() {
	cetakJudul("EDIT PETUGAS")
	id := inputInt("Masukkan ID Petugas : ")
	idx := binarySearchPetugasByID(id)
	if idx == -1 {
		fmt.Println("[!] Petugas tidak ditemukan.")
		tekanEnterLanjut()
		return
	}
	if arrPetugas[idx].Username == "admin" {
		fmt.Println("[!] Akun admin tidak bisa diubah.")
		tekanEnterLanjut()
		return
	}

	fmt.Printf("Nama lama     : %s\n", arrPetugas[idx].Nama)
	nama := inputString("Nama baru (kosong=skip)     : ")
	password := inputString("Password baru (kosong=skip) : ")

	if nama != "" {
		arrPetugas[idx].Nama = nama
	}
	if password != "" {
		arrPetugas[idx].Password = password
	}
	fmt.Println("[✓] Data petugas berhasil diperbarui.")
	tekanEnterLanjut()
}

func hapusPetugas() {
	cetakJudul("HAPUS PETUGAS")
	id := inputInt("Masukkan ID Petugas : ")
	idx := binarySearchPetugasByID(id)
	if idx == -1 {
		fmt.Println("[!] Petugas tidak ditemukan.")
		tekanEnterLanjut()
		return
	}
	if arrPetugas[idx].Username == "admin" {
		fmt.Println("[!] Akun admin tidak bisa dihapus.")
		tekanEnterLanjut()
		return
	}

	konfirmasi := inputString(fmt.Sprintf("Hapus petugas '%s'? (y/n) : ", arrPetugas[idx].Nama))
	if strings.ToLower(konfirmasi) == "y" {
		// geser array
		for i := idx; i < jumlahPetugas-1; i++ {
			arrPetugas[i] = arrPetugas[i+1]
		}
		arrPetugas[jumlahPetugas-1] = Petugas{}
		jumlahPetugas--
		fmt.Println("[✓] Petugas berhasil dihapus.")
	} else {
		fmt.Println("[i] Penghapusan dibatalkan.")
	}
	tekanEnterLanjut()
}

func daftarPetugas() {
	cetakJudul("DAFTAR PETUGAS")
	fmt.Println("Urutkan berdasarkan:")
	fmt.Println("  1. Nama A-Z")
	fmt.Println("  2. Nama Z-A")
	fmt.Println("  3. ID (default)")
	pilihan := inputInt("Pilihan : ")

	// buat salinan slice
	salinan := make([]Petugas, jumlahPetugas)
	for i := 0; i < jumlahPetugas; i++ {
		salinan[i] = arrPetugas[i]
	}

	switch pilihan {
	case 1:
		insertionSortPetugasByNama(salinan, true)
	case 2:
		insertionSortPetugasByNama(salinan, false)
	}

	cetakGaris("-", 60)
	fmt.Printf("%-5s %-15s %-20s %-6s\n", "ID", "Username", "Nama", "Aktif")
	cetakGaris("-", 60)
	for i := 0; i < len(salinan); i++ {
		p := salinan[i]
		aktif := "Ya"
		if !p.Aktif {
			aktif = "Tidak"
		}
		fmt.Printf("%-5d %-15s %-20s %-6s\n", p.ID, p.Username, p.Nama, aktif)
	}
	cetakGaris("-", 60)
	tekanEnterLanjut()
}

func cariPetugasByID() {
	cetakJudul("CARI PETUGAS (Binary Search by ID)")
	id := inputInt("Masukkan ID Petugas : ")
	idx := binarySearchPetugasByID(id)
	if idx == -1 {
		fmt.Println("[!] Petugas dengan ID tersebut tidak ditemukan.")
	} else {
		p := arrPetugas[idx]
		fmt.Printf("\nID       : %d\n", p.ID)
		fmt.Printf("Nama     : %s\n", p.Nama)
		fmt.Printf("Username : %s\n", p.Username)
		aktif := "Aktif"
		if !p.Aktif {
			aktif = "Tidak Aktif"
		}
		fmt.Printf("Status   : %s\n", aktif)
	}
	tekanEnterLanjut()
}

// ============================================================
// MODUL PETUGAS — Kelola Transaksi Parkir
// ============================================================

func tampilMenuPetugas() {
	for {
		cetakJudul(fmt.Sprintf("MENU PETUGAS — %s", petugasLogin.Nama))
		fmt.Println("  1. Kendaraan Masuk (Tambah Transaksi)")
		fmt.Println("  2. Kendaraan Keluar (Selesaikan Transaksi)")
		fmt.Println("  3. Edit Transaksi")
		fmt.Println("  4. Hapus Transaksi")
		fmt.Println("  5. Cari Kendaraan (by Nopol)")
		fmt.Println("  6. Cetak Daftar Kendaraan")
		fmt.Println("  7. Total Uang Parkir Hari Ini")
		fmt.Println("  0. Logout")
		cetakGaris("-", 60)
		pilihan := inputInt("Pilihan : ")
		switch pilihan {
		case 1:
			kendaraanMasuk()
		case 2:
			kendaraanKeluar()
		case 3:
			editTransaksi()
		case 4:
			hapusTransaksi()
		case 5:
			cariKendaraanByNopol()
		case 6:
			cetakDaftarKendaraan()
		case 7:
			totalUangHariIni()
		case 0:
			logout()
			return
		default:
			fmt.Println("[!] Pilihan tidak valid.")
			tekanEnterLanjut()
		}
	}
}

func kendaraanMasuk() {
	cetakJudul("KENDARAAN MASUK")
	if jumlahTransaksi >= MAX_TRANSAKSI {
		fmt.Println("[!] Data transaksi penuh.")
		tekanEnterLanjut()
		return
	}

	nopol := strings.ToUpper(inputString("Nomor Polisi : "))
	// cek apakah sudah ada transaksi aktif untuk nopol ini
	if sequentialSearchTransaksiByNopol(nopol) != -1 {
		fmt.Println("[!] Kendaraan dengan nomor polisi tersebut masih ada di dalam.")
		tekanEnterLanjut()
		return
	}

	fmt.Println("Jenis Kendaraan:")
	fmt.Println("  1. Motor")
	fmt.Println("  2. Mobil")
	jenis := inputInt("Pilihan : ")
	jenisStr := "motor"
	if jenis == 2 {
		jenisStr = "mobil"
	}

	t := Transaksi{
		ID: nextIDTransaksi,
		Kendaraan: Kendaraan{
			NomorPolisi: nopol,
			JenisKend:   jenisStr,
		},
		WaktuMasuk: time.Now(),
		Selesai:    false,
		IDPetugas:  petugasLogin.ID,
	}
	arrTransaksi[jumlahTransaksi] = t
	jumlahTransaksi++
	nextIDTransaksi++

	fmt.Printf("\n[✓] Kendaraan %s (%s) masuk pukul %s\n", nopol, jenisStr, formatWaktu(t.WaktuMasuk))
	fmt.Printf("    ID Transaksi: %d\n", t.ID)
	tekanEnterLanjut()
}

func kendaraanKeluar() {
	cetakJudul("KENDARAAN KELUAR")
	nopol := strings.ToUpper(inputString("Nomor Polisi : "))
	idx := sequentialSearchTransaksiByNopol(nopol)
	if idx == -1 {
		fmt.Println("[!] Kendaraan tidak ditemukan atau sudah keluar.")
		tekanEnterLanjut()
		return
	}

	waktuKeluar := time.Now()
	biaya := hitungBiaya(
		arrTransaksi[idx].Kendaraan.JenisKend,
		arrTransaksi[idx].WaktuMasuk,
		waktuKeluar,
	)

	arrTransaksi[idx].WaktuKeluar = waktuKeluar
	arrTransaksi[idx].BiayaParkir = biaya
	arrTransaksi[idx].Selesai = true

	durasi := waktuKeluar.Sub(arrTransaksi[idx].WaktuMasuk)
	jam := int(durasi.Hours())
	menit := int(durasi.Minutes()) % 60

	fmt.Println("\n====== STRUK PARKIR ======")
	fmt.Printf("Nopol       : %s\n", arrTransaksi[idx].Kendaraan.NomorPolisi)
	fmt.Printf("Jenis       : %s\n", arrTransaksi[idx].Kendaraan.JenisKend)
	fmt.Printf("Waktu Masuk : %s\n", formatWaktu(arrTransaksi[idx].WaktuMasuk))
	fmt.Printf("Waktu Keluar: %s\n", formatWaktu(waktuKeluar))
	fmt.Printf("Durasi      : %d jam %d menit\n", jam, menit)
	fmt.Printf("Biaya       : %s\n", formatRupiah(biaya))
	fmt.Println("==========================")
	tekanEnterLanjut()
}

func editTransaksi() {
	cetakJudul("EDIT TRANSAKSI")
	id := inputInt("ID Transaksi : ")
	idx := binarySearchTransaksiByID(id)
	if idx == -1 {
		fmt.Println("[!] Transaksi tidak ditemukan.")
		tekanEnterLanjut()
		return
	}
	if arrTransaksi[idx].Selesai {
		fmt.Println("[!] Transaksi sudah selesai, tidak bisa diedit.")
		tekanEnterLanjut()
		return
	}

	fmt.Printf("Nopol lama : %s\n", arrTransaksi[idx].Kendaraan.NomorPolisi)
	nopol := strings.ToUpper(inputString("Nopol baru (kosong=skip) : "))
	if nopol != "" {
		arrTransaksi[idx].Kendaraan.NomorPolisi = nopol
	}

	fmt.Println("Jenis (1=Motor, 2=Mobil, 0=Skip) :")
	jenis := inputInt("Pilihan : ")
	if jenis == 1 {
		arrTransaksi[idx].Kendaraan.JenisKend = "motor"
	} else if jenis == 2 {
		arrTransaksi[idx].Kendaraan.JenisKend = "mobil"
	}

	fmt.Println("[✓] Transaksi berhasil diperbarui.")
	tekanEnterLanjut()
}

func hapusTransaksi() {
	cetakJudul("HAPUS TRANSAKSI")
	id := inputInt("ID Transaksi : ")
	idx := binarySearchTransaksiByID(id)
	if idx == -1 {
		fmt.Println("[!] Transaksi tidak ditemukan.")
		tekanEnterLanjut()
		return
	}

	konfirmasi := inputString(fmt.Sprintf("Hapus transaksi ID %d (nopol: %s)? (y/n) : ",
		id, arrTransaksi[idx].Kendaraan.NomorPolisi))
	if strings.ToLower(konfirmasi) == "y" {
		for i := idx; i < jumlahTransaksi-1; i++ {
			arrTransaksi[i] = arrTransaksi[i+1]
		}
		arrTransaksi[jumlahTransaksi-1] = Transaksi{}
		jumlahTransaksi--
		fmt.Println("[✓] Transaksi berhasil dihapus.")
	} else {
		fmt.Println("[i] Penghapusan dibatalkan.")
	}
	tekanEnterLanjut()
}

func cariKendaraanByNopol() {
	cetakJudul("CARI KENDARAAN (Sequential Search by Nopol)")
	nopol := strings.ToUpper(inputString("Nomor Polisi : "))

	ketemu := false
	for i := 0; i < jumlahTransaksi; i++ {
		if strings.EqualFold(arrTransaksi[i].Kendaraan.NomorPolisi, nopol) {
			t := arrTransaksi[i]
			status := "Masih di dalam"
			if t.Selesai {
				status = "Sudah keluar"
			}
			fmt.Println(strings.Repeat("-", 50))
			fmt.Printf("ID Transaksi : %d\n", t.ID)
			fmt.Printf("Nopol        : %s\n", t.Kendaraan.NomorPolisi)
			fmt.Printf("Jenis        : %s\n", t.Kendaraan.JenisKend)
			fmt.Printf("Waktu Masuk  : %s\n", formatWaktu(t.WaktuMasuk))
			fmt.Printf("Waktu Keluar : %s\n", formatWaktu(t.WaktuKeluar))
			fmt.Printf("Status       : %s\n", status)
			if t.Selesai {
				fmt.Printf("Biaya        : %s\n", formatRupiah(t.BiayaParkir))
			}
			ketemu = true
		}
	}

	if !ketemu {
		fmt.Println("[!] Tidak ada transaksi dengan nomor polisi tersebut.")
	}
	tekanEnterLanjut()
}

func cetakDaftarKendaraan() {
	cetakJudul("CETAK DAFTAR KENDARAAN")
	fmt.Println("Filter:")
	fmt.Println("  1. Semua kendaraan")
	fmt.Println("  2. Hanya Mobil")
	fmt.Println("  3. Hanya Motor")
	filter := inputInt("Pilihan filter : ")

	fmt.Println("\nUrutkan berdasarkan:")
	fmt.Println("  1. Waktu masuk (terlama → terbaru)")
	fmt.Println("  2. Waktu masuk (terbaru → terlama)")
	fmt.Println("  3. Nomor Polisi (A-Z)")
	fmt.Println("  4. Nomor Polisi (Z-A)")
	fmt.Println("  5. Biaya (terkecil → terbesar)")
	fmt.Println("  6. Biaya (terbesar → terkecil)")
	urutan := inputInt("Pilihan urutan : ")

	// Kumpulkan data sesuai filter
	var data []Transaksi
	for i := 0; i < jumlahTransaksi; i++ {
		t := arrTransaksi[i]
		masuk := true
		if filter == 2 && t.Kendaraan.JenisKend != "mobil" {
			masuk = false
		}
		if filter == 3 && t.Kendaraan.JenisKend != "motor" {
			masuk = false
		}
		if masuk {
			data = append(data, t)
		}
	}

	// Sort
	switch urutan {
	case 1:
		selectionSortTransaksiByWaktu(data, true)
	case 2:
		selectionSortTransaksiByWaktu(data, false)
	case 3:
		insertionSortTransaksiByNopol(data, true)
	case 4:
		insertionSortTransaksiByNopol(data, false)
	case 5:
		selectionSortTransaksiByBiaya(data, true)
	case 6:
		selectionSortTransaksiByBiaya(data, false)
	}

	cetakGaris("-", 80)
	fmt.Printf("%-5s %-12s %-7s %-17s %-17s %-8s %-12s\n",
		"ID", "Nopol", "Jenis", "Waktu Masuk", "Waktu Keluar", "Status", "Biaya")
	cetakGaris("-", 80)

	var totalUang float64
	for _, t := range data {
		status := "Di dalam"
		biayaStr := "-"
		if t.Selesai {
			status = "Keluar"
			biayaStr = formatRupiah(t.BiayaParkir)
			totalUang += t.BiayaParkir
		}
		fmt.Printf("%-5d %-12s %-7s %-17s %-17s %-8s %-12s\n",
			t.ID,
			t.Kendaraan.NomorPolisi,
			t.Kendaraan.JenisKend,
			formatWaktu(t.WaktuMasuk),
			formatWaktu(t.WaktuKeluar),
			status,
			biayaStr,
		)
	}
	cetakGaris("-", 80)
	fmt.Printf("Total kendaraan : %d | Total pendapatan (selesai): %s\n", len(data), formatRupiah(totalUang))
	tekanEnterLanjut()
}

func totalUangHariIni() {
	cetakJudul("TOTAL UANG PARKIR HARI INI")
	sekarang := time.Now()
	tahun, bulan, hari := sekarang.Date()

	var totalMobil, totalMotor float64
	var countMobil, countMotor int

	for i := 0; i < jumlahTransaksi; i++ {
		t := arrTransaksi[i]
		if !t.Selesai {
			continue
		}
		ty, tb, td := t.WaktuKeluar.Date()
		if ty == tahun && tb == bulan && td == hari {
			if t.Kendaraan.JenisKend == "mobil" {
				totalMobil += t.BiayaParkir
				countMobil++
			} else {
				totalMotor += t.BiayaParkir
				countMotor++
			}
		}
	}

	tglStr := sekarang.Format("02 January 2006")
	fmt.Printf("Tanggal       : %s\n", tglStr)
	cetakGaris("-", 40)
	fmt.Printf("Mobil  : %3d kendaraan → %s\n", countMobil, formatRupiah(totalMobil))
	fmt.Printf("Motor  : %3d kendaraan → %s\n", countMotor, formatRupiah(totalMotor))
	cetakGaris("-", 40)
	fmt.Printf("TOTAL  : %3d kendaraan → %s\n", countMobil+countMotor, formatRupiah(totalMobil+totalMotor))
	tekanEnterLanjut()
}

// ============================================================
// MAIN
// ============================================================

func main() {
	inisialisasiData()

	for {
		if petugasLogin == nil {
			berhasil := login()
			if !berhasil {
				tekanEnterLanjut()
				continue
			}
		}

		if isAdmin {
			// Admin dapat akses menu admin DAN menu petugas
			cetakJudul("PILIH MODE")
			fmt.Println("  1. Menu Admin (Kelola Petugas)")
			fmt.Println("  2. Menu Petugas (Transaksi Parkir)")
			fmt.Println("  0. Logout")
			cetakGaris("-", 60)
			pilihan := inputInt("Pilihan : ")
			switch pilihan {
			case 1:
				tampilMenuAdmin()
			case 2:
				tampilMenuPetugas()
			case 0:
				logout()
			default:
				fmt.Println("[!] Pilihan tidak valid.")
				tekanEnterLanjut()
			}
		} else {
			tampilMenuPetugas()
		}
	}
}
