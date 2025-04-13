package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	port          = 8080
	videoDir      = "/data/animes/Devil May Cry 2025/Season 1"
	segmentLength = 6                // Segment length in seconds
	segmentDir    = "./segments"     // Directory to store segments
	playlistDir   = "./playlists"    // Directory to store playlists
	maxSegmentAge = 30 * time.Minute // How long to keep segments before cleanup
)

type VideoSegment struct {
	Path      string
	CreatedAt time.Time
}

// Global segment cache
var (
	segmentCache     = make(map[string]map[int]*VideoSegment) // videoId -> segmentIndex -> segment
	segmentCacheLock sync.RWMutex
)

func main() {
	// Create cache directories
	os.MkdirAll(segmentDir, 0755)
	os.MkdirAll(playlistDir, 0755)

	// Start cleanup routine for old segments
	go segmentCleanupRoutine()

	// Create HTTP handlers for streaming
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/hls/", hlsHandler)

	// Start the server
	fmt.Printf("Streaming server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Main handler that creates and serves the master playlist
func streamHandler(w http.ResponseWriter, r *http.Request) {
	// Extract video path from URL
	// pathParts := strings.SplitN(r.URL.Path, "/stream/", 2)
	// if len(pathParts) < 2 || pathParts[1] == "" {
	// 	http.Error(w, "Invalid video path", http.StatusBadRequest)
	// 	return
	// }

	// videoPath := pathParts[1]
	// sourcePath := filepath.Join(videoDir, videoPath)

	// Add CORS headers
	enableCORS(w)

	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	videoPath := "/Devil May Cry (2025) - S01E01 - Hell WEBDL-1080p.mkv"
	sourcePath := filepath.Join(videoDir, videoPath)

	// Security check
	if !strings.HasPrefix(sourcePath, videoDir) {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	// Check if file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Get quality parameters
	width := r.URL.Query().Get("width")
	if width == "" {
		width = "1920" // 720p default
	}

	bitrate := r.URL.Query().Get("bitrate")
	if bitrate == "" {
		bitrate = "5000k" // Default bitrate
	}

	const videoUuid = "video_uuid" // Placeholder for video UUID

	// Generate a unique ID for this video+settings combination
	videoId := fmt.Sprintf("%s_%s_%s",
		// strings.ReplaceAll(filepath.Base(videoPath), ".", "|"),
		videoUuid,
		width,
		strings.ReplaceAll(bitrate, "k", ""))

	// Create master playlist
	masterPlaylist := createMasterPlaylist(videoId)

	// Write the master playlist to a temporary file
	playlistPath := filepath.Join(playlistDir, videoId+".m3u8")
	err := ioutil.WriteFile(playlistPath, []byte(masterPlaylist), 0644)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	log.Printf("Created master playlist for %s at %s", videoPath, playlistPath)

	// Get video duration to set playlist parameters
	duration, err := getVideoDuration(sourcePath)
	if err != nil {
		log.Printf("Warning: couldn't get video duration: %v", err)
	}

	// Create variant playlist
	variantPlaylist := createVariantPlaylist(videoId, int(duration), segmentLength)
	variantPlaylistPath := filepath.Join(playlistDir, videoId+"_variant.m3u8")
	err = ioutil.WriteFile(variantPlaylistPath, []byte(variantPlaylist), 0644)
	if err != nil {
		http.Error(w, "Failed to create variant playlist", http.StatusInternalServerError)
		return
	}

	// Store video information in segment cache
	segmentCacheLock.Lock()
	segmentCache[videoId] = make(map[int]*VideoSegment)
	segmentCacheLock.Unlock()

	// Redirect to HLS
	http.Redirect(w, r, "/hls/"+videoId+".m3u8", http.StatusFound)
}

// Handler for HLS requests (playlists and segments)
func hlsHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	enableCORS(w)

	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	pathParts := strings.SplitN(r.URL.Path, "/hls/", 2)
	if len(pathParts) < 2 || pathParts[1] == "" {
		http.Error(w, "Invalid HLS path", http.StatusBadRequest)
		return
	}

	requestPath := pathParts[1]

	// Handle master playlist request
	if strings.HasSuffix(requestPath, ".m3u8") {
		servePlaylist(w, r, requestPath)
		return
	}

	// Handle segment request
	if strings.HasSuffix(requestPath, ".ts") {
		serveSegment(w, r, requestPath)
		return
	}

	http.Error(w, "Invalid HLS request", http.StatusBadRequest)
}

// Serve playlist files (master or variant)
func servePlaylist(w http.ResponseWriter, r *http.Request, playlistPath string) {
	fullPath := filepath.Join(playlistDir, playlistPath)

	// Check if playlist exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")

	http.ServeFile(w, r, fullPath)
}

// Serve video segments, generating them on demand if needed
func serveSegment(w http.ResponseWriter, r *http.Request, segmentPath string) {
	// Parse segment information from filename
	// Expected format: videoId_segmentNumber.ts
	parts := strings.Split(strings.TrimSuffix(segmentPath, ".ts"), "_")
	if len(parts) < 3 {
		http.Error(w, "Invalid segment path", http.StatusBadRequest)
		return
	}

	// Last part is segment number
	segmentNum, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "Invalid segment number", http.StatusBadRequest)
		return
	}

	// Reconstruct video ID
	videoId := strings.Join(parts[:len(parts)-1], "_")

	// Check if segment exists in cache
	segmentCacheLock.RLock()
	segmentMap, exists := segmentCache[videoId]
	var segment *VideoSegment
	if exists {
		segment = segmentMap[segmentNum]
	}
	segmentCacheLock.RUnlock()

	if segment != nil && fileExists(segment.Path) {
		// Update access time
		segmentCacheLock.Lock()
		segmentCache[videoId][segmentNum].CreatedAt = time.Now()
		segmentCacheLock.Unlock()

		// Serve existing segment
		w.Header().Set("Content-Type", "video/MP2T")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, segment.Path)
		return
	}

	// Need to generate segment
	segmentFilePath := filepath.Join(segmentDir, fmt.Sprintf("%s_%d.ts", videoId, segmentNum))

	// Parse video ID to get original parameters
	idParts := strings.Split(videoId, "_")
	if len(idParts) < 3 {
		http.Error(w, "Invalid video ID format", http.StatusBadRequest)
		return
	}

	// Extract parameters from the ID
	var (
		// filename = strings.Join(idParts[:len(idParts)-2], "_") // Reconstruct original filename
		width   = idParts[len(idParts)-2]
		bitrate = idParts[len(idParts)-1] + "k"
	)
	filename := "Devil May Cry (2025) - S01E01 - Hell WEBDL-1080p.mkv"
	// Replace underscores in the filename with periods
	// filename = strings.ReplaceAll(filename, "|", ".")

	// Full path to the source video
	sourcePath := filepath.Join(videoDir, filename)

	// Calculate start and end times for the segment
	startTime := segmentNum * segmentLength
	// endTime := startTime + segmentLength

	// Generate the segment using FFmpeg
	err = generateSegment(sourcePath, segmentFilePath, startTime, segmentLength, width, bitrate, 2)
	if err != nil {
		log.Printf("Failed to generate segment: %v", err)
		http.Error(w, "Failed to generate segment", http.StatusInternalServerError)
		return
	}

	// Add to cache
	segmentCacheLock.Lock()
	if segmentCache[videoId] == nil {
		segmentCache[videoId] = make(map[int]*VideoSegment)
	}
	segmentCache[videoId][segmentNum] = &VideoSegment{
		Path:      segmentFilePath,
		CreatedAt: time.Now(),
	}
	segmentCacheLock.Unlock()

	// Serve the segment
	w.Header().Set("Content-Type", "video/MP2T")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, segmentFilePath)
}

// Create a master playlist file
func createMasterPlaylist(videoId string) string {
	return fmt.Sprintf(
		`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=3000000,RESOLUTION=1280x720
/hls/%s_variant.m3u8
`, videoId)
}

// Create a variant playlist file with segments
func createVariantPlaylist(videoId string, durationSeconds, segmentLength int) string {
	var buffer bytes.Buffer

	buffer.WriteString("#EXTM3U\n")
	buffer.WriteString("#EXT-X-VERSION:3\n")
	buffer.WriteString("#EXT-X-ALLOW-CACHE:NO\n")
	buffer.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", segmentLength))
	buffer.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")
	buffer.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")

	// Calculate number of segments
	numSegments := (durationSeconds + segmentLength - 1) / segmentLength

	// Cap at reasonable max if duration is unknown
	if durationSeconds <= 0 || numSegments > 3000 {
		numSegments = 3000
	}

	// Add segment entries
	for i := 0; i < numSegments; i++ {
		// Last segment might be shorter
		var segmentDuration = segmentLength
		if i == numSegments-1 && durationSeconds > 0 {
			lastSegmentDuration := durationSeconds % segmentLength
			if lastSegmentDuration > 0 {
				segmentDuration = lastSegmentDuration
			}
		}
		if i != 0 {
			buffer.WriteString("#EXT-X-DISCONTINUITY\n") // Add this before each segment (except the first)
		}

		buffer.WriteString(fmt.Sprintf("#EXTINF:%d.000000,\n", segmentDuration))
		buffer.WriteString(fmt.Sprintf("/hls/%s_%d.ts\n", videoId, i))
	}

	buffer.WriteString("#EXT-X-ENDLIST\n")

	return buffer.String()
}

// // Generate a video segment using ffmpeg
// func generateSegment(sourcePath, outputPath string, startTime, duration int, width, bitrate string) error {
// 	log.Printf("Generating segment: %s (start: %d, duration: %d)",
// 		filepath.Base(outputPath), startTime, duration)

// 	// // Calculate PTS offset in milliseconds
// 	// ptsOffset := startTime * 1000

// 	err := ffmpeg.Input(sourcePath, ffmpeg.KwArgs{
// 		"ss":            startTime, // Seek to position
// 		"accurate_seek": "",        // More accurate seeking
// 	}).
// 		Output(outputPath, ffmpeg.KwArgs{
// 			"t": duration,
// 			// "c:v": "libx264",
// 			"c:v": "copy",
// 			"c:a": "aac",
// 			// "c:a": "copy",
// 			// "b:v": bitrate,
// 			"b:a": "128k",
// 			"ac":  "2",
// 			"map": ["0:v:0", "0:a:2"],
// 			// "vf":                fmt.Sprintf("scale=%s:-2", width),
// 			"preset":           "ultrafast",
// 			"tune":             "zerolatency",
// 			"f":                "mpegts",
// 			"force_key_frames": "expr:gte(t,0)", // Force keyframe at start of segment
// 			// "x264-params":       "scenecut=0:open_gop=0", // Prevent scene detection and open GOPs
// 			"avoid_negative_ts": "make_zero",
// 			"muxdelay":          "0",
// 			"mpegts_flags":      "resend_headers", // Include necessary headers in each segment
// 			"async":             "1",              // Add audio sync parameter
// 			// "vsync":             "1",              // Add video sync parameter
// 			// "map":               "0:v:0,0:a:0,0:s:0",
// 		}).
// 		ErrorToStdOut().
// 		Run()

// 	return err
// }

// Generate a video segment using ffmpeg with custom audio stream selection
func generateSegment(sourcePath, outputPath string, startTime, duration int, width, bitrate string, audioStreamIndex int) error {
	log.Printf("Generating segment: %s (start: %d, duration: %d, audio stream: %d)",
		filepath.Base(outputPath), startTime, duration, audioStreamIndex)

	// Create input with proper seeking
	input := ffmpeg.Input(sourcePath, ffmpeg.KwArgs{
		"ss":            startTime,
		"accurate_seek": "",
	})

	// Use AddOutputOption for -map options since they need to be separate
	cmd := input.Output(outputPath, ffmpeg.KwArgs{
		"t":   duration,
		"c:v": "copy",   // Copy video stream
		"c:a": "aac",    // Transcode audio to AAC
		"c:s": "webvtt", // Copy subtitle stream
		"b:a": "128k",   // Audio bitrate
		"ac":  "2",      // 2 audio channels
		// "ar":                "44100", // Standard audio sample rate
		"preset":            "ultrafast",
		"tune":              "zerolatency",
		"f":                 "mpegts",
		"force_key_frames":  "expr:gte(t,0)",
		"avoid_negative_ts": "make_zero",
		"muxdelay":          "0",
		"mpegts_flags":      "resend_headers",
		// "async":             "1",
		"map": []string{"0:v:0", "0:a:2", "0:s:1"},
	})

	// Run the command
	err := cmd.ErrorToStdOut().Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v", err)
	}

	return nil
}

// Get video duration using ffprobe
func getVideoDuration(filePath string) (float64, error) {
	data, err := ffmpeg.Probe(filePath)
	if err != nil {
		return 0, err
	}

	// Parse the JSON output to get duration
	type ProbeFormat struct {
		Duration string `json:"duration"`
	}

	type ProbeData struct {
		Format ProbeFormat `json:"format"`
	}

	var probeData ProbeData
	err = json.Unmarshal([]byte(data), &probeData)
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(probeData.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// Cleanup old segments to prevent filling disk
func segmentCleanupRoutine() {
	for {
		time.Sleep(5 * time.Minute)

		now := time.Now()
		var segmentsToRemove []string

		segmentCacheLock.Lock()

		// Find old segments
		for videoId, segments := range segmentCache {
			for segmentNum, segment := range segments {
				if now.Sub(segment.CreatedAt) > maxSegmentAge {
					// Mark for removal
					segmentsToRemove = append(segmentsToRemove, segment.Path)
					delete(segments, segmentNum)
				}
			}

			// If all segments are removed, remove the video entry
			if len(segments) == 0 {
				delete(segmentCache, videoId)
			}
		}

		segmentCacheLock.Unlock()

		// Delete the files
		for _, path := range segmentsToRemove {
			os.Remove(path)
		}

		if len(segmentsToRemove) > 0 {
			log.Printf("Cleaned up %d old segments", len(segmentsToRemove))
		}
	}
}

// Check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Range")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
}