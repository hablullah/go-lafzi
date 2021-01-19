package lafzi

const (
	hamza              = '\u0621'
	alefWithMaddaAbove = '\u0622'
	alefWithHamzaAbove = '\u0623'
	wawWithHamzaAbove  = '\u0624'
	alefWithHamzaBelow = '\u0625'
	yehWithHamzaAbove  = '\u0626'
	alef               = '\u0627'
	beh                = '\u0628'
	tehMarbuta         = '\u0629'
	teh                = '\u062A'
	theh               = '\u062B'
	jeem               = '\u062C'
	hah                = '\u062D'
	khah               = '\u062E'
	dal                = '\u062F'
	thal               = '\u0630'
	reh                = '\u0631'
	zain               = '\u0632'
	seen               = '\u0633'
	sheen              = '\u0634'
	sad                = '\u0635'
	dad                = '\u0636'
	tah                = '\u0637'
	zah                = '\u0638'
	ain                = '\u0639'
	ghain              = '\u063A'
	feh                = '\u0641'
	qaf                = '\u0642'
	kaf                = '\u0643'
	lam                = '\u0644'
	meem               = '\u0645'
	noon               = '\u0646'
	heh                = '\u0647'
	waw                = '\u0648'
	alefMaksura        = '\u0649'
	yeh                = '\u064A'
	fathatan           = '\u064B'
	dammatan           = '\u064C'
	kasratan           = '\u064D'
	fatha              = '\u064E'
	damma              = '\u064F'
	kasra              = '\u0650'
	shadda             = '\u0651'
	sukun              = '\u0652'
)

var phonetics = map[rune]rune{
	hamza:              'x',
	alefWithMaddaAbove: 'x',
	alefWithHamzaAbove: 'x',
	wawWithHamzaAbove:  'x',
	alefWithHamzaBelow: 'x',
	yehWithHamzaAbove:  'x',
	alef:               'x',
	beh:                'b',
	tehMarbuta:         't',
	teh:                't',
	theh:               's',
	jeem:               'z',
	hah:                'h',
	khah:               'h',
	dal:                'd',
	thal:               'z',
	reh:                'r',
	zain:               'z',
	seen:               's',
	sheen:              's',
	sad:                's',
	dad:                'd',
	tah:                't',
	zah:                'z',
	ain:                'x',
	ghain:              'g',
	feh:                'f',
	qaf:                'k',
	kaf:                'k',
	lam:                'l',
	meem:               'm',
	noon:               'n',
	heh:                'h',
	waw:                'w',
	yeh:                'y',

	fatha: 'a',
	damma: 'u',
	kasra: 'i',
	sukun: '0',
}
