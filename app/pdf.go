package app

import (
	"fmt"
	"image"
	_ "image/jpeg" // Import for JPEG decoding
	_ "image/png"  // Import for PNG decoding
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// PackageChapterToPDF creates a PDF from images in a directory.
// It sorts images alphabetically and scales them to fit A4 pages.
func PackageChapterToPDF(imageDir, mangaTitle, chapterHash string) error {
	files, err := os.ReadDir(imageDir)
	if err != nil {
		return fmt.Errorf("failed to read image directory %s: %w", imageDir, err)
	}

	var imageFiles []string
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				imageFiles = append(imageFiles, filepath.Join(imageDir, file.Name()))
			}
		}
	}

	if len(imageFiles) == 0 {
		return fmt.Errorf("no image files found in %s", imageDir)
	}

	// Sort files alphabetically to ensure correct page order
	sort.Strings(imageFiles)

	pdf := gofpdf.New("P", "mm", "A4", "") // Portrait, mm, A4
	pdf.SetTitle(fmt.Sprintf("%s - Chapter %s", mangaTitle, chapterHash), true)
	pdf.SetAuthor("manga-g", true) // Set author generically

	for _, imgPath := range imageFiles {
		// Get image dimensions without loading the whole image into memory
		file, err := os.Open(imgPath)
		if err != nil {
			fmt.Printf("Warning: Skipping file %s due to open error: %v\n", imgPath, err)
			continue
		}
		imgConfig, _, err := image.DecodeConfig(file)
		file.Close() // Close file immediately after getting config
		if err != nil {
			fmt.Printf("Warning: Skipping file %s due to decode config error: %v\n", imgPath, err)
			continue
		}

		// Use A4 and scale the image to fit the page width
		pageWidth, pageHeight := pdf.GetPageSize()
		margin := 10.0 // 10mm margin
		drawableWidth := pageWidth - (2 * margin)
		drawableHeight := pageHeight - (2 * margin)

		// Convert image pixels to mm (approx 96 DPI -> 1 px = 0.264583 mm)
		// Using DPI is tricky, let's scale based on aspect ratio relative to drawable area
		imgRatio := float64(imgConfig.Width) / float64(imgConfig.Height)
		pageRatio := drawableWidth / drawableHeight

		var scaledWidth, scaledHeight float64
		if imgRatio > pageRatio {
			// Image is wider than page area relative to height -> fit width
			scaledWidth = drawableWidth
			scaledHeight = scaledWidth / imgRatio
		} else {
			// Image is taller than page area relative to width -> fit height
			scaledHeight = drawableHeight
			scaledWidth = scaledHeight * imgRatio
		}

		pdf.AddPage()

		// Center the image on the page
		x := (pageWidth - scaledWidth) / 2
		y := (pageHeight - scaledHeight) / 2

		// Register image and place it
		pdf.ImageOptions(imgPath, x, y, scaledWidth, scaledHeight, false, gofpdf.ImageOptions{ReadDpi: false}, 0, "")
	}

	// Save the PDF to the parent directory
	pdfPath := filepath.Join(filepath.Dir(imageDir), fmt.Sprintf("%s_chapter_%s.pdf", mangaTitle, chapterHash))
	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to save PDF to %s: %w", pdfPath, err)
	}

	fmt.Printf("PDF generated successfully: %s\n", pdfPath) // Add success message here
	return nil
}
