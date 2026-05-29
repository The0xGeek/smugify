package utils

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// regex cache ==========
var (
	stringRegex = regexp.MustCompile(`"(\\.|[^"\\])*"|'(\\.|[^'\\])*'`)
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// options ==========
type Options struct {
	EncodeStrings    bool
	HexEscape        bool
	RemoveComments   bool
	MinifyCode       bool
	AdvancedEncoding bool
}

//obfuscator ==========
type Obfuscator struct {
	options Options
	nameMap map[string]string
	counter int
}

// constructor ==========
func NewObfuscator(options Options) *Obfuscator {
	return &Obfuscator{
		options: options,
		nameMap: make(map[string]string),
		counter: 0,
	}
}

// main pipeline ==========
func (o *Obfuscator) Obfuscate(code string) string {
	transformers := []func(string) string{}

	if o.options.RemoveComments {
		transformers = append(transformers, o.removeComments)
	}

	if o.options.EncodeStrings {
		transformers = append(transformers, o.encodeStrings)
	}

	if o.options.HexEscape {
		transformers = append(transformers, o.hexEscapeStrings)
	}

	if o.options.AdvancedEncoding {
		transformers = append(transformers, o.advancedEncode)
	}

	if o.options.MinifyCode {
		transformers = append(transformers, o.minify)
	}

	result := code

	for _, transformer := range transformers {
		result = transformer(result)
	}

	return result
}

// remove comments ==========
func (o *Obfuscator) removeComments(code string) string {

	var result strings.Builder

	inString := false
	stringChar := byte(0)

	for i := 0; i < len(code); i++ {

		if inString {

			if code[i] == stringChar && code[i-1] != '\\' {
				inString = false
			}

			result.WriteByte(code[i])
			continue
		}

		if code[i] == '"' || code[i] == '\'' || code[i] == '`' {
			inString = true
			stringChar = code[i]

			result.WriteByte(code[i])
			continue
		}

		if i+1 < len(code) {

			// single-line
			if code[i] == '/' && code[i+1] == '/' {

				for i < len(code) && code[i] != '\n' {
					i++
				}

				continue
			}

			// multi-line
			if code[i] == '/' && code[i+1] == '*' {

				i += 2

				for i+1 < len(code) &&
					!(code[i] == '*' && code[i+1] == '/') {
					i++
				}

				i++
				continue
			}
		}

		result.WriteByte(code[i])
	}

	return result.String()
}

// string encoder ==========
func (o *Obfuscator) encodeStrings(code string) string {
	found := false

	result := stringRegex.ReplaceAllStringFunc(code, func(match string) string {
		found = true

		content := match[1 : len(match)-1]

		encoded := base64.StdEncoding.EncodeToString([]byte(content))

		return `_d("` + encoded + `")`
	})

	if found {
		result = `function _d(s){return atob(s)};` + result
	}

	return result
}

// hex escaping ==========
func (o *Obfuscator) hexEscapeStrings(code string) string {
	return stringRegex.ReplaceAllStringFunc(code, func(match string) string {
		content := match[1 : len(match)-1]

		var builder strings.Builder

		for _, c := range content {
			builder.WriteString(fmt.Sprintf("\\x%02x", c))
		}

		return `"` + builder.String() + `"`
	})
}

// advanced multi-layer encoding ==========
func (o *Obfuscator) advancedEncode(code string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(code))

	var hex strings.Builder

	for _, c := range encoded {
		hex.WriteString(fmt.Sprintf("%02x", c))
	}

	return fmt.Sprintf(
		`(()=>{const d=h=>atob(h.match(/../g).map(x=>String.fromCharCode(parseInt(x,16))).join(""));eval(d("%s"))})();`,
		hex.String(),
	)
}

// minifier ==========
func (o *Obfuscator) minify(code string) string {
	code = whitespaceRegex.ReplaceAllString(code, " ")

	replacer := strings.NewReplacer(
		" {", "{",
		"{ ", "{",
		" }", "}",
		"} ", "}",
		"; ", ";",
		" ;", ";",
	)

	return replacer.Replace(strings.TrimSpace(code))
}