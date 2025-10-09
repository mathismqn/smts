package pdf

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func (p *PDF) AddWatermark(text string, posX, posY float64) error {
	wm, _ := pdfcpu.ParseTextWatermarkDetails(
		text,
		fmt.Sprintf("pos:bl, off:%f %f, scale:0.4 abs, op:1, rot:0, font:Helvetica, color:#000000", posX, posY),
		true,
		types.POINTS,
	)

	return api.AddWatermarksFile(p.path, p.path, []string{"l"}, wm, nil)
}

func (p *PDF) AddSignature(sigPath string) error {
	wm, _ := pdfcpu.ParseImageWatermarkDetails(
		sigPath,
		"pos:bl, off:80 0, scale:0.3 abs, op:1, rot:0",
		true,
		types.POINTS,
	)

	return api.AddWatermarksFile(p.path, p.path, []string{"l"}, wm, nil)
}
