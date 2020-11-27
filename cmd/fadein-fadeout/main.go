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
		log.Printf("fadein-fadeout subject sound_file(.DXX) move_width move_velocity end_angle outdir\n")
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
	return FadeinFadeout(subject, soundName, moveWidth, moveVelocity, endAngle, outDir)
}

func FadeinFadeout(subject, soundName string, moveWidth, moveVelocity, endAngle int, outDir string) error {
	// サンプリング周波数 [sample/sec]
	const samplingFreq = 48000
	// 移動時間 [sec]
	var moveTime float64 = float64(moveWidth) / float64(moveVelocity)
	// 移動時間 [sample]
	var moveSamples int = int(moveTime * samplingFreq)

	// 0.1度動くのに必要なサンプル数
	// [sec]*[sample/sec] / [0.1deg] = [sample/0.1deg]
	var dwellingSamples int = int(moveTime * samplingFreq / float64(moveWidth))
	var durationSamples int = dwellingSamples * 63 / 64
	var overlapSamples int = dwellingSamples * 1 / 64

	fadeinFilter, fadeoutFilter := spat.GenerateFadeinFadeoutFilter(overlapSamples)

	// 音データの読み込み
	sound, err := spat.ReadFile(soundName)
	if err != nil {
		return err
	}

	for _, direction := range []string{"c", "cc"} {
		for _, LR := range []string{"L", "R"} {
			moveOut := make([]float64, dwellingSamples, moveSamples)
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

				// SLTFの読み込み（十分に長い白色雑音を想定）
				SLTFName := fmt.Sprintf("%s/SLTF/SLTF_%d_%s.DDB", subject, (endAngle+dataAngle)%3600, LR)
				SLTF, err := spat.ReadFile(SLTFName)
				if err != nil {
					return err
				}

				// Fadein-Fadeout
				// 音データと伝達関数の畳込み
				cutSound := sound[angle*(durationSamples+overlapSamples) : durationSamples*2+angle*(durationSamples+overlapSamples)+len(SLTF)*3+1]
				soundSLTF := spat.LinearConvolutionTimeDomain(cutSound, SLTF)
				// 無音区間の切り出し
				soundSLTF = soundSLTF[len(SLTF)*2 : len(soundSLTF)-len(SLTF)*2]
				// 前の角度のfadeout部と現在の角度のfadein部の加算
				fadein := make([]float64, overlapSamples)
				for i := range fadein {
					fadein[i] = soundSLTF[i] * fadeinFilter[i]
					moveOut[(durationSamples+overlapSamples)*angle+i] += fadein[i]
				}

				// 持続時間
				moveOut = append(moveOut, soundSLTF[overlapSamples:len(soundSLTF)-overlapSamples]...)

				// fadeout
				fadeout := make([]float64, overlapSamples)
				for i := range fadeout {
					fadeout[i] = soundSLTF[len(soundSLTF)-overlapSamples+i] * fadeoutFilter[i]
				}
				moveOut = append(moveOut, fadeout...)
			}

			// 先頭のFadein部をカット
			out := moveOut[overlapSamples:]

			// DDBへ出力
			outName := fmt.Sprintf("%s/move_judge_w%03d_mt%03d_%s_%d_%s.DDB", outDir, moveWidth, moveVelocity, direction, endAngle, LR)
			if err := spat.WriteFile(outName, out); err != nil {
				return err
			}
			_, err = fmt.Fprintf(os.Stderr, "%s: length=%d\n", outName, len(out))
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
