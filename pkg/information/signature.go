package information

import "math/rand"

// 签名，由于没有给出签名上传上传接口，所以在此给出固定的几条签名信息用于随机选取
const (
	Signature0 string = "Golang是世界上最好的语言"
	Signature1 string = "C++好难学"
	Signature2 string = "Java好难学"
	Signature3 string = "C#好难学"
	Signature4 string = "python好难学"
	Signature5 string = "MySQL好难学"
	Signature6 string = "Redis好难学"
	Signature7 string = "Hertz好好用"
	Signature8 string = "青训营快乐出发"
)

const MaxSignatureCnt = 9

func GetSignature() string {
	idx := rand.Intn(MaxSignatureCnt)
	switch idx {
	case 0:
		return Signature0
	case 1:
		return Signature1
	case 2:
		return Signature2
	case 3:
		return Signature3
	case 4:
		return Signature4
	case 5:
		return Signature5
	case 6:
		return Signature6
	case 7:
		return Signature7
	case 8:
		return Signature8
	default:
		return "希望字节给我发个Offer"
	}
}
