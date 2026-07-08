package cardengine

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/rivo/uniseg"
	"golang.org/x/image/font"
)

const (
	DefaultEmptyText = "输入一个标题"
	Ellipsis         = "…"
)

var (
	BreakAfter   = makeSet("，", "。", "！", "？", "、", "；", "：", ",", ".", "!", "?", ";", ":", " ")
	BadLineStart = makeSet("，", "。", "！", "？", "、", "；", "：", ",", ".", "!", "?", ";", ":", ")", "）", "]", "】", "》", "」", "』")
	BadLineEnd   = makeSet("(", "（", "[", "【", "《", "「", "『")
	WidthRatios  = []float64{1, 0.97, 0.94, 0.91, 0.88, 0.85, 0.82, 0.79, 0.76, 0.73}
)

var spaceRegex = regexp.MustCompile(`[ \t]+`)

// --- 1. 数据结构定义 ---

type TextBox struct {
	Width         float64
	Height        float64
	X             float64
	Y             float64
	MaxLines      int
	MinFontSize   float64
	MaxFontSize   float64
	LineHeight    float64
	Align         string // "left", "center", "right"
	LetterSpacing float64
}

type CardTemplate struct {
	TextBox TextBox
}

type TextLayout struct {
	Lines        []string
	FontSize     float64
	LineHeightPx float64
	TotalHeight  float64
}

type HighlightRect struct {
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Rotation float64
}

type LayoutLineValue struct {
	Text  string
	Width float64
}

type LayoutCandidate struct {
	Lines     []LayoutLineValue
	WrapWidth float64
	Score     float64
}

// --- 2. 文本测量器 (抽象隔离层) ---

// Measurer 接口允许你注入具体的字体测量逻辑
type Measurer interface {
	MeasureText(text string, fontSize float64, letterSpacing float64) float64
}

// FontMeasurer 是基于 x/image/font 的具体实现，带有 LRU/并发缓存
type FontMeasurer struct {
	// 实际应用中，这里应该持有一个根据 fontSize 获取 font.Face 的工厂或缓存
	// 这里为了引擎完整性，抽象为一个获取 Face 的函数
	FaceProvider func(fontSize float64) font.Face
	cache        sync.Map // 简单的并发缓存: map[string]float64
}

func (m *FontMeasurer) MeasureText(text string, fontSize float64, letterSpacing float64) float64 {
	if text == "" {
		return 0
	}
	cacheKey := text + "|" + floatToStr(fontSize) + "|" + floatToStr(letterSpacing)
	if val, ok := m.cache.Load(cacheKey); ok {
		return val.(float64)
	}

	face := m.FaceProvider(fontSize)
	advance := font.MeasureString(face, text)
	// font.MeasureString 返回 26.6 定点数，右移 6 位转为 float64 像素
	width := float64(advance) / 64.0

	// 加上 LetterSpacing
	graphemeCount := GraphemeLength(text)
	if graphemeCount > 1 {
		width += float64(graphemeCount-1) * letterSpacing
	}

	m.cache.Store(cacheKey, width)
	return width
}

func floatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// --- 3. 字符串与字形工具 ---

func makeSet(items ...string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, item := range items {
		m[item] = struct{}{}
	}
	return m
}

func NormalizeText(value string) string {
	s := strings.ReplaceAll(value, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = spaceRegex.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func GetGraphemes(value string) []string {
	var graphemes []string
	gr := uniseg.NewGraphemes(value)
	for gr.Next() {
		graphemes = append(graphemes, gr.Str())
	}
	return graphemes
}

func GraphemeLength(value string) int {
	return uniseg.GraphemeClusterCount(value)
}

func FirstGrapheme(value string) string {
	gr := uniseg.NewGraphemes(value)
	if gr.Next() {
		return gr.Str()
	}
	return ""
}

func LastGrapheme(value string) string {
	gr := uniseg.NewGraphemes(value)
	last := ""
	for gr.Next() {
		last = gr.Str()
	}
	return last
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Min(maximum, math.Max(minimum, value))
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// --- 4. 核心排版与打分算法 ---

// LayoutWithLines 模拟 Pretext 的贪心断行逻辑
func LayoutWithLines(text string, wrapWidth float64, fontSize float64, ls float64, m Measurer) []LayoutLineValue {
	graphemes := GetGraphemes(text)
	var lines []LayoutLineValue
	var currentLine string
	var currentWidth float64

	for _, g := range graphemes {
		w := m.MeasureText(g, fontSize, 0)
		spacing := 0.0
		if currentLine != "" {
			spacing = ls
		}

		if currentWidth+spacing+w > wrapWidth && currentWidth > 0 {
			lines = append(lines, LayoutLineValue{Text: currentLine, Width: currentWidth})
			currentLine = g
			currentWidth = w
		} else {
			currentLine += g
			currentWidth += spacing + w
		}
	}
	if currentLine != "" {
		lines = append(lines, LayoutLineValue{Text: currentLine, Width: currentWidth})
	}
	return lines
}

func createCandidateWidths(text string, fontSize float64, ls float64, template CardTemplate, m Measurer) []float64 {
	boxWidth := template.TextBox.Width
	minWidth := boxWidth * 0.68
	widthSet := make(map[float64]struct{})

	for _, ratio := range WidthRatios {
		widthSet[round1(boxWidth*ratio)] = struct{}{}
	}

	naturalWidth := m.MeasureText(text, fontSize, ls)
	maxTargetLines := float64(template.TextBox.MaxLines)
	if maxTargetLines > 6 {
		maxTargetLines = 6
	}

	for targetLines := 2.0; targetLines <= maxTargetLines; targetLines++ {
		for _, factor := range []float64{0.94, 1, 1.06} {
			width := clamp((naturalWidth/targetLines)*factor, minWidth, boxWidth)
			widthSet[round1(width)] = struct{}{}
		}
	}

	// map 转 slice 并从大到小排序
	var widths []float64
	for w := range widthSet {
		widths = append(widths, w)
	}
	// Sort descending
	for i := 0; i < len(widths)-1; i++ {
		for j := i + 1; j < len(widths); j++ {
			if widths[i] < widths[j] {
				widths[i], widths[j] = widths[j], widths[i]
			}
		}
	}
	return widths
}

func getIdealLineCount(text string, maxLines int) int {
	length := GraphemeLength(strings.ReplaceAll(text, " ", ""))
	if length <= 9 {
		return 2
	}
	if length <= 18 {
		return 3
	}
	if length <= 32 {
		return 4
	}
	if 5 < maxLines {
		return 5
	}
	return maxLines
}

func scoreLinePunctuation(line string, isLast boolean) float64 {
	first := FirstGrapheme(line)
	last := LastGrapheme(line)
	score := 0.0

	if _, ok := BadLineStart[first]; ok {
		score += 1800
	}
	if _, ok := BadLineEnd[last]; ok {
		score += 1100
	}
	if !isLast {
		if _, ok := BreakAfter[last]; ok {
			score -= 95
		}
	}
	return score
}

type boolean bool

func scoreLayout(lines []LayoutLineValue, wrapWidth, boxWidth float64, maxLines int, text string) float64 {
	if len(lines) == 0 {
		return math.Inf(1)
	}

	var fillRatios []float64
	sumFill := 0.0
	for _, line := range lines {
		fill := clamp(line.Width/wrapWidth, 0, 1.5)
		fillRatios = append(fillRatios, fill)
		sumFill += fill
	}

	avgFill := sumFill / float64(len(fillRatios))
	fillVarianceSum := 0.0
	for _, fill := range fillRatios {
		fillVarianceSum += math.Pow(fill-avgFill, 2)
	}
	fillVariance := fillVarianceSum / float64(len(fillRatios))

	score := 0.0
	score += fillVariance * 1450
	score += math.Pow(math.Max(0, 0.67-avgFill), 2) * 1200
	score += math.Pow(1-wrapWidth/boxWidth, 2) * 360

	idealLineCount := getIdealLineCount(text, maxLines)
	score += math.Pow(float64(len(lines)-idealLineCount), 2) * 34

	for i, line := range lines {
		isLast := i == len(lines)-1
		length := GraphemeLength(line.Text)
		fill := fillRatios[i]

		score += scoreLinePunctuation(line.Text, boolean(isLast))

		if !isLast && fill < 0.54 {
			score += math.Pow(0.54-fill, 2) * 1700
		}

		if length <= 2 {
			if isLast {
				score += 780
			} else {
				score += 1250
			}
		} else if length == 3 && isLast {
			score += 160
		}
	}

	lastFill := fillRatios[len(fillRatios)-1]
	if lastFill < 0.38 {
		score += math.Pow(0.38-lastFill, 2) * 2500
	}

	if len(lines) >= 2 {
		prevFill := fillRatios[len(fillRatios)-2]
		if lastFill > prevFill+0.28 {
			score += math.Pow(lastFill-prevFill, 2) * 520
		}
	}

	return score
}

func findBestCandidateAtSize(rawText string, fontSize float64, template CardTemplate, m Measurer) *LayoutCandidate {
	text := NormalizeText(rawText)
	if text == "" {
		text = DefaultEmptyText
	}

	box := template.TextBox
	lineHeightPx := fontSize * box.LineHeight
	ls := box.LetterSpacing

	candidateWidths := createCandidateWidths(text, fontSize, ls, template, m)

	var best *LayoutCandidate

	for _, wrapWidth := range candidateWidths {
		lines := LayoutWithLines(text, wrapWidth, fontSize, ls, m)

		if len(lines) == 0 || len(lines) > box.MaxLines || float64(len(lines))*lineHeightPx > box.Height {
			continue
		}

		score := scoreLayout(lines, wrapWidth, box.Width, box.MaxLines, text)

		if best == nil || score < best.Score {
			best = &LayoutCandidate{
				Lines:     lines,
				WrapWidth: wrapWidth,
				Score:     score,
			}
		}
	}

	return best
}

func fitWithEllipsis(value string, maxWidth float64, fontSize float64, template CardTemplate, m Measurer) string {
	graphemes := GetGraphemes(value)
	ls := template.TextBox.LetterSpacing

	if m.MeasureText(value+Ellipsis, fontSize, ls) <= maxWidth {
		return value + Ellipsis
	}

	low := 0
	high := len(graphemes)
	best := Ellipsis

	for low <= high {
		middle := (low + high) / 2
		candidate := strings.Join(graphemes[:middle], "") + Ellipsis

		if m.MeasureText(candidate, fontSize, ls) <= maxWidth {
			best = candidate
			low = middle + 1
		} else {
			high = middle - 1
		}
	}
	return best
}

func createFallbackCandidate(rawText string, fontSize float64, template CardTemplate, m Measurer) *LayoutCandidate {
	text := NormalizeText(rawText)
	if text == "" {
		text = DefaultEmptyText
	}

	box := template.TextBox
	ls := box.LetterSpacing
	allLines := LayoutWithLines(text, box.Width, fontSize, ls, m)

	visibleCount := box.MaxLines
	if len(allLines) < visibleCount {
		visibleCount = len(allLines)
	}

	var visibleLines []LayoutLineValue
	if visibleCount == 0 {
		visibleLines = append(visibleLines, LayoutLineValue{
			Text:  DefaultEmptyText,
			Width: m.MeasureText(DefaultEmptyText, fontSize, ls),
		})
	} else {
		visibleLines = make([]LayoutLineValue, visibleCount)
		copy(visibleLines, allLines[:visibleCount])
	}

	if len(allLines) > box.MaxLines {
		lastIndex := len(visibleLines) - 1
		lastText := visibleLines[lastIndex].Text
		truncated := fitWithEllipsis(lastText, box.Width, fontSize, template, m)
		visibleLines[lastIndex] = LayoutLineValue{
			Text:  truncated,
			Width: m.MeasureText(truncated, fontSize, ls),
		}
	}

	return &LayoutCandidate{
		Lines:     visibleLines,
		WrapWidth: box.Width,
		Score:     math.Inf(1),
	}
}

// BuildTextLayout 核心入口：构建标题卡片的最终文字布局
func BuildTextLayout(rawText string, template CardTemplate, m Measurer) TextLayout {
	text := NormalizeText(rawText)
	if text == "" {
		text = DefaultEmptyText
	}

	box := template.TextBox

	low := box.MinFontSize
	high := box.MaxFontSize
	maxFittingSize := box.MinFontSize

	maxCandidate := findBestCandidateAtSize(text, box.MinFontSize, template, m)

	// 二分寻找能放下的最大字号
	for high-low > 0.8 {
		fontSize := (low + high) / 2
		candidate := findBestCandidateAtSize(text, fontSize, template, m)

		if candidate != nil {
			maxFittingSize = fontSize
			maxCandidate = candidate
			low = fontSize
		} else {
			high = fontSize
		}
	}

	// 比较附近的缩小字号
	nearbySizes := []float64{
		maxFittingSize,
		maxFittingSize * 0.97,
		maxFittingSize * 0.94,
		maxFittingSize * 0.91,
		maxFittingSize * 0.88,
	}

	selectedSize := maxFittingSize
	selectedCandidate := maxCandidate
	selectedScore := math.Inf(1)
	if maxCandidate != nil {
		selectedScore = maxCandidate.Score
	}

	for _, rawSize := range nearbySizes {
		fontSize := clamp(rawSize, box.MinFontSize, box.MaxFontSize)
		candidate := findBestCandidateAtSize(text, fontSize, template, m)

		if candidate == nil {
			continue
		}

		shrinkRatio := 0.0
		if maxFittingSize > 0 {
			shrinkRatio = (maxFittingSize - fontSize) / maxFittingSize
		}

		adjustedScore := candidate.Score + shrinkRatio*620

		if selectedCandidate == nil || adjustedScore < selectedScore {
			selectedSize = fontSize
			selectedCandidate = candidate
			selectedScore = adjustedScore
		}
	}

	if selectedCandidate == nil {
		selectedSize = box.MinFontSize
		selectedCandidate = createFallbackCandidate(text, selectedSize, template, m)
	}

	fontSize := math.Floor(selectedSize*10) / 10
	lineHeightPx := fontSize * box.LineHeight

	var finalLines []string
	for _, l := range selectedCandidate.Lines {
		finalLines = append(finalLines, l.Text)
	}

	return TextLayout{
		Lines:        finalLines,
		FontSize:     fontSize,
		LineHeightPx: lineHeightPx,
		TotalHeight:  float64(len(finalLines)) * lineHeightPx,
	}
}

// FindHighlightRects 计算高亮区域坐标
func FindHighlightRects(layout TextLayout, keyword string, template CardTemplate, m Measurer) []HighlightRect {
	normalized := strings.TrimSpace(keyword)
	if normalized == "" {
		return nil
	}

	var rectangles []HighlightRect
	box := template.TextBox
	ls := box.LetterSpacing

	for lineIndex, line := range layout.Lines {
		from := 0
		for from < len(line) {
			index := strings.Index(line[from:], normalized)
			if index < 0 {
				break
			}
			index += from

			before := line[:index]
			matched := line[index : index+len(normalized)]

			beforeWidth := m.MeasureText(before, layout.FontSize, ls)
			matchedWidth := m.MeasureText(matched, layout.FontSize, ls)
			lineWidth := m.MeasureText(line, layout.FontSize, ls)

			lineStart := box.X
			switch box.Align {
			case "center":
				lineStart = box.X + (box.Width-lineWidth)/2
			case "right":
				lineStart = box.X + box.Width - lineWidth
			}

			rectangles = append(rectangles, HighlightRect{
				X:        lineStart + beforeWidth - 8,
				Y:        box.Y + float64(lineIndex)*layout.LineHeightPx + layout.FontSize*0.55,
				Width:    matchedWidth + 16,
				Height:   layout.FontSize * 0.5,
				Rotation: float64(((lineIndex+index)%3)-1) * 1.8,
			})

			from = index + len(normalized)
		}
	}
	return rectangles
}

// --- 5. 辅助功能 ---

func InferKeyword(text string) string {
	trimmed := NormalizeText(text)

	// 匹配引号内的内容
	re := regexp.MustCompile(`[“「『《](.{2,8}?)[”」』》]`)
	if match := re.FindStringSubmatch(trimmed); len(match) > 1 {
		return match[1]
	}

	// 简单的标点分割兜底策略
	splitRe := regexp.MustCompile(`[，。！？、；：,\s!?]+`)
	parts := splitRe.Split(trimmed, -1)

	var fragments []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if GraphemeLength(p) >= 2 {
			fragments = append(fragments, p)
		}
	}

	if len(fragments) == 0 {
		return ""
	}

	last := fragments[len(fragments)-1]
	if GraphemeLength(last) <= 6 {
		return last
	}

	gr := GetGraphemes(last)
	if len(gr) >= 4 {
		return strings.Join(gr[len(gr)-4:], "")
	}
	return last
}

func HashText(value string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(value); i++ {
		hash ^= uint32(value[i])
		hash *= 16777619
	}
	return hash
}

func SeededValue(seed float64, index float64) float64 {
	val := math.Sin(seed*12.9898+index*78.233) * 43758.5453
	return val - math.Floor(val)
}
