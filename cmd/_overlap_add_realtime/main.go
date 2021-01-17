package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/tetsuzawa/spat"
)

const (
	port = 9999

	NumInputChannels  = 0
	NumOutputChannels = 2
	SamplingFreq      = 48000
	//FramesPerBuffer   = 1024
	FramesPerBuffer = 0

	AngleNum = 3600
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("Usage: %s subject sound_file(.DXX)\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	if err := run(); err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	if flag.NArg() != 2 {
		return errors.New("invalid arguments")
	}
	args := flag.Args()
	subject := args[0]
	soundName := args[1]

	// declare a channel to receive signal
	sigCh := make(chan os.Signal, 1)
	// receive
	signal.Notify(sigCh, os.Interrupt)
	//pass the channel to the main processing function
	ctx, cancel := context.WithCancel(context.Background())

	// ******** signal handler ********
	//wg := &sync.WaitGroup{}

	// wait a signal in another gotoutine
	go func() {
		// block until func receive a signal
		sig := <-sigCh
		log.Printf("got signal %v\n", sig)
		log.Println("exiting...")
		cancel()
	}()
	// ********************************

	angleCh := make(chan int)

	go updateAngleBySocketConn(ctx, sigCh, angleCh, port)

	return OverlapAdd(ctx, angleCh, subject, soundName)
}

func updateAngleBySocketConn(ctx context.Context, sigCh chan os.Signal, angleCh chan int, port int) {
	log.Printf("udp server is runnning on port: %v\n", port)
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("udp error:", err)
		sigCh <- os.Interrupt
	}
	defer conn.Close()
	buffer := make([]byte, 1024)

	const dataByteSize = 8

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Println(err)
				continue
			}
			length, _, err := conn.ReadFrom(buffer)
			if err != nil {
				log.Println("connection read error: ", err)
				continue
			}
			for i := 0; i < length/dataByteSize; i++ {
				angle, err := spat.BytesToInt64(buffer[dataByteSize*i : dataByteSize*(i+1)])
				if err != nil {
					log.Println("type conversion error:")
				}
				angleCh <- int(angle)
			}
		}
	}

}

var SLTFs [][][]float32
var dataAngle = 0
var sound []float32
var bufferLong [][]float32
var soundSLTF [][]float32

func OverlapAdd(ctx context.Context, angleCh chan int, subject, soundName string) error {
	// サンプリング周波数 [sample/sec]

	// 0.1度動くのに必要なサンプル数
	// [sec]*[sample/sec] / [0.1deg] = [sample/0.1deg]
	//var moveSamplesPerDeg int = moveSamples / moveWidth

	// 音データの読み込み
	//sound, err := spat.ReadDXXFile(soundName)
	//if err != nil {
	//	return err
	//}
	//sound = spat.Float64sToFloat32s(spat.RandNoise(FramesPerBuffer))

	channels := []string{"L", "R"}
	SLTFs = make([][][]float32, len(channels))
	for i, channel := range channels {
		SLTFs[i] = make([][]float32, AngleNum)
		for j := 0; j < AngleNum; j++ {
			SLTFName := fmt.Sprintf("%s/SLTF/SLTF_%d_%s.DDB", subject, j, channel)
			SLTF, err := spat.ReadDXXFile(SLTFName)
			SLTFF32s := spat.Float64sToFloat32s(SLTF)
			spat.NormFloat32s(SLTFF32s)
			SLTFs[i][j] = SLTFF32s
			if err != nil {
				return err
			}
		}
	}

	err := portaudio.Initialize()
	if err != nil {
		return err
	}
	defer portaudio.Terminate()
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	//outputDevice, err := portaudio.DefaultOutputDevice()
	//if err != nil {
	//	return err
	//}
	outputDevice := h.DefaultOutputDevice

	estimatedFramesPerBuffer := 2048
	sound = spat.Float64sToFloat32s(spat.RandNoise(estimatedFramesPerBuffer))

	streamParameters := portaudio.StreamParameters{
		Output: portaudio.StreamDeviceParameters{
			Device:   outputDevice,
			Channels: 2,
			Latency:  outputDevice.DefaultHighOutputLatency,
		},
		SampleRate: SamplingFreq,
		//FramesPerBuffer: FramesPerBuffer,
		FramesPerBuffer: FramesPerBuffer,
		//Flags:           portaudio.NoFlag,
		Flags: portaudio.ClipOff | portaudio.DitherOff,
	}
	//streamParameters := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)
	log.Println("streamParameters.Output.Channels:", streamParameters.Output.Channels)

	log.Println(streamParameters.Output.Device.Name)
	log.Println(streamParameters.Output.Device.MaxOutputChannels)
	log.Println(streamParameters.Output.Device.DefaultSampleRate)
	log.Printf("%#v\n", streamParameters.Output.Device.HostApi)
	log.Println(streamParameters.Output.Device.DefaultLowInputLatency)
	log.Println(streamParameters.Output.Device.DefaultHighInputLatency)

	//buffer := make([][]float32, len(channels))
	//for i := 0; i < len(channels); i++ {
	//	buffer[i] = make([]float32, FramesPerBuffer)
	//}
	//
	//bufferLong := make([][]float32, FramesPerBuffer+len(SLTFs[0][0]), (FramesPerBuffer+len(SLTFs[0][0]))*10)
	//for i := 0; i < len(channels); i++ {
	//	bufferLong[i] = make([]float32, FramesPerBuffer+len(SLTFs[0][0]), (FramesPerBuffer+len(SLTFs[0][0]))*10)
	//}
	bufferLong = make([][]float32, len(channels))
	soundSLTF = make([][]float32, len(channels))
	for i := 0; i < len(channels); i++ {
		//bufferLong[i] = make([]float32, FramesPerBuffer+len(SLTFs[0][0])-1, (FramesPerBuffer+len(SLTFs[0][0])-1)*10)
		//soundSLTF[i] = make([]float32, FramesPerBuffer+len(SLTFs[0][0])-1, (FramesPerBuffer+len(SLTFs[0][0])-1)*10)
		bufferLong[i] = make([]float32, estimatedFramesPerBuffer+len(SLTFs[0][0])-1)
		soundSLTF[i] = make([]float32, estimatedFramesPerBuffer+len(SLTFs[0][0])-1)
	}

	//buffer := make([]float32, len(channels)*FramesPerBuffer)
	//bufferLong := make([]float32, len(channels)*(FramesPerBuffer+len(SLTFs[0][0])), len(channels)*(FramesPerBuffer+len(SLTFs[0][0]))*10)

	//stream, err := portaudio.OpenStream(streamParameters, buffer)
	stream, err := portaudio.OpenStream(streamParameters, callback)
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.Stop()
	defer stream.Close()

	log.Println("playing...")
	err = stream.Start()
	if err != nil {
		return err
	}

	//dataAngle := 0
	for {
		select {
		case <-ctx.Done():
			if err != nil {
				return err
			}
			log.Println("stop stream")
			log.Println("done portaudio")
			return nil

		case dataAngle = <-angleCh:

		default:
			//for i := range channels {
			//	SLTF := SLTFs[i][dataAngle]
			////	音データと伝達関数の畳込み
			//soundSLTF := spat.LinearConvolutionTimeDomainFloat32s(sound, SLTF)

			// Overlap-Add
			//for j, v := range soundSLTF {
			//	bufferLong[i][j] += v
			//}
			//for j, v := range soundSLTF {
			//	bufferLong[len(channels)*j+i] += v
			//}
			//}
			//copy(buffer, bufferLong[:FramesPerBuffer*len(channels)])
			//shift(bufferLong, FramesPerBuffer*len(channels))
			//log.Println(bufferLong)
			//log.Println(len(bufferLong))
			//log.Println(buffer)

			//err = stream.Write()
			//if err != nil {
			//	return err
			//}
		}
	}
}

func callback(out [][]float32) {
	var err error
	fmt.Println(len(soundSLTF[0]), len(sound), )
	for i := range soundSLTF {
		soundSLTF[i], err = spat.LinearConvolutionTimeDomainFloat32s(soundSLTF[i], sound, SLTFs[i][dataAngle])
		spat.NormFloat32s(soundSLTF[i])
		if err != nil {
			log.Println(err)
		}
	}
	for i, soundSLTFLR := range soundSLTF {
		for j, v := range soundSLTFLR {
			bufferLong[i][j] += v
		}
	}
	//for i := range out[0] {
	//	bufferLong[0][i] += sound[i]
	//	bufferLong[1][i] += sound[i]
	//}
	for i, o := range out {
		for j := range o {
			out[i][j] = bufferLong[i][j]
		}
	}
	for i := range bufferLong {
		shift(bufferLong[i], len(out[i]))
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
