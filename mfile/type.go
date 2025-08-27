package mfile

import (
	"mime"
	"net/http"
	"strings"
)

func ContentToExtName(lType string) string {
	ext := ""
	switch lType {

	case "image/bmp":
		ext = "bmp"

	case "image/gif":
		ext = "gif"

	case "image/jpeg":
		ext = "jpeg"

	case "image/webp":
		ext = "webp"

	case "image/png":
		ext = "png"

	case "text/html":
		ext = "html"

	case "text/plain":
		ext = "txt"

	case "application/vnd.visio":
		ext = "vsd"

	case "application/vnd.ms-powerpoint":
		ext = "pptx"

	case "application/msword":
		ext = "docx"

	case "application/msexcel":
		ext = "xlsx"

	case "application/csv":
		ext = "csv"

	case "text/xml":
		ext = "xml"

	case "video/mp4":
		ext = "mp4"

	case "video/x-msvideo":
		ext = "avi"

	case "video/quicktime":
		ext = "mov"

	case "video/mpeg":
		ext = "mpeg"

	case "video/x-ms-wmv":
		ext = "wm"

	case "video/x-flv":
		ext = "flv"

	case "video/x-matroska":
		ext = "mkv"

	}

	if strings.Contains(lType, "text/html") {
		ext = "html"
	}

	return ext
}

// MimeToExt 返回不带点的扩展名（例如 "png"），若未知返回空字符串。
func MimeToExt(ct string) string {
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "" {
		return ""
	}
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = strings.TrimSpace(ct[:idx])
	}
	if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
		// exts 如 [".png"], 去掉点
		return strings.TrimPrefix(exts[0], ".")
	}
	// 兜底映射（只列出常见的）
	m := map[string]string{
		"image/jpeg":               "jpg",
		"image/png":                "png",
		"image/gif":                "gif",
		"text/plain":               "txt",
		"text/html":                "html",
		"application/pdf":          "pdf",
		"application/vnd.ms-excel": "xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": "xlsx",
		"application/msword": "doc",
		// ... 根据需要补充
	}
	if v, ok := m[ct]; ok {
		return v
	}
	return ""
}

// DetectMime 根据字节内容返回 MIME 类型，使用 http.DetectContentType
func DetectMime(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	// DetectContentType 只需要前512字节
	n := 512
	if len(content) < 512 {
		n = len(content)
	}
	return http.DetectContentType(content[:n])
}

// ExtByContent 根据内容推断合适的文件后缀（含点），若无法推断则返回空字符串
func ExtByContent(content []byte) string {
	mt := DetectMime(content)
	if mt == "" {
		return ""
	}
	// 去掉可能的 charset
	if idx := strings.Index(mt, ";"); idx != -1 {
		mt = strings.TrimSpace(mt[:idx])
	}
	exts, _ := mime.ExtensionsByType(mt)
	if len(exts) > 0 {
		return exts[0]
	}
	// 兜底：根据少数已知 mime 做映射
	switch mt {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "application/pdf":
		return ".pdf"
	case "text/plain":
		return ".txt"
	case "text/html":
		return ".html"
	}
	return ""
}
