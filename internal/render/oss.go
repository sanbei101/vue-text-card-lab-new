package render

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/purus-dev/aqua"
)

type R2Storage struct {
	client       *aqua.Client
	publicDomain string
}

func NewR2Storage(cfg aqua.Config, publicDomain string) *R2Storage {
	return &R2Storage{
		client:       aqua.NewClient(&cfg),
		publicDomain: publicDomain,
	}
}

func (s *R2Storage) UploadWebP(ctx context.Context, objectKey string, webpBytes []byte) (string, error) {
	reader := bytes.NewReader(webpBytes)

	err := s.client.UploadFile(ctx, objectKey, reader, "image/webp")
	if err != nil {
		return "", fmt.Errorf("上传 WebP 到 R2 失败: %w", err)
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(s.publicDomain, "/"), objectKey), nil

}
