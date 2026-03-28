package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"

	"github.com/beevik/etree"
)

func main() {
	letter := flag.String("letter", "ABCDEFGHIJKLMNOPQRSTUVWXYZÑ", "letter to export")
	template := flag.String("template", "letter-template-01.svg", "Template SVG file to use")
	letterid := flag.String("letterid", "flowPara2993", "Id inside the SVG of the letter")
	id := flag.String("id", "use3075", "Id inside the SVG of section to export")
	flag.Parse()

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(*template); err != nil {
		log.Fatalf("failed to parse SVG: %v", err)
	}

	elements := doc.FindElements("//*[@id='" + *letterid + "']")
	if len(elements) == 0 {
		log.Fatalf("no element found with id %q", *letterid)
	}
	textNode := elements[0]

	var generated []string

	for _, letra := range *letter {
		char := string(letra)
		textNode.SetText(char)

		tmpFile, err := os.CreateTemp("", "letterspritegen-*.svg")
		if err != nil {
			log.Fatalf("failed to create temp file: %v", err)
		}
		tmpName := tmpFile.Name()

		if _, err := doc.WriteTo(tmpFile); err != nil {
			log.Fatalf("failed to write SVG: %v", err)
		}
		tmpFile.Close()

		outputName := fmt.Sprintf("gen-%s.png", char)
		fmt.Printf("- Generated temp svg with letter %s  (%s)\n", char, tmpName)

		cmd := exec.Command("inkscape",
			"--export-id="+*id,
			"--export-filename="+outputName,
			"--export-dpi=180",
			tmpName,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("inkscape error for %q: %v", char, err)
			continue
		}
		fmt.Printf("  saved %s\n", outputName)
		generated = append(generated, outputName)
	}

	if len(generated) == 0 {
		return
	}

	// Find the maximum dimension across all generated PNGs.
	maxSize := 0
	for _, name := range generated {
		f, err := os.Open(name)
		if err != nil {
			log.Printf("could not open %s to measure: %v", name, err)
			continue
		}
		cfg, err := png.DecodeConfig(f)
		f.Close()
		if err != nil {
			log.Printf("could not decode config for %s: %v", name, err)
			continue
		}
		if cfg.Width > maxSize {
			maxSize = cfg.Width
		}
		if cfg.Height > maxSize {
			maxSize = cfg.Height
		}
	}

	fmt.Printf("\nPadding all images to %dx%d (transparent background, centered)...\n", maxSize, maxSize)

	// Pad each PNG to maxSize×maxSize, centered, transparent background.
	for _, name := range generated {
		if err := padImage(name, maxSize); err != nil {
			log.Printf("failed to pad %s: %v", name, err)
		} else {
			fmt.Printf("  padded %s\n", name)
		}
	}
}

func padImage(path string, size int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	src, err := png.Decode(f)
	f.Close()
	if err != nil {
		return err
	}

	dst := image.NewNRGBA(image.Rect(0, 0, size, size))
	// dst is zero-initialized → fully transparent

	b := src.Bounds()
	offsetX := (size - b.Dx()) / 2
	offsetY := (size - b.Dy()) / 2
	draw.Draw(dst, b.Add(image.Pt(offsetX, offsetY)), src, b.Min, draw.Src)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, dst)
}
