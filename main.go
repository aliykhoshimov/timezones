package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// Timezone represents a structured timezone entry
type Timezone struct {
	Label      string `json:"label"`
	Identifier string `json:"identifier"`
	UTCOffset  string `json:"utc_offset"`
}

// Predefined list of timezones
var timezones = []Timezone{
	{"UTC (Coordinated Universal Time)", "UTC", "UTC + 0:00"},
	{"GMT (Greenwich Mean Time)", "Europe/London", "UTC + 0:00"},
	{"CET (Central European Time)", "Europe/Berlin", "UTC + 1:00"},
	{"EET (Eastern European Time)", "Europe/Athens", "UTC + 2:00"},
	{"MSK (Moscow Standard Time)", "Europe/Moscow", "UTC + 3:00"},
	{"GST (Gulf Standard Time)", "Asia/Dubai", "UTC + 4:00"},
	{"IST (Indian Standard Time)", "Asia/Kolkata", "UTC + 5:30"},
	{"BST (Bangladesh Standard Time)", "Asia/Dhaka", "UTC + 6:00"},
	{"ICT (Indochina Time)", "Asia/Bangkok", "UTC + 7:00"},
	{"CST (China Standard Time)", "Asia/Shanghai", "UTC + 8:00"},
	{"AWST (Australian Western Standard Time)", "Australia/Perth", "UTC + 8:00"},
	{"JST (Japan Standard Time)", "Asia/Tokyo", "UTC + 9:00"},
	{"KST (Korea Standard Time)", "Asia/Seoul", "UTC + 9:00"},
	{"ACST (Australia Central Standard Time)", "Australia/Adelaide", "UTC + 9:30"},
	{"AEST (Australia Eastern Standard Time)", "Australia/Sydney", "UTC + 10:00"},
	{"ChST (Chamorro Standard Time)", "Pacific/Guam", "UTC + 10:00"},
	{"NZST (New Zealand Standard Time)", "Pacific/Auckland", "UTC + 12:00"},
	{"SST (Samoa Standard Time)", "Pacific/Pago_Pago", "UTC-11:00"},
	{"HST (Hawaii Standard Time)", "Pacific/Honolulu", "UTC-10:00"},
	{"AKST (Alaska Standard Time)", "America/Anchorage", "UTC - 9:00"},
	{"PST (Pacific Standard Time)", "America/Los_Angeles", "UTC - 8:00"},
	{"MST (Mountain Standard Time)", "America/Denver", "UTC - 7:00"},
	{"CST (Central Standard Time)", "America/Chicago", "UTC - 6:00"},
	{"EST (Eastern Standard Time)", "America/New_York", "UTC - 5:00"},
	{"AST (Atlantic Standard Time)", "America/Puerto_Rico", "UTC - 4:00"},
	{"NST (Newfoundland Standard Time)", "America/St_Johns", "UTC - 3:30"},
	{"BRT (Brasília Time)", "America/Sao_Paulo", "UTC - 3:00"},
	{"ART (Argentina Time)", "America/Argentina/Buenos_Aires", "UTC - 3:00"},
}

// UserTimezoneMap stores user ID and their selected timezone
var (
	userTimezoneMap = make(map[string]string)
	mu              sync.Mutex
)

// GetTimezonesHandler returns available timezones
func GetTimezonesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timezones)
}

// SaveTimezoneRequest represents the payload to save a timezone
type SaveTimezoneRequest struct {
	UserID   string `json:"user_id"`
	Timezone string `json:"timezone"`
}

// SaveTimezoneHandler validates and saves a user's timezone
func SaveTimezoneHandler(w http.ResponseWriter, r *http.Request) {
	var req SaveTimezoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate timezone
	valid := false
	var tzIdentifier string
	for _, tz := range timezones {
		if req.Timezone == tz.Label {
			valid = true
			tzIdentifier = tz.Identifier
			break
		}
	}

	if !valid {
		http.Error(w, "Invalid timezone selection", http.StatusBadRequest)
		return
	}

	// Save to user timezone map
	mu.Lock()
	userTimezoneMap[req.UserID] = tzIdentifier
	mu.Unlock()

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Timezone saved successfully",
		"timezone": tzIdentifier,
	})
}

// GetCurrentTimeHandler returns the current time for a user
func GetCurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	tzIdentifier, exists := userTimezoneMap[userID]
	mu.Unlock()

	if !exists {
		http.Error(w, "User timezone not found", http.StatusNotFound)
		return
	}

	// Get the current time in user's timezone
	loc, err := time.LoadLocation(tzIdentifier)
	if err != nil {
		http.Error(w, "Invalid timezone", http.StatusInternalServerError)
		return
	}
	timeInZone := time.Now().In(loc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id":      userID,
		"timezone":     tzIdentifier,
		"current_time": timeInZone.Format(time.RFC3339),
	})
}

func main() {
	http.HandleFunc("/timezones", GetTimezonesHandler)
	http.HandleFunc("/users/timezone", SaveTimezoneHandler)
	http.HandleFunc("/users/current_time", GetCurrentTimeHandler)

	log.Println("Server started on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
