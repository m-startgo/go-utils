package mfile

import (
	"mime"
	"net/http"
	"strings"
)

// ContentToExtName 根据给定的 content-type 返回不带点的扩展名（例如 "png"）。
// 函数对输入做基本规范化（Trim + ToLower + 去掉参数），并使用内部映射作兜底。
func ContentToExtName(lType string) string {
	ct := strings.TrimSpace(strings.ToLower(lType))
	if ct == "" {
		return ""
	}
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = strings.TrimSpace(ct[:idx])
	}

	// 直接使用 map，比长 switch 更易维护
	m := map[string]string{
		"image/bmp":             "bmp",
		"image/gif":             "gif",
		"image/jpeg":            "jpg", // 优先返回最常见的 "jpg"
		"image/webp":            "webp",
		"image/png":             "png",
		"text/html":             "html",
		"text/plain":            "txt",
		"application/vnd.visio": "vsd",
		// 一些历史/常见 Office MIME
		"application/vnd.ms-powerpoint":                                             "ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"application/msword":                                                        "doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "docx",
		"application/vnd.ms-excel":                                                  "xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
		"application/csv":  "csv",
		"text/xml":         "xml",
		"video/mp4":        "mp4",
		"video/x-msvideo":  "avi",
		"video/quicktime":  "mov",
		"video/mpeg":       "mpeg",
		"video/x-ms-wmv":   "wm",
		"video/x-flv":      "flv",
		"video/x-matroska": "mkv",
	}

	if v, ok := m[ct]; ok {
		return v
	}
	return ""
}

// MimeToExt 返回不带点的扩展名（例如 "png"），若未知返回空字符串。
// 优先使用标准库 mime.ExtensionsByType 获取扩展名，失败后回退到 ContentToExtName 的映射。
func MimeToExt(ct string) string {
	ct = strings.TrimSpace(strings.ToLower(ct))
	if ct == "" {
		return ""
	}
	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = strings.TrimSpace(ct[:idx])
	}
	if exts, _ := mime.ExtensionsByType(ct); len(exts) > 0 {
		return strings.TrimPrefix(exts[0], ".")
	}
	// 回退到内部映射
	return ContentToExtName(ct)
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
