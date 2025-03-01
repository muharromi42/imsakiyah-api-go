package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type JadwalImsakiyah struct {
	Tanggal string `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

func ScrapeJadwalImsakiyah() ([]JadwalImsakiyah, error) {
	var jadwal []JadwalImsakiyah

	url := "https://www.kompas.com/ramadhan/jadwal-imsakiyah/"
	resp, err := http.Get(url)
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

	// Cari elemen tabel yang berisi jadwal
	doc.Find(".jadwal-imsakiyah tbody tr").Each(func(i int, s *goquery.Selection) {
		var row JadwalImsakiyah
		s.Find("td").Each(func(j int, td *goquery.Selection) {
			text := td.Text()
			switch j {
			case 0:
				row.Tanggal = text
			case 1:
				row.Imsak = text
			case 2:
				row.Subuh = text
			case 3:
				row.Dhuha = text
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

	log.Println("Server running on port 3050")
	r.Run(":3050")
}
