package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var skipDirs = map[string]bool {
	"node_modules": true,
    ".git": true,
    "vendor": true,
    "dist": true,
    "build": true,
    ".idea": true,
    ".vscode": true,
    ".cache": true,
    ".npm": true,
    "target": true,
    "bin": true,
    "obj": true,
}

// App struct
type App struct {
	ctx context.Context
}

type SearchOptions struct {
	Pattern       string   `json:"pattern"`
	Directory     string   `json:"directory"`
	// CaseSensitive bool     `json:"caseSensitive"`
	IncludeHidden bool     `json:"includeHidden"`
	// FileTypes     []string `json:"fileTypes"`
	// MaxDepth      int      `json:"maxDepth"`
	// Extensions    []string `json:"extensions"`
}

type FileResult struct {
	ID	      int32     `json:"id"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Extension 	  string   `json:"extension"`
	// Directory string    `json:"directory"`
	IsDir     bool      `json:"isDir"`
	Size      int64     `json:"size"`
	// ModTime   time.Time `json:"modTime"`
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


func (a *App) GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

func walkDirRecursive(path string, fn func(path string, d fs.DirEntry, err error) error) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fn(path, nil, err)
	}

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())

		if err := fn(childPath, entry, nil); err != nil {
			if err == filepath.SkipDir {
				if entry.IsDir() {
					continue
				}
				return nil
			}
			return err
		}

		if entry.IsDir() {
			if err := walkDirRecursive(childPath, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func walkDir(root string, fn func(path string, d fs.DirEntry, err error) error) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	sem := make(chan struct{}, runtime.NumCPU()*2)

	for _, entry := range entries {
		if !entry.IsDir() {
			if err := fn(filepath.Join(root, entry.Name()), entry, nil); err != nil {
				return err
			}
			continue
		}

		subdir := filepath.Join(root, entry.Name())
		wg.Add(1)
		sem <- struct{}{}

		go func(path string, d fs.DirEntry) {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			if err := fn(path, d, nil); err != nil {
				if err != filepath.SkipDir {
					select {
					case errChan <- err:
					default:
					}
				}
				return
			}

			if err := walkDirRecursive(path, fn); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(subdir, entry)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err, ok := <-errChan; ok {
		return err
	}

	return nil
}

func containsRegexChars(path string) bool {
	specialChars := []byte{'*', '+', '?', '[', ']', '(', ')', '{', '}', '^', '$', '.', '|', '\\'}

	for _, char := range specialChars {
		if strings.ContainsRune(path, rune(char)) {
			return true
		}
	}

	return false
}

func (a *App) Search(options SearchOptions) ([]FileResult, error) {	
	root, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting root director %w", err)
	}

	isSimpleChar := containsRegexChars(options.Pattern)
	var matcher func(string) bool

	if isSimpleChar {
		pattern := strings.ToLower(options.Pattern)
		matcher = func(name string) bool {
			return strings.Contains(strings.ToLower(name), pattern )
		}
	} else {
		regex, err := regexp.Compile("(?i)" + options.Pattern)
		if err != nil {
			return nil, fmt.Errorf("error compiling pattern to regex %w", err)
		}
		matcher = regex.MatchString
	}

	var id int32

	err = walkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
	
		if strings.HasPrefix(d.Name(), ".") || skipDirs[d.Name()] {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
	
		if matcher(d.Name()) {
			id++
			wailsRuntime.EventsEmit(a.ctx, "entry:found", FileResult{
				ID: id,
				Path: path,
				Name: d.Name(),
				Extension: filepath.Ext(path),
				IsDir: d.IsDir(),
			})
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *App) SearchSubdirectory(path string) ([]FileResult, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading subdir %w", err)
	}

	var results []FileResult
	var id int32

	for _, entry := range entries {
		path := filepath.Join(path, entry.Name())
		id++

		results = append(results, FileResult{
			ID: id,
			Path: path,
			Name: entry.Name(),
			Extension: filepath.Ext(path),
			IsDir: entry.IsDir(),
		})
	}

	return results, nil
}

func (a *App) OpenFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	case "darwin":
		cmd = exec.Command("open", path)
	}

	return cmd.Start()
}
