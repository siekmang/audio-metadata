package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/wizzymore/tinyfiledialogs"
)

type AudioMetadata struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Year        int64  `json:"year"`
	Genre       string `json:"genre"`
	Duration		string `json:"duration"`
	Channels		int64	 `json:"channels"`
	SampleRate	int64  `json:"sample_rate"`
	BitRate     int64  `json:"bit_rate"`
	Track       int64  `json:"track"`
	TotalTracks int64	 `json:"total_tracks"`
	FileType    string `json:"file_type"`
}

// Helper function calling exiftool's GetString function
//
// file - exiftool.FileMetadata - The metadata pulled from a file
//
// key - string - the key for the metadata string we want to pull
//
// returns a string that represents a certain piece of metadata about an audio file
func getString(file exiftool.FileMetadata, key string) string {
    v, _ := file.GetString(key)
    return v
}

// Helper function calling exiftool's GetInt function
//
// file - exiftool.FileMetadata - The metadata pulled from a file
//
// key - string - the key for the metadata string we want to pull
//
// returns an int64 that represents a certain piece of metadata about an audio file
func getInt(file exiftool.FileMetadata, key string) int64 {
    v, _ := file.GetInt(key)
    return v
}

// Prompts user to select a batch of audio files, uses exiftool to pull
// metadata from those files, then maps that data and ouputs a JSON file
//
// If this is the first time the user is using this, it will prompt them
// to select an output location, and save that to a config file in the
// config directory
//
// return is void
func main() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("No config path found.")
		return
	}

	configFile := filepath.Join(configPath, "metadata_output_dir")

	dirBytes, err := os.ReadFile(configFile)
	fileDir := string(dirBytes)

	if err != nil || fileDir == "" {
		selected, ok := tinyfiledialogs.SelectFolderDialog(
			"Select Output Location",
			"",
		)

		if !ok {
			fmt.Println("Selection canceled.")
			return
		}

		os.WriteFile(configFile, []byte(selected), 0644)
		fileDir = selected
	}

	timestamp := time.Now().Format("2006-01-02-150405")
	fileName := fmt.Sprintf("metadata_%s.json", timestamp)
	filePath := filepath.Join(fileDir, fileName)

	var metadata []AudioMetadata

	et, err := exiftool.NewExiftool()
	if err != nil {
		log.Fatal(err)
	}

	defer et.Close()
	paths, ok := tinyfiledialogs.OpenFileDialog(
		"Select Audio Files",
		"",
		[]string{"*.mp3", "*.flac", "*.wav", "*.m4a"},
		"Audio Files",
		true,
	)

	if !ok {
		fmt.Println("Selection canceled.")
		return
	}

	fileList := strings.Split(paths, "|")
	data := et.ExtractMetadata(fileList...)

	fmt. Println("Selected Files:")
	for _, item := range data {

		if item.Err != nil{
			log.Printf("Skipping %s due to error: %v\n", item.File, item.Err)
    	continue
    }

		meta := AudioMetadata{
		Title:       getString(item, "Title"),
		Artist:      getString(item, "Artist"),
		Album:       getString(item, "Album"),
		Year:        getInt(item, "Year"),
		Duration:    getString(item, "TrackDuration"),
		Channels:    getInt(item, "AudioChannels"),
		SampleRate:  getInt(item, "AudioSampleRate"),
		BitRate:     getInt(item, "AudioCurrentTargetBitRate"),
		Genre:       getString(item, "Genre"),
		Track:       getInt(item, "CDDBTrackNumber"),
		TotalTracks: getInt(item, "CDDBDiscTracks"),
		FileType:    getString(item, "FileType"),
		}

		fmt.Println("-", meta.Title, "by", meta.Artist)

		metadata = append(metadata, meta)
	}

	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Println("Error writing file: ", err)
		return
	}

	fmt.Println("Metadata exported successfully to", filePath)
}
