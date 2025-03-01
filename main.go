package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

type JadwalImsakiyah struct {
	Hari    string `json:"hari"`
	Tanggal string `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

func ScrapeJadwalImsakiyah() ([]JadwalImsakiyah, error) {
	var jadwal []JadwalImsakiyah

	url := "https://www.kompas.com/ramadhan/jadwal-imsakiyah/kota-palembang"

	// Buat request dengan User-Agent agar tidak diblokir
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Scrape data dari tabel dengan class "w-full"
	doc.Find("table.w-full tbody tr").Each(func(i int, s *goquery.Selection) {
		var row JadwalImsakiyah
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			text := strings.TrimSpace(td.Text())
			switch j {
			case 0:
				row.Hari = text
			case 1:
				row.Tanggal = text
			case 2:
				row.Imsak = text
			case 3:
				row.Subuh = text
			case 4:
				row.Dzuhur = text
			case 5:
				row.Ashar = text
			case 6:
				row.Maghrib = text
			case 7:
				row.Isya = text
			}
		})
		jadwal = append(jadwal, row)
	})

	return jadwal, nil
}

// Fungsi untuk mendapatkan jadwal hari ini
func ScrapeJadwalHariIni() (*JadwalImsakiyah, error) {
	jadwal, err := ScrapeJadwalImsakiyah()
	if err != nil {
		return nil, err
	}

	// Format tanggal hari ini sesuai dengan format di situs (misal: "Sabtu, 01/03")
	today := time.Now().Format("Monday, 02/01")
	today = convertIndonesianDay(today) // Konversi nama hari ke bahasa Indonesia

	// Cari data yang sesuai dengan tanggal hari ini
	for _, j := range jadwal {
		if j.Tanggal == today {
			return &j, nil
		}
	}

	return nil, fmt.Errorf("Jadwal hari ini tidak ditemukan")
}

// Konversi nama hari dari English ke Indonesian (karena format tanggal di web Kompas)
func convertIndonesianDay(dateStr string) string {
	indonesianDays := map[string]string{
		"Monday":    "Senin",
		"Tuesday":   "Selasa",
		"Wednesday": "Rabu",
		"Thursday":  "Kamis",
		"Friday":    "Jumat",
		"Saturday":  "Sabtu",
		"Sunday":    "Minggu",
	}

	for en, id := range indonesianDays {
		dateStr = strings.ReplaceAll(dateStr, en, id)
	}

	return dateStr
}

func main() {
	r := gin.Default()

	r.GET("/imsakiyah", func(c *gin.Context) {
		jadwal, err := ScrapeJadwalImsakiyah()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, jadwal)
	})

	// Endpoint untuk mendapatkan jadwal hari ini
	r.GET("/imsakiyah/today", func(c *gin.Context) {
		jadwal, err := ScrapeJadwalHariIni()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, jadwal)
	})

	log.Println("Server running on port 3050")
	r.Run(":3050")
}
