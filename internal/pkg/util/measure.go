package util

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// MeasureArticleLength 测量文章长度并返回字符数
// 参数:
//
//	article - 需要测量长度的文章字符串
//
// 返回值:
//
//	int64 - 返回文章的字符数长度
func MeasureArticleLength(article string) int64 {
	// 使用 utf8.RuneCountInString() 来获取准确的字符数
	// 这对于包含多字节字符（如中文）的字符串特别重要
	return int64(utf8.RuneCountInString(article))
}

// 专门处理Markdown的摘要生成
func GenerateBriefFromMarkdown(markdown string) string {
	// 1. 移除代码块
	reCodeBlock := regexp.MustCompile("(?s)```.*?```")
	text := reCodeBlock.ReplaceAllString(markdown, "")

	// 2. 移除行内代码
	reInlineCode := regexp.MustCompile("`.*?`")
	text = reInlineCode.ReplaceAllString(text, "")

	// 3. 移除图片和链接
	reImage := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	text = reImage.ReplaceAllString(text, "")

	reLink := regexp.MustCompile(`\[.*?\]\(.*?\)`)
	text = reLink.ReplaceAllString(text, "$1") // 保留链接文本

	// 4. 移除标题标记
	reHeaders := regexp.MustCompile(`^#{1,6}\s+`)
	lines := strings.Split(text, "\n")

	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = reHeaders.ReplaceAllString(line, "")
		cleanLines = append(cleanLines, line)
	}

	// 5. 生成摘要
	cleanText := strings.Join(cleanLines, " ")
	return truncateToLength(cleanText, 200)
}

func truncateToLength(cleanText string, maxLength int) string {
	runes := []rune(cleanText)
	if len(runes) <= maxLength {
		return string(runes)
	}
	brief := string(runes[:maxLength])
	if lastSpace := strings.LastIndex(brief, " "); lastSpace > maxLength-20 {
		brief = brief[:lastSpace]
	}
	return brief + "..."
}
