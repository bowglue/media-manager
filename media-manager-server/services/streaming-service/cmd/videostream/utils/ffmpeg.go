package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"streaming-service/cmd/videostream/config"
	"streaming-service/cmd/videostream/models"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// RunFFmpegCommand runs an FFmpeg command with the given arguments
func RunFFmpegCommand(args []string) {
	cmd := exec.Command("ffmpeg", args...)
	
	// Capture output
	var logBuffer bytes.Buffer
	cmd.Stderr = &logBuffer
	
	// Run command and wait for completion
	err := cmd.Run()
	if err != nil {
		log.Printf("FFmpeg error: %v\n%s", err, logBuffer.String())
		return
	}
	
	log.Printf("FFmpeg command completed successfully")
}

// GetFileInfo gets information about a media file using ffprobe
func GetFileInfo(filePath string) (*models.FileInfo, error) {
	data, err := ffmpeg.Probe(filePath)
	if err != nil {
		return nil, err
	}
	
	var probeData map[string]interface{}
	err = json.Unmarshal([]byte(data), &probeData)
	if err != nil {
		return nil, err
	}
	
	info := &models.FileInfo{}
	
	// Get duration
	format, ok := probeData["format"].(map[string]interface{})
	if ok {
		if duration, ok := format["duration"].(string); ok {
			info.Duration = duration
		}
	}
	
	// Get stream info
	streams, ok := probeData["streams"].([]interface{})
	if ok {
		for _, stream := range streams {
			streamMap, ok := stream.(map[string]interface{})
			if !ok {
				continue
			}
			
			codecType, ok := streamMap["codec_type"].(string)
			if !ok {
				continue
			}
			
			if codecType == "video" && info.VideoCodec == "" {
				if codec, ok := streamMap["codec_name"].(string); ok {
					info.VideoCodec = codec
				}
			} else if codecType == "audio" && info.AudioCodec == "" {
				if codec, ok := streamMap["codec_name"].(string); ok {
					info.AudioCodec = codec
				}
			}
		}
	}
	
	return info, nil
}

// TranscodeFullVideo performs full transcoding for both video and audio
func TranscodeFullVideo(sourcePath, width, bitrate string, audioStream int) {
	log.Printf("Starting full transcoding of: %s", sourcePath)
	
	// Create the segment pattern
	// segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts")
	segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts.temp")
	segmentList := filepath.Join(config.OutputDir, "segment_list.csv")
	startIndex := 0 / config.SegmentLength
	
	// Build FFmpeg command with parameters to ensure exact segment times
	args := []string{
		"-i", sourcePath,
		"-c:v", "libx264",                           // H.264 video codec
		"-preset", "ultrafast",                      // Fastest encoding
		"-tune", "fastdecode",                       // Optimize for fast decoding
		// "-level", "3.0",                             // Compatible level
		// "-b:v", bitrate,                             // Video bitrate
		// "-maxrate", bitrate,                         // Maximum bitrate
		// "-bufsize", bitrate,                         // Buffer size
		// "-vf", fmt.Sprintf("scale=%s:-2", width),    // Scale video
		
		// Parameters to ensure consistent segment timing
		// "-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", config.SegmentLength), // Force keyframes at each segment boundary
		// "-sc_threshold", "0",                        // Disable scene change detection
		// "-sn",                                       // Disable subtitles
		
		// Audio parameters
		"-c:a", "aac",                               // AAC audio codec
		"-b:a", "128k",                              // Audio bitrate
		"-ac", "2",                                  // Stereo audio
		"-ar", "44100",                              // Standard audio sample rate
		
		// Stream mapping
		"-map", "0:v:0",                             // First video stream
		"-map", "0:a:0",                             // First audio stream
		
		// Segmentation parameters
		"-f", "segment",                             // Output format: segmented
		"-segment_time", fmt.Sprintf("%d", config.SegmentLength),
		"-segment_format", "mpegts",                 // MPEG-TS container for segments
		// "-segment_list", filepath.Join(config.OutputDir, "debug.m3u8"), // Debug playlist
		// "-segment_list_flags", "+live",              // Add the #EXT-X-ENDLIST tag at end of processing
		
		// Critical params for exact segment timing
		//"-reset_timestamps", "1",                    // Reset timestamps at each segment start
		// "-break_non_keyframes", "0",                 // Only break at keyframes
		// "-avoid_negative_ts", "make_zero",           // Handle negative timestamps by making them zero
		"-segment_list", segmentList,
		"-segment_list_type", "csv",
		"-segment_start_number", fmt.Sprintf("%d", startIndex),
		segmentPattern,
	}
	
	// Debug: Print the command being executed
	cmdString := "ffmpeg"
	for _, arg := range args {
		cmdString += " " + arg
	}
	log.Printf("Executing FFmpeg command: %s", cmdString)
	
	RunFFmpegCommand(args)
}

// TranscodeVideoOnly performs video-only transcoding (copy audio)
func TranscodeVideoOnly(sourcePath, width, bitrate string, audioStream int) {
	log.Printf("Starting video-only transcoding of: %s", sourcePath)
	
	// Create the segment pattern
	// segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts")
	segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts.temp")
	segmentList := filepath.Join(config.OutputDir, "segment_list.csv")
	
	// Build FFmpeg command
	args := []string{
		"-i", sourcePath,
		"-c:v", "libx264",
		"-c:a", "copy", // Copy original audio
		"-b:v", bitrate,
		"-vf", fmt.Sprintf("scale=%s:-2", width),
		"-preset", "veryfast",
		"-g", "48", // GOP size
		"-keyint_min", "48",
		"-sc_threshold", "0",
		"-force_key_frames", "expr:gte(t,n_forced*2)",
		"-map", "0:v:0",
		"-map", fmt.Sprintf("0:a:%d", audioStream),
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", config.SegmentLength),
		"-segment_format", "mpegts",
		"-segment_list", segmentList,
		"-segment_list_type", "csv",
		segmentPattern,
	}
	
	RunFFmpegCommand(args)
}

// TranscodeAudioOnly performs audio-only transcoding (copy video)
func TranscodeAudioOnly(sourcePath string, audioStream int) {
	log.Printf("Starting audio-only transcoding of: %s", sourcePath)
	
	// Create the segment pattern
	// segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts")

	segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts.temp")
	segmentList := filepath.Join(config.OutputDir, "segment_list.csv")
	startIndex := 0 / config.SegmentLength
	// Build FFmpeg command
	args := []string{
		// "-ss", "00:02:30", // Start at 2 minutes 30 seconds
		"-i", sourcePath,
		"-c:v", "copy", // Copy original video
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2", // Stereo audio
		"-map", "0:v:0",
		"-map", "0:a:0",
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", config.SegmentLength),
		"-segment_format", "mpegts",
		"-segment_start_number", fmt.Sprintf("%d", startIndex),
		"-segment_list", segmentList,
		"-segment_list_type", "csv",
		segmentPattern,
	}
	
	// Debug: Print the command being executed
	cmdString := "ffmpeg"
	for _, arg := range args {
		cmdString += " " + arg
	}
	log.Printf("Executing FFmpeg command: %s", cmdString)
	
	RunFFmpegCommand(args)
}

// SegmentVideo performs direct play with segmentation
func SegmentVideo(sourcePath string) {
	log.Printf("Segmenting video for direct play: %s", sourcePath)
	
	// Create the segment pattern
	segmentPattern := filepath.Join(config.OutputDir, "segment_%05d.ts")
	
	// Build FFmpeg command
	args := []string{
		"-i", sourcePath,
		"-c", "copy", // Copy both video and audio
		"-map", "0",
		"-f", "segment",
		"-segment_time", fmt.Sprintf("%d", config.SegmentLength),
		"-segment_format", "mpegts",
		segmentPattern,
	}
	
	RunFFmpegCommand(args)
}

// ExtractSegmentNumber extracts segment number from filename (e.g., "segment_00001.ts" -> 1)
func ExtractSegmentNumber(filename string) int {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0
	}
	
	numStr := strings.TrimSuffix(parts[1], ".ts")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	
	return num
}

// SortSegmentFiles sorts segment files by sequence number
func SortSegmentFiles(files []string) []string {
	// Extract segment numbers
	segments := make([]models.SegmentFile, 0, len(files))
	for _, file := range files {
		seqNum := ExtractSegmentNumber(filepath.Base(file))
		segments = append(segments, models.SegmentFile{Path: file, SeqNumber: seqNum})
	}
	
	// Sort segments by sequence number
	for i := 0; i < len(segments)-1; i++ {
		for j := i + 1; j < len(segments); j++ {
			if segments[i].SeqNumber > segments[j].SeqNumber {
				segments[i], segments[j] = segments[j], segments[i]
			}
		}
	}
	
	// Return sorted paths
	result := make([]string, len(segments))
	for i, segment := range segments {
		result[i] = segment.Path
	}
	
	return result
} 



