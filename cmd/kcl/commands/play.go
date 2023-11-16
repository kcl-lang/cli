package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"kcl-lang.io/kcl-playground/pkg/play"
)

const (
	playDesc = `
This command opens the kcl playground in the browser.
`
	playExample = `  # Open in the localhost:80 and open the browser.
  kcl play --open

  # Open with the addr
  kcl play --addr 127.0.0.1:80 --open

  # Only listen the address and do not open the the browser.
  kcl play
  `
)

var (
	addr string = ":80"
)

// NewPlayCmd returns the run command.
func NewPlayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "play",
		Short:   "Open the kcl playground in the browser.",
		Long:    playDesc,
		Example: playExample,
		RunE: func(*cobra.Command, []string) error {
			return runPlayground(addr)
		},
		Aliases:      []string{"p"},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&addr, "addr", ":80", "set listen address, in the form 'host:port'")

	return cmd
}

func runPlayground(addr string) error {
	opts := play.Options{
		PlayMode:   true,
		AllowShare: true,
	}
	fmt.Printf("[Info] Playground listens at %s\n", addr)
	go func() {
		time.Sleep(time.Second * 2)
		openBrowser(addr)
	}()
	return play.Run(addr, &opts)
}

func openBrowser(addr string) error {
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", addr).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", addr).Start()
	case "darwin":
		return exec.Command("open", addr).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
