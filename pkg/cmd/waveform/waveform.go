package waveform

import (
	"errors"
	"fmt"
	"os"

	"github.com/amaretto/go-waveform/pkg/waveform"
	"github.com/spf13/cobra"
)

var src string
var dst string

// NewCommand create command
func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "waveform path/to/mp3",
		Short: "A waveform generator",
		Long:  `A waveform generator`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires source mp3 path")
			}
			src = args[0]
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := saveWaveImage(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	c.PersistentFlags().StringVarP(&dst, "dest", "d", "barchart.png", "image file path(default is ./barchart.png")
	return c
}

func saveWaveImage() error {
	w := waveform.NewWaveformer()
	w.MusicPath = src
	if err := w.SaveWaveImage(dst); err != nil {
		return err
	}
	return nil
}
