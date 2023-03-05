package information

import (
	"byte_dance_5th/pkg/config"
	"fmt"
	"math/rand"
	"strconv"
)

const BackgroundCnt int = 6

func GetBackGroundUrl() string {
	i := rand.Intn(100)
	fileName := strconv.Itoa(i%BackgroundCnt) + ".jpg"
	var url string = fmt.Sprintf("http://%s:%d/static/background/%s", config.Conf.SE.IP, config.Conf.SE.Port, fileName)
	return url
}
