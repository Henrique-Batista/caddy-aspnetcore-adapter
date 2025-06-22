package middleware

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Config é preenchido via Caddyfile ou linha de comando:
// caddy run --process_path <path> --process_args "MyApp.dll --urls http://127.0.0.1:5000"
type Config struct {
	ProcessPath  string // caminho do dotnet ou exe
	Arguments    string // argumentos para o runtime
	RestartDelay time.Duration
}

// StartProcess inicia o loop de execução do processo ASP.NET Core
func StartProcess(ctx context.Context, cfg Config) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		parts := strings.Fields(cfg.Arguments)
		cmd := exec.CommandContext(ctx, cfg.ProcessPath, parts...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Iniciando processo: %s %s\n", cfg.ProcessPath, cfg.Arguments)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("falha ao iniciar ASP.NET Core: %w", err)
		}

		err := cmd.Wait()
		fmt.Printf("Processo terminou (%v), reiniciando em %v...\n", err, cfg.RestartDelay)
		time.Sleep(cfg.RestartDelay)
	}
}
