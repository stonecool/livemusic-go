package geekbang

// geekBangRespDataInfoPrice struct used in geekbang
type geekBangRespDataInfoPrice struct {
	Market     int `json:"market"`
	Sale       int `json:"sale"`
	Sale_type  int `json:"sale_type"`
	Start_time int `json:"start_time"`
	End_time   int `json:"end_time"`
}

type geekBangRespDataInfoAuthor struct {
	name       string
	intro      string
	info       string
	avatar     string
	brief_html string
	brief      string
}

type geekBangRespDataInfoCover struct {
	square             string
	rectangle          string
	horizontal         string
	lecture_horizontal string
	learn_horizontal   string
	transparent        string
	color              string
}

type geekBangRespDataInfoArticle struct {
	id                  int
	count               int
	count_req           int
	count_pub           int
	total_length        int
	first_article_id    int
	first_article_title string
}

type geekBangRespDataInfoSeo struct {
	keywords []string
}

type geekBangRespDataInfoShare struct {
	title   string
	content string
	cover   string
	poster  string
}

type geekBangRespDataInfo struct {
	ID             int `json:"id"`
	Spu            int `json:"spu"`
	ctime          int
	utime          int
	begin_tme      int
	end_time       int
	Price          geekBangRespDataInfoPrice `json:"price"`
	is_onboard     bool
	is_sale        bool
	is_groupbuy    bool
	is_promo       bool
	is_shareget    bool
	is_sharesale   bool
	Info_type      string `json:"type"`
	is_column      bool
	is_core        bool
	is_video       bool
	is_audio       bool
	is_dailylesson bool
	is_university  bool
	is_opencourse  bool
	is_qconp       bool
	nav_id         int
	time_not_sale  int
	title          string
	subtitle       string
	intro          string
	intro_html     string
	ucode          string
	is_finish      bool
	Author         geekBangRespDataInfoAuthor  `json:"author"`
	Cover          geekBangRespDataInfoCover   `json:"cover"`
	Article        geekBangRespDataInfoArticle `json:"article"`
	Seo            geekBangRespDataInfoSeo     `json:"seo"`
	Share          geekBangRespDataInfoShare   `json:"share"`
	Labels         []int                       `json:"labels"`
}

type geekBangRespData struct {
	Infos []geekBangRespDataInfo `json:"infos"`
}

type GeekBangResp struct {
	Code int              `json:"code"`
	Data geekBangRespData `json:"data"`
}
