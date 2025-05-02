# MangaDex API Debug Tool

A command-line utility for testing and debugging the MangaDex API integration in the manga-g application.

## Features

- Search for manga by title
- Get manga details by ID
- Get chapter details and image server information by ID
- Advanced manga chapter image downloading:
  - Multi-threaded downloads for speed
  - Resume capability for interrupted downloads
  - Organized by manga title and chapter
  - Duplicate prevention
  - Progress tracking
- Run with mock data for offline testing
- View API endpoints
- Interactive manga and chapter selection

## Usage

```bash
# Show available API endpoints
go run cmd/debug/main.go -test

# Search for manga by title (with interactive selection)
go run cmd/debug/main.go -search "one piece"
# Then select a manga by number to view details and chapters
# Then select a chapter by number to view chapter details
# Choose whether to download images when prompted
# Specify number of download threads (1-10)

# Get manga details and chapters by ID
go run cmd/debug/main.go -manga "some-manga-id"

# Get chapter details and image server information
go run cmd/debug/main.go -chapter "some-chapter-id"

# Use verbose mode to see raw JSON responses
go run cmd/debug/main.go -search "one piece" -v

# Use mock data (no API calls)
go run cmd/debug/main.go -mock -search "any title"
go run cmd/debug/main.go -mock -manga "any-id"

# Use a custom API URL
go run cmd/debug/main.go -api "https://custom-api-url.com" -search "one piece"
```

## Options

- `-search <title>`: Search for manga by title
- `-manga <id>`: Get manga details by ID
- `-chapter <id>`: Get chapter details by ID
- `-test`: Show available API endpoints (no API calls)
- `-mock`: Use mock data instead of real API calls
- `-v`: Verbose output (print raw JSON)
- `-api <url>`: Use custom API URL (default: https://api.mangadex.org)

## Interactive Navigation

The debug tool supports interactive navigation through the MangaDex API:

1. Search for manga: `go run cmd/debug/main.go -search "your search"`
2. Select a manga from the numbered results
3. View manga details and available chapters
4. Select a chapter to see chapter details and image server information
5. Choose to download chapter images when prompted
6. Specify number of download threads for parallel downloading

This flow mimics the user experience in the main application and helps with testing the entire process.

## Advanced Download Features

The tool includes a robust download system with the following features:

### Multi-threaded Downloads
- Specify number of download threads (1-10)
- Dramatically speeds up downloading large chapters
- Each thread handles downloading individual images in parallel

### Resume Capability
- Downloads can be paused at any time with Ctrl+C
- Progress is automatically saved to a state file
- Run the same command again to be prompted to resume
- No duplicate downloads - picks up exactly where it left off

### Organized File Structure
- Downloads are organized by manga title and chapter
- Directory structure: `debug_downloads/[Manga Title]/chapter_[hash]`
- Files are named sequentially: `001_image.png`, `002_image.jpg`, etc.
- Sanitizes filenames to be compatible with all operating systems

### Error Handling
- Automatic retries for failed downloads (up to 3 attempts)
- Temporary files used during download to prevent corrupt images
- Detailed progress tracking with percentage and file counts
- Graceful handling of network interruptions

## Image Downloads

The tool can download manga chapter images to your local machine:

- Images are saved to the `debug_downloads/[Manga Title]/chapter_[hash]` directory
- Files are named with sequential numbers and original filenames (`001_image.png`, etc.)
- A progress bar shows download status with percentage and file counts
- In mock mode, placeholder files are created to simulate the download process
- Press Ctrl+C at any time to pause downloading (can be resumed later)

## Purpose

This tool helps developers:

1. Understand the MangaDex API structure and endpoints
2. Test API requests and responses
3. Debug parsing issues with the JSON responses
4. Test application code without making actual API calls
5. Explore the available data from the MangaDex API
6. Test the image download functionality before implementing it in the GUI
7. Develop and refine multi-threaded download capabilities

## API Endpoints

The MangaDex API v5 endpoints used by this tool include:

- **Search Manga**: `/manga?title={title}&limit={limit}`
- **Manga Details**: `/manga/{id}`
- **Manga Chapters**: `/manga/{id}/feed?translatedLanguage[]={lang}&limit={limit}`
- **Chapter Details**: `/chapter/{id}`
- **At-Home Server**: `/at-home/server/{id}`

#### 4. Download Chapter Images

Downloads chapter images using the `mangadex-at-home` server information.

- **Function**: `DownloadChapterImages(atHome *md.AtHomeServerResponse, outputDir string) error`
- **Input**: `md.AtHomeServerResponse` (from `GetMangaChapterAtHome`), `outputDir` (base directory for saving)
- **Output**: `error` (if any issues occur during download)
- **Directory structure**: `manga/[Manga Title]/chapter_[hash]`
- **Description**: Fetches each page image URL from `atHome.Chapter.DataSaver`, downloads the image, and saves it to the specified directory. Handles potential download errors.

#### 5. Get Manga Volumes & Chapters

Retrieves volume and chapter information for a specific manga.

- **Command**: `go run main.go download --url <chapter_url>`
- **Example**: `go run main.go download --url https://mangadex.org/chapter/some-chapter-uuid`
- **Description**:
  - Parses the chapter URL to get the chapter UUID.
  - Fetches chapter details using the UUID.
  - Retrieves the manga title using the manga UUID from the chapter details.
  - Fetches the `mangadex-at-home` server information for the chapter.
  - Downloads the chapter images.
- **Output**:
  - Logs download progress for each page.
  - Images are saved to the `manga/[Manga Title]/chapter_[hash]` directory

### State Management

The download command uses a `download_state.json` file (in the `manga` directory by default) to keep track of downloaded chapters and pages. This allows the command to:

- Keep track of downloaded chapters and pages
- Resume downloads from where they left off
- Prevent duplicate downloads
- Provide detailed progress tracking 