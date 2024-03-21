package views

import (
	"fmt"
	"testing"
)

func TestXXX(t *testing.T) {
	status := getOutputDirectly("python3", "../scripts/save_geojson.py", "get_geo_json", "下载", "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/dataset/下载")
	fmt.Println(string(status))
}
