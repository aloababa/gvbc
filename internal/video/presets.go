package video

type Preset struct {
	Name    string
	Display Display
}

func Presets() []Preset {
	return []Preset{
		{
			Name: "8k Display",
			Display: Display{
				Width:       7680,
				Height:      4320,
				RefreshRate: 60,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "4k ESport",
			Display: Display{
				Width:       3840,
				Height:      2160,
				RefreshRate: 240,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "4k Gaming",
			Display: Display{
				Width:       3840,
				Height:      2160,
				RefreshRate: 144,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "2k ESport",
			Display: Display{
				Width:       2560,
				Height:      1440,
				RefreshRate: 360,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "2k Gaming",
			Display: Display{
				Width:       2560,
				Height:      1440,
				RefreshRate: 180,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "1080p ESport",
			Display: Display{
				Width:       1920,
				Height:      1080,
				RefreshRate: 480,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
		{
			Name: "1080p Gaming",
			Display: Display{
				Width:       1920,
				Height:      1080,
				RefreshRate: 240,
				ColorDepth:  colorDepth10bit,
				Timing:      cvtrbv2,
			},
		},
	}
}
