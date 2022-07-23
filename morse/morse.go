package morse

import (
	"strings"

	hangul "github.com/suapapa/go_hangul"
)

const (
	LangEnglish = "EN"
	LangKorean  = "KR"
)

var EnglishRuneToMorseMap = map[rune]string{
	'a': ".-",
	'b': "-...",
	'c': "-.-.",
	'd': "-..",
	'e': ".",
	'f': "..-.",
	'g': "--.",
	'h': "....",
	'i': "..",
	'j': ".---",
	'k': "-.-",
	'l': ".-..",
	'm': "--",
	'n': "-.",
	'o': "---",
	'p': ".--.",
	'q': "--.-",
	'r': ".-.",
	's': "...",
	't': "-",
	'u': "..-",
	'v': "...-",
	'w': ".--",
	'x': "-..-",
	'y': "-.--",
	'z': "--..",
	'1': ".----",
	'2': "..---",
	'3': "...--",
	'4': "....-",
	'5': ".....",
	'6': "-....",
	'7': "--...",
	'8': "---..",
	'9': "----.",
	'0': "-----",
	'?': "..--..",
	'!': "-.-.--",
	'.': ".-.-.-",
	',': "--..--",
	';': "-.-.-.",
	':': "---...",
	'+': ".-.-.",
	'-': "-....-",
	'/': "-..-.",
	'=': "-...-",
}
var 한글RuneToMorseMap = map[rune]string{
	'ㄱ': ".-..",
	'ㄴ': "..-.",
	'ㄷ': "-...",
	'ㄹ': "...-",
	'ㅁ': "--",
	'ㅂ': ".--",
	'ㅅ': "--.",
	'ㅇ': "-.-",
	'ㅈ': ".--.",
	'ㅊ': "-.-.",
	'ㅋ': "-..-",
	'ㅌ': "--..",
	'ㅍ': "---",
	'ㅎ': ".---",
	'ㅏ': ".",
	'ㅑ': "..",
	'ㅓ': "-",
	'ㅕ': "...",
	'ㅗ': ".-",
	'ㅛ': "-.",
	'ㅜ': "....",
	'ㅠ': ".-.",
	'ㅡ': "-..",
	'ㅣ': "..-",
	'ㅔ': "-.--",
	'ㅐ': "--.-",
	'ㅖ': "... ..-",
	'ㅒ': ".. ..-",
}
var (
	englishMorseToRuneMap = map[string]rune{}
	한글MorseToRuneMap      = map[string]rune{}
)

func init() {
	for k, v := range EnglishRuneToMorseMap {
		englishMorseToRuneMap[v] = k
	}

	for k, v := range 한글RuneToMorseMap {
		한글MorseToRuneMap[v] = k
	}

}

func runeToMorse(r rune) (result string, found bool) {
	if hangul.IsHangul(r) {
		result, found = 한글RuneToMorseMap[r]
	} else {
		result, found = EnglishRuneToMorseMap[r]
	}

	return
}

func ToMorseCode(text string) string {
	textFixed := &strings.Builder{}
	for _, c := range text {
		if hangul.IsHangul(c) {
			x, y, z := hangul.SplitCompat(c)
			textFixed.WriteRune(x)
			if y > 0 {
				textFixed.WriteRune(y)
			}
			if z > 0 {
				textFixed.WriteRune(z)
			}
		} else {
			textFixed.WriteRune(c)
		}
	}
	text = textFixed.String()

	sb := strings.Builder{}
	for i, c := range strings.ToLower(text) {
		if c == ' ' {
			sb.WriteRune('/')
			sb.WriteRune(' ')
			continue
		}
		if morse, ok := runeToMorse(c); ok {
			sb.WriteString(morse)
			if i < len(text)-1 {
				sb.WriteRune(' ')
			}
		}

	}

	return strings.TrimSuffix(sb.String(), " ")
}

func FromMorseCode(text string, lang string) string {
	var dict map[string]rune
	switch lang {
	case LangEnglish:
		dict = englishMorseToRuneMap
	case LangKorean:
		dict = 한글MorseToRuneMap
	}

	sb := strings.Builder{}

	current := ""
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '.', '-':
			current += string(text[i])
		case ' ':
			if r, ok := dict[current]; ok {
				sb.WriteRune(r)
			}
			current = ""
		case '/':
			sb.WriteRune(' ')
		}
	}
	if _, ok := dict[current]; ok {
		sb.WriteRune(dict[current])
	}

	return strings.ToUpper(sb.String())
}
