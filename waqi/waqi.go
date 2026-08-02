package waqi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"absen/redis" // Package redis kamu

	"github.com/gin-gonic/gin"
)

// ======================================================
// 1. DTO (Data Transfer Object)
// ======================================================

type CleanAQIResponse struct {
	Location          string             `json:"location"`
	AQI               int                `json:"aqi"`
	Status            string             `json:"status"`
	BadgeColor        string             `json:"badge_color"`
	DominantPollutant string             `json:"dominant_pollutant"`
	HealthAdvice      string             `json:"health_advice"`
	Weather           WeatherInfo        `json:"weather"`
	Coordinates       map[string]float64 `json:"coordinates"`
	UpdatedAt         string             `json:"updated_at"`
	Source            string             `json:"source"`
}

type WeatherInfo struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

type StationDropdown struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"` // Slug nama stasiun / kota untuk query feed
}

// WAQI Bounds Raw Struct (Untuk menangkap Big Data Global)
type waqiBoundsResponse struct {
	Status string `json:"status"`
	Data   []struct {
		UID     int    `json:"uid"`
		Aqi     string `json:"aqi"`
		Station struct {
			Name string `json:"name"`
		} `json:"station"`
	} `json:"data"`
}

// ======================================================
// 2. HTTP HANDLERS
// ======================================================

// GetCitiesDropdownHandler: AMBIL BIG DATA STASIUN GLOBAL (CACHED 5 JAM DI REDIS)
func GetCitiesDropdownHandler(c *gin.Context) {
	cacheKey := "waqi:global_stations_dropdown"

	// A. Cek apakah Big Data Stasiun ada di Redis Cache
	cachedData, err := redis.Getc(cacheKey)
	if err == nil && cachedData != "" {
		var cachedStations []StationDropdown
		if json.Unmarshal([]byte(cachedData), &cachedStations) == nil {
			c.JSON(http.StatusOK, gin.H{
				"source":   "⚡ Redis Cache (TTL 5 Jam)",
				"total":    len(cachedStations),
				"stations": cachedStations,
			})
			return
		}
	}

	// B. Jika Cache Miss: Fetch Big Data Global dari WAQI Bounds API
	token := os.Getenv("WAQI_TOKEN")
	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WAQI_TOKEN belum dikonfigurasi!"})
		return
	}

	// Bounding Box Seluruh Dunia (-90,-180 ke 90,180)
	globalBoundsURL := fmt.Sprintf("https://api.waqi.info/map/bounds/?token=%s&latlng=-90,-180,90,180", token)

	resp, err := http.Get(globalBoundsURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung ke WAQI API"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	var boundsResp waqiBoundsResponse
	if err := json.Unmarshal(bodyBytes, &boundsResp); err != nil || boundsResp.Status != "ok" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses Big Data dari WAQI"})
		return
	}

	// C. Transform ke Format Ringkas untuk Dropdown Frontend
	var stations []StationDropdown
	for _, item := range boundsResp.Data {
		if item.Station.Name != "" {
			stations = append(stations, StationDropdown{
				ID:   item.UID,
				Name: item.Station.Name,
				City: fmt.Sprintf("@%d", item.UID), // Menggunakan UID stasiun agar 100% presisi & tidak nyasar
			})
		}
	}

	// D. Simpan Big Data ke Redis Cache selama 5 JAM
	jsonToCache, errMarshal := json.Marshal(stations)
	if errMarshal == nil {
		_ = redis.SetC(cacheKey, string(jsonToCache), 5*time.Hour)
	}

	c.JSON(http.StatusOK, gin.H{
		"source":   "🌐 WAQI API Official (Big Data Fresh Fetch)",
		"total":    len(stations),
		"stations": stations,
	})
}

// AirQualityHandler: DIRECT FETCH TANPA CACHE SAAT USER MEMILIH LOKASI
func AirQualityHandler(c *gin.Context) {
	token := os.Getenv("WAQI_TOKEN")
	if token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WAQI_TOKEN belum dikonfigurasi di environment!"})
		return
	}

	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	city := c.Query("city")

	// 1. Buat URL Request (Pakai ID Stasiun / Koordinat GPS)
	apiURL, err := buildWAQIUrl(latStr, lngStr, city, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Direct Fetch ke WAQI API (Tanpa Redis)
	waqiData, err := fetchFromWAQI(apiURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Format Response
	response := formatAQIResponse(waqiData)

	c.JSON(http.StatusOK, response)
}

// ======================================================
// 3. SERVICE & CLIENT LOGIC
// ======================================================

type waqiRawResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI         int    `json:"aqi"`
		DominentPol string `json:"dominentpol"`
		City        struct {
			Name string    `json:"name"`
			Geo  []float64 `json:"geo"`
		} `json:"city"`
		IAQI struct {
			T struct{ V float64 `json:"v"` } `json:"t"`
			H struct{ V float64 `json:"v"` } `json:"h"`
		} `json:"iaqi"`
		Time struct {
			S string `json:"s"`
		} `json:"time"`
	} `json:"data"`
}

func buildWAQIUrl(latStr, lngStr, city, token string) (string, error) {
	// Mode 1: GPS
	if latStr != "" && lngStr != "" {
		if _, errLat := strconv.ParseFloat(latStr, 64); errLat != nil {
			return "", fmt.Errorf("koordinat latitude tidak valid")
		}
		if _, errLng := strconv.ParseFloat(lngStr, 64); errLng != nil {
			return "", fmt.Errorf("koordinat longitude tidak valid")
		}
		return fmt.Sprintf("https://api.waqi.info/feed/geo:%s;%s/?token=%s", latStr, lngStr, token), nil
	}

	// Mode 2: Pilih Stasiun / Kota dari Dropdown
	if city != "" {
		return fmt.Sprintf("https://api.waqi.info/feed/%s/?token=%s", city, token), nil
	}

	// Mode 3: Default (Yogyakarta)
	return fmt.Sprintf("https://api.waqi.info/feed/yogyakarta/?token=%s", token), nil
}

func fetchFromWAQI(url string) (*waqiRawResponse, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal terhubung ke layanan WAQI")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca data dari server")
	}

	var waqiResp waqiRawResponse
	if err := json.Unmarshal(bodyBytes, &waqiResp); err != nil || waqiResp.Status != "ok" {
		return nil, fmt.Errorf("data kualitas udara tidak ditemukan untuk lokasi ini")
	}

	return &waqiResp, nil
}

// ======================================================
// 4. HELPER LOGIC
// ======================================================

func formatAQIResponse(raw *waqiRawResponse) CleanAQIResponse {
	status, color, advice := getAQIMetadata(raw.Data.AQI)

	var latVal, lngVal float64
	if len(raw.Data.City.Geo) >= 2 {
		latVal = raw.Data.City.Geo[0]
		lngVal = raw.Data.City.Geo[1]
	}

	tempStr := "N/A"
	if raw.Data.IAQI.T.V != 0 {
		tempStr = fmt.Sprintf("%.1f°C", raw.Data.IAQI.T.V)
	}

	humidityStr := "N/A"
	if raw.Data.IAQI.H.V != 0 {
		humidityStr = fmt.Sprintf("%.0f%%", raw.Data.IAQI.H.V)
	}

	return CleanAQIResponse{
		Location:          raw.Data.City.Name,
		AQI:               raw.Data.AQI,
		Status:            status,
		BadgeColor:        color,
		DominantPollutant: raw.Data.DominentPol,
		HealthAdvice:      advice,
		Weather: WeatherInfo{
			Temperature: tempStr,
			Humidity:    humidityStr,
		},
		Coordinates: map[string]float64{
			"latitude":  latVal,
			"longitude": lngVal,
		},
		UpdatedAt: raw.Data.Time.S,
		Source:    "🌐 WAQI API Official (Live Direct)",
	}
}

func getAQIMetadata(aqi int) (status string, color string, advice string) {
	switch {
	case aqi <= 50:
		return "Baik", "#2ECC71", "Udara sangat bersih. Sangat aman untuk beraktivitas di luar ruangan!"
	case aqi <= 100:
		return "Sedang", "#F39C12", "Udara cukup baik. Kelompok sensitif sebaiknya kurangi aktivitas outdoor berat."
	case aqi <= 150:
		return "Tidak Sehat (Kelompok Sensitif)", "#E67E22", "Kelompok sensitif (anak/lansia/asma) disarankan memakai masker."
	case aqi <= 200:
		return "Tidak Sehat", "#E74C3C", "Kualitas udara buruk! Disarankan memakai masker jika keluar ruangan."
	case aqi <= 300:
		return "Sangat Tidak Sehat", "#8E44AD", "Peringatan kesehatan! Hindari semua aktivitas di luar ruangan."
	default:
		return "Berbahaya", "#7E5109", "Kondisi darurat polusi! Tetap di dalam ruangan."
	}
}