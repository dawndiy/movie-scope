package main

// API 参数
type apiOptions map[string]string

// 热映影片
type Movie struct {
	ID          string
	Name        string
	Type        string
	ReleaseDate string
	Nation      string
	Starring    string
	Length      string
	Picture     string
	Score       string
	Director    string
	Tags        string
	Message     string
	IsIMAX      string
	IsNew       string
	Word        string
}

//
type MovieCinema struct {
	ID        string
	Name      string
	Telephone string
	Location  map[string]float64
	Address   string
	Rating    string
	TimeTable []map[string]string
}

// 周围影院
type NearbyCinema struct {
	ID        string
	Name      string
	Telephone string
	Location  map[string]float64
	Address   string
	Distance  string
	Rating    string
	Review    []map[string]string
	Movies    []NearbyCinemaMovie
}

// 周围影院影片
type NearbyCinemaMovie struct {
	Movie
	Description string
	TimeTable   []map[string]string
}
