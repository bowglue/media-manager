package models

import (
	"os/exec"
	"sync"
	"time"
)

// FileInfo stores information about a media file
type FileInfo struct {
	VideoCodec string
	AudioCodec string
	Duration   string
}

// SegmentFile represents a segment file with a sequence number
type SegmentFile struct {
	Path      string
	SeqNumber int
}

// TranscodeSession stores information about an active transcoding session
type TranscodeSession struct {
	ID             string
	SourcePath     string
	Width          string
	Bitrate        string
	AudioStream    int
	StartTime      time.Time
	LastPaused     time.Time
	TotalPauseTime time.Duration
	CurrentSegment int
	IsActive       bool
	Cmd            *exec.Cmd
	OutputDir      string
	Command        []string // Store the original command for resuming
}

// TranscodeManager manages multiple transcoding sessions
type TranscodeManager struct {
	Sessions map[string]*TranscodeSession
	Lock     sync.RWMutex
} 