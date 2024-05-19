package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/imgconv"
	"github.com/urfave/cli/v2"
)

func hasFilePath(path string) bool {
	return strings.Contains(path, ".")
}

func getFilenameNoExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func getPathToFile(path string) string {
	return strings.TrimSuffix(path, filepath.Base(path))
}

func getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		log.SetFlags(0)
		log.Fatal(err)
	}
	return wd
}

func isDir(path string) bool {
	if !(strings.Contains(path, `\`) || strings.Contains(path, "/")) {
		return false
	}

	stat, err := os.Stat(path)
	if err != nil {
		log.SetFlags(0)
		log.Fatal(err)
	}
	if stat.IsDir() {
		return true
	} else {
		return false
	}
}

func getDecodedImage(path string) image.Image {
	if !hasFilePath(filepath.Base(path)) {
		log.SetFlags(0)
		log.Fatal("Error: Make sure you include the full filename (including extension).")
	}

	src, err := imgconv.Open(path)
	if err != nil {
		log.SetFlags(0)
		log.Fatal("Error: Could not save image, make sure you aren't using single backslashes in the filepath.")
	}

	return src
}

func convertImage(path, desiredFormat, output string) {
	imgData := getDecodedImage(path)

	if output == "" {
		if strings.Contains(path, `\`) || strings.Contains(path, "/") {
			output = getPathToFile(path)
		}
	}

	newFilename := filepath.Join(output, filepath.Base(getFilenameNoExt(path))+"."+desiredFormat)
	desiredFormatAsFormat, err := imgconv.FormatFromExtension(desiredFormat)

	if err != nil {
		log.SetFlags(0)
		log.Fatal("Error: Desired format not supported. Supported Formats:\n\nJPEG (JPG)\nPNG\nGIF\nTIFF\nBMP\nWEBP (only from webp to another format)\nPDF")
	}

	fmtOption := imgconv.FormatOption{Format: desiredFormatAsFormat}
	imgconv.Save(newFilename, imgData, &fmtOption)

	if isDir(getPathToFile(newFilename)) {
		fmt.Println("\n Saved new image as", newFilename)
	} else {
		fmt.Println("\n Saved new image as", filepath.Join(getWorkingDirectory(), newFilename))
	}
}

func main() {
	app := &cli.App{
		Version: "v1.0.0",
		Authors: []*cli.Author{
			{Name: "Brayden O'Neil", Email: "oneilb123@gmail.com"},
		},
		Usage: "imgc conv [--image or -i] filepath [--to or -t] format",
		Commands: []*cli.Command{
			{
				Name:    "conv",
				Aliases: []string{"c"},
				Usage:   "conv [--image or -i] filepath [--to or -t] format",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "image",
						Aliases: []string{"i"},
						Usage:   "The image to be converted.",
					},
					&cli.StringFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "The format to convert the specified image to.",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Value:   "",
						Usage:   "The desired output path. Set to the same directory as the image by default",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.String("image")
					format := c.String("to")
					output := c.String("output")

					convertImage(path, format, output)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
