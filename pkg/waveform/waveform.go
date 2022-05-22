package waveform

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	sampleInterval    = 800
	smoothSampleCount = 3
	heightMax         = 20
)

// Waveformer treats waveforms
type Waveformer struct {
	MusicPath         string
	SampleInterval    int
	SmoothSampleCount int
	HeightMax         int
}

// Waveform has information of wave thined out
type Waveform struct {
	MusicTitle string
	Wave       []int
}

// NewWaveformer returns instance of Waveformer
func NewWaveformer(heightMax, sampleInterval, smoothSampleCount int) *Waveformer {
	w := new(Waveformer)
	w.HeightMax = heightMax
	w.SampleInterval = sampleInterval
	w.SmoothSampleCount = smoothSampleCount
	return w
}

// SaveWaveImage generate waveform and output it as images
func (w *Waveformer) SaveWaveImage(destFile string) error {
	wave, err := w.GenWaveForm()
	if err != nil {
		return err
	}

	parray := make([]float64, len(wave.Wave))
	for i, n := range wave.Wave {
		parray[i] = float64(n)
	}

	var group plotter.Values
	group = parray
	fmt.Println(parray)
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Waveform"
	p.X.Label.Text = "sample"
	p.Y.Label.Text = "math of volume"
	p.X.Min = 0
	p.X.Max = 10000
	p.Y.Min = 0
	p.Y.Max = 20

	width := vg.Points(1)

	bars, err := plotter.NewBarChart(group, width)
	if err != nil {
		return err
	}

	bars.Color = plotutil.Color(1)
	p.Add(bars)

	if err := p.Save(13*vg.Inch, 3*vg.Inch, destFile); err != nil {
		panic(err)
	}

	return nil
}

// GenWaveForm generate and normalize waveform from mp3 file
func (w *Waveformer) GenWaveForm() (*Waveform, error) {
	f, err := os.Open(w.MusicPath)
	if err != nil {
		return nil, err
	}
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	defer streamer.Close()

	// gen,smoothing,normalize
	rwave := genRawWave(streamer, w.SampleInterval)
	smoothRawWave(rwave, w.SmoothSampleCount)
	wave := normalizeRawWave(rwave, w.HeightMax)

	wf := &Waveform{"", wave}

	return wf, nil
}

// GenRawWave generate raw waveform values([]float64) from volume values
func genRawWave(streamer beep.StreamSeeker, sampleInterval int) []float64 {
	var tmp [2][2]float64
	var count, ncount int
	var rwave []float64
	// ToDo : currently, raw wave limit is 100000
	rwave = make([]float64, 100000)

	for {
		// check EOF
		if sn, sok := streamer.Stream(tmp[:1]); sn == 0 && !sok {
			break
		}
		samplel := tmp[0][0]
		sampler := tmp[0][1]

		sumSquare := math.Pow(samplel, 2)
		sumSquare += math.Pow(sampler, 2)
		value := math.Sqrt(sumSquare)

		if count%sampleInterval == 0 {
			rwave[ncount] = value
			ncount++
		}

		count++
	}

	rwave = rwave[:ncount]
	return rwave
}

// SmoothRawWave make raw wave values smoothly for visualizaiton
func smoothRawWave(rwave []float64, smoothSampleCount int) {
	var sum float64
	for i := 0; i < len(rwave); i++ {
		if i < len(rwave)-smoothSampleCount {
			sum = 0
			for j := 0; j < smoothSampleCount; j++ {
				sum += rwave[i+j]
			}
			rwave[i] = sum / float64(smoothSampleCount)
		} else {
			rwave[i] = rwave[i-1]
		}
	}
}

// NormalizeRawWave arrange wave values utilizing heightMax as height
func normalizeRawWave(rwave []float64, heightMax int) []int {
	var max float64
	var limit float64
	// ToDo : need to check(use sqrt(2)?)
	max = 1.0
	limit = float64(heightMax)

	var r []int
	r = make([]int, len(rwave))
	for i, num := range rwave {
		r[i] = int(math.Ceil(limit * num / max))
	}

	return r
}
