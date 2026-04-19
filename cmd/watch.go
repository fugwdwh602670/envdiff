package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/envdiff/internal/diff"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch <fileA> <fileB>",
	Short: "Watch two .env files and print diffs on change",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileA, fileB := args[0], args[1]

		opts := diff.WatchOptions{
			Interval: time.Duration(watchInterval) * time.Millisecond,
			OnChange: func(results []diff.Result) {
				out, err := diff.FormatResults(results, "text")
				if err != nil {
					fmt.Fprintf(os.Stderr, "format error: %v\n", err)
					return
				}
				fmt.Println("--- diff update ---")
				fmt.Print(out)
			},
		}

		stop := make(chan struct{})
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sig
			close(stop)
		}()

		fmt.Fprintf(os.Stderr, "Watching %s and %s (interval %dms)...\n", fileA, fileB, watchInterval)
		return diff.Watch(fileA, fileB, opts, stop)
	},
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2000, "Poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
