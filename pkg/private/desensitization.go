package private

import "regexp"

func InsensitiveMobile(mobile string) string {
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return mobile
	}
	return mobile[:3] + "****" + mobile[7:]
}