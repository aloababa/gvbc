package main

import (
	"flag"
	"io"
	"log"
	"log/slog"

	"github.com/aloababa/gvbc/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fps := flag.Int("fps", 60, "The max fps at which the renderer should run (min 1, max 120)")
	debug := flag.Bool("debug", false, "Enable debug to log file")
	debugLogFile := flag.String("log-file", "debug.log", "The path to debug log file")
	flag.Parse()
	err := run(*debug, *debugLogFile, *fps)
	if err != nil {
		log.Fatal(err)
	}
}

func run(debug bool, debugLogFile string, fps int) error {
	if debug {
		f, err := tea.LogToFile(debugLogFile, "debug")
		if err != nil {
			return err
		}
		defer f.Close()
		slog.SetDefault(slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	}
	m := tui.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithFPS(fps))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
