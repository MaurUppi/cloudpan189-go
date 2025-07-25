// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"compress/gzip"
	"flag"
	"io"
	"io/ioutil"
	"net/http/cookiejar"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// TrimPathPrefix 去除目录的前缀
func TrimPathPrefix(path, prefixPath string) string {
	if prefixPath == "/" {
		return path
	}
	return strings.TrimPrefix(path, prefixPath)
}

// ContainsString 检测字符串是否在字符串数组里
func ContainsString(ss []string, s string) bool {
	for k := range ss {
		if ss[k] == s {
			return true
		}
	}
	return false
}

// GetURLCookieString 返回cookie字串
func GetURLCookieString(urlString string, jar *cookiejar.Jar) string {
	u, _ := url.Parse(urlString)
	cookies := jar.Cookies(u)
	cookieString := ""
	for _, v := range cookies {
		cookieString += v.String() + "; "
	}
	cookieString = strings.TrimRight(cookieString, "; ")
	return cookieString
}

// DecompressGZIP 对 io.Reader 数据, 进行 gzip 解压
func DecompressGZIP(r io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	gzipReader.Close()
	return ioutil.ReadAll(gzipReader)
}

// FlagProvided 检测命令行是否提供名为 name 的 flag, 支持多个name(names)
func FlagProvided(names ...string) bool {
	if len(names) == 0 {
		return false
	}
	var targetFlag *flag.Flag
	for _, name := range names {
		targetFlag = flag.Lookup(name)
		if targetFlag == nil {
			return false
		}
		if targetFlag.DefValue == targetFlag.Value.String() {
			return false
		}
	}
	return true
}

// Trigger 用于触发事件
func Trigger(f func()) {
	if f == nil {
		return
	}
	go f()
}

// TriggerOnSync 用于触发事件, 同步触发
func TriggerOnSync(f func()) {
	if f == nil {
		return
	}
	f()
}

func ParseVersionNum(versionStr string) int {
	versionStr = strings.ReplaceAll(versionStr, "-dev", "")
	versionStr = strings.ReplaceAll(versionStr, "v", "")
	versionParts := strings.Split(versionStr, ".")
	verNum := parseInt(versionParts[0])*1e4 + parseInt(versionParts[1])*1e2 + parseInt(versionParts[2])
	return verNum
}
func parseInt(numStr string) int {
	num, e := strconv.Atoi(numStr)
	if e != nil {
		return 0
	}
	return num
}

// IsExcludeFile 是否是指定排除的文件
func IsExcludeFile(filePath string, excludeNames *[]string) bool {
	if excludeNames == nil || len(*excludeNames) == 0 {
		return false
	}

	for _, pattern := range *excludeNames {
		fileName := path.Base(strings.ReplaceAll(filePath, "\\", "/"))
		m, _ := regexp.MatchString(pattern, fileName)
		if m {
			return true
		}
	}
	return false
}
