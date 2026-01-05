package img

import (
	"bytes"
)

// IsImage 通过文件魔数判断是否为图片
func IsImage(data []byte) (bool, string) {
	if len(data) < 4 {
		return false, ""
	}

	// 检查常见图片格式的魔数
	// JPEG: FF D8 FF
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true, "image/jpeg"
	}

	// PNG: 89 50 4E 47
	if len(data) >= 4 && data[0] == 0x89 && data[1] == 0x50 &&
		data[2] == 0x4E && data[3] == 0x47 {
		return true, "image/png"
	}

	// GIF: 47 49 46 38 (GIF8)
	if len(data) >= 4 && data[0] == 0x47 && data[1] == 0x49 &&
		data[2] == 0x46 && data[3] == 0x38 {
		return true, "image/gif"
	}

	// WebP: 需要检查 RIFF 头部
	if len(data) >= 12 && bytes.Equal(data[0:4], []byte("RIFF")) &&
		bytes.Equal(data[8:12], []byte("WEBP")) {
		return true, "image/webp"
	}

	// BMP: 42 4D
	if len(data) >= 2 && data[0] == 0x42 && data[1] == 0x4D {
		return true, "image/bmp"
	}

	return false, ""
}
