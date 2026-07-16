package render

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var sfntBufferPool = sync.Pool{
	New: func() any {
		return &sfnt.Buffer{}
	},
}

// TextToPath 把字符串按 sfnt 字形轮廓转成一条 SVG path d 属性,并返回总宽度。
func TextToPath(f *sfnt.Font, text string, fontSize float64) (string, float64, error) {
	var d strings.Builder
	d.Grow(len(text) * 200)

	var x fixed.Int26_6
	var prev sfnt.GlyphIndex
	var hasPrev bool

	ppem := fixed.Int26_6(fontSize * 64)

	buf := sfntBufferPool.Get().(*sfnt.Buffer)
	defer sfntBufferPool.Put(buf)

	for _, r := range text {
		idx, err := f.GlyphIndex(buf, r)
		if err != nil {
			idx = 0
		}

		segments, err := f.LoadGlyph(buf, idx, ppem, nil)
		if err == nil {
			appendSegments(&d, segments, x)
		}

		if hasPrev {
			kern, err := f.Kern(buf, prev, idx, ppem, font.HintingNone)
			if err == nil {
				x += kern
			}
		}

		advance, err := f.GlyphAdvance(buf, idx, ppem, font.HintingNone)
		if err == nil {
			x += advance
		}

		prev = idx
		hasPrev = true
	}

	return d.String(), float64(x) / 64.0, nil
}

func appendSegments(d *strings.Builder, segs sfnt.Segments, offset fixed.Int26_6) {
	for _, seg := range segs {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			p := seg.Args[0]
			d.WriteString("M")
			writePoint(d, addX(p, offset), p.Y)
		case sfnt.SegmentOpLineTo:
			p := seg.Args[0]
			d.WriteString("L")
			writePoint(d, addX(p, offset), p.Y)
		case sfnt.SegmentOpQuadTo:
			c := seg.Args[0]
			p := seg.Args[1]
			d.WriteString("Q")
			writePoint(d, addX(c, offset), c.Y)
			writePoint(d, addX(p, offset), p.Y)
		case sfnt.SegmentOpCubeTo:
			c1 := seg.Args[0]
			c2 := seg.Args[1]
			p := seg.Args[2]
			d.WriteString("C")
			writePoint(d, addX(c1, offset), c1.Y)
			writePoint(d, addX(c2, offset), c2.Y)
			writePoint(d, addX(p, offset), p.Y)
		}
	}
}

func addX(p fixed.Point26_6, dx fixed.Int26_6) fixed.Int26_6 {
	return fixed.Int26_6(int64(p.X) + int64(dx))
}

func writePoint(d *strings.Builder, x, y fixed.Int26_6) {
	d.WriteByte(' ')
	writeFixed(d, x)
	d.WriteByte(' ')
	writeFixed(d, y)
}

func writeFixed(d *strings.Builder, v fixed.Int26_6) {
	val := math.Round(float64(v)*10/64.0) / 10.0
	var buf [32]byte
	d.Write(strconv.AppendFloat(buf[:0], val, 'f', -1, 64))
}
