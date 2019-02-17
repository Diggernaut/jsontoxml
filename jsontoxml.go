package jsontoxml

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/buger/jsonparser"
)

var cleanKeys map[string]string
var re *regexp.Regexp
var sre *strings.Replacer
var ident int
var doident bool

// SetKeysForClean you can use this func to set key that
// you want to replace in result xml
func SetKeysForClean(keys map[string]string) {}

// SetSymbolsForClean you can use this func to set symbols that
// you want to replace in result xml
func SetSymbolsForClean(symbols map[string]string) {}

func init() {
	re = regexp.MustCompile(`^[A-Za-z]+`)
	sre = strings.NewReplacer(":", "_",
		"@", "_", "/", "_", "\\", "_", "?", "_", "$", "_", ".", "_")
	cleanKeys = make(map[string]string)
	cleanKeys["html"] = "html_safe"
	cleanKeys["link"] = "url"
	cleanKeys["caption"] = "caption_safe"
	cleanKeys["body"] = "body_safe"
	cleanKeys["area"] = "area_safe"
	cleanKeys["base"] = "base_safe"
	cleanKeys["br"] = "br_safe"
	cleanKeys["col"] = "col_safe"
	cleanKeys["command"] = "command_safe"
	cleanKeys["embed"] = "embed_safe"
	cleanKeys["hr"] = "hr_safe"
	cleanKeys["img"] = "img_safe"
	cleanKeys["input"] = "input_safe"
	cleanKeys["keygen"] = "keygen_safe"
	cleanKeys["link"] = "link_safe"
	cleanKeys["meta"] = "meta_safe"
	cleanKeys["param"] = "param_safe"
	cleanKeys["source"] = "source_safe"
	cleanKeys["track"] = "track_safe"
	cleanKeys["wbr"] = "wbr_safe"
	cleanKeys["image"] = "image_safe"
	doident = true
<<<<<<< HEAD
	ident = 0
=======
>>>>>>> 5a63bc002c901951d2bf9e879767e98010471d11
}
func clear(b []byte) string {
	key := string(b)
	key = SpaceMap(key)
	key = strings.TrimSpace(key)
	key = sre.Replace(key)
	key = strings.ToLower(key)
	if !re.MatchString(key) {
		key = "safe_" + key
	}
	if v, ok := cleanKeys[key]; ok {
		return v
	}

	return key
}

// Convert - call it with you json to get xml string
func Convert(json []byte, k string) ([]byte, error) {
	json = []byte(strings.TrimSpace(string(json)))
	d, err := handleArray(json, "body")
	if err != nil {
		d, err := handleObject(json, "body")
		if err != nil {
			return nil, err
		}
		return d, nil
	}
	return d, nil
}
func handleArray(data []byte, key string) ([]byte, error) {
	var b bytes.Buffer
	var err error
	handlerArray := func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var s []byte
		switch dataType {
		case jsonparser.Object:
			if string(key) == "" {
				key = "element"
			}
			s, err = handleObject(value, key)
			if err != nil {
				fmt.Println(err)
			}

			b.Write(s)
		case jsonparser.Array:
			if string(key) == "" {
				key = "element"
			}
			s, err = handleArray(value, key)
			if err != nil {
				fmt.Println(err)
			}
			b.Write(s)
		default:
			b.Write(handleScalar(value, []byte(clear([]byte(key)))))
		}
	}
	err = jsonparser.ArrayEach(data, handlerArray)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
func handleObject(data []byte, key string) ([]byte, error) {
	var b bytes.Buffer

	if key != "" {
		if doident {
			b.WriteString("\n")
			b.WriteString(strings.Repeat("\t", ident))
			ident++
		}
		b.WriteString("<")
		b.WriteString(clear([]byte(key)))
		b.WriteString(">")

	} else {
		ident--
	}
	handlerObject := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch dataType {
		case jsonparser.Object:
			v, err := handleObject(value, string(key))
			if err != nil {
				return err
			}
			b.Write(v)
		case jsonparser.Array:
			if string(key) == "" {
				key = []byte("element")
			}
			v, err := handleArray(value, string(key))
			if err != nil {
				return err
			}
			b.Write(v)
		default:
			b.Write(handleScalar(value, key))
		}
		return nil
	}
	err := jsonparser.ObjectEach(data, handlerObject)
	if err != nil {
		return nil, err
	}
	if doident {
		ident--
		if ident < 0 {
			ident = 0
		}
		b.WriteString("\n")
		b.WriteString(strings.Repeat("\t", ident))

	}
	if key != "" {
		b.WriteString("</")
		b.WriteString(clear([]byte(key)))
		b.WriteString(">")
	}

	return b.Bytes(), nil
}
func handleScalar(data []byte, key []byte) []byte {
	var b string
	if doident {
		b += "\n"
		b += strings.Repeat("\t", ident)
	}
	b += "<"
	b += clear([]byte(key))
	b += ">"
	b += string(data)
	b += "</"
	b += clear([]byte(key))
	b += ">"

	return []byte(b)
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
