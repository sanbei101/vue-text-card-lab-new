# 小蓝书卡片引擎 · Blue Card Engine

> 把文字变成小红书类型的封面图

基于 **Go + ConnectRPC** 的服务,给定一段文字,自动套用 12 套精心设计的模板,生成高质量 SVG / WebP 卡片图,并上传至云端存储。


### TemplateList - 获取模板列表

**请求:**

```jsonc
// POST /cardengine.v1.CardEngineService/TemplateList
// Content-Type: application/json
```

**响应:**

```jsonc
{
  "templates": [
    {
      "id": "question-blue",
      "name": "清透问号",
      "kind": "question"
    },
    {
      "id": "yellow-memo",
      "name": "便签日记",
      "kind": "memo"
    }
    // ... 共 12 个模板
  ]
}
```

### 2. Cards - 生成卡片

根据标题生成所有模板的卡片图,自动转换为 WebP 并上传到 R2 存储。

**请求:**

```jsonc
// POST /cardengine.v1.CardEngineService/Cards
// Content-Type: application/json

{
  "title": "为什么我们总是怀念夏天?",   // 必填,最长 100 字符
  "keyword": "夏天"                      // 可选,关键词用于文字高亮
}
```

**响应:**

```jsonc
{
  "templates": [
    {
      "id": "question-blue",
      "name": "清透问号",
      "url": "https://image-bed.sanbei.codes/cards/question-blue/abc12345.webp"
    },
    {
      "id": "yellow-memo",
      "name": "便签日记",
      "url": "https://image-bed.sanbei.codes/cards/yellow-memo/def67890.webp"
    }
    // ... 每个模板一张卡片
  ]
}
```

#### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `title` | `string` | 卡片文字内容,最长 100 个字符(Unicode 计数) |
| `keyword` | `string` | 可选项。为空时自动从标题中提取关键词;指定后会在卡片中高亮匹配的文字 |
| `templates[].id` | `string` | 模板唯一 ID |
| `templates[].name` | `string` | 模板中文名 |
| `templates[].url` | `string` | 生成的 WebP 图片公开访问 URL |



## 作为库使用

如果不需要运行服务,只想在本地生成 SVG:

```go
import (
    "github.com/sanbei101/blue-card-engine/internal/render"
    "github.com/sanbei101/blue-card-engine/internal/templates"
)

// 加载字体库
lib, err := fonts.NewLibrary()
if err != nil { /* handle */ }

// 加载模板
reg, err := templates.Load()
if err != nil { /* handle */ }

// 获取某个模板
tpl, _ := reg.ByKind("mono")

// 生成 SVG
svgBytes, err := render.Card(tpl, "少即是多。", "设计", lib)
// svgBytes 可直接写入 .svg 文件,或传给 webp.SVGToWebP 转 WebP
```

## 项目结构

```
blue-card-engine/
├── cmd/server/                 # 服务入口
│   └── main.go                  # HTTP 服务器启动
├── internal/
│   ├── cardengine/              # 文字排版引擎(换行、字号自适应、高亮)
│   ├── fonts/                   # 字体库(5 款嵌入中文字体)
│   ├── render/                  # SVG 渲染 + R2 上传
│   │   ├── render.go            # Card() - SVG 生成核心
│   │   ├── textpath.go          # 文字 → SVG path
│   │   ├── oss.go               # R2 存储上传
│   │   ├── shared/              # 共享 SVG 模板组件
│   │   └── templates/           # 12 套 SVG 模板
│   ├── render/webp/             # SVG → WebP 转换(依赖 libvips)
│   ├── server/                  # ConnectRPC 请求处理
│   └── templates/               # 模板注册表 + 12 套模板配置
├── gen/                         # 自动生成的 Protobuf / Connect 代码
├── proto/                       # .proto 接口定义
├── testdata/                    # 测试产出的 SVG / PNG
└── .github/workflows/
    └── render_preview.yml       # 自动预览:生成所有模板 → PNG → 上传
```
