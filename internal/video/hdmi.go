package video

import (
	"strings"

	"github.com/hekmon/cunits/v3"
)

type HDMI struct {
	Version string
	DSC     bool
	HDR     bool
	Modes   []TransmissionMode
}

func (h HDMI) CanHDR(colorDepth ColorDepth) bool {
	if h.HDR {
		if colorDepth >= colorDepth10bit {
			return true
		}
	}
	return false
}

func HDMIVersions() []HDMI {
	return hdmiVersions
}

var hdmiVersions = []HDMI{
	{
		Version: "2.2",
		DSC:     true,
		HDR:     true,
		Modes: []TransmissionMode{
			frl9g, frl18g, frl24g, frl32g, frl40g, frl48g, frl64g, frl80g, frl96g,
		},
	},
	{
		Version: "2.1",
		DSC:     true,
		HDR:     true,
		Modes: []TransmissionMode{
			frl9g, frl18g, frl24g, frl32g, frl40g, frl48g,
		},
	},
	{
		Version: "2.0",
		DSC:     false,
		HDR:     true,
		Modes: []TransmissionMode{
			tmds165, tmds340, tmds600,
		},
	},
	{
		Version: "1.4",
		DSC:     false,
		HDR:     false,
		Modes: []TransmissionMode{
			tmds165, tmds340,
		},
	},
	{
		Version: "1.3",
		DSC:     false,
		HDR:     false,
		Modes: []TransmissionMode{
			tmds165, tmds340,
		},
	},
	{
		Version: "1.2",
		DSC:     false,
		HDR:     false,
		Modes: []TransmissionMode{
			tmds165,
		},
	},
	{
		Version: "1.1",
		DSC:     false,
		HDR:     false,
		Modes: []TransmissionMode{
			tmds165,
		},
	},
	{
		Version: "1.0",
		DSC:     false,
		HDR:     false,
		Modes: []TransmissionMode{
			tmds165,
		},
	},
}

var _ TransmissionMode = HDMITransmissionMode{}

type HDMITransmissionMode struct {
	Name         string
	MaxBandwidth cunits.Speed
}

func (m HDMITransmissionMode) GetName() string {
	return m.Name
}

func (m HDMITransmissionMode) GetBandwidth() cunits.Speed {
	return m.MaxBandwidth
}

func (m HDMITransmissionMode) EffectiveBandwidth() cunits.Speed {
	if strings.Contains(m.Name, "FRL") {
		return cunits.Speed{Bits: cunits.Bits(float64(m.MaxBandwidth.Bits) * (88.8 / 100.0))}
	}
	return cunits.Speed{Bits: cunits.Bits(float64(m.MaxBandwidth.Bits) * (80.0 / 100.0))}
}

func (m HDMITransmissionMode) MaxCompressedBandwidth(colorDepth ColorDepth) cunits.Speed {
	return cunits.Speed{Bits: m.EffectiveBandwidth().Bits * cunits.Bits(float64(colorDepth)/8)}
}

func (m HDMITransmissionMode) Usage(bandwidth cunits.Speed) float64 {
	return float64(bandwidth.Bits*100) / float64(m.EffectiveBandwidth().Bits)
}

var (
	tmds165 = HDMITransmissionMode{
		Name:         "TMDS (165 MHz)",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(4.95)},
	}
	tmds340 = HDMITransmissionMode{
		Name:         "TMDS (340 MHz)",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(10.20)},
	}
	tmds600 = HDMITransmissionMode{
		Name:         "TMDS (600 MHz)",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(18.00)},
	}
	frl9g = HDMITransmissionMode{
		Name:         "FRL 1",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(9.00)},
	}
	frl18g = HDMITransmissionMode{
		Name:         "FRL 2",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(18.00)},
	}
	frl24g = HDMITransmissionMode{
		Name:         "FRL 3",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(24.00)},
	}
	frl32g = HDMITransmissionMode{
		Name:         "FRL 4",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(32.00)},
	}
	frl40g = HDMITransmissionMode{
		Name:         "FRL 5",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(40.00)},
	}
	frl48g = HDMITransmissionMode{
		Name:         "FRL 6",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(48.00)},
	}
	frl64g = HDMITransmissionMode{
		Name:         "FRL 7",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(64.00)},
	}
	frl80g = HDMITransmissionMode{
		Name:         "FRL 8",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(80.00)},
	}
	frl96g = HDMITransmissionMode{
		Name:         "FRL 9",
		MaxBandwidth: cunits.Speed{Bits: cunits.ImportInGb(96.00)},
	}
)
