package merge_ip2region

import (
	"github.com/orestonce/Ip2regionTool"
)

func ConvertDbToTxt(req Ip2regionTool.ConvertDbReq) (errMsg string) {
	return Ip2regionTool.ConvertDb(req)
}
