export type TemplateKind =
  | 'question'
  | 'memo'
  | 'editorial'
  | 'lavender'
  | 'notebook'
  | 'stamp'
  | 'grid'
  | 'night'
  | 'burst'
  | 'paper'
  | 'bubble'
  | 'mono'

export interface TextBoxSpec {
  x: number
  y: number
  width: number
  height: number
  minFontSize: number
  maxFontSize: number
  maxLines: number
  lineHeight: number
  align: 'left' | 'center' | 'right'
}

export interface CardTemplate {
  id: string
  name: string
  kind: TemplateKind
  background: string
  foreground: string
  accent: string
  muted: string
  fontFamily: string
  fontWeight: number
  textBox: TextBoxSpec
  radius: number
  frame?: string
}

export interface TextLayout {
  lines: string[]
  fontSize: number
  lineHeightPx: number
  totalHeight: number
}

export interface HighlightRect {
  x: number
  y: number
  width: number
  height: number
  rotation: number
}
