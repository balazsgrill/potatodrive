package ui

import (
	"log"

	"github.com/balazsgrill/potatodrive/core"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type widthDPI struct {
	width int // in native pixels
	dpi   int
}

type textWidthDPI struct {
	text  string
	width int // in native pixels
	dpi   int
}

type Styler struct {
	lb     **walk.ListBox
	canvas *walk.Canvas
	model  *StatusList

	font                *walk.Font
	dpi2StampSize       map[int]walk.Size
	widthDPI2WsPerLine  map[widthDPI]int
	textWidthDPI2Height map[textWidthDPI]int // in native pixels
	stateicons          map[core.FileSyncStateEnum]*walk.Icon
}

func (s *Styler) ItemHeightDependsOnWidth() bool {
	return true
}

func (s *Styler) DefaultItemHeight() int {
	dpi := (*s.lb).DPI()
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)

	return s.StampSize().Height + marginV*2
}

const (
	marginH96dpi int = 6
	marginV96dpi int = 2
	lineW96dpi   int = 1
)

const (
	imageres_dll_mun = "C:\\Windows\\SystemResources\\imageres.dll.mun"
	iconid_done      = 233
	iconid_working   = 229
	iconid_error     = 230
	iconid_pending   = 232

	status_icon_size = 16
)

func (s *Styler) loadIcons() {
	var err error
	s.stateicons[core.FileSyncStateDone], err = walk.NewIconExtractedFromFileWithSize(imageres_dll_mun, iconid_done, status_icon_size)
	if err != nil {
		log.Fatal(err)
	}
	s.stateicons[core.FileSyncStateDeleted] = s.stateicons[core.FileSyncStateDone]
	s.stateicons[core.FileSyncStateDirty], err = walk.NewIconExtractedFromFileWithSize(imageres_dll_mun, iconid_pending, status_icon_size)
	if err != nil {
		log.Fatal(err)
	}
	s.stateicons[core.FileSyncStatePending] = s.stateicons[core.FileSyncStateDirty]
	s.stateicons[core.FileSyncStateDownloading], err = walk.NewIconExtractedFromFileWithSize(imageres_dll_mun, iconid_working, status_icon_size)
	if err != nil {
		log.Fatal(err)
	}
	s.stateicons[core.FileSyncStateUploading] = s.stateicons[core.FileSyncStateDownloading]
	s.stateicons[core.FileSyncStateError], err = walk.NewIconExtractedFromFileWithSize(imageres_dll_mun, iconid_error, status_icon_size)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Styler) ItemHeight(index, width int) int {
	dpi := (*s.lb).DPI()
	marginH := walk.IntFrom96DPI(marginH96dpi, dpi)
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
	lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

	msg := s.model.statuses[index].Path

	twd := textWidthDPI{msg, width, dpi}

	if height, ok := s.textWidthDPI2Height[twd]; ok {
		return height + marginV*2
	}

	canvas, err := s.Canvas()
	if err != nil {
		return 0
	}

	stampSize := s.StampSize()

	wd := widthDPI{width, dpi}
	wsPerLine, ok := s.widthDPI2WsPerLine[wd]
	if !ok {
		bounds, _, err := canvas.MeasureTextPixels("W", (*s.lb).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
		if err != nil {
			return 0
		}
		wsPerLine = (width - marginH*4 - lineW - stampSize.Width) / bounds.Width
		s.widthDPI2WsPerLine[wd] = wsPerLine
	}

	if len(msg) <= wsPerLine {
		s.textWidthDPI2Height[twd] = stampSize.Height
		return stampSize.Height + marginV*2
	}

	bounds, _, err := canvas.MeasureTextPixels(msg, (*s.lb).Font(), walk.Rectangle{Width: width - marginH*4 - lineW - stampSize.Width, Height: 255}, walk.TextEditControl|walk.TextWordbreak|walk.TextEndEllipsis)
	if err != nil {
		return 0
	}

	s.textWidthDPI2Height[twd] = bounds.Height

	return bounds.Height + marginV*2
}

func (s *Styler) StyleItem(style *walk.ListItemStyle) {
	if canvas := style.Canvas(); canvas != nil {
		if style.Index()%2 == 1 && style.BackgroundColor == walk.Color(win.GetSysColor(win.COLOR_WINDOW)) {
			style.BackgroundColor = walk.Color(win.GetSysColor(win.COLOR_BTNFACE))
			if err := style.DrawBackground(); err != nil {
				return
			}
		}

		pen, err := walk.NewCosmeticPen(walk.PenSolid, style.LineColor)
		if err != nil {
			return
		}
		defer pen.Dispose()

		dpi := (*s.lb).DPI()
		marginH := walk.IntFrom96DPI(marginH96dpi, dpi)
		marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
		lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

		b := style.BoundsPixels()
		b.X += marginH
		b.Y += marginV

		item := s.model.statuses[style.Index()]

		canvas.DrawImagePixels(s.stateicons[item.State], walk.Point{b.X, b.Y})

		stampSize := s.StampSize()

		x := b.X + stampSize.Width + marginH + lineW
		canvas.DrawLinePixels(pen, walk.Point{x, b.Y - marginV}, walk.Point{x, b.Y - marginV + b.Height})

		b.X += stampSize.Width + marginH*2 + lineW
		b.Width -= stampSize.Width + marginH*4 + lineW

		style.DrawText(item.Path, b, walk.TextEditControl|walk.TextWordbreak|walk.TextEndEllipsis)
	}
}

func (s *Styler) StampSize() walk.Size {
	dpi := (*s.lb).DPI()

	stampSize, ok := s.dpi2StampSize[dpi]
	if !ok {
		s.dpi2StampSize[dpi] = walk.SizeTo96DPI(walk.Size{Width: 32, Height: 32}, dpi)
	}

	return stampSize
}

func (s *Styler) Canvas() (*walk.Canvas, error) {
	if s.canvas != nil {
		return s.canvas, nil
	}

	canvas, err := (*s.lb).CreateCanvas()
	if err != nil {
		return nil, err
	}
	s.canvas = canvas
	(*s.lb).AddDisposable(canvas)

	return canvas, nil
}
