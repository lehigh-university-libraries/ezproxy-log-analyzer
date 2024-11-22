package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EzpaarseLog struct {
	Date               string `json:"date"`
	Datetime           string `json:"datetime"`
	Domain             string `json:"domain"`
	City               string `json:"geoip-city"`
	Coordinates        string `json:"geoip-coordinates"`
	Country            string `json:"geoip-country"`
	Latitude           string `json:"geoip-latitude"`
	Longitude          string `json:"geoip-longitude"`
	Region             string `json:"geoip-region"`
	Host               string `json:"host"`
	Identd             string `json:"identd"`
	LogID              string `json:"log_id"`
	Login              string `json:"login"`
	Middlewares        string `json:"middlewares"`
	MiddlewaresDate    string `json:"middlewares_date"`
	MiddlewaresVersion string `json:"middlewares_version"`
	OnCampus           string `json:"on_campus"`
	Platform           string `json:"platform"`
	PlatformName       string `json:"platform_name"`
	PlatformsDate      string `json:"platforms_date"`
	PlatformsVersion   string `json:"platforms_version"`
	PublisherName      string `json:"publisher_name"`
	Robot              string `json:"robot"`
	Size               string `json:"size"`
	Status             string `json:"status"`
	Timestamp          int    `json:"timestamp"`
	URL                string `json:"url"`
}

type Quicksight struct {
	Date          string  `json:"date"`
	Domain        string  `json:"domain"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lng"`
	Region        string  `json:"region"`
	Platform      string  `json:"platform"`
	PlatformName  string  `json:"platform_name"`
	PublisherName string  `json:"publisher_name"`
	Department    string  `json:"department"`
	Role          string  `json:"role"`
}

// Function to check if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Function to save the response body to a file
func saveToFile(filename string, data io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	return err
}

// Function to extract description
func extractDescription(filename string) (string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", "", err
	}
	// Use regex to capture everything in the second <td> tag
	re := regexp.MustCompile(`<tr><td>Description</td><td>([^<]+)</td></tr>`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) < 2 {
		return "", "", fmt.Errorf("description not found")
	}
	// Now, we have everything inside the second <td>
	fullDescription := matches[1]

	// Try to split it on " - " to separate category and description
	parts := strings.SplitN(fullDescription, " - ", 2)
	category := parts[0]
	description := "unknown"
	if len(parts) == 2 {
		description = parts[1]
	}

	return category, description, nil
}

func main() {
	ldapSearchServer := os.Getenv("LDAP_SEARCH_SERVER")

	file, err := os.Open("logs/ezpaarse.json")
	if err != nil {
		slog.Error("Failed to open ezpaarse JSON. Did you run process.sh?", "err", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode the JSON file into a slice of EzpaarseLog
	var logs []EzpaarseLog
	if err := json.NewDecoder(file).Decode(&logs); err != nil {
		slog.Error("Failed to decode JSON", "err", err)
		os.Exit(1)
	}

	// Iterate over the logs and print some fields
	visitors := []Quicksight{}
	for _, log := range logs {
		visitor := Quicksight{
			Date:          log.Date,
			Domain:        log.Domain,
			City:          log.City,
			Country:       log.Country,
			Region:        log.Region,
			Platform:      log.Platform,
			PlatformName:  log.PlatformName,
			PublisherName: log.PublisherName,
		}
		lat, err := strconv.ParseFloat(log.Latitude, 64)
		if err != nil {
			slog.Warn("Error converting string to float", "err", err, "lng", log.Latitude)
			continue
		}
		lng, err := strconv.ParseFloat(log.Longitude, 64)
		if err != nil {
			slog.Warn("Error converting string to float", "err", err, "lat", log.Latitude)
			continue
		}

		visitor.Latitude = lat
		visitor.Longitude = lng

		email := fmt.Sprintf("%s@lehigh.edu", strings.ToLower(log.Login))
		filename := fmt.Sprintf("/tmp/%s.html", log.Login)
		// Check if the file already exists
		if !fileExists(filename) {
			slog.Info("Querying LDAP")
			// If not, make the HTTP POST request
			resp, err := http.PostForm(ldapSearchServer, map[string][]string{
				"cn": {email},
			})
			if err != nil {
				slog.Error("Error making request", "err", err)
				return
			}
			defer resp.Body.Close()

			// Save response to STR.html
			err = saveToFile(filename, resp.Body)
			if err != nil {
				slog.Error("Error saving file", "err", err)
				return
			}
			time.Sleep(2 * time.Second)
		}

		// Extract and parse the description
		role, department, err := extractDescription(filename)
		if err != nil {
			slog.Warn("Error extracting description", "login", log.Login, "err", err)
		}
		// Output the extracted values
		visitor.Department = department
		visitor.Role = role
		visitors = append(visitors, visitor)
	}

	f, err := os.Create("ezproxy.json")
	if err != nil {
		slog.Error("failed to create ezproxy JSON", "err", err)
	}
	defer f.Close()

	// Marshal the logs slice into JSON
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(visitors); err != nil {
		slog.Error("failed to encode JSON", "err", err)
	}
}
