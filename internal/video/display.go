package video

import (
	"fmt"
	"math"

	"github.com/hekmon/cunits/v3"
)

type Display struct {
	Width       int
	Height      int
	RefreshRate int
	ColorDepth  ColorDepth
	Timing      Timing
}

func (d Display) String() string {
	return fmt.Sprintf("%dx%d@%dHz, color depth: %s, timing: %s",
		d.Width, d.Height, d.RefreshRate, d.ColorDepth, d.Timing.String())
}

func (d Display) FrameSize() int {
	return d.Width * d.Height
}

func (d Display) EffectiveFrameSize() int {
	return d.Timing.EffectiveFrameSize(d)
}

func (d Display) EffectivePixelRate() int {
	return d.EffectiveFrameSize() * d.RefreshRate
}

func (d Display) Bandwidth() cunits.Speed {
	return cunits.Speed{Bits: cunits.Bits(d.EffectivePixelRate() * int(d.ColorDepth))}
}

func (d Display) DSC() cunits.Speed {
	return cunits.Speed{Bits: cunits.Bits(float64(d.Bandwidth().Bits) / (float64(d.ColorDepth) / 8))}
}

type ColorDepth int

func (c ColorDepth) String() string {
	switch c {
	case colorDepth8bit:
		return "8 bpc (24 bit/px)"
	case colorDepth10bit:
		return "10 bpc (30 bit/px)"
	case colorDepth12bit:
		return "12 bpc (36 bit/px)"
	case colorDepth16bit:
		return "16 bpc (48 bit/px)"
	}
	return ""
}

func ColorDepths() []ColorDepth {
	return []ColorDepth{
		colorDepth8bit,
		colorDepth10bit,
		colorDepth12bit,
		colorDepth16bit,
	}
}

func ColorDepth8bit() ColorDepth {
	return colorDepth8bit
}

func ColorDepth10bit() ColorDepth {
	return colorDepth10bit
}

func ColorDepth12bit() ColorDepth {
	return colorDepth12bit
}

func ColorDepth16bit() ColorDepth {
	return colorDepth16bit
}

const (
	colorDepth8bit  ColorDepth = 24
	colorDepth10bit ColorDepth = 30
	colorDepth12bit ColorDepth = 36
	colorDepth16bit ColorDepth = 48
)

type Timing interface {
	EffectiveFrameSize(d Display) int
	String() string
}

func Timings() []Timing {
	return []Timing{cvtrb, cvtrbv2}
}

type CVTRBTiming struct {
	Name   string
	VMin   float64
	HBlank int
}

func (t CVTRBTiming) EffectiveFrameSize(d Display) int {
	return (d.Height + int(math.Ceil((float64(d.Height)*t.VMin)/(1/float64(d.RefreshRate)-t.VMin)))) * (d.Width + t.HBlank)
}

func (t CVTRBTiming) String() string {
	return t.Name
}

func CVTRB() Timing {
	return cvtrb
}

func CVTRBv2() Timing {
	return cvtrbv2
}

var (
	cvtrb = CVTRBTiming{
		Name:   "CVT-RB",
		VMin:   0.00046,
		HBlank: 160,
	}

	cvtrbv2 = CVTRBTiming{
		Name:   "CVT-RBv2",
		VMin:   0.00046,
		HBlank: 80,
	}
)

type TransmissionMode interface {
	GetName() string
	GetBandwidth() cunits.Speed
	EffectiveBandwidth() cunits.Speed
	MaxCompressedBandwidth(colorDepth ColorDepth) cunits.Speed
	Usage(bandwidth cunits.Speed) float64
}
