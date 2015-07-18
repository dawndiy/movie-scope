package main

import (
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

	// 地点, 失败时 location 可能为 nil
	location := metadata.Location()
	// 地点坐标参数
	paramLoc := "39.9289,116.3883"
	if location != nil {
		// paramLoc = strconv.FormatFloat(location.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(location.Latitude, 'f', -1, 64)
		paramLoc = geoConv(location.Longitude, location.Latitude, false)
		log.Println(paramLoc)
	}
	param := apiOptions{
		"ak":       APP_KEY,
		"location": paramLoc,
		"output":   "json",
	}

	dptId := query.DepartmentID()
	queryStr := query.QueryString()

	if queryStr == "" {
		// 无搜索
		switch dptId {
		case "":
			m.showHome(metadata, reply)
			nearbyCinemas := APINearbyCinema(param)
			m.showNearbyCinemas(nearbyCinemas, reply)
		case "hot_movie":
			movies := APIHotMovie(param)
			m.showHotMovies(movies, reply)
		case "nearby_cinema":
			nearbyCinemas := APINearbyCinema(param)
			m.showNearbyCinemas(nearbyCinemas, reply)
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

func (m *MovieScope) PerformAction(result *scopes.Result, metadata *scopes.ActionMetadata, widgetId, actionId string) (*scopes.ActivationResponse, error) {
	// handle the action and then tell the dash what to do next
	// through an ActivationResponse.
	log.Println("******************************************************")
	resp := scopes.NewActivationResponse(scopes.ActivationHideDash)
	return resp, nil
}

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

	// 地点, 失败时 location 可能为 nil
	location := metadata.Location()
	// 地点坐标参数
	paramLoc := "39.9289,116.3883"
	if location != nil {
		paramLoc = geoConv(location.Longitude, location.Latitude, false)
		log.Println(paramLoc)
	}
	param := apiOptions{
		"ak":       APP_KEY,
		"location": paramLoc,
		"output":   "json",
	}

	movies := APIHotMovie(param)
	category := reply.RegisterCategory("home_hot_movie", "热映影片", "", homeMovieTemplate)

	for _, movie := range movies {

		log.Println(movie.Name)
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "movie")
		result.Set("instance", movie)
		result.SetTitle(movie.Name)
		result.SetArt(movie.Picture)
		result.SetURI(movie.Picture)
		result.Set("subtitle", movie.ReleaseDate)

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}

// 热映影片页面
func (m *MovieScope) showHotMovies(movies []Movie, reply *scopes.SearchReply) {

	category := reply.RegisterCategory("hot_movie", "热映影片", "", hotMovieTemplate)

	for _, movie := range movies {

		log.Println(movie.Name)
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "movie")
		result.Set("instance", movie)
		result.SetTitle(movie.Name)
		result.SetArt(movie.Picture)
		result.SetURI(movie.Picture)
		result.Set("subtitle", movie.ReleaseDate)

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}

// 周围影院页面
func (m *MovieScope) showNearbyCinemas(nearbyCinemas []NearbyCinema, reply *scopes.SearchReply) {

	category := reply.RegisterCategory("nearby_cinema", "周围影院", "", nearbyCinemaTemplate)

	for _, cinema := range nearbyCinemas {

		log.Println(cinema.Name)
		result := scopes.NewCategorisedResult(category)

		result.Set("type", "cinema")
		result.Set("instance", cinema)
		result.SetTitle(cinema.Name)
		result.SetArt(m.base.ScopeDirectory() + "/local.png")
		result.SetURI(m.base.ScopeDirectory() + "/local.png")

		if err := reply.Push(result); err != nil {
			log.Println(err)
		}
	}
}

func (m *MovieScope) viewHotMovie(result *scopes.Result, reply *scopes.PreviewReply) {

	movie := *new(Movie)
	err := result.Get("instance", &movie)
	log.Println(movie, err)

	//opts := apiOptions{
	//	"ak":       APP_KEY,
	//	"location": "北京",
	//	"wd":       movie.Name,
	//	"output":   "json",
	//}
	//t := APISearchMovie(opts)
	//log.Println("==========")
	//log.Println(t)
	//log.Println("==========")

	layout1col := scopes.NewColumnLayout(1)
	//layout2col := scopes.NewColumnLayout(2)
	//layout3col := scopes.NewColumnLayout(3)

	// Single column layout
	layout1col.AddColumn("image", "header", "text", "description")

	// Two column layout
	//layout2col.AddColumn("image")
	//layout2col.AddColumn("header", "summary", "actions")

	// Three cokumn layout
	//layout3col.AddColumn("image")
	//layout3col.AddColumn("header", "summary", "actions")
	//layout3col.AddColumn()

	reply.RegisterLayout(layout1col)

	image := scopes.NewPreviewWidget("image", "image")
	image.AddAttributeValue("source", movie.Picture)
	image.AddAttributeValue("zoomable", true)

	header := scopes.NewPreviewWidget("header", "header")
	header.AddAttributeValue("title", movie.Name)
	header.AddAttributeValue("subtitle", movie.Tags)

	text := scopes.NewPreviewWidget("text", "text")
	text.AddAttributeValue("title", "信息")
	info := strings.Join([]string{
		"评分: " + movie.Score,
		"类型: " + movie.Type,
		"分类: " + movie.Tags,
		"时长: " + movie.Length + " 分钟",
		"上映时间: " + movie.ReleaseDate,
		"所属地区: " + movie.Nation,
		"导演: " + movie.Director,
		"演员: " + movie.Starring,
	}, "\n")
	text.AddAttributeValue("text", info)

	description := scopes.NewPreviewWidget("description", "text")
	description.AddAttributeValue("title", "描述")
	description.AddAttributeValue("text", movie.Message)

	reply.PushWidgets(image, header, text, description)

}

func (m *MovieScope) viewNearbyCinema(result *scopes.Result, reply *scopes.PreviewReply) {

	cinema := *new(NearbyCinema)
	err := result.Get("instance", &cinema)
	log.Println(cinema, err)

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
	header.AddAttributeValue("title", "<b>"+cinema.Name+"</b>")
	header.AddAttributeValue("subtitle", cinema.Address)

	// 地点, 失败时 location 可能为 nil
	// 地点坐标参数
	lng := cinema.Location["lng"]
	lat := cinema.Location["lat"]
	paramLoc := geoConv(lng, lat, true)
	param := apiOptions{
		"center":       paramLoc,
		"markers":      paramLoc,
		"markerStyles": "l,A,0xff0000",
	}
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
	log.Println("URL: ", staticImgUrl)
	log.Println("URL: ", panoramaImgUrl)

	// 地图图片
	mapImage := scopes.NewPreviewWidget("mapImage", "gallery")
	mapImage.AddAttributeValue("fallback", m.base.ScopeDirectory()+"/local.png")
	mapImage.AddAttributeValue("sources", []string{staticImgUrl, panoramaImgUrl})
	mapImage.AddAttributeValue("zoomable", true)

	info := scopes.NewPreviewWidget("info", "text")
	info.AddAttributeValue("title", "<b>信息</b>")
	info.AddAttributeValue(
		"text",
		"<b>电话: </b>"+cinema.Telephone+"<br/>"+
			"<b>距离: </b>"+cinema.Distance)

	// 按钮
	actions := scopes.NewPreviewWidget("actions", "actions")
	tuple1 := make(map[string]interface{})
	tuple1["id"] = "tel"
	tuple1["label"] = "电话咨询"
	tuple1["uri"] = "tel:///" + cinema.Telephone
	tuple2 := make(map[string]interface{})
	tuple2["id"] = "map"
	tuple2["label"] = "地图"
	tuple3 := make(map[string]interface{})
	tuple3["id"] = "three"
	tuple3["label"] = "THREE"
	actions.AddAttributeValue("actions", []interface{}{tuple1, tuple2, tuple3})

	// 上映的影片图片
	title := scopes.NewPreviewWidget("online", "text")
	title.AddAttributeValue("title", "<b>上映的影片</b>")
	movies := scopes.NewPreviewWidget("gallery", "gallery")
	array := []string{}
	for _, i := range cinema.Movies {
		array = append(array, i.Picture)
	}
	movies.AddAttributeValue("sources", array)

	// 影片场次
	onMovies := scopes.NewPreviewWidget("onMovies", "expandable")
	onMovies.AddAttributeValue("title", "<b>影片场次</b>")
	onMovies.AddAttributeValue("collapsed-widgets", 2)
	for i, m := range cinema.Movies {
		text := scopes.NewPreviewWidget("movie"+strconv.Itoa(i), "text")
		text.AddAttributeValue("title", "")
		table := ""
		for _, t := range m.TimeTable {
			if t["type"] == "" {
				t["type"] = "暂无数据"
			}
			if t["price"] == "" {
				t["price"] = "暂无数据"
			}
			table = table + "<br/> 时间:" + t["date"] + " " + t["time"] + "<br/>" +
				" 类型:" + t["lan"] + " " + t["type"] + "<br/>" +
				" 票价:" + t["price"] + "<br/>"
		}
		text.AddAttributeValue(
			"text",
			"<b>「"+m.Name+"」</b><br/>"+
				"<b>场次: </b>"+table)
		onMovies.AddWidget(text)
	}

	// 影院评论
	reviews := scopes.NewPreviewWidget("reviews", "expandable")
	reviews.AddAttributeValue("title", "<b>影院评论</b>")
	reviews.AddAttributeValue("collapsed-widgets", 1)
	for i, r := range cinema.Review {
		text := scopes.NewPreviewWidget("review"+strconv.Itoa(i), "text")
		//text.AddAttributeValue("title", r["content"])
		text.AddAttributeValue("title", "")
		text.AddAttributeValue(
			"text",
			"<b>"+r["date"]+"</b><br/>"+
				r["content"])
		reviews.AddWidget(text)
	}

	reply.PushWidgets(
		mapImage,
		header,
		info,
		title,
		movies,
		actions,
		onMovies,
		reviews)
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
