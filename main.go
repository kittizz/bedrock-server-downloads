package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/hashicorp/go-version"
)

// Constants
const (
	MinecraftBedrockAPIURL = "https://net-secondary.web.minecraft-services.net/api/v1.0/download/links"
	DataFilePath           = "bedrock-server-downloads.json"
	UserAgent              = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
)

// Data structures
type ServerData struct {
	Release map[string]VersionData `json:"release"`
	Preview map[string]VersionData `json:"preview"`
}

type VersionData struct {
	Windows PlatformData `json:"windows"`
	Linux   PlatformData `json:"linux"`
}

type PlatformData struct {
	URL string `json:"url"`
}

func main() {
	log.Println("Processing Minecraft Bedrock Server versions...")
	if err := processVersions(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("All versions processed successfully.")
}

func processVersions() error {
	// Fetch downloads from new API
	client := &http.Client{}
	req, err := http.NewRequest("GET", MinecraftBedrockAPIURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type apiResponse struct {
		Result struct {
			Links []struct {
				DownloadType string `json:"downloadType"`
				DownloadURL  string `json:"downloadUrl"`
			} `json:"links"`
		} `json:"result"`
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return err
	}

	urls := make(map[string]string)
	for _, link := range apiResp.Result.Links {
		urls[link.DownloadType] = link.DownloadURL
	}

	platforms := []string{"serverBedrockWindows", "serverBedrockLinux",
		"serverBedrockPreviewWindows", "serverBedrockPreviewLinux"}
	for _, p := range platforms {
		if urls[p] == "" {
			return fmt.Errorf("missing download URL for %s", p)
		}
	}

	// Load existing data
	data := ServerData{
		Release: make(map[string]VersionData),
		Preview: make(map[string]VersionData),
	}
	if fileData, err := os.ReadFile(DataFilePath); err == nil {
		json.Unmarshal(fileData, &data) // Ignore error, use empty if can't parse
	}

	// Process regular version
	if err := processVersion(&data, false,
		urls["serverBedrockWindows"], urls["serverBedrockLinux"]); err != nil {
		return err
	}

	// Process preview version
	if err := processVersion(&data, true,
		urls["serverBedrockPreviewWindows"], urls["serverBedrockPreviewLinux"]); err != nil {
		return err
	}

	// Save data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(DataFilePath, jsonData, 0644)
}

func processVersion(data *ServerData, isPreview bool, windowsURL, linuxURL string) error {
	// Extract versions
	windowsVersion, err := extractVersion(windowsURL)
	if err != nil {
		return err
	}
	linuxVersion, err := extractVersion(linuxURL)
	if err != nil {
		return err
	}

	// Check version consistency
	if windowsVersion != linuxVersion {
		versionType := "regular"
		if isPreview {
			versionType = "preview"
		}
		return fmt.Errorf("version mismatch in %s release: Windows=%s, Linux=%s",
			versionType, windowsVersion, linuxVersion)
	}

	// Normalize version
	v, err := version.NewVersion(windowsVersion)
	if err != nil {
		return err
	}
	segments := v.Segments()
	normalizedVersion := ""
	if len(segments) >= 2 {
		if len(segments) >= 3 {
			normalizedVersion = fmt.Sprintf("%d.%d.%d", segments[0], segments[1], segments[2])
		} else {
			normalizedVersion = fmt.Sprintf("%d.%d", segments[0], segments[1])
		}
	} else {
		return fmt.Errorf("invalid version format: %s", windowsVersion)
	}

	// Update data structure
	versionType := "Regular"
	versionMap := data.Release
	if isPreview {
		versionType = "Preview"
		versionMap = data.Preview
	}

	if _, exists := versionMap[normalizedVersion]; exists {
		log.Printf("%s version v%s is already up to date.", versionType, normalizedVersion)
	} else {
		log.Printf("Adding new %s version v%s to database.", versionType, normalizedVersion)
		versionMap[normalizedVersion] = VersionData{
			Windows: PlatformData{URL: windowsURL},
			Linux:   PlatformData{URL: linuxURL},
		}
		if isPreview {
			data.Preview = versionMap
		} else {
			data.Release = versionMap
		}
	}

	return nil
}

func extractVersion(url string) (string, error) {
	pattern := regexp.MustCompile(`bedrock-server-(\d+\.\d+(\.\d+){1,2})\.zip`)
	matches := pattern.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("unable to extract version from URL: %s", url)
	}
	return matches[1], nil
}
