package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/tetsuzawa/spat"
	"log"
	"math/rand"
	"time"
)

const sampleRate = 48000
const framesPerBuffer = 64

var SLTFL, SLTFR []float32
var longL, longR []float32
var soundSLTFL []float32
var soundSLTFR []float32

func init() {
	f, err := spat.ReadDXXFile("SLTF_450_L.DDB")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(f)
	SLTFL = spat.Float64sToFloat32s(f)
	NormFloat32s(SLTFL)
	f, err = spat.ReadDXXFile("SLTF_450_R.DDB")
	if err != nil {
		log.Fatalln(err)
	}
	SLTFR = spat.Float64sToFloat32s(f)
	NormFloat32s(SLTFR)
	fmt.Println("read sltf")

	//longL = make([]float32, framesPerBuffer+len(SLTFL)-1)
	//longR = make([]float32, framesPerBuffer+len(SLTFR)-1)
	longL = make([]float32, 2048+len(SLTFL)-1)
	longR = make([]float32, 2048+len(SLTFR)-1)
}

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	//s := newStereoSine(256, 320, sampleRate)
	s, err := portaudio.OpenDefaultStream(0, 2, sampleRate, framesPerBuffer, processAudio)
	chk(err)
	defer s.Close()
	chk(s.Start())
	time.Sleep(5 * time.Second)
	chk(s.Stop())
}

//type stereoSine struct {
//	*portaudio.Stream
//	stepL, phaseL float64
//	stepR, phaseR float64
//}

//func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
//	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
//	var err error
//	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
//chk(err)
//return s
//}

//func (g *stereoSine) processAudio(out [][]float32) {
//	for i := range out[0] {
//		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
//		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
//		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
//		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
//	}
//}


func processAudio(out [][]float32) {
	var err error
	//fmt.Println(len(out[0]))
	for i := range out[0] {
		r := rand.Float32()
		out[0][i] = r
		out[1][i] = r
	}
	//fmt.Println(len(SLTFL))
	soundSLTFL, err = spat.LinearConvolutionTimeDomainFloat32s(soundSLTFL, out[0], SLTFL)
	if err != nil {
		log.Fatalln(err)
	}
	soundSLTFR, err = spat.LinearConvolutionTimeDomainFloat32s(soundSLTFR, out[1], SLTFR)
	if err != nil {
		log.Fatalln(err)
	}
	NormFloat32s(soundSLTFL)
	NormFloat32s(soundSLTFR)
	for i := range soundSLTFL {
		longL[i] += soundSLTFL[i]
		longR[i] += soundSLTFR[i]
	}
	//fmt.Println(len(soundSLTFL))
	//fmt.Println(len(soundSLTFR))
	//fmt.Println(len(out[0]))
	for i := range out[0] {
		//out[0][i] = float32(rand.Float32())
		//out[1][i] = float32(rand.Float32())
		//out[0][i] = soundSLTFL[i]
		//out[1][i] = soundSLTFR[i]
		//out[0][i] = outBuf[0][i]
		//out[1][i] = outBuf[1][i]
		out[0][i] = longL[i]
		out[1][i] = longR[i]
	}
	shift(longL, len(out[0]))
	shift(longR, len(out[1]))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func d10000(x []float32) {
	for i, v := range x {
		x[i] = v / 10000.
	}
}

func NormFloat32s(data []float32) {
	_, min, max := spat.AbsMinMaxFloat32s(data)
	for i, v := range data {
		vv := (spat.AbsFloat32(v) - min) / (max - min)
		data[i] = vv
		if v < 0 {
			vv = -vv
		}
	}
}

func shift(x []float32, n int) {
	if n > len(x)-1 {
		return
	}
	for i := 0; i < len(x)-n; i++ {
		if i > len(x)-1 {
			return
		}
		x[i] = x[i+n]
	}
	for i := 0; i < n; i++ {
		x[len(x)-n+i] = 0
	}
}
