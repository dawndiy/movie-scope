package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"strings"
)

// 接口地址
const (
	APP_KEY            string = "GHFmaeFG2Ma2ryL2N614YbdD"
	API_BASE           string = "http://api.map.baidu.com/telematics/v3/movie"
	API_HOT_MOVIE      string = API_BASE + "?qt=hot_movie"
	API_SEARCH_MOVIE   string = API_BASE + "?qt=search_movie"
	API_SEARCH_CINEMA  string = API_BASE + "?qt=search_cinema"
	API_NEARBY_CINEMA  string = API_BASE + "?qt=nearby_cinema"
	API_STATIC_IMAGE   string = "http://api.map.baidu.com/staticimage?"
	API_PANORAMA_IMAGE string = "http://api.map.baidu.com/panorama?"
	API_GEO_CONV       string = "http://api.map.baidu.com/geoconv/v1/?"
)

// 构建URL
func makeURL(api string, kwargs map[string]interface{}) string {

	params := []string{api}

	for k, v := range kwargs {
		params = append(params, k+"="+fmt.Sprint(v))
	}

	return strings.Join(params, "&")
}

func GetAPIData(api string, kwargs map[string]interface{}) (*simplejson.Json, error) {

	params := []string{api}
	kwargs["ak"] = APP_KEY
	kwargs["output"] = "json"
	for k, v := range kwargs {
		params = append(params, k+"="+fmt.Sprint(v))
	}
	url := strings.Join(params, "&")
	log.Println("URL: ", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := simplejson.NewFromReader(res.Body)
	return data, err
}

// 静态图
func APIStaticImage(kwargs map[string]interface{}) string {

	url := makeURL(API_STATIC_IMAGE, kwargs)
	return url
}

// 街景图
func APIPanoramaImage(kwargs map[string]interface{}) string {

	url := makeURL(API_PANORAMA_IMAGE, kwargs)
	return url
}

// 坐标转换
func APIGeoconv(kwargs map[string]interface{}) (float64, float64) {

	url := makeURL(API_GEO_CONV, kwargs)

	res, err := http.Get(url)
	if err != nil {
		return 0, 0
	}
	defer res.Body.Close()

	data, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return 0, 0
	}

	lng := data.Get("result").GetIndex(0).Get("x").MustFloat64()
	lat := data.Get("result").GetIndex(0).Get("y").MustFloat64()

	return lng, lat

}
