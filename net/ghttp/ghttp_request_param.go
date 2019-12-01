// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"mime/multipart"
)

// Get is alias of GetRequest, which is one of the most commonly used functions for
// retrieving parameter.
// See GetRequest.
func (r *Request) Get(key string, def ...interface{}) interface{} {
	return r.GetRequest(key, def...)
}

// GetVar is alis of GetRequestVar.
// See GetRequestVar.
func (r *Request) GetVar(key string, def ...interface{}) *gvar.Var {
	return r.GetRequestVar(key, def...)
}

// GetRaw is alias of GetBody.
// See GetBody.
func (r *Request) GetRaw() []byte {
	return r.GetBody()
}

// GetRawString is alias of GetBodyString.
// See GetBodyString.
func (r *Request) GetRawString() string {
	return r.GetBodyString()
}

// GetRaw retrieves and returns request body content as bytes.
func (r *Request) GetBody() []byte {
	if r.bodyContent == nil {
		r.bodyContent, _ = ioutil.ReadAll(r.Body)
	}
	return r.bodyContent
}

// GetRawString retrieves and returns request body content as string.
func (r *Request) GetBodyString() string {
	return gconv.UnsafeBytesToStr(r.GetRaw())
}

// GetJson parses current request content as JSON format, and returns the JSON object.
// Note that the request content is read from request BODY, not from any field of FORM.
func (r *Request) GetJson() (*gjson.Json, error) {
	return gjson.LoadJson(r.GetRaw())
}

func (r *Request) GetString(key string, def ...interface{}) string {
	return r.GetRequestString(key, def...)
}

func (r *Request) GetBool(key string, def ...interface{}) bool {
	return r.GetRequestBool(key, def...)
}

func (r *Request) GetInt(key string, def ...interface{}) int {
	return r.GetRequestInt(key, def...)
}

func (r *Request) GetInt32(key string, def ...interface{}) int32 {
	return r.GetRequestInt32(key, def...)
}

func (r *Request) GetInt64(key string, def ...interface{}) int64 {
	return r.GetRequestInt64(key, def...)
}

func (r *Request) GetInts(key string, def ...interface{}) []int {
	return r.GetRequestInts(key, def...)
}

func (r *Request) GetUint(key string, def ...interface{}) uint {
	return r.GetRequestUint(key, def...)
}

func (r *Request) GetUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestUint32(key, def...)
}

func (r *Request) GetUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestUint64(key, def...)
}

func (r *Request) GetFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestFloat32(key, def...)
}

func (r *Request) GetFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestFloat64(key, def...)
}

func (r *Request) GetFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestFloats(key, def...)
}

func (r *Request) GetArray(key string, def ...interface{}) []string {
	return r.GetRequestArray(key, def...)
}

func (r *Request) GetStrings(key string, def ...interface{}) []string {
	return r.GetRequestStrings(key, def...)
}

func (r *Request) GetInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestInterfaces(key, def...)
}

func (r *Request) GetMap(def ...map[string]interface{}) map[string]interface{} {
	return r.GetRequestMap(def...)
}

func (r *Request) GetMapStrStr(def ...map[string]interface{}) map[string]string {
	return r.GetRequestMapStrStr(def...)
}

// GetToStruct is alias of GetRequestToStruct.
// See GetRequestToStruct.
func (r *Request) GetToStruct(pointer interface{}, mapping ...map[string]string) error {
	return r.GetRequestToStruct(pointer, mapping...)
}

// ParseQuery parses query string into r.queryMap.
func (r *Request) ParseQuery() {
	if r.parsedQuery {
		return
	}
	r.parsedQuery = true
	if r.URL.RawQuery != "" {
		var err error
		r.queryMap, err = gstr.Parse(r.URL.RawQuery)
		if err != nil {
			panic(err)
		}
	}
}

// ParseRaw parses the request raw data into r.rawMap.
func (r *Request) ParseBody() {
	if r.parsedBody {
		return
	}
	r.parsedBody = true
	if body := r.GetBodyString(); len(body) > 0 {
		r.bodyMap, _ = gstr.Parse(body)
	}
}

// ParseForm parses the request form for HTTP method PUT, POST, PATCH.
// The form data is pared into r.formMap.
func (r *Request) ParseForm() {
	if r.parsedForm {
		return
	}
	r.parsedForm = true
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		var err error
		if gstr.Contains(contentType, "multipart/") {
			// multipart/form-data, multipart/mixed
			if err = r.ParseMultipartForm(r.Server.config.FormParsingMemory); err != nil {
				panic(err)
			}
		} else if gstr.Contains(contentType, "form") {
			// application/x-www-form-urlencoded
			if err = r.Request.ParseForm(); err != nil {
				panic(err)
			}
		}
		if len(r.PostForm) > 0 {
			// Re-parse the form data using united parsing way.
			params := ""
			for name, values := range r.PostForm {
				if len(values) == 1 {
					if len(params) > 0 {
						params += "&"
					}
					params += name + "=" + gurl.Encode(values[0])
				} else {
					if len(name) > 2 && name[len(name)-2:] == "[]" {
						name = name[:len(name)-2]
						for _, v := range values {
							if len(params) > 0 {
								params += "&"
							}
							params += name + "[]=" + gurl.Encode(v)
						}
					} else {
						if len(params) > 0 {
							params += "&"
						}
						params += name + "=" + gurl.Encode(values[len(values)-1])
					}
				}
			}
			if r.formMap, err = gstr.Parse(params); err != nil {
				panic(err)
			}
		} else {
			r.ParseBody()
			if len(r.bodyMap) > 0 {
				r.formMap = r.bodyMap
			}
		}
	}
}

// GetMultipartForm parses and returns the form as multipart form.
func (r *Request) GetMultipartForm() *multipart.Form {
	r.ParseForm()
	return r.MultipartForm
}

// GetMultipartFiles returns the post files array.
// Note that the request form should be type of multipart.
func (r *Request) GetMultipartFiles(name string) []*multipart.FileHeader {
	form := r.GetMultipartForm()
	if form == nil {
		return nil
	}
	if v := form.File[name]; len(v) > 0 {
		return v
	}
	// Support "name[]" as array parameter.
	if v := form.File[name+"[]"]; len(v) > 0 {
		return v
	}
	return nil
}
