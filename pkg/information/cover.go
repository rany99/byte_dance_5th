//----------!!!Attention:请确保ffmpeg.exe已经置于GoPath路径下!!!----------

package information

import (
	"byte_dance_5th/pkg/errortype"
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
)

//videoPath 视频保存路径
//snapShotPath 截图保存路径
//frameNum 截图帧数

// SnapShotFromVideo 生成截图
func SnapShotFromVideo(videoPath, snapShotPath string, frameNum int) (err error) {

	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		//log.Fatal(errortype.SnapShotErr, err)
		return errors.New(errortype.SnapShotErr)
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		//log.Fatal(errortype.ImgDecodeErr, err)
		return errors.New(errortype.ImgDecodeErr)
	}

	err = imaging.Save(img, snapShotPath)
	if err != nil {
		//log.Fatal(errortype.SaveSnapErr, err)
		return errors.New(errortype.SaveSnapErr)
	}

	return nil
}
