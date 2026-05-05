package transcoder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	Resolution1080p = "1080p"
	Resolution720p  = "720p"
	Resolution480p  = "480p"
	Resolution360p  = "360p"
)

type Variant struct {
	Quality     string
	ManifestURL string
	Bitrate     int
	FileSize    int64
	Width       int
	Height      int
}

type Result struct {
	Qualities []Variant
	Duration  int // seconds
	ThumbPath string
}

// TranscodeToHLS converts a raw video file into HLS segments with multiple quality variants.
func TranscodeToHLS(ctx context.Context, inputPath string, outputDir string) (Result, error) {
	if _, err := os.Stat(inputPath); err != nil {
		return Result{}, fmt.Errorf("input file not found: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create output dir: %w", err)
	}

	duration, err := probeDuration(ctx, inputPath)
	if err != nil {
		return Result{}, fmt.Errorf("probe duration: %w", err)
	}

	variants := []struct {
		name   string
		width  int
		height int
		bitrate int
	}{
		{Resolution1080p, 1920, 1080, 5000},
		{Resolution720p, 1280, 720, 2800},
		{Resolution480p, 854, 480, 1400},
		{Resolution360p, 640, 360, 800},
	}

	var result Result
	result.Duration = duration

	for _, v := range variants {
		qualDir := filepath.Join(outputDir, v.name)
		if err := os.MkdirAll(qualDir, 0o755); err != nil {
			continue
		}

		manifestPath := filepath.Join(qualDir, "index.m3u8")
		segmentPattern := filepath.Join(qualDir, "segment_%03d.ts")

		args := []string{
			"-i", inputPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:v", strconv.Itoa(v.bitrate) + "k",
			"-maxrate", strconv.Itoa(v.bitrate) + "k",
			"-bufsize", strconv.Itoa(v.bitrate*2) + "k",
			"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", v.width, v.height),
			"-preset", "fast",
			"-g", "48",
			"-sc_threshold", "0",
			"-hls_time", "4",
			"-hls_list_size", "0",
			"-hls_segment_filename", segmentPattern,
			manifestPath,
		}

		cmd := exec.CommandContext(ctx, "ffmpeg", args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return Result{}, fmt.Errorf("ffmpeg %s failed: %w\n%s", v.name, err, string(output))
		}

		var fileSize int64
		entries, _ := os.ReadDir(qualDir)
		for _, e := range entries {
			if info, err := e.Info(); err == nil {
				fileSize += info.Size()
			}
		}

		result.Qualities = append(result.Qualities, Variant{
			Quality:     v.name,
			ManifestURL: filepath.Join(outputDir, v.name, "index.m3u8"),
			Bitrate:     v.bitrate,
			FileSize:    fileSize,
			Width:       v.width,
			Height:      v.height,
		})
	}

	// Generate thumbnail (seek to 1s — safe for short videos)
	thumbPath := filepath.Join(outputDir, "thumbnail.jpg")
	thumbArgs := []string{
		"-ss", "00:00:01",
		"-i", inputPath,
		"-vframes", "1",
		"-vf", "scale=320:180",
		thumbPath,
	}
	if err := exec.CommandContext(ctx, "ffmpeg", thumbArgs...).Run(); err == nil {
		result.ThumbPath = thumbPath
	}

	return result, nil
}

func probeDuration(ctx context.Context, path string) (int, error) {
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}
