package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
	"west2-video/common/config"
	"west2-video/common/img"
)

var Uploader = NewHTTPUploader(config.C.UpC.APIURL, time.Duration(config.C.UpC.Timeout)*time.Second)

// getExtensionFromMimeType 根据 MIME 类型返回文件扩展名
func getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

// HTTPUploader HTTP 上传客户端
type HTTPUploader struct {
	apiURL  string
	timeout time.Duration
	client  *http.Client
}

// UploadResponse 上传接口响应结构
type UploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	} `json:"data"`
}

// NewHTTPUploader 创建 HTTP 上传客户端
func NewHTTPUploader(apiURL string, timeout time.Duration) *HTTPUploader {
	if timeout <= 0 {
		timeout = 30 * time.Second // 默认30秒
	}

	return &HTTPUploader{
		apiURL:  apiURL,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// UploadImage 上传图片到服务器
// data: 图片二进制数据
// filename: 文件名（可选，用于设置 Content-Disposition）
// 返回: 图片访问URL
func (u *HTTPUploader) UploadImage(data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("图片数据为空")
	}

	// 验证图片格式（只支持 JPG、PNG、GIF、WEBP）
	isImage, mimeType := img.IsImage(data)
	if !isImage {
		return "", fmt.Errorf("只支持上传图片文件（JPG、PNG、GIF、WEBP）")
	}

	// 检查是否为支持的格式（排除 BMP 等其他格式）
	ext := getExtensionFromMimeType(mimeType)
	if ext == "" {
		return "", fmt.Errorf("只支持上传图片文件（JPG、PNG、GIF、WEBP）")
	}

	// 如果没有提供文件名，使用默认名称
	if filename == "" {
		filename = fmt.Sprintf("image_%d%s", time.Now().Unix(), ext)
	} else {
		// 确保文件名有正确的扩展名
		fileExt := strings.ToLower(filepath.Ext(filename))
		if fileExt == "" {
			// 如果文件名没有扩展名，则添加正确的扩展名
			filename = filename + ext
		} else if fileExt != ext {
			// 如果扩展名不匹配，则替换为正确的扩展名
			filenameWithoutExt := strings.TrimSuffix(filename, fileExt)
			filename = filenameWithoutExt + ext
		}
		// 如果扩展名已经匹配，保持原样
	}

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段，并手动设置 Content-Type
	// 使用 CreatePart 而不是 CreateFormFile，以便完全控制 Content-Type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("创建表单文件字段失败: %v", err)
	}

	// 写入文件数据
	_, err = part.Write(data)
	if err != nil {
		return "", fmt.Errorf("写入文件数据失败: %v", err)
	}

	// 关闭 writer（这会写入结束边界）
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭 multipart writer 失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", u.apiURL, body)
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("上传失败，HTTP状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析 JSON 响应
	var uploadResp UploadResponse
	err = json.Unmarshal(respBody, &uploadResp)
	if err != nil {
		return "", fmt.Errorf("解析响应 JSON 失败: %v, 响应内容: %s", err, string(respBody))
	}

	// 检查业务状态码
	if uploadResp.Code != http.StatusOK {
		return "", fmt.Errorf("上传失败，错误码: %d, 消息: %s", uploadResp.Code, uploadResp.Message)
	}

	// 检查 success 字段
	if !uploadResp.Success {
		return "", fmt.Errorf("上传失败: %s", uploadResp.Message)
	}

	// 检查返回的 URL 是否为空
	if uploadResp.Data.URL == "" {
		return "", fmt.Errorf("服务器返回的 URL 为空")
	}

	return uploadResp.Data.URL, nil
}

// UploadImageWithContext 带上下文的图片上传（支持取消）
func (u *HTTPUploader) UploadImageWithContext(ctx context.Context, data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("图片数据为空")
	}

	// 验证图片格式（只支持 JPG、PNG、GIF、WEBP）
	isImage, mimeType := img.IsImage(data)
	if !isImage {
		return "", fmt.Errorf("只支持上传图片文件（JPG、PNG、GIF、WEBP）")
	}

	// 检查是否为支持的格式（排除 BMP 等其他格式）
	ext := getExtensionFromMimeType(mimeType)
	if ext == "" {
		return "", fmt.Errorf("只支持上传图片文件（JPG、PNG、GIF、WEBP）")
	}

	// 处理文件名：确保有正确的扩展名
	originalFilename := filename
	if filename == "" {
		filename = fmt.Sprintf("image_%d%s", time.Now().Unix(), ext)
	} else {
		// 获取文件名的扩展名（小写）
		fileExt := strings.ToLower(filepath.Ext(filename))
		if fileExt == "" {
			// 如果文件名没有扩展名，则添加正确的扩展名
			filename = filename + ext
		} else {
			// 如果文件名已有扩展名，检查是否匹配检测到的图片类型
			// 如果不匹配，替换为正确的扩展名（确保服务器能正确识别）
			if fileExt != ext {
				filenameWithoutExt := strings.TrimSuffix(filename, fileExt)
				filename = filenameWithoutExt + ext
			}
			// 如果扩展名匹配，保持原样
		}
	}

	// 调试信息：记录文件名变化
	if originalFilename != filename {
		// 文件名被修改了，这在调试时很有用
	}

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段，并手动设置 Content-Type
	// 使用 CreatePart 而不是 CreateFormFile，以便完全控制 Content-Type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("创建表单文件字段失败: %v", err)
	}

	_, err = part.Write(data)
	if err != nil {
		return "", fmt.Errorf("写入文件数据失败: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭 multipart writer 失败: %v", err)
	}

	// 使用 context 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", u.apiURL, body)
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := u.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("上传失败，HTTP状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var uploadResp UploadResponse
	err = json.Unmarshal(respBody, &uploadResp)
	if err != nil {
		return "", fmt.Errorf("解析响应 JSON 失败: %v, 响应内容: %s", err, string(respBody))
	}

	if uploadResp.Code != http.StatusOK {
		return "", fmt.Errorf("上传失败，错误码: %d, 消息: %s", uploadResp.Code, uploadResp.Message)
	}

	if !uploadResp.Success {
		return "", fmt.Errorf("上传失败: %s", uploadResp.Message)
	}

	if uploadResp.Data.URL == "" {
		return "", fmt.Errorf("服务器返回的 URL 为空")
	}

	return uploadResp.Data.URL, nil
}
