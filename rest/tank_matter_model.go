package rest

type TankMatter struct {
	TankBase
	Puuid    string `json:"puuid"`
	UserUuid string `json:"userUuid"`
	Dir      bool   `json:"dir"`
	Alien    bool   `json:"alien"`
	Name     string `json:"name"`
	Md5      string `json:"md5"`
	Size     int64  `json:"size"`
	Privacy  bool   `json:"privacy"`
	Path     string `json:"path"`
	Url      string `json:"url"`
}
