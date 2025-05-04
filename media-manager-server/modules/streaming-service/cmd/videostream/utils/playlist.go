package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"streaming-service/cmd/videostream/config"
)

// CleanOutputDir removes all files from the output directory
func CleanOutputDir() {
	// Remove old files
	files, err := os.ReadDir(config.OutputDir)
	if err != nil {
		log.Printf("Error reading output directory: %v", err)
		return
	}
	
	for _, file := range files {
		os.Remove(filepath.Join(config.OutputDir, file.Name()))
	}
	
	log.Printf("Cleaned output directory")
}

// CreateMasterPlaylist creates the master HLS playlist
func CreateMasterPlaylist(width string) {
	content := fmt.Sprintf(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=3000000,RESOLUTION=%sx720
/segments/index.m3u8
`, width)

	// Write playlist file
	masterPath := filepath.Join(config.OutputDir, "master.m3u8")
	err := os.WriteFile(masterPath, []byte(content), 0644)
	if err != nil {
		log.Printf("Error writing master playlist: %v", err)
	}
}

// CreateIndexPlaylist creates the index HLS playlist with all segments based on duration
func CreateIndexPlaylist(duration float64) {
	// Calculate number of segments (round up)
	numSegments := int((duration + config.SegmentLength - 1) / float64(config.SegmentLength))
	if numSegments <= 0 {
		numSegments = 100 // Default to 100 segments if we couldn't determine duration
	}
	
	log.Printf("Creating playlist for video with duration %.2f seconds, %d segments", duration, numSegments)
	
	// Start building playlist content
	var playlist strings.Builder
	playlist.WriteString("#EXTM3U\n")
	playlist.WriteString("#EXT-X-VERSION:3\n")
	playlist.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", config.SegmentLength))
	playlist.WriteString("#EXT-X-ALLOW-CACHE:YES\n")
	playlist.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")
	playlist.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")
	
	// Add segment entries with consistent durations
	for i := 0; i < numSegments; i++ {
		segmentDuration := float64(config.SegmentLength)
		
		// Last segment might be shorter
		if i == numSegments-1 && duration > 0 {
			remainingDuration := duration - float64(i*config.SegmentLength)
			if remainingDuration < segmentDuration {
				segmentDuration = remainingDuration
				if segmentDuration <= 0 {
					segmentDuration = 1 // Ensure positive duration
				}
			}
		}


		if i != 0 {
			playlist.WriteString("#EXT-X-DISCONTINUITY\n") // Add this before each segment (except the first)
		}
		
		segmentName := fmt.Sprintf("segment_%05d.ts", i)
		playlist.WriteString(fmt.Sprintf("#EXTINF:%.6f,\n", segmentDuration))
		playlist.WriteString(fmt.Sprintf("/segments/%s\n", segmentName))
	}
	
	// End playlist
	playlist.WriteString("#EXT-X-ENDLIST\n")
	
	// Write the playlist file
	indexPath := filepath.Join(config.OutputDir, "index.m3u8")
	err := os.WriteFile(indexPath, []byte(playlist.String()), 0644)
	if err != nil {
		log.Printf("Error writing index playlist: %v", err)
	} else {
		log.Printf("Created index playlist with %d segments", numSegments)
	}
} 