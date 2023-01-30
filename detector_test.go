package chardet_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogs/chardet"
)

func TestDetector(t *testing.T) {
	type file_charset_language struct {
		File     string
		IsHtml   bool
		Charset  string
		Language string
	}
	var data = []file_charset_language{
		{"utf8.html", true, "UTF-8", ""},
		{"utf8_bom.html", true, "UTF-8", ""},
		{"8859_1_en.html", true, "ISO-8859-1", "en"},
		{"8859_1_da.html", true, "ISO-8859-1", "da"},
		{"8859_1_de.html", true, "ISO-8859-1", "de"},
		{"8859_1_es.html", true, "ISO-8859-1", "es"},
		{"8859_1_fr.html", true, "ISO-8859-1", "fr"},
		{"8859_1_pt.html", true, "ISO-8859-1", "pt"},
		{"shift_jis.html", true, "Shift_JIS", "ja"},
		{"gb18030.html", true, "GB18030", "zh"},
		{"euc_jp.html", true, "EUC-JP", "ja"},
		{"euc_kr.html", true, "EUC-KR", "ko"},
		{"big5.html", true, "Big5", "zh"},
	}

	textDetector := chardet.NewTextDetector()
	htmlDetector := chardet.NewHtmlDetector()
	buffer := make([]byte, 32<<10)
	for _, d := range data {
		f, err := os.Open(filepath.Join("testdata", d.File))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		size, _ := io.ReadFull(f, buffer)
		input := buffer[:size]
		var detector = textDetector
		if d.IsHtml {
			detector = htmlDetector
		}
		result, err := detector.DetectBest(input)
		if err != nil {
			t.Fatal(err)
		}
		if result.Charset != d.Charset {
			t.Errorf("Expected charset %s, actual %s", d.Charset, result.Charset)
		}
		if result.Language != d.Language {
			t.Errorf("Expected language %s, actual %s", d.Language, result.Language)
		}
	}

	// "ノエル" Shift JIS encoded
	test := []byte("\x83m\x83G\x83\x8b")

	result, err := textDetector.DetectAll(test)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 results, actual %d", len(result))
	}
	if result[0].Charset != "Shift_JIS" || result[1].Charset != "GB18030" || result[2].Charset != "Big5" {
		t.Errorf("DetectAll order is wrong: %v", result)
	}

	singleResult, err := textDetector.DetectBest(test)
	if err != nil {
		t.Fatal(err)
	}
	if singleResult.Charset != "Shift_JIS" {
		t.Errorf("DetectBest result is wrong: %v", singleResult)
	}
}

func BenchmarkDetectBest(b *testing.B) {
	textDetector := chardet.NewTextDetector()
	aaaa := bytes.Repeat([]byte("A"), 1024)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		textDetector.DetectBest(aaaa)
	}
}
