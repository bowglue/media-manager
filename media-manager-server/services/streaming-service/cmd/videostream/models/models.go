package models

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