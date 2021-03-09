package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/pkg/context"
	"github.com/fztcjjl/zim/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tiaotiao/mapstruct"
	"io/ioutil"
	"mime"
	"sort"
	"strings"
)

type AppHeader struct {
	Token string `json:"token"` // 令牌
	Nonce string `json:"nonce"` // 随机字符串
	Sign  string `json:"sign"`  // 签名
}

func CheckSign() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx := context.New(gctx)

		h := &AppHeader{}
		ctx.ShouldBindHeader(h)

		mapBody := make(map[string]interface{})

		ct := gctx.Request.Header.Get("Content-Type")
		ct, _, _ = mime.ParseMediaType(ct)
		if ct == "application/json" {
			body, _ := ioutil.ReadAll(gctx.Request.Body)
			gctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			log.Debug(string(body))
			dec := json.NewDecoder(bytes.NewReader(body))
			dec.UseNumber()
			dec.Decode(&mapBody)
			log.Debug(mapBody)
		} else {
			for k, vs := range ctx.GetForm() {
				v := ""
				if len(vs) > 0 {
					v = vs[0]
				}
				if v != "" {
					mapBody[k] = v
				}
			}
		}

		sign := sign(h, mapBody, "")
		log.Debugf("计算得到的sign:%s", sign)
		log.Debugf("客户端上传的sign:%s", h.Sign)
		if h.Sign != sign {
			ctx.ResponseError(errors.ErrSign)
			gctx.Abort()
		}

		gctx.Next()
	}
}

func sign(h *AppHeader, mapBody map[string]interface{}, secret string) string {
	mapAll := make(map[string]interface{})
	mapHeader := mapstruct.Struct2MapTag(h, "json")
	delete(mapHeader, "sign")

	var keys []string
	for k := range mapHeader {
		v := cast.ToString(mapHeader[k])
		if v != "" {
			keys = append(keys, k)
			mapAll[k] = v
		}
	}

	for k, v := range mapBody {
		keys = append(keys, k)
		mapAll[k] = v

	}
	sort.Strings(keys)

	var v string
	var plist []string
	for _, k := range keys {
		switch obj := mapAll[k].(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32,
			float32, float64:
			v = cast.ToString(obj)
			log.Debug(k, v)
			break
		case json.Number:
			v = string(obj)
			break
		case bool:
			v = cast.ToString(obj)
			break
		case string:
			v = obj
		case []interface{}:
			v = concatArray(obj)
			break
		case map[string]interface{}:
			v = concatMap(obj)
			break
		}

		if v != "" {
			plist = append(plist, k+"="+v)
		}
	}
	var src = strings.Join(plist, "&")
	log.Debug(src)
	src += secret
	bs := sha256.Sum256([]byte(src))

	return hex.EncodeToString(bs[:])

}

func concatArray(a []interface{}) string {
	if len(a) == 0 {
		return ""
	}
	str := "["
	first := true
	for _, o := range a {
		if !first {
			str += ","
		}

		first = false
		switch m := o.(type) {
		case map[string]interface{}:
			str += concatMap(m)
			break
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32,
			float32, float64:
			v := cast.ToString(m)
			str += v
			break
		case json.Number:
			str += string(m)
			break
		case bool:
			str += cast.ToString(m)
			break
		case string:
			str += m
			break
		default:
			log.Debug(m)
			break
		}
	}
	str += "]"
	return str
}

func concatMap(m map[string]interface{}) string {
	str := "{"
	var keys []string
	for key := range m {
		if m[key] != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var v string
	var plist []string
	for _, key := range keys {
		v = ""
		switch obj := m[key].(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32,
			float32, float64:
			v = cast.ToString(obj)
			break
		case json.Number:
			v = string(obj)
			break
		case bool:
			v = cast.ToString(obj)
			break
		case string:
			v = obj
		case []interface{}:
			v = concatArray(obj)
			break
		case map[string]interface{}:
			v = concatMap(obj)
			break
		}

		if v != "" {
			plist = append(plist, key+"="+v)
		}
	}

	str += strings.Join(plist, "&")
	str += "}"
	return str
}
