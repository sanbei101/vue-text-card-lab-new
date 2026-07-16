package cardengine

import (
	"hash/fnv"
	"math"
	"regexp"
	"slices"
	"strings"

	"github.com/phuslu/lru"
	"github.com/rivo/uniseg"
	"golang.org/x/image/font"
)

const (
	ScorePenaltyBadLineStart = 1800.0 // 行首出现标点符号的惩罚
	ScorePenaltyBadLineEnd   = 1100.0 // 行尾出现左括号等符号的惩罚
	ScoreBonusBreakAfter     = 95.0   // 在正常句读后断行的奖励
	ScoreWeightFillVariance  = 1450.0 // 整体填满率方差的权重
	ScoreWeightAvgFill       = 1200.0 // 平均填满率过低的惩罚权重
	ScoreWeightWrapWidth     = 360.0  // 实际排版宽度偏离目标宽度的惩罚
	ScoreWeightLineCount     = 34.0   // 偏离理想行数的惩罚

	ScorePenaltyEarlyLineFill = 1700.0 // 非最后一行填满率极低的惩罚
	ScorePenaltyShortLastLine = 780.0  // 最后一行字数极少的惩罚
	ScorePenaltyShortMidLine  = 1250.0 // 中间行字数极少的惩罚
	ScorePenaltyThreeCharLast = 160.0  // 最后一行仅3个字的轻微惩罚
	ScorePenaltyLastLineFill  = 2500.0 // 最后一行极度空旷的惩罚
	ScorePenaltyOrphanWord    = 520.0  // 导致最后一行看起来像"孤儿词"的差异惩罚
)

const (
	DefaultEmptyText = "输入一个标题"
	Ellipsis         = "…"
)

var (
	BreakAfter   = makeSet("，", "。", "！", "？", "、", "；", "：", ",", ".", "!", "?", ";", ":", " ")
	BadLineStart = makeSet(
		"，",
		"。",
		"！",
		"？",
		"、",
		"；",
		"：",
		",",
		".",
		"!",
		"?",
		";",
		":",
		")",
		"）",
		"]",
		"】",
		"》",
		"」",
		"』",
	)
	BadLineEnd  = makeSet("(", "（", "[", "【", "《", "「", "『")
	WidthRatios = []float64{1, 0.97, 0.94, 0.91, 0.88, 0.85, 0.82, 0.79, 0.76, 0.73}
)

var (
	spaceRegex = regexp.MustCompile(`[ \t]+`)
	quoteRegex = regexp.MustCompile(`[“「『《](.{2,8}?)[”」』》]`)
	splitRegex = regexp.MustCompile(`[，。！？、；：,\s!?]+`)
)

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
	Text          string
	Width         float64
	GraphemeCount int    // 提前记录字素数量
	First         string // 提前记录行首字素
	Last          string // 提前记录行尾字素
}

type LayoutCandidate struct {
	Lines     []LayoutLineValue
	WrapWidth float64
	Score     float64
}
type measureKey struct {
	text          string
	fontSize      float64
	letterSpacing float64
}

// 文本测量器

type Measurer interface {
	MeasureText(text string, fontSize, letterSpacing float64) float64
}

type FontMeasurer struct {
	FaceProvider func(fontSize float64) font.Face
	cache        *lru.LRUCache[measureKey, float64]
}

func NewFontMeasurer(faceProvider func(fontSize float64) font.Face, cacheSize int) *FontMeasurer {
	return &FontMeasurer{
		FaceProvider: faceProvider,
		cache:        lru.NewLRUCache[measureKey, float64](cacheSize),
	}
}

func (m *FontMeasurer) MeasureText(text string, fontSize, letterSpacing float64) float64 {
	cacheKey := measureKey{
		text:          text,
		fontSize:      fontSize,
		letterSpacing: letterSpacing,
	}
	if width, ok := m.cache.Get(cacheKey); ok {
		return width
	}

	face := m.FaceProvider(fontSize)
	advance := font.MeasureString(face, text)
	width := float64(advance) / 64.0

	graphemeCount := GraphemeLength(text)
	if graphemeCount > 1 {
		width += float64(graphemeCount-1) * letterSpacing
	}

	m.cache.Set(cacheKey, width)
	return width
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
	if value == "" {
		return ""
	}
	gr, _, _, _ := uniseg.StepString(value, -1)
	return gr
}

func LastGrapheme(value string) string {
	if value == "" {
		return ""
	}
	var last string
	state := -1
	for len(value) > 0 {
		last, value, _, state = uniseg.StepString(value, state)
	}
	return last
}

func clamp(value, minimum, maximum float64) float64 {
	return max(minimum, min(maximum, value))
}

func roundTo1Decimal(v float64) float64 {
	return math.Round(v*10) / 10
}

// LayoutWithLines 模拟 Pretext 的贪心断行逻辑
func LayoutWithLines(graphemes []string, wrapWidth, fontSize, ls float64, m Measurer) []LayoutLineValue {
	lines := make([]LayoutLineValue, 0, 6)

	start := 0
	currentWidth := 0.0

	for i, g := range graphemes {
		w := m.MeasureText(g, fontSize, 0)
		spacing := 0.0
		if i > start {
			spacing = ls
		}

		if currentWidth+spacing+w > wrapWidth && currentWidth > 0 {
			lines = append(lines, LayoutLineValue{
				Text:          strings.Join(graphemes[start:i], ""),
				Width:         currentWidth,
				GraphemeCount: i - start,
				First:         graphemes[start],
				Last:          graphemes[i-1],
			})
			start = i
			currentWidth = w
		} else {
			currentWidth += spacing + w
		}
	}
	if start < len(graphemes) {
		lines = append(lines, LayoutLineValue{
			Text:          strings.Join(graphemes[start:], ""),
			Width:         currentWidth,
			GraphemeCount: len(graphemes) - start,
			First:         graphemes[start],
			Last:          graphemes[len(graphemes)-1],
		})
	}
	return lines
}

func createCandidateWidths(text string, fontSize, letterSpace float64, template *CardTemplate, m Measurer) []float64 {
	boxWidth := template.TextBox.Width
	minWidth := boxWidth * 0.68
	widths := make([]float64, 0, 32)

	for _, ratio := range WidthRatios {
		widths = append(widths, roundTo1Decimal(boxWidth*ratio))
	}

	naturalWidth := m.MeasureText(text, fontSize, letterSpace)
	maxTargetLines := float64(template.TextBox.MaxLines)
	if maxTargetLines > 6 {
		maxTargetLines = 6
	}

	for targetLines := 2.0; targetLines <= maxTargetLines; targetLines++ {
		for _, factor := range []float64{0.94, 1, 1.06} {
			width := clamp((naturalWidth/targetLines)*factor, minWidth, boxWidth)
			widths = append(widths, roundTo1Decimal(width))
		}
	}

	slices.SortFunc(widths, func(a, b float64) int {
		if a < b {
			return 1
		}
		if a > b {
			return -1
		}
		return 0
	})
	return slices.Compact(widths)
}

func getIdealLineCount(text string, maxLines int) int {
	length := GraphemeLength(text) - strings.Count(text, " ")

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

func scoreLinePunctuation(line LayoutLineValue, isLast bool) float64 {
	score := 0.0
	if _, ok := BadLineStart[line.First]; ok {
		score += ScorePenaltyBadLineStart
	}
	if _, ok := BadLineEnd[line.Last]; ok {
		score += ScorePenaltyBadLineEnd
	}
	if !isLast {
		if _, ok := BreakAfter[line.Last]; ok {
			score -= ScoreBonusBreakAfter
		}
	}
	return score
}

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
		fillVarianceSum += (fill - avgFill) * (fill - avgFill)
	}
	fillVariance := fillVarianceSum / float64(len(fillRatios))

	score := 0.0
	score += fillVariance * ScoreWeightFillVariance
	score += math.Pow(max(0, 0.67-avgFill), 2) * ScoreWeightAvgFill
	score += (1 - wrapWidth/boxWidth) * (1 - wrapWidth/boxWidth) * ScoreWeightWrapWidth

	idealLineCount := getIdealLineCount(text, maxLines)
	score += math.Pow(float64(len(lines)-idealLineCount), 2) * ScoreWeightLineCount

	for i, line := range lines {
		isLast := i == len(lines)-1
		length := line.GraphemeCount
		fill := fillRatios[i]

		score += scoreLinePunctuation(line, isLast)

		if !isLast && fill < 0.54 {
			score += (0.54 - fill) * (0.54 - fill) * 1700
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
		score += (0.38 - lastFill) * (0.38 - lastFill) * 2500
	}

	if len(lines) >= 2 {
		prevFill := fillRatios[len(fillRatios)-2]
		if lastFill > prevFill+0.28 {
			score += (lastFill - prevFill) * (lastFill - prevFill) * 520
		}
	}

	return score
}

func findBestCandidateAtSize(rawText string, fontSize float64, template *CardTemplate, m Measurer) *LayoutCandidate {
	text := NormalizeText(rawText)
	if text == "" {
		text = DefaultEmptyText
	}

	box := template.TextBox
	lineHeightPx := fontSize * box.LineHeight
	ls := box.LetterSpacing

	candidateWidths := createCandidateWidths(text, fontSize, ls, template, m)

	graphemes := GetGraphemes(text)
	var best *LayoutCandidate

	for _, wrapWidth := range candidateWidths {
		lines := LayoutWithLines(graphemes, wrapWidth, fontSize, ls, m)

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

func fitWithEllipsis(value string, maxWidth, fontSize float64, template *CardTemplate, m Measurer) string {
	ls := template.TextBox.LetterSpacing

	if m.MeasureText(value+Ellipsis, fontSize, ls) <= maxWidth {
		return value + Ellipsis
	}

	var byteIndices []int
	state := -1
	rest := value
	currentByte := 0
	for len(rest) > 0 {
		var cluster string
		cluster, rest, _, state = uniseg.StepString(rest, state)
		currentByte += len(cluster)
		byteIndices = append(byteIndices, currentByte)
	}

	low := 0
	high := len(byteIndices)
	best := Ellipsis

	for low <= high {
		middle := (low + high) / 2
		var candidate string

		if middle == 0 {
			candidate = Ellipsis
		} else {
			candidate = value[:byteIndices[middle-1]] + Ellipsis
		}

		if m.MeasureText(candidate, fontSize, ls) <= maxWidth {
			best = candidate
			low = middle + 1
		} else {
			high = middle - 1
		}
	}

	return best
}

func createFallbackCandidate(rawText string, fontSize float64, template *CardTemplate, m Measurer) *LayoutCandidate {
	text := NormalizeText(rawText)
	if text == "" {
		text = DefaultEmptyText
	}

	box := template.TextBox
	ls := box.LetterSpacing
	graphemes := GetGraphemes(text)
	allLines := LayoutWithLines(graphemes, box.Width, fontSize, ls, m)

	visibleCount := min(len(allLines), box.MaxLines)

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

// BuildTextLayout 核心入口:构建标题卡片的最终文字布局
func BuildTextLayout(rawText string, template *CardTemplate, m Measurer) TextLayout {
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
func FindHighlightRects(layout TextLayout, keyword string, template *CardTemplate, m Measurer) []HighlightRect {
	normalized := strings.TrimSpace(keyword)
	if normalized == "" {
		return nil
	}

	var rectangles []HighlightRect
	box := template.TextBox
	ls := box.LetterSpacing
	matchedWidth := m.MeasureText(normalized, layout.FontSize, ls)

	for lineIndex, line := range layout.Lines {
		from := 0
		for from < len(line) {
			index := strings.Index(line[from:], normalized)
			if index < 0 {
				break
			}
			index += from

			before := line[:index]

			beforeWidth := m.MeasureText(before, layout.FontSize, ls)
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
	if match := quoteRegex.FindStringSubmatch(trimmed); len(match) > 1 {
		return match[1]
	}

	// 简单的标点分割兜底策略
	parts := splitRegex.Split(trimmed, -1)

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
	h := fnv.New32a()
	h.Write([]byte(value))
	return h.Sum32()
}

func SeededValue(seed, index float64) float64 {
	val := math.Sin(seed*12.9898+index*78.233) * 43758.5453
	return val - math.Floor(val)
}
