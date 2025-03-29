package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

type SearchOptions struct {
	Pattern       string   `json:"pattern"`
	Directory     string   `json:"directory"`
	CaseSensitive bool     `json:"caseSensitive"`
	IncludeHidden bool     `json:"includeHidden"`
	FileTypes     []string `json:"fileTypes"`
	MaxDepth      int      `json:"maxDepth"`
	Extensions    []string `json:"extensions"`
}

type FileResult struct {
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Directory string    `json:"directory"`
	IsDir     bool      `json:"isDir"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"modTime"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Search(options SearchOptions) ([]FileResult, error) {
	if !checkFdAvailable() {
		return nil, fmt.Errorf("fd command not found. Please install fd: https://github.com/sharkdp/fd")
	}

	args := []string{}

	args = append(args, options.Pattern)

	if options.Directory != "" {
		args = append(args, options.Directory)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			args = append(args, ".")
		} else {
			args = append(args, home)
		}
	}

	if !options.CaseSensitive {
		args = append(args, "--ignore-case")
	}

	if options.IncludeHidden {
		args = append(args, "--hidden")
	}
	
	for _, fileType := range options.FileTypes {
		if fileType == "f" || fileType == "d" {
			args = append(args, "--type", fileType)
		}
	}

	if options.MaxDepth > 0 {
		args = append(args, "--max-depth", fmt.Sprintf("%d", options.MaxDepth))
	}

	for _, ext := range options.Extensions {
		args = append(args, "--extension", ext)
	}

	cmd := exec.Command("fd", args...)
	fmt.Println(cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting fd command: %w", err)
	}

	results := []FileResult{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		path := scanner.Text()

		fmt.Println(path)

		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		result := FileResult{
			Path: path,
			Name: filepath.Base(path),
			Directory: filepath.Dir(path),
			IsDir: info.IsDir(),
			Size: info.Size(),
			ModTime: info.ModTime(),
		}

		results = append(results, result)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading fd output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("fd command error: %w", err)
	}

	return results, nil
}

func checkFdAvailable() bool {
	_, err := exec.LookPath("fd")
	return err == nil
}

