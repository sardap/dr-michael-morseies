package morse_test

import (
	"testing"

	"github.com/sardap/dr-michael-morseies/morse"
	"github.com/stretchr/testify/assert"
)

func TestToMorseCode(t *testing.T) {
	type scenario struct {
		input    string
		expceted string
	}

	scenarios := []scenario{
		{input: "A", expceted: ".-"},
		{input: "B", expceted: "-..."},
		{input: "C", expceted: "-.-."},
		{input: "D", expceted: "-.."},
		{input: "E", expceted: "."},
		{input: "F", expceted: "..-."},
		{input: "G", expceted: "--."},
		{input: "H", expceted: "...."},
		{input: "I", expceted: ".."},
		{input: "J", expceted: ".---"},
		{input: "K", expceted: "-.-"},
		{input: "L", expceted: ".-.."},
		{input: "M", expceted: "--"},
		{input: "N", expceted: "-."},
		{input: "O", expceted: "---"},
		{input: "P", expceted: ".--."},
		{input: "Q", expceted: "--.-"},
		{input: "R", expceted: ".-."},
		{input: "S", expceted: "..."},
		{input: "T", expceted: "-"},
		{input: "U", expceted: "..-"},
		{input: "V", expceted: "...-"},
		{input: "W", expceted: ".--"},
		{input: "X", expceted: "-..-"},
		{input: "Y", expceted: "-.--"},
		{input: "Z", expceted: "--.."},
		{
			input:    "a b c d e f g h i j k l m n o p q r s t u v w x y z",
			expceted: ".- / -... / -.-. / -.. / . / ..-. / --. / .... / .. / .--- / -.- / .-.. / -- / -. / --- / .--. / --.- / .-. / ... / - / ..- / ...- / .-- / -..- / -.-- / --..",
		},
		{input: "1234567890", expceted: ".---- ..--- ...-- ....- ..... -.... --... ---.. ----. -----"},
		{input: ".,?:", expceted: ".-.-.- --..-- ..--.. ---..."},
		{input: "한글", expceted: ".--- . ..-. .-.. -.. ...-"},
		{input: "구천육백사십일", expceted: ".-.. .... -.-. - ..-. -.- .-. .-.. .-- --.- .-.. --. . --. ..- .-- -.- ..- ...-"},
		{input: "사과를 사요", expceted: "--. . .-.. ...- -.. ...- / --. . -.- -."},
	}

	for _, scenario := range scenarios {
		out := morse.ToMorseCode(scenario.input)
		assert.Equalf(t, scenario.expceted, out, "Input %s", scenario.input)
	}
}

func TestFromMorseCode(t *testing.T) {
	type scenario struct {
		expected string
		input    string
		lang     string
	}

	scenarios := []scenario{
		{expected: "A", input: ".-", lang: morse.LangEnglish},
		{expected: "B", input: "-...", lang: morse.LangEnglish},
		{expected: "C", input: "-.-.", lang: morse.LangEnglish},
		{expected: "D", input: "-..", lang: morse.LangEnglish},
		{expected: "E", input: ".", lang: morse.LangEnglish},
		{expected: "F", input: "..-.", lang: morse.LangEnglish},
		{expected: "G", input: "--.", lang: morse.LangEnglish},
		{expected: "H", input: "....", lang: morse.LangEnglish},
		{expected: "I", input: "..", lang: morse.LangEnglish},
		{expected: "J", input: ".---", lang: morse.LangEnglish},
		{expected: "K", input: "-.-", lang: morse.LangEnglish},
		{expected: "L", input: ".-..", lang: morse.LangEnglish},
		{expected: "M", input: "--", lang: morse.LangEnglish},
		{expected: "N", input: "-.", lang: morse.LangEnglish},
		{expected: "O", input: "---", lang: morse.LangEnglish},
		{expected: "P", input: ".--.", lang: morse.LangEnglish},
		{expected: "Q", input: "--.-", lang: morse.LangEnglish},
		{expected: "R", input: ".-.", lang: morse.LangEnglish},
		{expected: "S", input: "...", lang: morse.LangEnglish},
		{expected: "T", input: "-", lang: morse.LangEnglish},
		{expected: "U", input: "..-", lang: morse.LangEnglish},
		{expected: "V", input: "...-", lang: morse.LangEnglish},
		{expected: "W", input: ".--", lang: morse.LangEnglish},
		{expected: "X", input: "-..-", lang: morse.LangEnglish},
		{expected: "Y", input: "-.--", lang: morse.LangEnglish},
		{expected: "Z", input: "--..", lang: morse.LangEnglish},
		{
			expected: "A B C D E F G H I J K L M N O P Q R S T U V W X Y Z",
			input:    ".- / -... / -.-. / -.. / . / ..-. / --. / .... / .. / .--- / -.- / .-.. / -- / -. / --- / .--. / --.- / .-. / ... / - / ..- / ...- / .-- / -..- / -.-- / --..",
			lang:     "EN",
		},
		{expected: "1234567890", input: ".---- ..--- ...-- ....- ..... -.... --... ---.. ----. -----", lang: morse.LangEnglish},
		{expected: ".,?:", input: ".-.-.- --..-- ..--.. ---...", lang: morse.LangEnglish},
		{expected: "ㅎㅏㄴㄱㅡㄹ", input: ".--- . ..-. .-.. -.. ...-", lang: morse.LangKorean},
	}

	for _, scenario := range scenarios {
		out := morse.FromMorseCode(scenario.input, scenario.lang)
		assert.Equalf(t, scenario.expected, out, "Input %s", scenario.input)
	}
}
