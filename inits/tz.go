package inits

import "os"

func TZInit() {
	os.Setenv("TZ", "Africa/Nairobi")
}
