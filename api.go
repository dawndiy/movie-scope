package main

import (
	"encoding/json"
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
func makeURL(api string, opts apiOptions) string {

	params := []string{api}

	for k, v := range opts {
		params = append(params, k+"="+v)
	}

	return strings.Join(params, "&")
}

// 热映电影
func APIHotMovie(opts apiOptions) []Movie {

	movies := []Movie{}

	url := makeURL(API_HOT_MOVIE, opts)

	res, err := http.Get(url)
	if err != nil {
		return movies
	}
	defer res.Body.Close()

	data, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return movies
	}

	// 成功
	if data.Get("error").MustInt() == 0 {
		for _, m := range data.Get("result").Get("movie").MustArray() {
			m := m.(map[string]interface{})
			movie := Movie{
				ID:          m["movie_id"].(string),
				Name:        m["movie_name"].(string),
				Type:        m["movie_type"].(string),
				ReleaseDate: m["movie_release_date"].(string),
				Nation:      m["movie_nation"].(string),
				Starring:    m["movie_starring"].(string),
				Length:      m["movie_length"].(string),
				Picture:     m["movie_picture"].(string),
				Score:       m["movie_score"].(string),
				Director:    m["movie_director"].(string),
				Tags:        m["movie_tags"].(string),
				Message:     m["movie_message"].(string),
				IsIMAX:      m["is_imax"].(json.Number).String(),
				IsNew:       m["is_new"].(string),
				Word:        m["movies_wd"].(string),
			}
			movies = append(movies, movie)
		}
	}

	//log.Println(movies)
	return movies

}

// 影片影讯检索
func APISearchMovie(opts apiOptions) []MovieCinema {

	movieCinemas := []MovieCinema{}

	url := makeURL(API_SEARCH_MOVIE, opts)
	log.Println("URL: ", url)

	res, err := http.Get(url)
	if err != nil {
		return movieCinemas
	}

	defer res.Body.Close()

	data, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return movieCinemas
	}

	if data.Get("error").MustInt() == 0 {
		for _, item := range data.Get("result").MustArray() {
			item := item.(map[string]interface{})

			lng, _ := item["location"].(map[string]interface{})["lng"].(json.Number).Float64()
			lat, _ := item["location"].(map[string]interface{})["lat"].(json.Number).Float64()

			timeTable := []map[string]string{}
			for _, table := range item["time_table"].([]interface{}) {
				table := table.(map[string]interface{})
				timeTable = append(timeTable, map[string]string{
					"time":  table["time"].(string),
					"date":  table["date"].(string),
					"lan":   table["lan"].(string),
					"type":  table["type"].(string),
					"price": table["price"].(json.Number).String(),
				})
			}

			movieCinema := MovieCinema{
				ID:        item["uid"].(string),
				Name:      item["name"].(string),
				Telephone: item["telephone"].(string),
				Location: map[string]float64{
					"lng": lng,
					"lat": lat,
				},
				TimeTable: timeTable,
			}

			movieCinemas = append(movieCinemas, movieCinema)
		}
	}

	return movieCinemas

}

// 影院影讯检索
func APISearchCinema() {
}

//周边影院检索
func APINearbyCinema(opts apiOptions) []NearbyCinema {

	url := makeURL(API_NEARBY_CINEMA, opts)
	log.Println(url)

	nearbyCinemas := []NearbyCinema{}

	res, err := http.Get(url)
	if err != nil {
		return nearbyCinemas
	}

	defer res.Body.Close()

	data, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return nearbyCinemas
	}

	if data.Get("error").MustInt() == 0 {
		for _, item := range data.Get("result").MustArray() {
			item := item.(map[string]interface{})

			// 坐标
			lng, _ := item["location"].(map[string]interface{})["lng"].(json.Number).Float64()
			lat, _ := item["location"].(map[string]interface{})["lat"].(json.Number).Float64()

			// 评论
			reviews := []map[string]string{}
			if _, ok := item["review"].([]interface{}); ok {
				for _, review := range item["review"].([]interface{}) {
					review := review.(map[string]interface{})
					reviews = append(reviews, map[string]string{
						"content": review["content"].(string),
						"date":    review["date"].(string),
					})
				}
			}

			// 影片
			movies := []NearbyCinemaMovie{}
			onlineMovies, ok := item["movies"].([]interface{})
			if ok {
				for _, movie := range onlineMovies {
					movie := movie.(map[string]interface{})

					// 场次
					timeTable := []map[string]string{}
					for _, table := range movie["time_table"].([]interface{}) {
						table := table.(map[string]interface{})
						price, ok := table["price"].(string)
						if !ok {
							price = table["price"].(json.Number).String()
						}
						timeTable = append(timeTable, map[string]string{
							"time": table["time"].(string),
							"date": table["date"].(string),
							"lan":  table["lan"].(string),
							"type": table["type"].(string),
							//"price": table["price"].(json.Number).String(),
							"price": price,
						})
					}

					movies = append(movies, NearbyCinemaMovie{
						Movie: Movie{
							Name:        movie["movie_name"].(string),
							Type:        movie["movie_type"].(string),
							Nation:      movie["movie_nation"].(string),
							Director:    movie["movie_director"].(string),
							Starring:    movie["movie_starring"].(string),
							ReleaseDate: movie["movie_release_date"].(string),
							Picture:     movie["movie_picture"].(string),
							Length:      movie["movie_length"].(string),
						},
						Description: movie["movie_description"].(string),
						TimeTable:   timeTable,
					})

				}
			}

			nc := NearbyCinema{
				ID:        item["uid"].(string),
				Name:      item["name"].(string),
				Telephone: item["telephone"].(string),
				Location: map[string]float64{
					"lng": lng,
					"lat": lat,
				},
				Address:  item["address"].(string),
				Distance: item["distance"].(json.Number).String(),
				Rating:   item["rating"].(string),
				Review:   reviews,
				Movies:   movies,
			}
			nearbyCinemas = append(nearbyCinemas, nc)

		}
	}

	return nearbyCinemas

}

// 静态图
func APIStaticImage(opts apiOptions) string {

	url := makeURL(API_STATIC_IMAGE, opts)
	return url
}

func APIPanoramaImage(opts apiOptions) string {

	url := makeURL(API_PANORAMA_IMAGE, opts)
	return url
}

// 坐标转换
func APIGeoconv(opts apiOptions) (float64, float64) {

	url := makeURL(API_GEO_CONV, opts)

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
