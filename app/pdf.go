package app

import (
    "path/filepath"

    pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

// PackToPDF packs chapter to .pdf
func PackToPDF(images []string, destination string) (string, error) {
    destination += ".pdf"

    err := RemoveIfExists(destination)
    if err != nil {
        return "", err
    }

    // Create parent directory since pdfcpu have some troubles when it doesn't exist
    if exists, err := Afero.Exists(filepath.Dir(destination)); err != nil {
        return "", err
    } else if !exists {
        if err := Afero.MkdirAll(filepath.Dir(destination), 0777); err != nil {
            return "", err
        }
    }

    if err := pdfcpu.ImportImagesFile(images, destination, nil, nil); err != nil {
        return "", err
    }

    return destination, nil
}
