package object

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/clbanning/mxj/v2"
	"strings"
)

func Str2Xml(in string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(in))
	return b.String()
}

func Map2Xml(obj *HashMap, isSub bool) (strXML string) {

	for k, v := range *obj {
		switch v.(type) {
		case string:
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k)
			break
		case int, int8, int16, int32, int64:
			strXML = strXML + fmt.Sprintf("<%s>%d</%s>", k, v, k)
			break
		case float32:
		case float64:
			strXML = strXML + fmt.Sprintf("<%s>%f</%s>", k, v, k)
			break
		case []*HashMap:
			strXML = strXML + fmt.Sprintf("<%s>", k)
			for _, subV := range v.([]*HashMap) {
				strXML = strXML + fmt.Sprintf("%s", Map2Xml(subV, true))
			}
			strXML = strXML + fmt.Sprintf("</%s>", k)
			break
		case []HashMap:
			strXML = strXML + fmt.Sprintf("<%s>", k)
			for _, subV := range v.([]HashMap) {
				strXML = strXML + fmt.Sprintf("%s", Map2Xml(&subV, true))
			}
			strXML = strXML + fmt.Sprintf("</%s>", k)
			break
		case *HashMap:
			strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(v.(*HashMap), true), k)
			break
		case HashMap:
			val := v.(HashMap)
			strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(&val, true), k)
			break
		case interface{}:
			b, _ := json.Marshal(v)
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, string(b), k)
			break
		}
	}
	if isSub {
		return strXML
	} else {
		return "<xml>" + strXML + "</xml>"
	}
}

func StringMap2Xml(obj *StringMap) (strXML string) {

	for k, v := range *obj {
		strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k)
	}
	return "<xml>" + strXML + "</xml>"
}

func Xml2HashMap(b []byte) (m map[string]interface{}, err error) {
	mv, err := mxj.NewMapXml(b) // unmarshal
	m = map[string]interface{}(mv)

	// 我觉得这个判断加上是不是更好一点，因为我需要的数据其实是xml对象里面的，不然的话使用者还需要额外处理一下，感觉没必要
	if _, ok := m["xml"]; ok {
		m = m["xml"].(map[string]interface{})
	}

	return m, err
}

func Xml2Map(b []byte) (m HashMap, err error) {

	decoder := xml.NewDecoder(bytes.NewReader(b))
	m = make(HashMap)
	tag := ""
	for {
		token, err := decoder.Token()

		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local != "xml" {
				tag = t.Name.Local
			} else {
				tag = ""
			}
			break
		case xml.EndElement:
			break
		case xml.CharData:
			data := strings.TrimSpace(string(t))
			fmt2.Dump(data)
			if len(data) != 0 {
				m[tag] = data
			}
			break
		}
	}
	return m, err
}
