package go_wm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	pfont "gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/plotter"
)

type WFont struct {
	Uri      string
	LocTTF   string
	TypeFace pfont.Typeface

	Size   int64
	Weight font.Weight
	Style  font.Style
}

// download font from debian
// url = "http://http.debian.net/debian/pool/main/f/fonts-ipafont/fonts-ipafont_00303.orig.tar.gz"
// ttf, err := untargz("IPAfont00303/ipam.ttf", resp.Body)
func (f *WFont) GetFont() (*pfont.Face, error) {
	if f.Uri == "" {
		cache := pfont.NewCache(liberation.Collection())
		fontStyle := pfont.Font{
			Typeface: "Liberation",
			Variant:  "Sans",
			Weight:   f.Weight,
			Style:    f.Style,
			Size:     pfont.Length(f.Size),
		}
		fnt := cache.Lookup(fontStyle, fontStyle.Size)

		return &fnt, nil
	}

	var ttf []byte
	if strings.Contains(f.Uri, "http") {
		resp, err := http.Get(f.Uri)
		if err != nil {
			log.Fatalf("could not download IPA font file: %+v", err)
			return nil, err
		}
		defer resp.Body.Close()

		font, err := untargz(f.LocTTF, resp.Body)
		if err != nil {
			log.Fatalf("could not untar archive: %+v", err)
			return nil, err
		}

		ttf = font
	} else {
		f, err := os.Open(f.Uri)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		b := bytes.Buffer{}
		b.ReadFrom(f)

		ttf = b.Bytes()
	}

	fontTTF, err := opentype.Parse(ttf)
	if err != nil {
		return nil, err
	}

	font := pfont.Font{
		Typeface: f.TypeFace,
		Weight:   f.Weight,
		Size:     pfont.Length(f.Size),
		Style:    f.Style,
	}
	pfont.DefaultCache.Add([]pfont.Face{
		{
			Font: font,
			Face: fontTTF,
		},
	})
	if !pfont.DefaultCache.Has(font) {
		log.Fatalf("no font %q!", font.Typeface)
		return nil, fmt.Errorf("no font %q!", font.Typeface)
	}
	plot.DefaultFont = font
	plotter.DefaultFont = font

	cache := pfont.DefaultCache
	fc := cache.Lookup(font, font.Size)

	return &fc, nil
}

func untargz(name string, r io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("could not find %q in tar archive", name)
			}
			return nil, fmt.Errorf("could not extract header from tar archive: %v", err)
		}

		if hdr == nil || hdr.Name != name {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, tr)
		if err != nil {
			return nil, fmt.Errorf("could not extract %q file from tar archive: %v", name, err)
		}
		return buf.Bytes(), nil
	}
}
