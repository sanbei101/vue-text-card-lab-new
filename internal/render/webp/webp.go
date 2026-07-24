package webp

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

func SVGToWebP(svgData []byte) ([]byte, error) {
	params := vips.NewImportParams()
	img, err := vips.LoadImageFromBuffer(svgData, params)
	if err != nil {
		return nil, fmt.Errorf("加载 SVG 失败: %w", err)
	}
	defer img.Close()

	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = 85
	exportParams.Lossless = false

	webpBytes, _, err := img.ExportWebp(exportParams)
	if err != nil {
		return nil, fmt.Errorf("导出 WebP 失败: %w", err)
	}
	return webpBytes, nil
}
