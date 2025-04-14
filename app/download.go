package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// DownloadProgressInfo represents generic progress update.
type DownloadProgressInfo struct {
	Completed int
	Total     int
	Percent   float64
	Error     error  // Used to signal errors during progress or final
	Done      bool   // Used to signal completion
	ChapterID string // Identify which download this belongs to
}

// DownloadState defines the structure for tracking download progress and metadata.
type DownloadState struct {
	MangaID        string   `json:"manga_id"`
	MangaTitle     string   `json:"manga_title"`
	ChapterID      string   `json:"chapter_id"`
	ChapterHash    string   `json:"chapter_hash"`
	BaseURL        string   `json:"base_url"`
	Files          []string `json:"files"`
	Downloaded     []bool   `json:"downloaded"`
	LastUpdated    string   `json:"last_updated"`
	TotalFiles     int      `json:"total_files"`
	CompletedFiles int      `json:"completed_files"`
}

// Downloader manages the state and execution of a chapter download.
type Downloader struct {
	state         DownloadState
	mutex         sync.Mutex
	active        bool
	taskChan      chan int
	wg            sync.WaitGroup
	NumWorkers    int                       // Configurable number of workers
	StateFilePath string                    // Configurable path for state file
	OutputDirBase string                    // Configurable base path for downloads
	ProgressChan  chan DownloadProgressInfo // Changed to bidirectional chan
}

// NewDownloader creates a new Downloader instance with configuration.
func NewDownloader(numWorkers int, stateFilePath, outputDirBase string, progressChan chan DownloadProgressInfo) *Downloader {
	if stateFilePath == "" {
		stateFilePath = "debug_downloads/download_state.json" // Default if empty
	}
	if outputDirBase == "" {
		outputDirBase = "debug_downloads" // Default if empty
	}
	return &Downloader{
		NumWorkers:    numWorkers,
		StateFilePath: stateFilePath,
		OutputDirBase: outputDirBase,
		ProgressChan:  progressChan, // Assign the bidirectional channel
	}
}

// --- State Management Methods ---

// saveState writes the current download state to the configured JSON file.
// IMPORTANT: Caller MUST hold the Downloader's mutex before calling this method.
func (d *Downloader) saveState() error {
	d.state.LastUpdated = time.Now().Format(time.RFC3339)

	stateDir := filepath.Dir(d.StateFilePath)
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	stateData, err := json.MarshalIndent(d.state, "", "  ")
	if err != nil {
		return err
	}

	// Use a temporary file for atomic write
	tmpStateFile := d.StateFilePath + ".tmp"
	err = os.WriteFile(tmpStateFile, stateData, 0644)
	if err != nil {
		os.Remove(tmpStateFile)
		return err
	}

	return os.Rename(tmpStateFile, d.StateFilePath)
}

// loadState attempts to load a previous download state from the configured file.
// Returns the loaded state, a bool indicating if found, and any error.
func (d *Downloader) loadState() (loadedState DownloadState, found bool, err error) {
	if _, statErr := os.Stat(d.StateFilePath); os.IsNotExist(statErr) {
		return DownloadState{}, false, nil
	}

	stateData, readErr := os.ReadFile(d.StateFilePath)
	if readErr != nil {
		return DownloadState{}, false, readErr
	}

	// Temporary state struct for unmarshalling without locking the main state yet
	var tempState DownloadState
	if unmarshalErr := json.Unmarshal(stateData, &tempState); unmarshalErr != nil {
		return DownloadState{}, false, unmarshalErr
	}

	return tempState, true, nil
}

// --- Download Control Methods ---

// setupSignalHandler creates a handler for interruption signals (Ctrl+C).
// It sets the downloader's active flag to false upon receiving a signal.
func (d *Downloader) setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nDownload paused by signal. State should be saved by workers.")
		d.mutex.Lock()
		d.active = false
		d.mutex.Unlock()
		// Don't exit here, let the main download function handle cleanup.
	}()
}

// monitorProgress updates the progress display periodically.
func (d *Downloader) monitorProgress() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		d.mutex.Lock()
		if !d.active {
			d.mutex.Unlock()
			return
		}
		completed := d.state.CompletedFiles
		total := d.state.TotalFiles
		chapterID := d.state.ChapterID
		d.mutex.Unlock()

		// Send generic progress message
		if d.ProgressChan != nil {
			percent := 0.0
			if total > 0 {
				percent = float64(completed) / float64(total)
			}
			d.ProgressChan <- DownloadProgressInfo{
				Completed: completed,
				Total:     total,
				Percent:   percent,
				ChapterID: chapterID,
			}
		}
	}
}

// DownloadChapter manages the multi-threaded download of chapter images.
// Runs synchronously and returns an error on failure.
func (d *Downloader) DownloadChapter(atHome *AtHomeResponse, mangaTitle string, chapterID string) error {
	var downloadErr error // Use local error variable

	// Try to load existing state for this specific chapter
	state, found, err := d.loadState()
	if err != nil {
		downloadErr = fmt.Errorf("error loading download state: %w", err)
		// Send error immediately if channel exists
		if d.ProgressChan != nil {
			d.ProgressChan <- DownloadProgressInfo{Error: downloadErr, Done: true, ChapterID: chapterID}
		}
		return downloadErr
	}

	// Initialize or resume state
	d.mutex.Lock()
	if found && state.ChapterID == chapterID && state.ChapterHash == atHome.Chapter.Hash {
		d.state = state // Resume using loaded state
		fmt.Printf("Resuming download for manga: %s, chapter: %s\n", d.state.MangaTitle, d.state.ChapterID)
	} else {
		// Start new download state
		d.state = DownloadState{
			MangaID:        "TBD", // TODO: Pass MangaID if available
			MangaTitle:     SanitizeFilename(mangaTitle),
			ChapterID:      chapterID,
			ChapterHash:    atHome.Chapter.Hash,
			BaseURL:        atHome.BaseURL,
			Files:          atHome.Chapter.Data,
			Downloaded:     make([]bool, len(atHome.Chapter.Data)),
			TotalFiles:     len(atHome.Chapter.Data),
			CompletedFiles: 0,
		}
	}
	d.mutex.Unlock()

	// Setup directory
	sanitizedTitle := d.state.MangaTitle
	outputDir := filepath.Join(d.OutputDirBase, sanitizedTitle, fmt.Sprintf("chapter_%s", d.state.ChapterHash[:8]))

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		downloadErr = fmt.Errorf("error creating directory %s: %w", outputDir, err)
		if d.ProgressChan != nil {
			d.ProgressChan <- DownloadProgressInfo{Error: downloadErr, Done: true, ChapterID: chapterID}
		}
		return downloadErr
	}

	fmt.Printf("Downloading %d images to: %s\n", d.state.TotalFiles, outputDir)
	fmt.Println("Press Ctrl+C to pause download (can be resumed later)")

	// Setup download infrastructure
	d.taskChan = make(chan int, d.state.TotalFiles)
	d.mutex.Lock()
	d.active = true
	d.mutex.Unlock()

	d.setupSignalHandler()

	// Verify existing files and update state
	d.mutex.Lock()
	completedCount := 0
	for i, isDownloaded := range d.state.Downloaded {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%03d_%s", i+1, d.state.Files[i]))
		if isDownloaded && fileExists(outputPath) {
			completedCount++
		} else if fileExists(outputPath) {
			d.state.Downloaded[i] = true
			completedCount++
		} else {
			d.state.Downloaded[i] = false
		}
	}
	d.state.CompletedFiles = completedCount

	err = d.saveState()
	if err != nil {
		fmt.Printf("Warning: Error saving initial download state: %v\n", err)
	}
	d.mutex.Unlock()

	// Start monitoring & workers
	progressDone := make(chan struct{}) // Channel to signal progress monitor completion
	go func() {
		d.monitorProgress()
		close(progressDone)
	}()

	d.wg = sync.WaitGroup{}
	for w := 1; w <= d.NumWorkers; w++ {
		d.wg.Add(1)
		go d.downloadWorker(outputDir, w)
	}

	// Queue tasks
	var tasksQueued int
	for i := range d.state.Files {
		d.mutex.Lock()
		shouldQueue := !d.state.Downloaded[i]
		d.mutex.Unlock()
		if shouldQueue {
			d.taskChan <- i
			tasksQueued++
		}
	}
	fmt.Printf("Queued %d tasks for download.\n", tasksQueued)
	close(d.taskChan)

	// Wait for workers completion
	d.wg.Wait()

	// Stop progress monitor
	d.mutex.Lock()
	d.active = false // Ensure monitor stops
	d.mutex.Unlock()
	<-progressDone // Wait for monitor to finish

	// Final state update
	d.mutex.Lock()
	d.state.CompletedFiles = 0
	for _, downloaded := range d.state.Downloaded {
		if downloaded {
			d.state.CompletedFiles++
		}
	}
	d.state.CompletedFiles = d.state.CompletedFiles
	err = d.saveState()
	if err != nil {
		fmt.Printf("Warning: Error saving final download state: %v\n", err)
	}
	success := d.state.CompletedFiles == d.state.TotalFiles

	if success {
		fmt.Println("\nDownload complete!")
		os.Remove(d.StateFilePath)
	} else {
		fmt.Printf("\nWarning: Only %d out of %d files processed. Run again to resume.\n", d.state.CompletedFiles, d.state.TotalFiles)
	}
	d.mutex.Unlock()

	// Send final message via channel
	if d.ProgressChan != nil {
		percent := 0.0
		if d.state.TotalFiles > 0 {
			percent = float64(d.state.CompletedFiles) / float64(d.state.TotalFiles)
		}
		d.ProgressChan <- DownloadProgressInfo{
			Completed: d.state.CompletedFiles,
			Total:     d.state.TotalFiles,
			Percent:   percent,
			Done:      true,
			Error:     downloadErr,
			ChapterID: chapterID,
		}
	}

	// Return the error state
	return downloadErr
}

// --- Internal Worker & Helpers (now methods) ---

// downloadWorker is executed by goroutines to download individual images.
func (d *Downloader) downloadWorker(outputDir string, workerID int) {
	defer d.wg.Done()

	for taskIndex := range d.taskChan {
		d.mutex.Lock()
		if !d.active {
			d.mutex.Unlock()
			return
		}
		if taskIndex >= len(d.state.Files) {
			fmt.Printf("Worker %d: Invalid task index %d.\n", workerID, taskIndex)
			d.mutex.Unlock()
			continue
		}
		filename := d.state.Files[taskIndex]
		baseURL := d.state.BaseURL
		chapterHash := d.state.ChapterHash
		isAlreadyMarkedDone := d.state.Downloaded[taskIndex]
		d.mutex.Unlock()

		if isAlreadyMarkedDone {
			continue
		}

		imageURL := fmt.Sprintf("%s/data/%s/%s", baseURL, chapterHash, filename)
		outputPath := filepath.Join(outputDir, fmt.Sprintf("%03d_%s", taskIndex+1, filename))

		if fileExists(outputPath) {
			d.mutex.Lock()
			if !d.state.Downloaded[taskIndex] {
				d.state.Downloaded[taskIndex] = true
				d.state.CompletedFiles++
				d.saveState()
			}
			d.mutex.Unlock()
			continue
		}

		tmpPath := outputPath + ".tmp"
		os.MkdirAll(filepath.Dir(outputPath), 0755)
		success := d.downloadFileWithRetry(imageURL, tmpPath)

		d.mutex.Lock()
		if !d.active {
			d.state.Downloaded[taskIndex] = false
			os.Remove(tmpPath)
			d.saveState()
			d.mutex.Unlock()
			return
		}

		if success {
			err := os.Rename(tmpPath, outputPath)
			if err != nil {
				fmt.Printf("\nWorker %d: Error renaming temp file %s: %v\n", workerID, tmpPath, err)
				d.state.Downloaded[taskIndex] = false
				os.Remove(tmpPath)
			} else {
				d.state.Downloaded[taskIndex] = true
				d.state.CompletedFiles++
			}
		} else {
			fmt.Printf("\nWorker %d: Failed to download %s after retries\n", workerID, filename)
			d.state.Downloaded[taskIndex] = false
			os.Remove(tmpPath)
		}
		d.saveState()
		d.mutex.Unlock()
	}
}

// downloadFileWithRetry attempts to download a file with exponential backoff.
func (d *Downloader) downloadFileWithRetry(url, outputPath string) bool {
	for retry := 0; retry < 3; retry++ {
		d.mutex.Lock()
		isActive := d.active
		d.mutex.Unlock()
		if !isActive {
			return false
		}

		if downloadFile(url, outputPath) {
			return true
		}

		d.mutex.Lock()
		isActive = d.active
		d.mutex.Unlock()
		if !isActive {
			return false
		}
		time.Sleep(time.Second * time.Duration(1<<retry))
	}
	return false
}

// downloadFile performs a single download attempt (remains a package function).
func downloadFile(url, outputPath string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("Referer", "https://mangadex.org/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err == nil
}

// SanitizeFilename and fileExists remain package functions as they don't depend on Downloader state.
// SanitizeFilename removes invalid characters from filenames.
func SanitizeFilename(name string) string {
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	return strings.TrimSpace(result)
}

// fileExists checks if a file exists and is not empty.
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("\nWARNING: Unexpected error checking file existence for %s: %v\n", filePath, err)
		}
		return false
	}
	return !info.IsDir() && info.Size() > 0
}
