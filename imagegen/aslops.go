package imagegen

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	
	"github.com/nfnt/resize"
)



func AsciiGen(imageURL string,reqwidth int) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "[Image Unavailable]", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	var img image.Image
	switch {
	case strings.Contains(contentType, "png"):
		img, err = png.Decode(bytes.NewReader(body))
	case strings.Contains(contentType, "jpeg"), strings.Contains(contentType, "jpg"):
		img, err = jpeg.Decode(bytes.NewReader(body))
	default:
		img, _, err = image.Decode(bytes.NewReader(body))
	}
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	const newWidth = 2048
	img = resize.Resize(newWidth, 0, img, resize.Lanczos3)

	ascii := convertToAscii(img)
	ascii = downscaleAscii(ascii, newWidth, reqwidth)
	trimascii := trimAndPadAscii(ascii)
	return trimascii, nil
}

func convertToAscii(img image.Image) string {
	asciiChars := " .:-=+*#%@"
	var asciiArt strings.Builder
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			charIndex := int(brightness * float64(len(asciiChars)-1) / 65535)
			char := asciiChars[charIndex]
			asciiArt.WriteByte(char)
		}
		asciiArt.WriteString("\n")
	}

	return asciiArt.String()
}



func downscaleAscii(ascii string, originalWidth, targetWidth int) string {
	lines := strings.Split(ascii, "\n")
	var downscaled strings.Builder
	ratio := originalWidth / targetWidth

	for i := 0; i < len(lines); i += ratio {
		for j := 0; j < originalWidth; j += ratio {
			if j < len(lines[i]) {
				downscaled.WriteByte(lines[i][j])
			}
		}
		downscaled.WriteString("\n")
	}
	return downscaled.String()
}

func trimAndPadAscii(ascii string) string {
    lines := strings.Split(ascii, "\n")

    for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
        lines = lines[1:]
    }

    for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
        lines = lines[:len(lines)-1]
    }


    minLeadingSpaces := len(lines[0])
    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)
        if len(trimmedLine) > 0 {
            leadingSpaces := len(line) - len(strings.TrimLeft(line, " "))
            if leadingSpaces < minLeadingSpaces {
                minLeadingSpaces = leadingSpaces
            }
        }
    }

    for i, line := range lines {
        trimmed := line[minLeadingSpaces:]
        lines[i] = "  " + trimmed + "  " 
    }

    lines = append([]string{""}, lines...)
    lines = append(lines, "")

    return strings.Join(lines, "\n")
}
