//go:build !dev
// +build !dev

package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Renderer struct {
	rt  wazero.Runtime
	mod wazero.CompiledModule
}

func NewRenderer(ctx context.Context) (*Renderer, error) {
	rt := wazero.NewRuntime(ctx)

	_, err := wasi_snapshot_preview1.Instantiate(ctx, rt)
	if err != nil {
		return nil, fmt.Errorf("error instantiating WASI module: %w", err)
	}

	mod, err := rt.CompileModule(ctx, wasmBundle)
	if err != nil {
		return nil, fmt.Errorf("error compiling wasm module: %w", err)
	}

	return &Renderer{rt: rt, mod: mod}, nil
}

func (r *Renderer) Close(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping WASM runtime")
	return errors.Join(r.mod.Close(ctx), r.rt.Close(ctx))
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

	var stdin bytes.Buffer
	err := json.NewEncoder(&stdin).Encode(input)
	if err != nil {
		return fmt.Errorf("error rendering page %s: error marshalling data to JSON: %w", input.Page, err)
	}

	stderr := bytes.NewBuffer(make([]byte, 0, 1024))

	instance, err := r.rt.InstantiateModule(ctx, r.mod, wazero.NewModuleConfig().WithName("").WithStdin(&stdin).WithStdout(output).WithStderr(stderr))
	if err != nil {
		return fmt.Errorf("error rendering page %s: error instantiating module: %w: %s", input.Page, err, stderr.String())
	}

	if stderr.Len() != 0 {
		return fmt.Errorf("error rendering page %s: %s", input.Page, strings.TrimSpace(stderr.String()))
	}

	return instance.Close(ctx)
}
