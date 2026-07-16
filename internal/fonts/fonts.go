package fonts

import (
	"embed"
	"fmt"
	"sync"

	"github.com/phuslu/log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"

	"github.com/sanbei101/blue-card-engine/internal/cardengine"
)

//go:embed fonts/*.ttf
var fontFS embed.FS

type FontFamily string

const (
	FontZCOOLXiaoWei        FontFamily = "ZCOOL XiaoWei"
	FontMaShanZheng         FontFamily = "Ma Shan Zheng"
	FontZCOOLKuaiLe         FontFamily = "ZCOOL KuaiLe"
	FontLongCang            FontFamily = "Long Cang"
	FontZCOOLQingKeHuangYou FontFamily = "ZCOOL QingKe HuangYou"
)

var familyFilenames = map[FontFamily]string{
	FontZCOOLXiaoWei:        "ZCOOLXiaoWei-Regular.ttf",
	FontMaShanZheng:         "MaShanZheng-Regular.ttf",
	FontZCOOLKuaiLe:         "ZCOOLKuaiLe-Regular.ttf",
	FontLongCang:            "LongCang-Regular.ttf",
	FontZCOOLQingKeHuangYou: "ZCOOLQingKeHuangYou-Regular.ttf",
}

type Library struct {
	fonts map[FontFamily]*opentype.Font
	sfnts map[FontFamily]*sfnt.Font

	mu        sync.Mutex
	faceCache map[FontFamily]map[float64]font.Face
	measurers map[FontFamily]*cardengine.FontMeasurer
}

// NewLibrary 加载并解析所有嵌入字体。
func NewLibrary() (*Library, error) {
	lib := &Library{
		fonts:     make(map[FontFamily]*opentype.Font),
		sfnts:     make(map[FontFamily]*sfnt.Font),
		faceCache: make(map[FontFamily]map[float64]font.Face),
		measurers: make(map[FontFamily]*cardengine.FontMeasurer),
	}

	for family, filename := range familyFilenames {
		data, err := fontFS.ReadFile("fonts/" + filename)
		if err != nil {
			return nil, fmt.Errorf("read font %s: %w", filename, err)
		}

		f, err := opentype.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("parse opentype font %s: %w", filename, err)
		}

		sfntFont, err := sfnt.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("parse sfnt font %s: %w", filename, err)
		}

		lib.fonts[family] = f
		lib.sfnts[family] = sfntFont
		lib.faceCache[family] = make(map[float64]font.Face)
		lib.measurers[family] = cardengine.NewFontMeasurer(lib.FaceProvider(family), 4096)
	}

	return lib, nil
}

// FaceProvider 现在自带字号缓存
func (lib *Library) FaceProvider(fontFamily FontFamily) func(fontSize float64) font.Face {
	f := lib.fonts[fontFamily]

	return func(fontSize float64) font.Face {
		lib.mu.Lock()
		defer lib.mu.Unlock()

		cache := lib.faceCache[fontFamily]
		if face, ok := cache[fontSize]; ok {
			return face
		}

		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		})
		if err != nil {
			log.Panic().Err(err).Msg("failed to create font face")
			return nil
		}

		cache[fontSize] = face
		return face
	}
}

func (lib *Library) Measurer(fontFamily FontFamily) *cardengine.FontMeasurer {
	return lib.measurers[fontFamily]
}

func (lib *Library) SFNT(fontFamily FontFamily) *sfnt.Font {
	return lib.sfnts[fontFamily]
}
