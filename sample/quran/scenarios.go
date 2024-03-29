package main

type Scenario struct {
	Name      string
	Queries   []string
	Documents []string
}

var scenarios = []Scenario{
	{
		Name: "A1",
		Queries: []string{
			"ulul albaab", "uulul albab", "ulul albab", "uulul albaab", "uuluul albaab",
			"ululalbab", "walulalbab", "uululalbab", "uululalbaab", "ulul al bab",
			"uwlulalbaab", "uulu al albaab", "uuluulalbaab", "uwlul albaab",
		},
		Documents: []string{
			"38:29", "14:52", "13:19", "2:269", "39:18", "39:9", "3:7", "40:54",
			"2:179", "38:43", "3:190", "65:10", "5:100", "12:111", "39:21", "2:197",
		},
	}, {
		Name: "A2",
		Queries: []string{
			"alhamdulillah", "alhamdulillaah", "alhamdulilla", "alhmadulillah",
			"alhamdu lillaah",
		},
		Documents: []string{
			"37:182", "1:2", "6:45", "35:34", "40:65", "27:59", "27:93", "14:39", "31:25",
			"23:28", "18:1", "10:10", "27:15", "39:29", "39:74", "17:111", "39:75", "29:63",
			"34:1", "6:1", "16:75", "35:1", "7:43",
		},
	}, {
		Name: "A3",
		Queries: []string{
			"yaa aiyuhannabiyu", "yaa 'ayyuhaannabbiyyu", "ya ayyuhannabiyyu",
			"yaa ayyuhannabiyyu", "yaa ayyuhan nabiyyu", "yaa ayyuhannbiyyu",
			"yaa ayyuhannubiyyu", "yaa ayyuhan nabiiyyu", "yaa ayyuhannabiiyu",
			"yaa ayyuhann nabiyyu", "yaayyuhannabiyyu", "yaa ayyuhannabiiy",
			"ya ayyuhannabiy", "yaa ayyuhannabiyyuu", "ya ayyuhannabiyu",
			"yaaayyuahannabiyu", "yaa aiyuhannabiiyu", "ya ayyuhan nabiyyu",
			"ya ayyuhannas", "yaa ayyuhannabiy", "ya aiyuhan nabiyu",
			"yaa aiyyuhannabii", "yaa ayyuhaa an nabiyyu", "yaa ayyukhannabiyu",
			"yaaayyuhaannabiyyu", "yaaaiyuhannabiyyu", "yaa ayyuhaan nabiyyu",
			"yaa ayyuhannabiyy",
		},
		Documents: []string{
			"33:45", "8:64", "66:9", "9:73", "33:1", "33:28", "66:1", "8:70",
			"33:59", "8:65", "60:12", "65:1", "33:50",
		},
	}, {
		Name: "A4",
		Queries: []string{
			"rosulullah", "rasuulullah", "rasulullah", "rasuulullaah", "rosuulullah",
			"rosuululloh", "rosuulullaah", "rasululloh", "rosululloh", "rasulullaah",
			"rosuu lullooh", "rosululullah", "rasuwlullah", "rasuul allaah", "rasullullaah",
		},
		Documents: []string{
			"91:13", "63:5", "63:1", "61:5", "61:6", "7:158", "48:29", "4:171",
			"33:21", "33:40", "49:3", "63:7", "49:7", "9:81", "9:61", "4:157",
			"9:120", "33:53", "6:124",
		},
	}, {
		Name: "A5",
		Queries: []string{
			"kholidiinafiihaa abada", "kholidiinafiihaa'abadaa", "kholidiina fiihaa abada",
			"khaalidiina fiihaa abadaa", "khaalidiina fiihaa abada", "khoo lidiina fiihaa abada",
			"kha lidina fiha abada", "khoolidiina fiihaa abadaa", "khoolidiina fiiha abadaa",
			"khalidinafiha abada", "khoolidiinafiihaa abadaa", "khoolidiina fii haa abadaa",
			"kholidina fiiha abadaa", "khaalidiina fiiha abada", "khaalidiyna fiihaa abadaa",
			"khalidiina fiiha abada", "kholidiina fiha abadaa", "kholidiina fiihaa abadaa",
			"khalidiinafiihaabada", "ghaalidiina fiihaa abadaa", "kholidina fihaabada",
			"khoo lidiina fiiha abadaa", "kholidiina fiiha abada", "kholidiina fiiha abadaa",
			"kholidii nafiiha abadaa", "khoo lidiyna fiyhaa a bada", "kholidiniina fiiha abada",
			"kholidina fiha abada", "khalidiyna fiyha abadaa", "kholidiinafiihaa abadaa",
			"kholidinafiha abada", "khalidiina fiihaaa abada", "khaakudiina fiihaa abada",
			"khaalidiina fiiha abadaa", "khaalidiinafiihaaabadaa", "khalidina fiyha abada",
			"kholodina fiha abada", "khalidiina fiihaa abadaa", "khaa lidiina fihaa abadaa",
			"khaa lidiina fiihaa abadaa",
		},
		Documents: []string{
			"9:22", "33:65", "4:169", "72:23", "98:8", "4:122", "4:57", "5:119",
			"9:100", "64:9", "65:11",
		},
	}, {
		Name: "A6",
		Queries: []string{
			"asshobiriin", "ash shoobiriin", "asshabiriina", "ash-shaabiriin", "asshoobiriin",
			"shabirin", "assoobiriin", "ashshobiriin", "ashshoobiriin", "ashshabiriin",
			"shaabiriin", "ash shobiriin", "shoobiriin", "ash shaabiriina", "ashshaabiriyn",
			"ssobiriin", "ash-shobiriin", "ashshobirin", "asshaabiriin", "ashshaabiriin",
			"assobiriin", "ashobiriin", "ash shoobiriiin", "asyoobiriin", "ash shaabiriin",
			"ash shoobiriyn", "as shobirin", "asshobirin", "asshabiriin", "al shoobiriin",
			"ashoobiriin", "ashobirin", "assabirin",
		},
		Documents: []string{
			"21:85", "47:31", "16:126", "3:17", "2:153", "3:142", "22:35", "8:46",
			"2:155", "3:146", "37:102", "8:66", "33:35", "2:177", "2:249",
		},
	}, {
		Name: "A7",
		Queries: []string{
			"wailun yaumaizillilmukazzibiin", "wailun yaumaidzillilmukadzibiin",
			"way luyyau maidzi llilmukadzibiin", "wayluy yawma`idzillil mukadzdzibiin",
			"wayluy yaumaidzil lilmukadzdzibiin", "wayluyyawmaidzillil mukadzibiin",
			"wayluy yaumaizil lil mukazzibin", "wailuyyaumaidzillilmukadzdzibiin",
			"wayluyyaumaidzillilmukadzibiin", "wailuy yaumaidzil lil mukadzdzibiin",
			"wailuyyaumaizillilmukazzibiin", "wayluyauma i dzillilmukadzdzibiin",
			"wayluyaumaidzil lilmukadz dzibiin", "waeluy yaumaidzill lilmukadzibiin",
			"wayluyyaumaidzil lil mukadzdzibiin", "waylun yaumaidzin lilmukadzdzibiyn",
			"wayluy yauma idzil lilmukadzdzibiin", "waylun yawmaizillilmukazibiin",
			"wayluyyawmaidzil lil mukadzibiin", "wailuy yaumaidzil lilmukadzibiin",
			"wailuyyaumaidzillillmukadzdzibin", "wailuyyauma idzil lil mukadzdzibiin",
			"wayluy yaumaidzil lilmukaddzibiin", "wailuyyaumaidzil lilmukadz dzibiin",
			"wayluyyauma'idzillilmukadzdzibiin", "wailuy yaumaidzil lilmukadzdzibiin ",
			"wailuyaumaidzillilmukadzibiin", "wayluyawma-idzil lil mukadzdzibiin",
			"waylun yuma idzil lilmukadz dhibin", "wayluy yawmaidzil lilmukadzibiin",
			"wayluyyaw maijillilmukazzibiin ", "wailuy yaumaidzillilmukadzibiin",
			"wailun yaumaizil lilmukazzibiin", "wayluyyawma idzil lilmukadzibiin",
			"wailun yau maidzinillilmukadzdzibin", "wayluy yawmaidzinl lilmukadzdzibiyn",
			"wayluyyaumaidzillilmukadzzibiin", "wayluyyaumaizillilmukazzzibiin",
			"wayluy yauma idzil lilmukadzibin", "wailunyau maizil lilmukazzibin",
			"wayluyyaumaidzillilmukadzdzibiin", "waylun yaumaidzinil lilmukadzzibiin",
			"wailuyyaumaidillil mukadzdzibiin", "waylun yaumaidzin lilmukadzibiin",
			"wayluy yawmaidzin lilmukadzdzibiin", "wayluyaumaidzinnilmukadzibiin",
			"wayluy yaumaidzil lilmukadzzibiin", "wailuyyaumaidzillilmukadzibiin",
			"wailuyyaumaidzillilmukadzzibiin", "wayluy yaumaidzil lil mukadzibin",
			"wayluy yawma idzil lilmukadzzibiin", "wailun yaumaidzillilmukadzibin",
			"wailun yauma'idzillilmukadzibiin", "waylun yaumaidzililmukadzdzibiin",
			"wailun yau maidzil lilmukadzdzi biin", "wayluyyaw maidzillilmukadzzibiin",
		},
		Documents: []string{
			"83:10", "77:49", "77:47", "77:45", "77:40", "77:37", "77:34",
			"77:28", "77:24", "77:19", "77:15", "52:11",
		},
	}, {
		Name: "A8",
		Queries: []string{
			"'ala kullisyai inggqodiir", "'ala kulli syaiin qodiir", "'alaa kulli syay inqodiir",
			"'alaa kulli syay`in qadiir", "'alaa kulli syai in qodiir", "'alaa kulli syaiin kodiir",
			"'ala kulli syaiin qadir", "'ala kulli syai'in qodiir", "'ala kulli syaiinnqadiir",
			"'ala kulli syai ingqodiir", "'alaa kulli syai'in qadiir", "'ala kulli syayin qadiir",
			"'ala kulilli syaiingqodiir", "'alakullisyay in qodiir", "alaa kulli syay ing qodiir",
			"'alakullisyai in qodiir", "'ala kulli syaiing qodiir", "'ala kulli syai`in qadiir",
			"'alaa kulli syayin qadiir", "'ala kulli syaiin wadiir", "'ala kullisyainkhodir",
			"'alaa kulli syai'in qodiir", "'alaa kulli syai-in(g) qodiir", "a'lakullisyaiinkhdir",
			"'alaa kulli syaiin qodiir", "'ala kulli syaiinqodiir", "'ala kulli syai'in~ qadiir",
			"'ala kulli syai in qadiir", "ala kulli syay ingkodiir", "'ala kulli syai-in qodiir",
			"'alaa kulli syaiin qadiir", "'a lakulli syai innqodir", "'ala kulli syaii kodiir",
			"ala kulli syaiin qodiir", "'alaa kulli syai in qadiir", "'alakulli shai inqodiir",
			"'ala kulli syai in qodiyr", "a'la kulli syaiingkodiiir", "'ala kulli syaiin qodir",
			"'ala kulli syai inq qhodiir", "ala kulli syai'ing qodir", "ngala kulli syai in kodir",
			"'alaa kullisyai inq-nqadiir", "'alaa kulli syai-in qodiir", "'ala kulli syai in qodiir",
			"'alaa kulli syay in qadiir", "'ala kulli syaiing khodiir", "'alaa kulli syay in qodiir",
			"'alakullisyai ingkodiir", "'alaakullisyainkodiir", "'ala kulli syayin qadir",
			"'ala kulli syay in qodiir", "alakullisyaiin qodir", "'ala kulli syai'in qadiir",
			"'ala kulli syai in 'qadiir", "'ala kulli syain qodiir",
		}, Documents: []string{
			"11:4", "3:189", "5:120", "57:2", "22:6", "48:21", "67:1", "33:27",
			"6:17", "42:9", "2:106", "9:39", "29:20", "5:40", "30:50", "16:77",
			"2:148", "3:165", "3:29", "64:1", "46:33", "41:39", "59:6", "3:26",
			"65:12", "2:284", "2:20", "5:19", "24:45", "35:1", "2:109", "8:41",
			"5:17", "66:8", "2:259",
		},
	}, {
		Name: "A9",
		Queries: []string{
			"yastathii'uun", "yastathi'un", "yastathi'uun", "yastatii uun", "yastatii'uun",
			"yastathiy'uun", "yastati'uwn", "yastathiun", "yastathii'uuna", "yasta thii'uun",
			"yastii uun", "yas tatii uun", "yastathy'uwn", "yastathiiu'uun", "yastathii'un",
			"yastathingun", "yastatwii'uun", "yastat'ii 'uun", "yastatiuun", "yas tathii'uun",
		},
		Documents: []string{
			"26:211", "36:75", "7:192", "36:50", "68:42", "25:9", "17:48", "21:40",
			"18:101", "7:197", "21:43", "4:98", "16:73", "11:20", "2:273",
		},
	}, {
		Name: "A10",
		Queries: []string{
			"innalloha ghofururrohiim", "innallaha ghofuururrohiim",
			"innallaha ghofuururrahiim", "innallaaha ghafuurur rahiim",
			"innallaaha ghofuurur rohiim", "innallaaha ghafuururrohiim",
			"innallaha ghafururrahim", "innalloha ghofuururrohiim",
			"innallooha ghofuururrohiim", "innallaha ghofururrahiim",
			"innalloha ghafuururrahiim", "innallaha gofururrohiim",
			"innallaha ghofurur rohiim", "innallaha gafururrrahiim",
			"innaallaha ghafuururrahiim", "inna allaha ghafuurur rahiim",
			"innallaha ghafururrohim", "innalloha ghofuurur rohiim",
			"innallahaghfururrahim", "innaalloha ghafuururrahiim",
			"innalllaha ghofuurur rohiim", "innallaaha ghafuururrahiim",
			"innalloha gofuururrohiiim", "innallaha ghafuurun rahiim",
			"innallaha ghofuururrohim", "innallooha ghopuu rurrohiim",
			"innalloha gofurur rhoiim", "innallaha gafuurur rahiim",
			"innallaha ghafuururrahiim", "innallahaghofuururrahiim",
			"innaallaha ghofuwrurochiym", "innallaha ghofururrohiim",
			"innaalloha gofururrohim", "innallaha ghafururrahiim",
			"innalloha ghofurur rohim", "innallaha ghofururrohim",
			"innallaha ghafuurur-rahiim", "inna allaaha ghafuurrurrahiim",
			"innallaha gofururrahiim", "innallaha ghafuru rahim",
			"innallaha ghofuurur rohiim", "innallaha qhofururrohiim",
		},
		Documents: []string{
			"2:192", "24:5", "3:89", "8:69", "2:199", "5:39", "2:226", "2:182",
			"9:102", "64:14", "16:115", "58:12", "2:173", "49:14", "9:99", "9:5",
			"60:12", "24:62", "5:3", "73:20",
		},
	}, {
		Name: "A11",
		Queries: []string{
			"masyalullazii", "masyalulladzii", "matsalulladzii", "matsalul ladzii",
			"matsalul lazi", "matsalullazi", "matsaluladzii", "matsalulladzi",
			"matsalulladziy", "matsalulladziii", "masaluladzi", "matsalul ladzi",
			"masyalulladzi", "matsalullazii", "masalul ladzi", "matsaludzin",
			"mathalullazii", "mastalul ladziy", "masalullazii", "masalullazi",
			"matsalu al ladzii", "matsalul ladziy", "masalulladzi",
		}, Documents: []string{
			"2:171", "59:15", "2:17", "29:41", "14:18", "2:261", "62:5", "2:265",
			"2:214", "2:228",
		},
	}, {
		Name: "A12",
		Queries: []string{
			"qaumizzolimiin", "qoumizh zhooliman", "qow mizhoolimiin", "qawmizhzhaalimiin",
			"qowmidz dzoolimiin", "qaumidzhoolimiin", "qaumidhalimin", "qoumizhzhoolimiin",
			"qowmidzhoolimiin", "qoumidzhdzhoolimiin", "qaumidh dhoolimiin", "qaumitzalimiin",
			"qoumizzholimiin", "qoumizhoolimiin", "qoumidz dzoolimiin", "qoumizzoolimiin",
			"qaumidh dholimiin", "qaumidzh dzhalimiin", "qaumidzdzaalimiin", "qaumidzoolimiyn",
			"qawmizh zhaalimiin", "qawmidzdzolimiin", "qowmidzoolimiin", "qoumizh-zholimiin",
			"khumizhzholimin", "qoumadzdzoolimiin", "qoumi dzaalimiin", "qoumizh zhoolimiin",
			"qaumizhzhoolimiin", "koumizzolimiin", "qowmizh zhoolimiin", "qaumidhdhaalimiin",
			"qowmidh dholimiin", "qowmizhoolimiin", "qaumizh zhaalimiin", "qoumidzoolimiin",
			"qaumidzolimiin", "qowmi adhdhoolumiyn", "kowmittholimiin", "qoumizhzholimin",
			"qoumidhoolimiin", "qoumizzhoolimiin", "qaumizzhoolimiin", "qoumidh dholimin",
			"koumizzholimin", "qaumi dzhaalimiin", "qoumidhdholimiin", "qaumi al dzaalimiin",
			"khowmizoolimiin", "qoumidzhoolimiin", "qoumdzoolimiin", "qawmi zholimin",
			"qowmidz dzhoolimiin", "qoumizzholimin", "qaumizzaalimiin", "qaumizzhaalimiin",
			"qaumizh zhalimiin", "qoumidzzolimiin",
		},
		Documents: []string{
			"6:47", "23:94", "26:10", "10:85", "28:21", "23:41", "7:47", "23:28",
			"61:7", "11:44", "3:86", "46:10", "28:50", "66:11", "6:68", "62:5",
			"9:109", "5:51", "9:19", "28:25", "2:258", "7:150", "6:144",
		},
	}, {
		Name: "A13",
		Queries: []string{
			"dholliin", "dhooolllin", "dhaalliin", "dhoolliin", "dhallin", "dhoooooolliin",
			"dloolliin", "dhool liin", "dolliin", "dhalliin", "dhaalliyn", "dhollin",
			"zooliin", "dhaaalliiin", "dhlolliyn", "dholliiiin", "dlaalliiyn", "dhdhalliin",
			"dho lin", "dhaaaaaalliiiin", "dlolliin", "dzoolliin", "dhool liinn", "dhaal liin",
		}, Documents: []string{
			"68:26", "26:86", "26:20", "37:69", "56:92", "56:51", "83:32",
			"15:56", "1:7", "23:106", "3:90", "6:77", "2:198",
		},
	}, {
		Name: "A14",
		Queries: []string{
			"mimba'di maa jaa a", "mimba'dimaajaa", "mim ba'di maa jaa`a", "mim ba'di maa jaa a",
			"mim ba'di ma ja a", "mimba'di maa jaaaaaa'a", "mimba' dimaajaa'", "min ba'di maa jaa'a",
			"mim ba'di ma jaa", "mimba'dimaajaa a", "mim ba di maa jaa a", "mimba'di maajaa a",
			"min ba'di ma jaa`a", "min ba'di maa jaa a", "min ba'di maa jaa", "mim ba'di maa jaa-a",
			"minbakdimja'", "mimmba'di maa jaa a", "mimba'di maa jaa", "min ba'di maa jaa- a",
			"min ba'di maa ja a", "minba'dimaaja", "mimm ba'di maa jaa a", "mim ba'dimaajaaa'",
			"min ba'di maja-a", "mimba'dimaa jaa a", "mimba'di majaa a", "mimba'di ma ja'a",
			"minba'di ma ja a", "mimmba'kdimaajaa", "min ba'di maa jaa a'", "minba'di maajaa-a",
			"mamba'di maa jaa a", "mimba'dimaajaa'", "min ba'di majaa", "minba'dimaajaa a",
			"mimba'dimaajaa'a", "min' ba'di maa jaa a", "min ba'di maajaa a",
		},
		Documents: []string{
			"98:4", "2:209", "3:105", "2:211", "45:17", "3:19", "42:14", "3:61",
			"2:145", "4:153", "2:253", "2:213",
		},
	}, {
		Name: "A15",
		Queries: []string{
			"tanziil", "tangziil", "tanzil", "tan ziil", "tanziyl", "tandziyl", "tannziil",
			"tan nziil", "tangzil", "tanjiil", "tanzilyl", "tanzhil", "tangjziil", "tandziil",
		},
		Documents: []string{
			"36:5", "69:43", "56:80", "41:2", "26:192", "76:23", "46:2", "45:2",
			"40:2", "20:4", "32:2", "25:25", "17:106", "39:1", "41:42",
		},
	}, {
		Name: "A16",
		Queries: []string{
			"fa-ula-ikahum", "faulaaaikahum", "fa `ulaaika hum", "fa ulaa ikahum", "faulaaika hum",
			"fa ulaika hum", "fa ulaaaaaaikahum", "faulaikahum", "fa uulaikahum", "fa uu la ikahum",
			"fauulaikahum", "fauula ika hum", "fa`ulaaikahum", "faulaaika", "fauwlaikahum",
			"fauulaika hum", "fa ulaa ika hum", "faulaa ika hum", "fa ula ikahum", "fa uu la ika hum",
			"fa-uulaikahum", "fa ulaaika hum", "fa ula ika hum", "fau laika hum", "fa u laika hum",
			"fa'ula'ika hum", "fa-ula-ika hum", "fa uulaika hum", "faulaikakhum", "fa uulaikahun",
			"faulaika hum",
			"fauu laika hum",
		},
		Documents: []string{
			"23:102", "70:31", "23:7", "3:82", "7:178", "3:94", "24:52", "7:8",
			"5:47", "2:121", "63:9", "64:16", "30:39", "9:23", "60:9", "59:9",
			"5:45", "49:11", "24:55", "5:44", "2:229",
		},
	},
}
