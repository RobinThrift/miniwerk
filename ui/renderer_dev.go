//go:build dev
// +build dev

package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Renderer struct {
	rt wazero.Runtime

	mu      sync.RWMutex
	mod     wazero.CompiledModule
	modtime time.Time
}

func NewRenderer(ctx context.Context) (*Renderer, error) {
	rt := wazero.NewRuntime(ctx)
	_, err := wasi_snapshot_preview1.Instantiate(ctx, rt)
	if err != nil {
		return nil, fmt.Errorf("error instantiating WASI module: %w", err)
	}

	return &Renderer{rt: rt}, nil
}

func (r *Renderer) Close(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping WASM runtime")
	if r.mod != nil {
		err := r.mod.Close(ctx)
		if err != nil {
			return err
		}
	}
	return errors.Join(r.rt.Close(ctx))
}

type RenderInput struct {
	Page string `json:"page,omitempty"`
	Data any    `json:"data,omitempty"`
}

func (r *Renderer) Render(ctx context.Context, input RenderInput, output io.Writer) error {
	start := time.Now()
	defer func(ctx context.Context) {
		dur := float64(time.Since(start)) / float64(time.Millisecond)
		slog.DebugContext(ctx, fmt.Sprintf("rendered page %s in %fms", input.Page, dur))
	}(ctx)
	slog.DebugContext(ctx, fmt.Sprintf("rendering page %s", input.Page))

	mod, err := r._loadMod(ctx)
	if err != nil {
		return err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var stdin bytes.Buffer
	err = json.NewEncoder(&stdin).Encode(input)
	if err != nil {
		return fmt.Errorf("error rendering page %s: error marshalling data to JSON: %w", input.Page, err)
	}

	stderr := bytes.NewBuffer(make([]byte, 0, 1024))

	instance, err := r.rt.InstantiateModule(ctx, mod, wazero.NewModuleConfig().WithName("").WithStdin(&stdin).WithStdout(output).WithStderr(stderr))
	if err != nil {
		return fmt.Errorf("error rendering page %s: error instantiating module: %w: %s", input.Page, err, stderr.String())
	}

	if stderr.Len() != 0 {
		return fmt.Errorf("error rendering page %s: %s", input.Page, strings.TrimSpace(stderr.String()))
	}

	return instance.Close(ctx)
}

const wasmFilename = "./ui/build/server/index.wasm"

func (r *Renderer) _loadMod(ctx context.Context) (wazero.CompiledModule, error) {
	newer, modtime, err := r.isNewerWASMFileAvailable()
	if err != nil {
		return nil, err
	}
	if !newer {
		return r.mod, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	start := time.Now()
	defer func(ctx context.Context) {
		dur := float64(time.Since(start)) / float64(time.Millisecond)
		slog.DebugContext(ctx, fmt.Sprintf("compiled WASM module in %fms", dur))
	}(ctx)
	slog.DebugContext(ctx, "newer WASM bundle detected, (re)loading from disk")

	serverWasmBytes, err := os.ReadFile(wasmFilename)
	if err != nil {
		return nil, err
	}

	mod, err := r.rt.CompileModule(ctx, serverWasmBytes)
	if err != nil {
		return nil, fmt.Errorf("error compiling wasm module: %w", err)
	}

	r.mod = mod
	r.modtime = modtime

	return mod, nil
}

func (r *Renderer) isNewerWASMFileAvailable() (bool, time.Time, error) {
	stat, err := os.Stat(wasmFilename)
	if err != nil {
		return false, time.Time{}, err
	}

	modtime := stat.ModTime()

	return modtime.After(r.modtime), modtime, nil
}
