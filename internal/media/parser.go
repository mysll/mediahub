package media

import (
	"fmt"
	"mediahub/internal/utils"
	"regexp"
	"strconv"
	"strings"
)

const (
	ParseFlagCnName = 0x1
	ParseFlagName   = 0x2
	ParseFlagYear   = 0x4
)

const (
	ParseStatePart = iota + 1
	ParseStateCN
	ParseStateEN
	ParseSeasonEpisode
	ParseStateYear
	ParseStatePix
	ParseStateSeason
	ParseStateEpisode
	ParseStateSource
	ParseStateEffect
	ParseStateVideoEncode
	ParseStateAudioEncode
)

func CheckFlag(f1, f2 int) bool {
	return f1&f2 == f2
}

type Step interface {
	Run(ctx *Parser, token string) bool
	Complete() bool
	Done()
}

type ParseStep struct {
	done bool
	meta *Meta
}

func NewParseStep(meta *Meta) *ParseStep {
	return &ParseStep{
		meta: meta,
	}
}

func (p *ParseStep) Complete() bool {
	return p.done
}

func (p *ParseStep) Done() {
	p.done = true
}

type Parser struct {
	token      *utils.Tokenizer
	steps      []Step
	tokenState int
	flag       int
	unknown    string
	lastToken  string
	Source     string
	Effect     []string
}

func NewParser(token *utils.Tokenizer, step ...Step) *Parser {
	p := &Parser{
		token: token,
		steps: make([]Step, 0, len(step)),
	}

	p.steps = append(p.steps, step...)
	return p
}

func (p *Parser) Run(token string) {
	for _, s := range p.steps {
		if !s.Complete() {
			// 是否需要继续
			if !s.Run(p, token) {
				break
			}
		}
	}
}

func (p *Parser) CheckFlag(f int) bool {
	return CheckFlag(p.flag, f)
}

func (p *Parser) SetFlag(f int) {
	p.flag = p.flag | f
}

type ParseYear struct {
	*ParseStep
}

func NewParseYear(m *Meta) Step {
	return &ParseYear{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseYear) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}

	if !DigitRe.MatchString(token) || len(token) != 4 {
		return true
	}

	year, _ := strconv.Atoi(token)
	if year < 1900 || year > 2050 {
		return true
	}

	if p.meta.Year != 0 {
		if p.meta.EnName != "" {
			p.meta.EnName = fmt.Sprintf("%s %s", p.meta.EnName, token)
		} else if p.meta.CnName != "" {
			p.meta.CnName = fmt.Sprintf("%s %s", p.meta.CnName, token)
		}
	} else if ok, _ := regexp.MatchString(`(?i)SEASON$`, p.meta.EnName); ok && p.meta.EnName != "" {
		// 如果匹配到年，且英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
		p.meta.EnName += " "
	}
	p.meta.Year = year
	ctx.tokenState = ParseStateYear
	ctx.SetFlag(ParseFlagName)
	return false
}

type ParsePart struct {
	*ParseStep
}

func NewParsePart(m *Meta) Step {
	return &ParsePart{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParsePart) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}
	if p.meta.Year == 0 &&
		p.meta.ResourcePix == "" &&
		p.meta.ResourceType == "" &&
		p.meta.BeginSeason == 0 &&
		p.meta.BeginEpisode == 0 {
		return true
	}
	part := PartRe.FindString(token)
	if part != "" {
		if p.meta.Part == "" {
			p.meta.Part = part
		}
		if next, err := ctx.token.Peek(); err != nil {
			if DigitRe.MatchString(next) && (len(next) == 1 || len(next) == 2 && strings.HasPrefix(next, "0")) ||
				utils.InList(next, []string{"A", "B", "C", "I", "II", "III"}, true) {
				p.meta.Part = fmt.Sprintf("%s%s", p.meta.Part, next)
				// skip
				ctx.token.Next()
			}
		}

		ctx.tokenState = ParseStatePart
		return false
	}
	return true
}

type ParseName struct {
	*ParseStep
}

func NewParseName(m *Meta) Step {
	return &ParseName{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseName) Run(ctx *Parser, token string) bool {
	if ctx.unknown != "" {
		if p.meta.CnName == "" {
			if p.meta.EnName == "" {
				p.meta.EnName = ctx.unknown
			} else if ctx.unknown != strconv.Itoa(p.meta.Year) {
				p.meta.EnName = fmt.Sprintf("%s %s", p.meta.EnName, ctx.unknown)
			}
			ctx.tokenState = ParseStateEN
		}
		ctx.unknown = ""
	}
	if ctx.CheckFlag(ParseFlagName) {
		p.Done()
		return true
	}
	if strings.ToUpper(token) == "AKA" {
		p.Done()
		return false
	}
	if NameSeWords.MatchString(token) {
		ctx.tokenState = ParseSeasonEpisode
		return true
	}
	if utils.IsChinese(token) {
		ctx.tokenState = ParseStateCN
		if p.meta.CnName == "" {
			p.meta.CnName = token
		} else if !ctx.CheckFlag(ParseFlagCnName) {
			if !NameNoChineseRe.MatchString(token) && !NameSeWords.MatchString(token) {
				p.meta.CnName = fmt.Sprintf("%s %s", p.meta.CnName, token)
			}
			ctx.SetFlag(ParseFlagCnName)
		}
	} else {
		isRomanNumeral := RomanNumerals.MatchString(token)
		isNumber := DigitRe.MatchString(token)
		//阿拉伯数字或者罗马数字
		if isNumber || isRomanNumeral {
			// 第几集跳过
			if ctx.tokenState == ParseSeasonEpisode {
				return true
			}
			name := p.meta.GetName()
			if name != "" {
				// 名字后面以 0 开头的不要，极有可能是集
				if strings.HasPrefix(token, "0") {
					return true
				}
				// 中文名后面跟的数字不是年份的极有可能是集
				if isNumber {
					if ctx.tokenState == ParseStateCN {
						val, _ := strconv.Atoi(token)
						if val < 1900 {
							return true
						}
					}
				}
				if isNumber && len(token) < 4 || isRomanNumeral {
					// 4位以下的数字或者罗马数字，拼装到已有标题中
					if ctx.tokenState == ParseStateCN {
						p.meta.CnName = fmt.Sprintf("%s %s", p.meta.CnName, token)
					} else if ctx.tokenState == ParseStateEN {
						p.meta.EnName = fmt.Sprintf("%s %s", p.meta.EnName, token)
					}
					return false
				} else if isNumber && len(token) == 4 {
					// 4位数字，可能是年份，也可能真的是标题的一部分，也有可能是集
					if ctx.unknown == "" {
						ctx.unknown = token
					}
				}

			} else {
				// 名字未出现前的第一个数字，记下来
				if ctx.unknown == "" {
					ctx.unknown = token
				}
			}
		} else if NameSeWords.MatchString(token) {
			//# 如果匹配到季，英文名结尾为Season，说明Season属于标题，不应在后续作为干扰词去除
			if ok, _ := regexp.MatchString(`(?i)SEASON$`, token); ok && p.meta.EnName != "" {
				p.meta.EnName += " "
			}
			ctx.SetFlag(ParseFlagName)
		} else if EpisodeRe.MatchString(token) ||
			ResourcesTypeRe.MatchString(token) ||
			ResourcesPixRe.MatchString(token) {
			// 集、来源、版本等不要
			ctx.SetFlag(ParseFlagName)
		} else {
			// 后缀名不要
			if IsMediaFile(fmt.Sprintf("xxx.%s", token)) {
				return true
			}
			// 英文或者英文+数字，拼装起来
			if p.meta.EnName != "" {
				p.meta.EnName = fmt.Sprintf("%s %s", p.meta.EnName, token)
			} else {
				p.meta.EnName = token
			}
			ctx.tokenState = ParseStateEN
		}

	}
	return true
}

type ParseResourcePix struct {
	*ParseStep
}

func NewParseResourcePix(m *Meta) Step {
	return &ParseResourcePix{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseResourcePix) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}
	pixes := ResourcesPixRe.FindAllString(token, -1)
	if len(pixes) > 0 {
		for _, pix := range pixes {
			if pix != "" && p.meta.ResourcePix == "" {
				p.meta.ResourcePix = strings.ToLower(pix)
			}
		}
		if p.meta.ResourcePix != "" && DigitRe.MatchString(p.meta.ResourcePix) {
			p.meta.ResourcePix = fmt.Sprintf("%sp", p.meta.ResourcePix)
		}
		ctx.tokenState = ParseStatePix
		ctx.SetFlag(ParseFlagName)
		return false
	} else {
		pixes = ResourcesPixRe2.FindAllString(token, -1)
		if len(pixes) > 0 {
			if p.meta.ResourcePix == "" {
				p.meta.ResourcePix = strings.ToLower(pixes[0])
			}
			ctx.tokenState = ParseStatePix
			ctx.SetFlag(ParseFlagName)
			return false
		}
	}
	return true
}

type ParseSeason struct {
	*ParseStep
}

func NewParseSeason(m *Meta) Step {
	return &ParseSeason{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseSeason) Run(ctx *Parser, token string) bool {
	seasons := SeasonRe.FindStringSubmatch(token)
	if len(seasons) > 0 {
		for _, se := range seasons {
			if se == "" {
				continue
			}
			if DigitRe.MatchString(se) {
				se, _ := strconv.Atoi(se)
				if p.meta.BeginSeason == 0 {
					p.meta.BeginSeason = se
					p.meta.TotalSeasons = 1
				} else {
					if se > p.meta.BeginSeason {
						p.meta.EndSeason = se
						p.meta.TotalSeasons = p.meta.EndSeason - p.meta.BeginSeason + 1
						if p.meta.IsFile && p.meta.TotalSeasons > 1 {
							p.meta.EndSeason = 0
							p.meta.TotalSeasons = 1
						}
					}
				}
			}
		}
		p.meta.MediaType = MediaTypeTv
		ctx.SetFlag(ParseFlagName)
		ctx.tokenState = ParseStateSeason
		return true
	} else if DigitRe.MatchString(token) {
		if ctx.tokenState == ParseStateSeason && p.meta.EndSeason == 0 && len(token) < 3 {
			p.meta.BeginSeason, _ = strconv.Atoi(token)
			p.meta.TotalSeasons = 1
			p.meta.MediaType = MediaTypeTv
			ctx.SetFlag(ParseFlagName)
			ctx.tokenState = ParseStateSeason
			return false
		}
	} else if p.meta.EndSeason == 0 && strings.ToUpper(token) == "SEASON" {
		ctx.tokenState = ParseStateSeason
	}
	return true
}

type ParseEpisode struct {
	*ParseStep
}

func NewParseEpisode(m *Meta) Step {
	return &ParseEpisode{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseEpisode) Run(ctx *Parser, token string) bool {
	episodes := EpisodeRe.FindStringSubmatch(token)
	if len(episodes) > 0 {

		for _, ep := range episodes {
			if ep == "" {
				continue
			}
			if DigitRe.MatchString(ep) {
				se, _ := strconv.Atoi(ep)
				if p.meta.BeginEpisode == 0 {
					p.meta.BeginEpisode = se
					p.meta.TotalEpisodes = 1
				} else {
					if se > p.meta.BeginEpisode {
						p.meta.EndEpisode = se
						p.meta.TotalEpisodes = p.meta.EndEpisode - p.meta.BeginEpisode + 1
						if p.meta.IsFile && p.meta.TotalEpisodes > 2 {
							p.meta.EndEpisode = 0
							p.meta.TotalEpisodes = 1
						}
					}
				}
			}
		}
		p.meta.MediaType = MediaTypeTv
		ctx.SetFlag(ParseFlagName)
		ctx.tokenState = ParseStateEpisode
		return false
	} else if DigitRe.MatchString(token) {
		episode, _ := strconv.Atoi(token)
		if ctx.tokenState == ParseStateEpisode &&
			p.meta.BeginEpisode > 0 &&
			p.meta.EndEpisode == 0 &&
			len(token) < 5 &&
			episode > p.meta.BeginEpisode {
			p.meta.EndEpisode = episode
			p.meta.TotalEpisodes = p.meta.EndEpisode - p.meta.BeginEpisode + 1
			if p.meta.IsFile && p.meta.TotalEpisodes > 2 {
				p.meta.EndEpisode = 0
				p.meta.TotalEpisodes = 1
			}
			p.meta.MediaType = MediaTypeTv
			return false
		} else if p.meta.BeginEpisode == 0 &&
			episode > 1 && episode < 4 &&
			ctx.tokenState != ParseStateYear &&
			ctx.tokenState != ParseStateVideoEncode &&
			token != ctx.unknown {
			p.meta.BeginEpisode = episode
			p.meta.TotalEpisodes = 1
			ctx.tokenState = ParseStateEpisode
			ctx.SetFlag(ParseFlagName)
			p.meta.MediaType = MediaTypeTv
			return false
		} else if ctx.tokenState == ParseStateEpisode &&
			p.meta.BeginEpisode == 0 &&
			len(token) < 5 {
			p.meta.BeginEpisode = episode
			p.meta.TotalEpisodes = 1
			ctx.tokenState = ParseStateEpisode
			return false
		}
	} else if strings.ToUpper(token) == "EPISODE" {
		ctx.tokenState = ParseStateSeason
	}
	return true
}

type ParseResourceType struct {
	*ParseStep
}

func NewParseResourceType(m *Meta) Step {
	return &ParseResourceType{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseResourceType) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}
	source := SourceRe.FindString(token)
	if source != "" {
		ctx.tokenState = ParseStateSource
		if ctx.Source == "" {
			ctx.Source = source
			ctx.lastToken = strings.ToUpper(ctx.Source)
		}
		ctx.SetFlag(ParseFlagName)
		return false
	} else {
		switch strings.ToUpper(token) {
		case "DL":
			if ctx.tokenState == ParseStateSource && ctx.lastToken == "WEB" {
				ctx.Source = "WEB-DL"
				return false
			}
		case "RAY":
			if ctx.tokenState == ParseStateSource && ctx.lastToken == "BLU" {
				ctx.Source = "BluRay"
				return false
			}
		case "WEBDL":
			ctx.Source = "WEB-DL"
			return false
		}
	}
	effect := EffectRe.FindString(token)
	if effect != "" {
		ctx.tokenState = ParseStateEffect
		if ctx.Effect == nil {
			ctx.Effect = make([]string, 0, 4)
		}
		find := false
		for _, eff := range ctx.Effect {
			if eff == effect {
				find = true
			}
		}
		if !find {
			ctx.Effect = append(ctx.Effect, effect)
		}
		ctx.lastToken = effect
		ctx.SetFlag(ParseFlagName)
		return false
	}
	return true
}

type ParseVideoEncode struct {
	*ParseStep
}

func NewParseVideoEncode(m *Meta) Step {
	return &ParseVideoEncode{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseVideoEncode) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}
	if p.meta.Year == 0 &&
		p.meta.ResourcePix == "" &&
		p.meta.ResourceType == "" &&
		p.meta.BeginSeason == 0 &&
		p.meta.BeginEpisode == 0 {
		return true
	}
	encode := VideoEncodeRe.FindString(token)
	up := strings.ToUpper(token)
	if encode != "" {
		upEnc := strings.ToUpper(encode)
		ctx.tokenState = ParseStateVideoEncode
		if p.meta.VideoEncode == "" {
			p.meta.VideoEncode = upEnc
			ctx.lastToken = p.meta.VideoEncode
		} else if p.meta.VideoEncode == "10bit" {
			p.meta.VideoEncode = fmt.Sprintf("%s 10bit", upEnc)
			ctx.lastToken = upEnc
		}
		ctx.SetFlag(ParseFlagName)
		return false
	} else if len(up) == 1 && (up[0] == 'H' || up[0] == 'X') {
		ctx.tokenState = ParseStateVideoEncode
		if up[0] == 'H' {
			ctx.lastToken = "H"
		} else {
			ctx.lastToken = "x"
		}
		ctx.SetFlag(ParseFlagName)
		return false
	} else if up == "264" || up == "265" &&
		ctx.tokenState == ParseStateVideoEncode &&
		ctx.lastToken == "H" || ctx.lastToken == "x" {
		p.meta.VideoEncode = fmt.Sprintf("%s%s", ctx.lastToken, token)
	} else if DigitRe.MatchString(token) &&
		ctx.tokenState == ParseStateVideoEncode &&
		ctx.lastToken == "VC" || ctx.lastToken == "MPEG" {
		p.meta.VideoEncode = fmt.Sprintf("%s%s", ctx.lastToken, token)
	} else if strings.ToUpper(token) == "10BIT" {
		ctx.tokenState = ParseStateVideoEncode
		if p.meta.VideoEncode == "" {
			p.meta.VideoEncode = "10bit"
		} else {
			p.meta.VideoEncode = fmt.Sprintf("%s 10bit", p.meta.VideoEncode)
		}
	}
	return true
}

type ParseAudioEncode struct {
	*ParseStep
}

func NewParseAudioEncode(m *Meta) Step {
	return &ParseAudioEncode{
		ParseStep: NewParseStep(m),
	}
}

func (p *ParseAudioEncode) Run(ctx *Parser, token string) bool {
	if p.meta.GetName() == "" {
		return true
	}
	if p.meta.Year == 0 &&
		p.meta.ResourcePix == "" &&
		p.meta.ResourceType == "" &&
		p.meta.BeginSeason == 0 &&
		p.meta.BeginEpisode == 0 {
		return true
	}

	encode := AudioEncodeRe.FindString(token)
	if encode != "" {
		upEnc := strings.ToUpper(encode)
		if p.meta.AudioEncode == "" {
			p.meta.AudioEncode = encode
		} else {
			if strings.ToUpper(p.meta.AudioEncode) == "DTS" {
				p.meta.AudioEncode = fmt.Sprintf("%s-%s", p.meta.AudioEncode, encode)
			} else {
				p.meta.AudioEncode = fmt.Sprintf("%s %s", p.meta.AudioEncode, encode)
			}
		}
		ctx.tokenState = ParseStateAudioEncode
		ctx.lastToken = upEnc
		ctx.SetFlag(ParseFlagName)
		return false
	} else if DigitRe.MatchString(token) &&
		ctx.tokenState == ParseStateAudioEncode {
		if p.meta.AudioEncode != "" {
			n := len(p.meta.AudioEncode)
			if DigitRe.MatchString(ctx.lastToken) {
				p.meta.AudioEncode = fmt.Sprintf("%s.%s", p.meta.AudioEncode, token)
			} else if DigitRe.MatchString(p.meta.AudioEncode[n-1:]) {
				p.meta.AudioEncode = fmt.Sprintf("%s %s.%s", p.meta.AudioEncode[:n-1], p.meta.AudioEncode[n-1:], token)
			} else {
				p.meta.AudioEncode = fmt.Sprintf("%s %s", p.meta.AudioEncode, token)
			}
		}
		ctx.lastToken = token
	}
	return true
}
