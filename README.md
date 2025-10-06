# Go Video Bandwidth Calculator

## Overview

The **Go Video Bandwidth Calculator** is a command-line application written in Go that calculates the bandwidth requirements for various video displays based on resolution, refresh rate, color depth, and timing standards. It supports different digital interfaces such as DisplayPort and HDMI and provides a user-friendly text-based user interface (TUI) for input and output.

![Demo](./demo.gif)

## Features

- **Dynamic Input**: Users can specify video parameters interactively via a text-based interface.
- **Presets**: Predefined video display settings (e.g. 4K, 1080p) for quick selection.
- **Bandwidth Calculation**: Computes the required bandwidth for the specified video display parameters.
- **Compatibility Check**: Checks compatibility with DisplayPort and HDMI versions based on the calculated bandwidth and color depth.

## Installation

### Prerequisites

- Go 1.25.1 or later installed on your system.

### Install using go

```bash
go install github.com/aloababa/gvbc@latest
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
