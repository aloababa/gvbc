package video

import "github.com/hekmon/cunits/v3"

type DisplayPort struct {
	Version string
	DSC     bool
	HDR     bool
	Modes   []TransmissionMode
}

func (d DisplayPort) CanHDR(colorDepth ColorDepth) bool {
	if d.HDR {
		if colorDepth >= colorDepth10bit {
			return true
		}
	}
	return false
}

func DisplayPortVersions() []DisplayPort {
	return displayPortVersions
}

var displayPortVersions = []DisplayPort{
	{
		Version: "2.x",
		DSC:     true,
		HDR:     true,
		Modes:   []TransmissionMode{uhbr20, uhbr135, uhbr10},
	},
	{
		Version: "1.4",
		DSC:     true,
		HDR:     true,
		Modes:   []TransmissionMode{hbr3},
	},
	{
		Version: "1.3",
		DSC:     false,
		HDR:     false,
		Modes:   []TransmissionMode{hbr3},
	},
	{
		Version: "1.2",
		DSC:     false,
		HDR:     false,
		Modes:   []TransmissionMode{hbr2},
	},
	{
		Version: "1.1",
		DSC:     false,
		HDR:     false,
		Modes:   []TransmissionMode{hbr},
	},
	{
		Version: "1.0",
		DSC:     false,
		HDR:     false,
		Modes:   []TransmissionMode{hbr, rbr},
	},
}

var _ TransmissionMode = DisplayPortTransmissionMode{}

type DisplayPortTransmissionMode struct {
	Name         string
	MaxBandwidth cunits.Speed
}

func (m DisplayPortTransmissionMode) GetName() string {
	return m.Name
}

func (m DisplayPortTransmissionMode) GetBandwidth() cunits.Speed {
	return m.MaxBandwidth
}

func (m DisplayPortTransmissionMode) EffectiveBandwidth() cunits.Speed {
	switch m.Name {
	case "RBR", "HBR", "HBR2", "HBR3":
		return cunits.Speed{Bits: cunits.Bits(float64(m.MaxBandwidth.Bits) * (80.0 / 100.0))}
	}
	return cunits.Speed{Bits: cunits.Bits(float64(m.MaxBandwidth.Bits) * (96.7 / 100.0))}
}

func (m DisplayPortTransmissionMode) MaxCompressedBandwidth(colorDepth ColorDepth) cunits.Speed {
	return cunits.Speed{Bits: m.EffectiveBandwidth().Bits * cunits.Bits(float64(colorDepth)/8)}
}

func (m DisplayPortTransmissionMode) Usage(bandwidth cunits.Speed) float64 {
	return float64(bandwidth.Bits*100) / float64(m.EffectiveBandwidth().Bits)
}

var (
	rbr = DisplayPortTransmissionMode{
		Name:         "RBR",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(6.48)},
	}
	hbr = DisplayPortTransmissionMode{
		Name:         "HBR",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(10.8)},
	}
	hbr2 = DisplayPortTransmissionMode{
		Name:         "HBR2",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(21.6)},
	}
	hbr3 = DisplayPortTransmissionMode{
		Name:         "HBR3",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(32.4)},
	}
	uhbr10 = DisplayPortTransmissionMode{
		Name:         "UHBR10",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(40)},
	}
	uhbr135 = DisplayPortTransmissionMode{
		Name:         "UHBR13.5",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(54)},
	}
	uhbr20 = DisplayPortTransmissionMode{
		Name:         "UHBR20",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(80)},
	}
)
