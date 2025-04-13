package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"app/cmd/videostream/config"
	"app/cmd/videostream/utils"
)

// StreamHandler creates and serves the stream
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	utils.EnableCORS(w)

	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Use hardcoded path or get from request
	videoPath := "Devil May Cry (2025) - S01E01 - Hell WEBDL-1080p.mkv"
	//videoPath := "Alien Romulus (2024) WEBDL-2160p_HB.mkv"
	if pathParam := r.URL.Query().Get("path"); pathParam != "" {
		videoPath = pathParam
	}
	
	sourcePath := filepath.Join(config.VideoDir, videoPath)

	// Check if file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		http.Error(w, "Video not found: "+sourcePath, http.StatusNotFound)
		return
	}
	
	// Get stream parameters
	width := r.URL.Query().Get("width")
	if width == "" {
		width = "1280" // 720p default
	}
	
	bitrate := r.URL.Query().Get("bitrate")
	if bitrate == "" {
		bitrate = "3000k" // Default bitrate
	}
	
	// Get audio stream index
	audioStream := 2 // Default audio stream
	if audioParam := r.URL.Query().Get("audio"); audioParam != "" {
		if as, err := strconv.Atoi(audioParam); err == nil {
			audioStream = as
		}
	}
	
	// Clean output directory
	utils.CleanOutputDir()
	
	// Get video info for duration
	fileInfo, err := utils.GetFileInfo(sourcePath)
	if err != nil {
		http.Error(w, "Error analyzing video: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Parse duration for playlist creation
	duration, err := strconv.ParseFloat(fileInfo.Duration, 64)
	if err != nil {
		log.Printf("Error parsing duration: %v, using default", err)
		duration = 3600 // Default to 1 hour if parsing fails
	}
	
	// Generate playlists upfront
	utils.CreateMasterPlaylist(width)
	utils.CreateIndexPlaylist(duration)
	
	// Choose streaming method based on request
	forceTranscode := r.URL.Query().Get("transcode") == "true"
	
	// Determine if transcoding is needed
	needsVideoTranscode := forceTranscode || (fileInfo.VideoCodec != "h264" && fileInfo.VideoCodec != "hevc")
	needsAudioTranscode := forceTranscode || (fileInfo.AudioCodec != "aac" && fileInfo.AudioCodec != "mp3")
	
	// Generate segments
	if needsVideoTranscode && needsAudioTranscode {
		// Full transcode
		go utils.TranscodeFullVideo(sourcePath, width, bitrate, audioStream)
	} else if needsVideoTranscode {
		// Video only transcode
		go utils.TranscodeVideoOnly(sourcePath, width, bitrate, audioStream)
	} else if needsAudioTranscode {
		// Audio only transcode
		go utils.TranscodeAudioOnly(sourcePath, audioStream)
	} else {
		// Direct play with segmentation
		go utils.SegmentVideo(sourcePath)
	}
	
	// Redirect to master playlist
	http.Redirect(w, r, "/segments/master.m3u8", http.StatusFound)
}

// SegmentsHandler handles requests for segments and playlist files
func SegmentsHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	utils.EnableCORS(w)

	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get the requested file
	parts := strings.SplitN(r.URL.Path, "/segments/", 2)
	if len(parts) != 2 || parts[1] == "" {
		http.Error(w, "Invalid segments request", http.StatusBadRequest)
		return
	}
	
	filename := parts[1]
	filePath := filepath.Join(config.OutputDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found: "+filePath, http.StatusNotFound)
		return
	}
	
	// Fix paths in playlists if needed
	if filename == "debug.m3u8" {
		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading playlist", http.StatusInternalServerError)
			return
		}
		
		// Fix the paths in the playlist
		fixedContent := strings.ReplaceAll(string(content), 
			"segment_", "/segments/segment_")
		
		// Set headers
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache") // Don't cache playlists
		
		// Write fixed content
		w.Write([]byte(fixedContent))
		return
	}
	
	// Set appropriate content type
	if strings.HasSuffix(filename, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache") // Don't cache playlists
	} else if strings.HasSuffix(filename, ".ts") {
		w.Header().Set("Content-Type", "video/MP2T")
		w.Header().Set("Cache-Control", "max-age=86400") // Cache segments
	}
	
	http.ServeFile(w, r, filePath)
} 