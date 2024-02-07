package media

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mediahub/internal/utils"
	"regexp"
	"strconv"
	"strings"
)

var (
	NameSeWords     = regexp.MustCompile(`["共", "第", "季", "集", "话", "話", "期"]`)
	DigitRe         = regexp.MustCompile(`^[0-9]+$`)
	NameNoBeginRe   = regexp.MustCompile(`^\[.+?]`)
	NameNoChineseRe = regexp.MustCompile(`.*版|.*字幕`)
	NameNoStringRe  = regexp.MustCompile(`^PTS|^JADE|^AOD|^CHC|^[A-Z]{1,4}TV[\-0-9UVHDK]*` +
		`|HBO$|\s+HBO|\d{1,2}th|\d{1,2}bit|NETFLIX|AMAZON|IMAX|^3D|\s+3D|^BBC\s+|\s+BBC|BBC$|DISNEY\+?|XXX|\s+DC$` +
		`|[第\s共]+[0-9一二三四五六七八九十\-\s]+季` +
		`|[第\s共]+[0-9一二三四五六七八九十百零\-\s]+[集话話]` +
		`|连载|日剧|美剧|电视剧|动画片|动漫|欧美|西德|日韩|超高清|高清|蓝光|翡翠台|梦幻天堂·龙网|★?\d*月?新番` +
		`|最终季|合集|[多中国英葡法俄日韩德意西印泰台港粤双文语简繁体特效内封官译外挂]+字幕|版本|出品|台版|港版|\w+字幕组` +
		`|未删减版|UNCUT$|UNRATE$|WITH EXTRAS$|RERIP$|SUBBED$|PROPER$|REPACK$|SEASON$|EPISODE$|Complete$|Extended$|Extended Version$` +
		`|S\d{2}\s*-\s*S\d{2}|S\d{2}|\s+S\d{1,2}|EP?\d{2,4}\s*-\s*EP?\d{2,4}|EP?\d{2,4}|\s+EP?\d{1,4}` +
		`|CD[\s.]*[1-9]|DVD[\s.]*[1-9]|DISK[\s.]*[1-9]|DISC[\s.]*[1-9]` +
		`|[248]K|\d{3,4}[PIX]+` +
		`|CD[\s.]*[1-9]|DVD[\s.]*[1-9]|DISK[\s.]*[1-9]|DISC[\s.]*[1-9]`)
	NameStripYear         = regexp.MustCompile(`([\s.]+)(\d{4})-(\d{4})`)
	NameStripSize         = regexp.MustCompile(`(?i)[0-9.]+\s*[MGT]i?B([A-Z]*)`)
	NameStripYearMonthDay = regexp.MustCompile(`\d{4}[\s._-]\d{1,2}[\s._-]\d{1,2}`)
	RomanNumerals         = regexp.MustCompile(`^M*(C[MD]|D?C{0,3})(X[CL]|L?X{0,3})(I[XV]|V?I{0,3})$`)
	EpisodeRe             = regexp.MustCompile(`(?i)EP?(\d{2,4})$|^EP?(\d{1,4})$|^S\d{1,2}EP?(\d{1,4})$|S\d{2}EP?(\d{2,4})`)
	SeasonRe              = regexp.MustCompile(`(?i)S(\d{2})|^S(\d{1,2})$|S(\d{1,2})E`)

	SourceExp       = `(?i)^BLURAY$|^HDTV$|^UHDTV$|^HDDVD$|^WEBRIP$|^DVDRIP$|^BDRIP$|^BLU$|^WEB$|^BD$|^HDRip$`
	EffectExp       = `^REMUX$|^UHD$|^SDR$|^HDR\d*$|^DOLBY$|^DOVI$|^DV$|^3D$|^REPACK$`
	SourceRe        = regexp.MustCompile(SourceExp)
	EffectRe        = regexp.MustCompile(EffectExp)
	ResourcesTypeRe = regexp.MustCompile(fmt.Sprintf("(?i)%s|%s", SourceExp, SourceExp))
	ResourcesPixRe  = regexp.MustCompile(`(?i)^[SBUHD]*(\d{3,4}[PI]+)|\d{3,4}X(\d{3,4})`)
	ResourcesPixRe2 = regexp.MustCompile(`(?i)(^[248]+K)`)
	VideoEncodeRe   = regexp.MustCompile(`(?i)^[HX]26[45]$|^AVC$|^HEVC$|^VC\d?$|^MPEG\d?$|^Xvid$|^DivX$|^HDR\d*$`)
	AudioEncodeRe   = regexp.MustCompile(`(?i)^DTS\d?$|^DTSHD$|^DTSHDMA$|^Atmos$|^TrueHD\d?$|^AC3$|^\dAudios?$|^DDP\d?$|^DD\d?$|^LPCM\d?$|^AAC\d?$|^FLAC\d?$|^HD\d?$|^MA\d?$`)
	PartRe          = regexp.MustCompile(`(?i)(^PART[0-9ABI]{0,2}$|^CD[0-9]{0,2}$|^DVD[0-9]{0,2}$|^DISK[0-9]{0,2}$|^DISC[0-9]{0,2}$)`)
)

type MetaInfo interface {
	GetMeta() *Meta
	GetName() string
}

type Meta struct {
	OrgString      string   //原字符串
	RevString      string   // 识别词处理后字符串
	OrgTitle       string   // 原标题
	Title          string   // 媒体标题
	Subtitle       string   // 副标题
	MediaType      int      // 类型 电影、电视剧
	CnName         string   // 中文名
	EnName         string   // 英文名
	TotalSeasons   int      // 总季数
	BeginSeason    int      // 识别的开始季 数字
	EndSeason      int      // 识别的结束季 数字
	TotalEpisodes  int      // 总集数
	BeginEpisode   int      // 识别的开始集
	EndEpisode     int      // 识别的结束集
	Category       string   // 二级分类
	TmdbId         int      // TMDB ID
	ImdbId         string   // IMDB ID
	TvdbId         int      // TVDB ID
	DoubanId       int      // 豆瓣 ID
	Keyword        []string // 自定义搜索词
	ReleaseDate    string   // 媒体发行日期
	Runtime        int      // 播放时长
	Year           int      // 媒体年份
	ResourcePix    string   // 分辨率
	ResourceType   string   // 来源
	ResourceEffect string   // 特效
	VideoEncode    string   // 视频编码
	AudioEncode    string   // 音频编码
	Part           string
	ReplacedWords  []string // 识别辅助 替换词
	IgnoredWords   []string // 识别辅助 忽略词
	OffsetWords    []string // 识别辅助 集偏移词
	tokens         *utils.Tokenizer
	IsFile         bool
}

func (m *Meta) GetMeta() *Meta {
	return m
}

func (m *Meta) GetName() string {
	if m.CnName != "" && utils.IsChinese(m.CnName) {
		return m.CnName
	} else if m.EnName != "" {
		return m.EnName
	} else {
		return m.CnName
	}
}

func (m *Meta) GetTitle() string {
	if m.Title != "" {
		if m.Year != 0 {
			return fmt.Sprintf("%s (%d)", m.Title, m.Year)
		}
		return m.Title
	}
	name := m.GetName()
	if name != "" {
		if m.Year != 0 {
			return fmt.Sprintf("%s (%d)", name, m.Year)
		}
		return name
	}
	return ""
}

type MetaVideo struct {
	*Meta
}

func NewMetaVideo(title, subtitle string, isFile bool) *MetaVideo {
	self := &MetaVideo{
		Meta: &Meta{OrgTitle: title, Subtitle: subtitle, IsFile: isFile},
	}
	// 判断是否纯数字命名的文件
	name := utils.GetFileName(title)
	if IsMediaFile(title) && DigitRe.Match(utils.ToBytes(name)) && len(name) < 5 {
		self.Meta.BeginEpisode, _ = strconv.Atoi(name)
		self.Meta.MediaType = MediaTypeTv
		return self
	}
	// 去掉名称中第1个[]的内容
	title = NameNoBeginRe.ReplaceAllString(title, "")
	// 把xxxx-xxxx年份换成前一个年份，常出现在季集上
	title = NameStripYear.ReplaceAllString(title, "$1$2")
	// 把大小去掉
	title = NameStripSize.ReplaceAllString(title, "")
	// 把年月日去掉
	title = NameStripYearMonthDay.ReplaceAllString(title, "")
	self.parseTitle(title)
	return self
}

func (m *MetaVideo) parseTitle(title string) {
	tokens := utils.NewToken(title)
	p := NewParser(
		tokens,
		NewParsePart(m.GetMeta()),
		NewParseName(m.GetMeta()),
		NewParseYear(m.GetMeta()),
		NewParseResourcePix(m.GetMeta()),
		NewParseSeason(m.GetMeta()),
		NewParseEpisode(m.GetMeta()),
		NewParseResourceType(m.GetMeta()),
		NewParseVideoEncode(m.GetMeta()),
		NewParseAudioEncode(m.GetMeta()),
	)
	token, err := tokens.Cur()
	for err == nil {
		p.Run(token)
		token, err = tokens.Next()
	}
	if len(p.Effect) > 0 {
		m.ResourceEffect = strings.Join(p.Effect, " ")
	}
	if p.Source != "" {
		m.ResourceType = strings.TrimSpace(p.Source)
	}
	// 默认为电影
	if m.MediaType == MediaTypeUnknown {
		m.MediaType = MediaTypeMovie
	}
	if strings.ToUpper(m.Part) == "PART" {
		m.Part = ""
	}
}

func (m *MetaVideo) parsePart(token string) {

}

func (m *MetaVideo) parseName(token string) {

}

type MetaAnim struct {
	*Meta
}

func NewMetaAnim(title, subtitle string) *MetaAnim {
	return &MetaAnim{
		Meta: &Meta{Title: title, Subtitle: subtitle},
	}
}

func NewMeta(title, subtitle string, mediaType int, isFile bool) MetaInfo {
	if title == "" {
		return nil
	}
	orgTitle := title
	var meta MetaInfo
	var err error
	var info utils.ProcessInfo
	title, info, err = utils.ProcessTitle(title)
	if err != nil {
		log.Warnf("process title failed, %s", err.Error())
	}
	subtitle, _, _ = utils.ProcessTitle(subtitle)

	if mediaType == MediaAnim || IsAnim(title) {
		meta = NewMetaAnim(title, subtitle)
	} else {
		meta = NewMetaVideo(title, subtitle, isFile)
	}
	meta.GetMeta().OrgString = orgTitle
	meta.GetMeta().RevString = title
	meta.GetMeta().IgnoredWords = info.Ignored
	meta.GetMeta().ReplacedWords = info.Replaced
	meta.GetMeta().OffsetWords = info.Offset
	return meta
}

func IsAnim(title string) bool {
	return false
}
