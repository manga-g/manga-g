package app

import (
	"github.com/signintech/gopdf"
)

// createPDF creates a PDF from a list of image file paths, using the specified
// page size, font path, font name, and font size. It saves the PDF to a file
// named "compiled.pdf".
func createPDF(imageFilePaths []string, pageSize gopdf.Rect, fontPath, fontName string, fontSize float64) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: pageSize})

	// Add the specified font to the PDF.
	err := addFont(pdf, fontPath, fontName)
	if err != nil {
		return err
	}

	// Add the images from the file paths provided in the imageFilePaths array
	// to the PDF, followed by a new page.
	err = addImages(pdf, imageFilePaths)
	if err != nil {
		return err
	}

	// Save the PDF to a file.
	err = pdf.WritePdf("compiled.pdf")
	if err != nil {
		return err
	}

	return nil
}

// addFont adds the specified font to the PDF.
func addFont(pdf gopdf.GoPdf, fontPath, fontName string) error {
	return pdf.AddTTFFont(fontName, fontPath)
}

// addImages adds the images from the file paths provided in the imageFilePaths
// array to the PDF, followed by a new page.
func addImages(pdf gopdf.GoPdf, imageFilePaths []string) error {
	// Loop through the list of image file paths.
	for _, filePath := range imageFilePaths {
		pdf.AddPage()

		// Load the image.
		err := pdf.Image(filePath, 0, 0, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
