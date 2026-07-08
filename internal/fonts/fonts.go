package fonts

import (
	"embed"
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/*.ttf
var fontFS embed.FS

// familyNameMap 把模板里使用的 CSS 字体名映射到磁盘文件名。
var familyNameMap = map[string]string{
	"ZCOOL XiaoWei":         "ZCOOLXiaoWei-Regular.ttf",
	"Ma Shan Zheng":         "MaShanZheng-Regular.ttf",
	"ZCOOL KuaiLe":          "ZCOOLKuaiLe-Regular.ttf",
	"Long Cang":             "LongCang-Regular.ttf",
	"ZCOOL QingKe HuangYou": "ZCOOLQingKeHuangYou-Regular.ttf",
}

// Library 保存已解析的字体。
type Library struct {
	fonts    map[string]*opentype.Font
	fallback string
}

// NewLibrary 加载并解析所有嵌入字体。
func NewLibrary() (*Library, error) {
	lib := &Library{fonts: make(map[string]*opentype.Font)}

	for family, filename := range familyNameMap {
		data, err := fontFS.ReadFile("fonts/" + filename)
		if err != nil {
			return nil, fmt.Errorf("read font %s: %w", filename, err)
		}
		f, err := opentype.Parse(data)
		if err != nil {
			return nil, fmt.Errorf("parse font %s: %w", filename, err)
		}
		lib.fonts[family] = f
		if lib.fallback == "" {
			lib.fallback = family
		}
	}

	return lib, nil
}

// Resolve 返回模板 fontFamily 中第一个可用的字体家族名。
func (lib *Library) Resolve(fontFamily string) string {
	// 模板里的 fontFamily 形如 '"ZCOOL XiaoWei", "PingFang SC", ...'
	// 我们简单按逗号拆分，取引号内的名字或第一个非空 token。
	start := 0
	for i := 0; i < len(fontFamily); i++ {
		c := fontFamily[i]
		if c == ',' {
			name := cleanName(fontFamily[start:i])
			if _, ok := lib.fonts[name]; ok {
				return name
			}
			start = i + 1
		}
	}
	name := cleanName(fontFamily[start:])
	if _, ok := lib.fonts[name]; ok {
		return name
	}
	return lib.fallback
}

func cleanName(s string) string {
	// 去掉首尾空白、引号。
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '"' || s[0] == '\'') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
		s = s[:len(s)-1]
	}
	return s
}

// FaceProvider 返回一个 closure，用于给 FontMeasurer 提供 font.Face。
func (lib *Library) FaceProvider(fontFamily string) func(fontSize float64) font.Face {
	family := lib.Resolve(fontFamily)
	f := lib.fonts[family]

	return func(fontSize float64) font.Face {
		face, err := opentype.NewFace(f, &opentype.FaceOptions{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		})
		if err != nil {
			// opentype.NewFace 理论上不会失败，但为了安全起见返回空 Face。
			face, _ = opentype.NewFace(f, &opentype.FaceOptions{Size: fontSize, DPI: 72})
		}
		return face
	}
}
