package main

import (
	"fmt"
	"launchpad.net/go-unityscopes/v2"
	"log"
	"strconv"
	"strings"
	//"net/http"
)

const homeMovieTemplate = `{
	"schema-version": 1,
	"template": {
		"category-layout": "carousel",
		"card-layout": "vertical",
		"card-size": "large",
		"overlay": true
	},
	"components": {
		"title": "title",
		"art": {
			"field": "art",
			"aspect-ratio": 0.6
		},
		"subtitle": "subtitle"
	}
}`

const hotMovieTemplate = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-layout": "vertical",
		"card-size": "medium",
		"overlay": true
	},
	"components": {
		"title": "title",
		"art": {
			"field": "art",
			"aspect-ratio": 0.8
		},
		"subtitle": "subtitle"
	}
}`

const nearbyCinemaTemplate = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-layout": "vertical",
		"card-size": "medium",
		"overlay": true
	},
	"components": {
		"title": "title",
		"art": {
			"field": "art",
			"aspect-ratio": 1
		},
		"subtitle": "subtitle"
	}
}`

// Movie Scope
type MovieScope struct {
	base *scopes.ScopeBase
}

func (m *MovieScope) SetScopeBase(base *scopes.ScopeBase) {
	m.base = base
}

func (m *MovieScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {

	// 创建分类
	department := m.createDepartments(query, metadata, reply)
	reply.RegisterDepartments(department)

	dptId := query.DepartmentID()   // 分类ID
	queryStr := query.QueryString() // 搜索词

	if queryStr == "" {
		// 无搜索
		switch dptId {
		case "":
			m.showHome(metadata, reply)
			m.showNearbyCinemas(metadata, reply)
		case "hot_movie":
			m.showHotMovies(metadata, reply)
		case "nearby_cinema":
			m.showNearbyCinemas(metadata, reply)
		}
	} else {
		// 有搜索
		log.Println("[搜索]: ", queryStr)
	}

	return nil
}

func (m *MovieScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {

	resType := *new(string)
	result.Get("type", &resType)

	if resType == "movie" {
		// 查看电影
		m.viewHotMovie(result, reply)
	} else if resType == "cinema" {
		// 查看附近影院
		m.viewNearbyCinema(result, reply)
	}
	return nil
}

// func (m *MovieScope) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
// 	// handle the action and then tell the dash what to do next
// 	// through an ActivationResponse.
// 	log.Println("******************************************************")
// 	resp := scopes.NewActivationResponse(scopes.ActivationHideDash)
// 	return resp, nil
// }

// 建立分类菜单
func (m *MovieScope) createDepartments(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply) *scopes.Department {

	index, _ := scopes.NewDepartment("", query, "首页")

	hotMovie, _ := scopes.NewDepartment("hot_movie", query, "热映影片")
	if hotMovie != nil {
		index.AddSubdepartment(hotMovie)
	}

	nearbyCinema, _ := scopes.NewDepartment("nearby_cinema", query, "周边影院")
	if nearbyCinema != nil {
		index.AddSubdepartment(nearbyCinema)
	}

	return index
}

// 主页影片
func (m *MovieScope) showHome(metadata *scopes.SearchMetadata, reply *scopes.SearchReply) {

	location := metadata.Location()
	loc := "39.9289,116.3883"
	if location != nil {
		loc = geoConv(location.Longitude, location.Latitude, false)
	}
	log.Println("LOC: ", location, loc)

	args := map[string]interface{}{
		"location": loc,
	}

	json, err := GetAPIData(API_HOT_MOVIE, args)
	var movies []interface{}
	if err == nil {
		movies = json.GetPath("result", "movie").MustArray()
	}

	category := reply.RegisterCategory("hot_movie", "热映影片", "", homeMovieTemplate)

	for _, movie := range movies {

		movie := movie.(map[string]interface{})
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "movie")
		result.Set("map", movie)
		result.SetTitle(fmt.Sprint(movie["movie_name"]))
		result.SetArt(fmt.Sprint(movie["movie_picture"]))
		result.SetURI(fmt.Sprint(movie["movie_picture"]))
		result.Set("subtitle", movie["movie_release_date"])

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}

// 热映影片页面
func (m *MovieScope) showHotMovies(metadata *scopes.SearchMetadata, reply *scopes.SearchReply) {

	location := metadata.Location()
	loc := "39.9289,116.3883"
	if location != nil {
		loc = geoConv(location.Longitude, location.Latitude, false)
	}
	log.Println("LOC: ", location, loc)

	args := map[string]interface{}{
		"location": loc,
	}

	json, err := GetAPIData(API_HOT_MOVIE, args)
	var movies []interface{}
	if err == nil {
		movies = json.GetPath("result", "movie").MustArray()
	}

	category := reply.RegisterCategory("hot_movie", "热映影片", "", hotMovieTemplate)

	for _, movie := range movies {

		movie := movie.(map[string]interface{})
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "movie")
		result.Set("map", movie)
		result.SetTitle(fmt.Sprint(movie["movie_name"]))
		result.SetArt(fmt.Sprint(movie["movie_picture"]))
		result.SetURI(fmt.Sprint(movie["movie_picture"]))
		result.Set("subtitle", movie["movie_release_date"])

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}

// 周围影院页面
func (m *MovieScope) showNearbyCinemas(metadata *scopes.SearchMetadata, reply *scopes.SearchReply) {

	location := metadata.Location()
	loc := "39.9289,116.3883"
	if location != nil {
		loc = geoConv(location.Longitude, location.Latitude, false)
	}
	loc = "116.43229664960474,40.04485553445356"
	log.Println("LOC: ", location, loc)

	args := map[string]interface{}{
		"location": loc,
	}

	json, err := GetAPIData(API_NEARBY_CINEMA, args)
	var cinemas []interface{}
	if err == nil {
		cinemas = json.GetPath("result").MustArray()
	}

	category := reply.RegisterCategory("nearby_cinema", "周围影院", "", nearbyCinemaTemplate)

	for _, cinema := range cinemas {

		cinema := cinema.(map[string]interface{})
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "cinema")
		result.Set("map", cinema)
		result.SetTitle(fmt.Sprint(cinema["name"]))
		result.SetArt(m.base.ScopeDirectory() + "/local.png")
		result.SetURI(m.base.ScopeDirectory() + "/local.png")
		result.Set("subtitle", cinema["rating"])

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}
func (m *MovieScope) viewHotMovie(result *scopes.Result, reply *scopes.PreviewReply) {

	var movie map[string]interface{}
	err := result.Get("map", &movie)
	// log.Println(movie, err)
	log.Println(movie, err)

	layout1col := scopes.NewColumnLayout(1)
	layout1col.AddColumn(
		"header",
		"image",
		"text",
		"description",
	)
	reply.RegisterLayout(layout1col)

	image := scopes.NewPreviewWidget("image", "image")
	image.AddAttributeValue("source", movie["movie_big_picture"])
	image.AddAttributeValue("fallback", movie["movie_picture"])
	image.AddAttributeValue("zoomable", true)

	header := scopes.NewPreviewWidget("header", "header")
	header.AddAttributeValue("title", "<b>"+fmt.Sprint(movie["movie_name"])+"</b>")

	text := scopes.NewPreviewWidget("text", "text")
	text.AddAttributeValue("title", "信息")
	info := strings.Join([]string{
		"评分: \t" + fmt.Sprint(movie["movie_score"]),
		"类型: \t" + fmt.Sprint(movie["movie_type"]),
		"分类: \t" + fmt.Sprint(movie["movie_tags"]),
		"时长: \t" + fmt.Sprint(movie["movie_length"]) + " 分钟",
		"上映: \t" + fmt.Sprint(movie["movie_release_date"]),
		"地区: \t" + fmt.Sprint(movie["movie_nation"]),
		"导演: \t" + fmt.Sprint(movie["movie_director"]),
		"演员: \t" + fmt.Sprint(movie["movie_starring"]),
	}, "\n")
	text.AddAttributeValue("text", info)

	description := scopes.NewPreviewWidget("description", "text")
	description.AddAttributeValue("title", "描述")
	description.AddAttributeValue("text", movie["movie_message"])

	reply.PushWidgets(
		header,
		image,
		text,
		description,
	)

}

func (m *MovieScope) viewNearbyCinema(result *scopes.Result, reply *scopes.PreviewReply) {

	var cinema map[string]interface{}
	err := result.Get("map", &cinema)
	// log.Println(cinema, err)
	log.Println(err)

	layout1col := scopes.NewColumnLayout(1)
	//layout2col := scopes.NewColumnLayout(2)
	//layout3col := scopes.NewColumnLayout(3)

	// Single column layout
	layout1col.AddColumn(
		"header",
		"mapImage",
		"info",
		"actions",
		"online",
		"gallery",
		"onMovies",
		"reviews")

	// Two column layout
	//layout2col.AddColumn("image")
	//layout2col.AddColumn("header", "summary", "actions")

	// Three cokumn layout
	//layout3col.AddColumn("image")
	//layout3col.AddColumn("header", "summary", "actions")
	//layout3col.AddColumn()

	reply.RegisterLayout(layout1col)

	// 标题头
	header := scopes.NewPreviewWidget("header", "header")
	header.AddAttributeValue("title", "<b>"+fmt.Sprint(cinema["name"])+"</b>")
	header.AddAttributeValue("subtitle", cinema["address"])

	// 地点坐标参数
	location := cinema["location"].(map[string]interface{})
	lng := location["lng"].(float64)
	lat := location["lat"].(float64)
	paramLoc := geoConv(lng, lat, true)
	param := apiOptions{
		"center":       paramLoc,
		"markers":      paramLoc,
		"markerStyles": "l,o,0xff0000",
		"zoom":         "17",
	}
	baiduLngLat := strings.Split(paramLoc, ",")
	baiduMap := fmt.Sprintf("http://api.map.baidu.com/marker?location=%s,%s&title=影院地址&content=%s&output=html&src=Movie-Scope", baiduLngLat[1], baiduLngLat[0], cinema["name"])

	// 静态地图
	staticImgUrl := APIStaticImage(param)
	param = apiOptions{
		"ak":       APP_KEY,
		"location": paramLoc,
		"width":    "1000",
		"fov":      "270",
	}
	// 全景地图
	panoramaImgUrl := APIPanoramaImage(param)

	// 地图图片
	mapImage := scopes.NewPreviewWidget("mapImage", "gallery")
	mapImage.AddAttributeValue("fallback", m.base.ScopeDirectory()+"/local.png")
	mapImage.AddAttributeValue("sources", []string{staticImgUrl, panoramaImgUrl})
	mapImage.AddAttributeValue("zoomable", true)

	info := scopes.NewPreviewWidget("info", "text")
	info.AddAttributeValue("title", "<b>信息</b>")
	info.AddAttributeValue(
		"text",
		"<b>电话: </b>"+fmt.Sprint(cinema["telephone"])+"<br/>"+
			"<b>距离: </b>"+fmt.Sprint(cinema["distance"]))

	// 按钮
	actions := scopes.NewPreviewWidget("actions", "actions")
	tuple1 := make(map[string]interface{})
	tuple1["id"] = "tel"
	tuple1["label"] = "电话咨询"
	tuple1["uri"] = "tel:///" + fmt.Sprint(cinema["telephone"])
	tuple2 := make(map[string]interface{})
	tuple2["id"] = "map"
	tuple2["label"] = "地图"
	tuple2["uri"] = baiduMap
	actions.AddAttributeValue("actions", []interface{}{tuple1, tuple2})

	// 上映的影片图片
	title := scopes.NewPreviewWidget("online", "text")
	title.AddAttributeValue("title", "<b>上映的影片</b>")
	movies := scopes.NewPreviewWidget("gallery", "gallery")
	array := []string{}
	for _, i := range cinema["movies"].([]interface{}) {
		i := i.(map[string]interface{})
		array = append(array, fmt.Sprint(i["movie_picture"]))
	}
	movies.AddAttributeValue("sources", array)

	// // 影片场次
	// onMovies := scopes.NewPreviewWidget("onMovies", "expandable")
	// onMovies.AddAttributeValue("title", "<b>影片场次</b>")
	// onMovies.AddAttributeValue("collapsed-widgets", 2)
	// for i, m := range cinema.Movies {
	// 	text := scopes.NewPreviewWidget("movie"+strconv.Itoa(i), "text")
	// 	text.AddAttributeValue("title", "")
	// 	table := ""
	// 	for _, t := range m.TimeTable {
	// 		if t["type"] == "" {
	// 			t["type"] = "暂无数据"
	// 		}
	// 		if t["price"] == "" {
	// 			t["price"] = "暂无数据"
	// 		}
	// 		table = table + "<br/> 时间:" + t["date"] + " " + t["time"] + "<br/>" +
	// 			" 类型:" + t["lan"] + " " + t["type"] + "<br/>" +
	// 			" 票价:" + t["price"] + "<br/>"
	// 	}
	// 	text.AddAttributeValue(
	// 		"text",
	// 		"<b>「"+m.Name+"」</b><br/>"+
	// 			"<b>场次: </b>"+table)
	// 	onMovies.AddWidget(text)
	// }

	// // 影院评论
	// reviews := scopes.NewPreviewWidget("reviews", "expandable")
	// reviews.AddAttributeValue("title", "<b>影院评论</b>")
	// reviews.AddAttributeValue("collapsed-widgets", 1)
	// for i, r := range cinema.Review {
	// 	text := scopes.NewPreviewWidget("review"+strconv.Itoa(i), "text")
	// 	//text.AddAttributeValue("title", r["content"])
	// 	text.AddAttributeValue("title", "")
	// 	text.AddAttributeValue(
	// 		"text",
	// 		"<b>"+r["date"]+"</b><br/>"+
	// 			r["content"])
	// 	reviews.AddWidget(text)
	// }

	reply.PushWidgets(
		mapImage,
		header,
		info,
		title,
		movies,
		actions,
		//		onMovies,
		//		reviews,
	)
}

func geoConv(lng, lat float64, t bool) string {

	paramLoc := strconv.FormatFloat(lng, 'f', -1, 64) + "," + strconv.FormatFloat(lat, 'f', -1, 64)

	if t == false {
		return paramLoc
	}

	param := apiOptions{
		"ak":     APP_KEY,
		"coords": paramLoc,
		"from":   "3",
		"to":     "5",
	}

	newLng, newLat := APIGeoconv(param)
	paramLoc = strconv.FormatFloat(newLng, 'f', -1, 64) + "," + strconv.FormatFloat(newLat, 'f', -1, 64)

	return paramLoc
}

func Test() {
	param := map[string]string{
		"ak": APP_KEY,
		//"wd":       "末日崩塌",
		"radius":   "3000m",
		"location": "116.3883,39.9289",
		"output":   "json",
	}
	out := APINearbyCinema(param)
	log.Println(out)
}

func main() {
	if err := scopes.Run(&MovieScope{}); err != nil {
		log.Fatalln(err)
	}
	//Test()
}
