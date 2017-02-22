package tokenizer

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
	"tktg/shared"
	"tktg/util"
	"unicode/utf8"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"github.com/acidd15/go-scala-util/src/scala"
)

type KoreanChunk struct {
	text   string
	offset int
	length int
}

type ChunkMatch struct {
	start int
	end   int
	text  string
	pos   util.KoreanPos
}

func (t ChunkMatch) disjoint(that interface{}) bool {
	v := that.(ChunkMatch)
	return (v.start < t.start && v.end <= t.start) ||
		(v.start >= t.end && v.end > t.end)
}

var posPatterns map[util.KoreanPos]pcre.Regexp

var chunkOrder []util.KoreanPos

func init() {
	gtLds := []string{"abogado", "academy", "accountants", "active", "actor", "adult", "aero", "agency", "airforce", "allfinanz", "alsace", "android", "aquarelle", "archi", "army", "arpa", "asia", "associates", "attorney", "auction", "audio", "autos", "axa", "band", "bar", "bargains", "bayern", "beer", "berlin", "best", "bid", "bike", "bio", "biz", "black", "blackfriday", "bloomberg", "blue", "bmw", "bnpparibas", "boo", "boutique", "brussels", "budapest", "build", "builders", "business", "buzz", "bzh", "cab", "cal", "camera", "camp", "cancerresearch", "capetown", "capital", "caravan", "cards", "care", "career", "careers", "cartier", "casa", "cash", "cat", "catering", "center", "ceo", "cern", "channel", "cheap", "christmas", "chrome", "church", "citic", "city", "claims", "cleaning", "click", "clinic", "clothing", "club", "coach", "codes", "coffee", "college", "cologne", "com", "community", "company", "computer", "condos", "construction", "consulting", "contractors", "cooking", "cool", "coop", "country", "credit", "creditcard", "cricket", "crs", "cruises", "cuisinella", "cymru", "dad", "dance", "dating", "day", "deals", "degree", "delivery", "democrat", "dental", "dentist", "desi", "diamonds", "diet", "digital", "direct", "directory", "discount", "dnp", "domains", "durban", "dvag", "eat", "edu", "education", "email", "emerck", "energy", "engineer", "engineering", "enterprises", "equipment", "esq", "estate", "eurovision", "eus", "events", "everbank", "exchange", "expert", "exposed", "fail", "farm", "fashion", "feedback", "finance", "financial", "firmdale", "fish", "fishing", "fitness", "flights", "florist", "flsmidth", "fly", "foo", "forsale", "foundation", "frl", "frogans", "fund", "furniture", "futbol", "gal", "gallery", "gbiz", "gent", "gift", "gifts", "gives", "glass", "gle", "global", "globo", "gmail", "gmo", "gmx", "google", "gop", "gov", "graphics", "gratis", "green", "gripe", "guide", "guitars", "guru", "hamburg", "haus", "healthcare", "help", "here", "hiphop", "hiv", "holdings", "holiday", "homes", "horse", "host", "hosting", "house", "how", "ibm", "immo", "immobilien", "industries", "info", "ing", "ink", "institute", "insure", "int", "international", "investments", "irish", "jetzt", "jobs", "joburg", "juegos", "kaufen", "kim", "kitchen", "kiwi", "koeln", "krd", "kred", "lacaixa", "land", "latrobe", "lawyer", "lds", "lease", "legal", "lgbt", "life", "lighting", "limited", "limo", "link", "loans", "london", "lotto", "ltda", "luxe", "luxury", "madrid", "maison", "management", "mango", "market", "marketing", "media", "meet", "melbourne", "meme", "memorial", "menu", "miami", "mil", "mini", "mobi", "moda", "moe", "monash", "money", "mormon", "mortgage", "moscow", "motorcycles", "mov", "museum", "nagoya", "name", "navy", "net", "network", "neustar", "new", "nexus", "ngo", "nhk", "ninja", "nra", "nrw", "nyc", "okinawa", "ong", "onion", "onl", "ooo", "org", "organic", "otsuka", "ovh", "paris", "partners", "parts", "party", "pharmacy", "photo", "photography", "photos", "physio", "pics", "pictures", "pink", "pizza", "place", "plumbing", "pohl", "poker", "porn", "post", "praxi", "press", "pro", "prod", "productions", "prof", "properties", "property", "pub", "qpon", "quebec", "realtor", "recipes", "red", "rehab", "reise", "reisen", "reit", "ren", "rentals", "repair", "report", "republican", "rest", "restaurant", "reviews", "rich", "rio", "rip", "rocks", "rodeo", "rsvp", "ruhr", "ryukyu", "saarland", "samsung", "sarl", "sca", "scb", "schmidt", "schule", "science", "scot", "services", "sexy", "shiksha", "shoes", "singles", "social", "software", "sohu", "solar", "solutions", "soy", "space", "spiegel", "supplies", "supply", "support", "surf", "surgery", "suzuki", "sydney", "systems", "taipei", "tatar", "tattoo", "tax", "technology", "tel", "tienda", "tips", "tirol", "today", "tokyo", "tools", "top", "town", "toys", "trade", "training", "travel", "trust", "tui", "university", "uno", "uol", "vacations", "vegas", "ventures", "vermögensberater", "vermögensberatung", "versicherung", "vet", "viajes", "villas", "vision", "vlaanderen", "vodka", "vote", "voting", "voto", "voyage", "wales", "wang", "watch", "webcam", "website", "wed", "wedding", "whoswho", "wien", "wiki", "williamhill", "wme", "work", "works", "world", "wtc", "wtf", "xxx", "xyz", "yachts", "yandex", "yoga", "yokohama", "youtube", "zip", "zone", "дети", "москва", "онлайн", "орг", "рус", "сайт", "بازار", "شبكة", "موقع", "संगठन", "みんな", "グーグル", "世界", "中信", "中文网", "企业", "佛山", "八卦", "公司", "公益", "商城", "商店", "商标", "在线", "广东", "我爱你", "手机", "政务", "机构", "游戏", "移动", "组织机构", "网址", "网店", "网络", "谷歌", "集团", "삼성"}
	ctLds := []string{"ac", "ad", "ae", "af", "ag", "ai", "al", "am", "an", "ao", "aq", "ar", "as", "at", "au", "aw", "ax", "az", "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi", "bj", "bl", "bm", "bn", "bo", "bq", "br", "bs", "bt", "bv", "bw", "by", "bz", "ca", "cc", "cd", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn", "co", "cr", "cu", "cv", "cw", "cx", "cy", "cz", "de", "dj", "dk", "dm", "do", "dz", "ec", "ee", "eg", "eh", "er", "es", "et", "eu", "fi", "fj", "fk", "fm", "fo", "fr", "ga", "gb", "gd", "ge", "gf", "gg", "gh", "gi", "gl", "gm", "gn", "gp", "gq", "gr", "gs", "gt", "gu", "gw", "gy", "hk", "hm", "hn", "hr", "ht", "hu", "id", "ie", "il", "im", "in", "io", "iq", "ir", "is", "it", "je", "jm", "jo", "jp", "ke", "kg", "kh", "ki", "km", "kn", "kp", "kr", "kw", "ky", "kz", "la", "lb", "lc", "li", "lk", "lr", "ls", "lt", "lu", "lv", "ly", "ma", "mc", "md", "me", "mf", "mg", "mh", "mk", "ml", "mm", "mn", "mo", "mp", "mq", "mr", "ms", "mt", "mu", "mv", "mw", "mx", "my", "mz", "na", "nc", "ne", "nf", "ng", "ni", "nl", "no", "np", "nr", "nu", "nz", "om", "pa", "pe", "pf", "pg", "ph", "pk", "pl", "pm", "pn", "pr", "ps", "pt", "pw", "py", "qa", "re", "ro", "rs", "ru", "rw", "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sj", "sk", "sl", "sm", "sn", "so", "sr", "ss", "st", "su", "sv", "sx", "sy", "sz", "tc", "td", "tf", "tg", "th", "tj", "tk", "tl", "tm", "tn", "to", "tp", "tr", "tt", "tv", "tw", "tz", "ua", "ug", "uk", "um", "us", "uy", "uz", "va", "vc", "ve", "vg", "vi", "vn", "vu", "wf", "ws", "ye", "yt", "za", "zm", "zw", "бел", "мкд", "мон", "рф", "срб", "укр", "қаз", "հայ", "الاردن", "الجزائر", "السعودية", "المغرب", "امارات", "ایران", "بھارت", "تونس", "سودان", "سورية", "عراق", "عمان", "فلسطين", "قطر", "مصر", "مليسيا", "پاکستان", "भारत", "বাংলা", "ভারত", "ਭਾਰਤ", "ભારત", "இந்தியா", "இலங்கை", "சிங்கப்பூர்", "భారత్", "ලංකා", "ไทย", "გე", "中国", "中國", "台湾", "台灣", "新加坡", "香港", "한국"}
	urlValidGtld := `(?:(?:` + strings.Join(gtLds, "|") + `)(?=[^[:alnum:]@]|$))`;
	urlValidCtld := `(?:(?:` + strings.Join(ctLds, "|") + `)(?=[^[:alnum:]@]|$))`;
	urlValidDomain := `(?:(?>(?:[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}][[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\-_]*)?[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\.)+(?:(?:[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}][[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\-]*)?[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\.)(?:` + urlValidGtld + `|` + urlValidCtld + `|` + `(?:xn--[0-9a-z]+)` + `)` + `)` + `|(?:` + `(?:(?:[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}][[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\-]*)?[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\.)` + `(?:` + urlValidGtld + `|` + `(?:xn--[0-9a-z]+)` + `|` + `(?:(?:co|tv)(?=[^[:alnum:]@]|$))` + `)` + `)` + `|(?:` + `(\Khttps?://)` + `(?:` + `(?:` + `(?:(?:[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}][[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\-]*)?[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\.)` + urlValidCtld + `)` + `|(?:` + `[.[^[:punct:]\s\p{Z}]]` + `+\.` + `(?:` + urlValidGtld + `|` + urlValidCtld + `)` + `)` + `)` + `)` + `|(?:` + `(?:(?:[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}][[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\-]*)?[[:alnum:]\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]\.)` + urlValidCtld + `(?=/)` + `)`
	validUrl := `(((?:[^A-Z0-9@＠$#＃\x{202a}-\x{202e}]|^))((https?://)?(` + urlValidDomain + `)` + `(?::(` + `[0-9]++` + `))?` + `(/` + `(?:(?:[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*(?:\((?:[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]+|(?:[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*\([a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]+\)[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*))\)[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*)*[a-z0-9=_#/\-\+\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]|(?:\((?:[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]+|(?:[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*\([a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]+\)[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]*))\)))|(?:@[a-z0-9!\*\';:=\+,.\$/%#\[\]\-_~\|&@\x{00c0}-\x{00d6}\x{00d8}-\x{00f6}\x{00f8}-\x{00ff}\x{0100}-\x{024f}\x{0253}\x{0254}\x{0256}\x{0257}\x{0259}\x{025b}\x{0263}\x{0268}\x{026f}\x{0272}\x{0289}\x{028b}\x{02bb}\x{0300}-\x{036f}\x{1e00}-\x{1eff}]+/))` + `*+` + `)?` + `(\?` + `[a-z0-9!?\*\'\(\);:&=\+\$/%#\[\]\-_\.,~\|@]` + `*` + `[a-z0-9_&=#/]` + `)?` + `)` + `)`

	posPatterns = map[util.KoreanPos]pcre.Regexp{
		util.Korean: pcre.MustCompile(`([가-힣]+)`, pcre.MULTILINE|pcre.UTF8),
		util.Alpha:  pcre.MustCompile(`([[:alnum:]]+)`, pcre.MULTILINE|pcre.UTF8),
		util.Number: pcre.MustCompile(`(\$?[[:digit:]]+(,[[:digit:]]{3})*([/~:\.-][[:digit:]]+)?(천|만|억|조)*(%|원|달러|위안|옌|엔|유로|등|년|월|일|회|시간|시|분|초)?)`, pcre.MULTILINE|pcre.UTF8),
		util.KoreanParticle: pcre.MustCompile(`([ㄱ-ㅣ]+)`, pcre.MULTILINE|pcre.UTF8),
		util.Punctuation:    pcre.MustCompile(`([[:punct:]·…’]+)`, pcre.MULTILINE|pcre.UTF8),
		util.URL:            pcre.MustCompile(validUrl, pcre.MULTILINE|pcre.UTF8),
		util.Email:          pcre.MustCompile(`([[:alnum:]\.\-_]+@[[:alnum:]\.]+)`, pcre.MULTILINE|pcre.UTF8),
		util.Hashtag:        pcre.MustCompile(`(^|[^&\p{L}\p{M}\p{Nd}_\x{200c}\x{a67e}\x{05be}\x{05f3}\x{05f4}\x{309b}\x{309c}\x{30a0}\x{30fb}\x{3003}\x{0f0b}\x{0f0c}\x{0f0d}])(#|＃)([\p{L}\p{M}\p{Nd}_\x{200c}\x{a67e}\x{05be}\x{05f3}\x{05f4}\x{309b}\x{309c}\x{30a0}\x{30fb}\x{3003}\x{0f0b}\x{0f0c}\x{0f0d}]*[\p{L}\p{M}][\p{L}\p{M}\p{Nd}_\x{200c}\x{a67e}\x{05be}\x{05f3}\x{05f4}\x{309b}\x{309c}\x{30a0}\x{30fb}\x{3003}\x{0f0b}\x{0f0c}\x{0f0d}]*)`, pcre.MULTILINE|pcre.UTF8),
		util.ScreenName:     pcre.MustCompile(`([^a-z0-9_!#$%&*@＠]|^|RT:?)([@＠]+)([a-z0-9_]{1,20})(/[a-z][a-z0-9_\-]{0,24})?`, pcre.MULTILINE|pcre.UTF8),
		util.CashTag:        pcre.MustCompile(`(^|[\x{0009}-\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{180E}\x{2000}-\x{200a}\x{2028}\x{2029}\x{202F}\x{205F}\x{3000}])(\$)([a-z]{1,6}(?:[._][a-z]{1,2})?)(?=$|\s|[[[:punct:]]])`, pcre.MULTILINE|pcre.UTF8),
		util.Space:          pcre.MustCompile(`\s+`, pcre.MULTILINE|pcre.UTF8),
	}

	chunkOrder = []util.KoreanPos{
		util.URL,
		util.Email,
		util.ScreenName,
		util.Hashtag,
		util.CashTag,
		util.Number,
		util.Korean,
		util.KoreanParticle,
		util.Alpha,
		util.Punctuation,
	}
}

func getChunks(input string, keepSpace bool) []string {
	return shared.ConvertSliceOfInterfaceToSliceOfString(
		scala.Map(chunk(input), func(key interface{}, value interface{}) interface{} {
			return value.(KoreanToken).Text
		}).([]interface{}))
}

func splitBySpaceKeepingSpace(s string) []string {
	space := regexp.MustCompile(`\s+`)

	loc := space.FindAllStringIndex(s, -1)

	tokens := []string{}
	index := 0
	for _, v := range loc {
		if index < v[0] {
			tokens = append(tokens, s[index:v[0]])
		}
		tokens = append(tokens, s[v[0]:v[1]])
		index = v[1]
	}

	l := len(s)
	if index < l {
		tokens = append(tokens, s[index:l])
	}

	return tokens
}

type ByOccur []ChunkMatch

func (t ByOccur) Len() int {
	return len(t)
}

func (t ByOccur) Swap(i int, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByOccur) Less(i int, j int) bool {
	return t[i].start < t[j].start
}

func splitChunks(text string) []ChunkMatch {
	chunkMatches := []ChunkMatch{}
	tl := len(text)
	if []rune(text)[0] == ' ' {
		return append(chunkMatches, ChunkMatch{0, tl, text, util.Space})
	} else {
		chunksBuf := []ChunkMatch{}
		matchedLen := 0

		for _, v := range chunkOrder {
			if matchedLen < tl {
				start := 0

				re := posPatterns[v]

				for {
					offset := (&re).FindIndex([]byte(text[start:]), 0)

					if offset == nil {
						break
					}

					s := start + offset[0]
					e := start + offset[1]

					cm := ChunkMatch{s, e, text[s:e], v}

					if scala.ForAll(chunksBuf, cm.disjoint) {
						chunksBuf = append(chunksBuf, cm)
						matchedLen += cm.end - cm.start
					}

					start = e
				}
			}
		}

		sort.Sort(ByOccur(chunksBuf))

		return fillInUnmatched(text, chunksBuf, util.Foreign)
	}
}

type tmpFillInUnmatchedTuple struct {
	l       []ChunkMatch
	prevEnd util.KoreanPos
}

func fillInUnmatched(text string, chunks []ChunkMatch, pos util.KoreanPos) []ChunkMatch {
	v := (scala.FoldLeft(chunks, tmpFillInUnmatchedTuple{[]ChunkMatch{}, 0},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			v := folded.(tmpFillInUnmatchedTuple)
			cm := value.(ChunkMatch)
			if util.KoreanPos(cm.start) == v.prevEnd {
				t := []ChunkMatch{cm}
				t = append(t, v.l...)
				return tmpFillInUnmatchedTuple{t, util.KoreanPos(cm.end)}
			} else if util.KoreanPos(cm.start) > v.prevEnd {
				t := []ChunkMatch{cm}
				t = append(t, ChunkMatch{int(v.prevEnd), cm.start, text[int(v.prevEnd):cm.start], pos})
				t = append(t, v.l...)
				return tmpFillInUnmatchedTuple{t, util.KoreanPos(cm.end)}
			} else {
				panic("Non-disjoint chunk matches found.")
			}
		}).(interface{})).(tmpFillInUnmatchedTuple)

	l := len(text)

	output := []ChunkMatch{}
	if int(v.prevEnd) < l {
		output = append(output, ChunkMatch{int(v.prevEnd), l, text[int(v.prevEnd):l], pos})
		output = append(output, v.l...)
	} else {
		output = append(output, v.l...)
	}

	return ConvertSliceOfInterfaceToSliceChunkMatch(scala.Reverse(output))
}

func getChunksByPos(input string, pos util.KoreanPos) []KoreanToken {
	return ConvertSliceOfInterfaceToSliceOfKoreanToken(scala.Filter(chunk(input), func(value interface{}) bool {
		return value.(KoreanToken).Pos == pos
	}).([]interface{}))
}

type tmpChunkTuple struct {
	l []KoreanToken
	i int
}

func chunk(input string) []KoreanToken {
	flatMap := func(s []string, f func(value interface{}) interface{}) []ChunkMatch {
		m := []ChunkMatch{}
		for _, v := range s {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				for _, d := range x.([]ChunkMatch) {
					m = append(m, d)
				}
			} else {
				panic("Kind must be array or slice")
			}
		}
		return m
	}

	m := flatMap(splitBySpaceKeepingSpace(input), func(value interface{}) interface{} {
		return splitChunks(value.(string))
	})

	v := (scala.FoldLeft(m, tmpChunkTuple{[]KoreanToken{}, 0},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			v := folded.(tmpChunkTuple)
			cm := value.(ChunkMatch)

			segStart := strings.Index(input[v.i:], cm.text)

			if segStart != -1 {
				segStart = v.i + segStart
			}

			ss := utf8.RuneCountInString(input[0:segStart])

			l := utf8.RuneCountInString(cm.text)

			t := []KoreanToken{{cm.text, cm.pos, ss, l, false}}
			t = append(t, v.l...)

			return tmpChunkTuple{t, segStart + l}
		}).(interface{})).(tmpChunkTuple)

	return ConvertSliceOfInterfaceToSliceOfKoreanToken(scala.Reverse(v.l).([]interface{}))
}
