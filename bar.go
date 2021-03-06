package termui

import "fmt"

type BarChart struct {
	Block
	BarColor   Attribute
	TextColor  Attribute
	NumColor   Attribute
	Data       []int
	DataLabels []string
	BarWidth   int
	BarGap     int
	labels     [][]rune
	dataNum    [][]rune
	numBar     int
	scale      float64
	max        int
}

func NewBarChart() *BarChart {
	bc := &BarChart{Block: *NewBlock()}
	bc.BarColor = ColorCyan
	bc.NumColor = ColorWhite
	bc.TextColor = ColorWhite
	bc.BarGap = 1
	bc.BarWidth = 3
	return bc
}

func (bc *BarChart) layout() {
	bc.numBar = bc.innerWidth / (bc.BarGap + bc.BarWidth)
	bc.labels = make([][]rune, bc.numBar)
	bc.dataNum = make([][]rune, len(bc.Data))

	for i := 0; i < bc.numBar && i < len(bc.DataLabels) && i < len(bc.Data); i++ {
		bc.labels[i] = trimStr2Runes(bc.DataLabels[i], bc.BarWidth)
		n := bc.Data[i]
		s := fmt.Sprint(n)
		bc.dataNum[i] = trimStr2Runes(s, bc.BarWidth)
	}

	bc.max = bc.Data[0] //  what if Data is nil?
	for i := 0; i < len(bc.Data); i++ {
		if bc.max < bc.Data[i] {
			bc.max = bc.Data[i]
		}
	}
	bc.scale = float64(bc.max) / float64(bc.innerHeight-1)
}

func (bc *BarChart) Buffer() []Point {
	ps := bc.Block.Buffer()
	bc.layout()

	for i := 0; i < bc.numBar && i < len(bc.Data) && i < len(bc.DataLabels); i++ {
		h := int(float64(bc.Data[i]) / bc.scale)
		oftX := i * (bc.BarWidth + bc.BarGap)
		// plot bar
		for j := 0; j < bc.BarWidth; j++ {
			for k := 0; k < h; k++ {
				p := Point{}
				p.Ch = ' '
				p.Bg = bc.BarColor
				p.X = bc.innerX + i*(bc.BarWidth+bc.BarGap) + j
				p.Y = bc.innerY + bc.innerHeight - 2 - k
				ps = append(ps, p)
			}
		}
		// plot text
		for j := 0; j < len(bc.labels[i]); j++ {
			p := Point{}
			p.Ch = bc.labels[i][j]
			p.Bg = bc.BgColor
			p.Fg = bc.TextColor
			p.Y = bc.innerY + bc.innerHeight - 1
			p.X = bc.innerX + oftX + j
			ps = append(ps, p)
		}
		// plot num
		for j := 0; j < len(bc.dataNum[i]); j++ {
			p := Point{}
			p.Ch = bc.dataNum[i][j]
			p.Fg = bc.NumColor
			p.Bg = bc.BarColor
			if h == 0 {
				p.Bg = bc.BgColor
			}
			p.X = bc.innerX + oftX + (bc.BarWidth-len(bc.dataNum[i]))/2 + j
			p.Y = bc.innerY + bc.innerHeight - 2
			ps = append(ps, p)
		}
	}

	return ps
}
