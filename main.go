package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wizzymore/tinyfiledialogs"
	"github.com/barasher/go-exiftool"
)

type AudioMetadata struct {
	Title       string `json:"title,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Album       string `json:"album,omitempty"`
	Year        int64    `json:"year,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Duration		string `json:"duration,omitempty"`
	Channels		int64		 `json:"channels,omitempty"`
	SampleRate	int64    `json:"sample_rate,omitempty"`
	BitRate     int64    `json:"bit_rate,omitempty"`
	Track       int64    `json:"track,omitempty"`
	TotalTracks int64		 `json:"total_tracks,omitempty"`
	FileType    string `json:"file_type,omitempty"`
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

// Prompts user to select a batch of audio files, uses exiftool to pull metadata from those files, then maps that data and ouputs a JSON file
//
// return is void
func main() {
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

	err = os.WriteFile("metadata.json", jsonData, 0644)
	if err != nil {
		log.Println("Error writing file: ", err)
		return
	}

	fmt.Println("Metadata exported successfully to metadata.json.")
}
