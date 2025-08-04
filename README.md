# Logging component for [bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework
If you want to duplicate logs of your application to TUI and any other number of io.Writers.
To integrate it to slog:
```go
	slog.SetDefault(slog.New(slog.NewTextHandler(m)))
```
