package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tetsuzawa/spat"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Printf("Usage of %s:\n", os.Args[0])
		log.Printf("overlap-add subject sound_file(.DXX) move_width move_velocity end_angle outdir\n")
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
	if flag.NArg() != 6 {
		return errors.New("invalid arguments")
	}
	args := flag.Args()
	subject := args[0]
	soundName := args[1]
	moveWidth, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	moveVelocity, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}
	endAngle, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}
	outDir := args[5]
	return OverlapAdd(subject, soundName, moveWidth, moveVelocity, endAngle, outDir)
}

func OverlapAdd(subject, soundName string, moveWidth, moveVelocity, endAngle int, outDir string) error {
	// サンプリング周波数 [sample/sec]
	const samplingFreq = 48000
	// 移動時間 [sec]
	var moveTime float64 = float64(moveWidth) / float64(moveVelocity)
	// 移動時間 [sample]
	var moveSamples int = int(moveTime * samplingFreq)

	// 0.1度動くのに必要なサンプル数
	// [sec]*[sample/sec] / [0.1deg] = [sample/0.1deg]
	var moveSamplesPerDeg int = moveSamples / moveWidth

	// 音データの読み込み
	sound, err := spat.ReadDXXFile(soundName)
	if err != nil {
		return err
	}

	SLTFName := fmt.Sprintf("%s/SLTF/SLTF_%d_%s.DDB", subject, 0, "L")
	SLTF, err := spat.ReadDXXFile(SLTFName)
	if err != nil {
		return err
	}

	for _, direction := range []string{"c", "cc"} {
		for _, LR := range []string{"L", "R"} {
			moveOut := make([]float64, moveSamples+len(SLTF)-1)
			usedAngles := make([]int, moveWidth)

			for angle := 0; angle < moveWidth; angle++ {
				// 畳み込むSLTFの角度を決定
				dataAngle := angle % (moveWidth * 2)
				if dataAngle > moveWidth {
					dataAngle = moveWidth*2 - dataAngle
				}
				if direction == "cc" {
					dataAngle = -dataAngle
				}
				dataAngle = dataAngle
				if dataAngle < 0 {
					dataAngle += 3600
				}
				// 使用した角度を記録（ログ出力用）
				usedAngles[angle] = (endAngle + dataAngle) % 3600

				// SLTFの読み込み
				SLTFName := fmt.Sprintf("%s/SLTF/SLTF_%d_%s.DDB", subject, (endAngle+dataAngle)%3600, LR)
				SLTF, err := spat.ReadDXXFile(SLTFName)
				if err != nil {
					return err
				}

				// 音データと伝達関数の畳込み
				cutSound := sound[moveSamplesPerDeg*angle : moveSamplesPerDeg*(angle+1)]
				soundSLTF := spat.LinearConvolutionTimeDomain(cutSound, SLTF)
				// Overlap-Add
				for i, v := range soundSLTF {
					moveOut[moveSamplesPerDeg*angle+i] += v
				}
			}

			// DDBへ出力
			outName := fmt.Sprintf("%s/move_judge_w%03d_mt%03d_%s_%d_%s.DDB", outDir, moveWidth, moveVelocity, direction, endAngle, LR)
			if err := spat.WriteDXXFile(outName, moveOut); err != nil {
				return err
			}
			_, err = fmt.Fprintf(os.Stderr, "%s: length=%d\n", outName, len(moveOut))
			if err != nil {
				return err
			}
			_, err := fmt.Fprintf(os.Stderr, "used angle:%v\n", usedAngles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
