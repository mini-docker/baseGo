// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utility

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	r "math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	UNDERLINE      = "_"
	BEFORE_STRINGS = "s6h2it"
	REAR_STRINGS   = "h9l5d"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// PasswordEncrypt
func PasswordEncrypt(account, password string, siteId string) string {
	return Md5(siteId + account + password)
}

func Md5_16(str string) []byte {
	h := md5.New()
	h.Write([]byte(str))
	//strd := hex.EncodeToString(h.Sum(nil))
	//strd = strd[8:24]
	return h.Sum(nil)
}

func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// base64 encode
func Base64Encode_(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}

// base64 decode
func Base64Decode_(str string) ([]byte, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return s, e
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

func DesEncrypt(src, key []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	//src = ZeroPadding(src, bs)
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}

	return Base64Encode_(out), nil
}

func DesDecrypt(src, key []byte) ([]byte, error) {
	src, err := Base64Decode_(string(src))
	if err != nil {
		return nil, err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	//out = ZeroUnPadding(out)
	out = PKCS5UnPadding(out)

	return out, nil
}

func Encrypt3(key1 string, src []byte) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(key1)
	iv, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAA=")
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = PKCS5Padding(src, bs)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return Base64Encode(string(dst)), nil
}

func Decrypt3(key, src []byte) ([]byte, error) {
	vi, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAA=")
	src, err := Base64Decode_(string(src))
	if err != nil {
		return []byte{}, err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return []byte{}, err
	}
	blockMode := cipher.NewCBCDecrypter(block, vi)
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func DesEncryptIv(origData, key, iv []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	//origData = ZeroPadding(origData, )
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return Base64Encode_(crypted), nil
}
func DesDecryptIv(crypted, key, iv []byte) ([]byte, error) {
	crypted, err := Base64Decode_(string(crypted))
	if err != nil {
		return nil, err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// AESEncrypt encrypts text and given key with AES.
func AESEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// AESEncrypt encrypts text and given key with AES.
func AESEncryptCBC(key, text []byte) (string, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	ciphertext := PKCS5Padding(text, block.BlockSize())
	crypted := make([]byte, len(ciphertext))
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(crypted, ciphertext)
	return base64.StdEncoding.EncodeToString(crypted), base64.StdEncoding.EncodeToString(iv), nil
}

// AESDecrypt decrypts text and given key with AES.
func AESDecrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// IsLetter returns true if the 'l' is an English letter.
func IsLetter(l uint8) bool {
	n := (l | 0x20) - 'a'
	if n >= 0 && n < 26 {
		return true
	}
	return false
}

// Expand replaces {k} in template with match[k] or subs[atoi(k)] if k is not in match.
func Expand(template string, match map[string]string, subs ...string) string {
	var p []byte
	var i int
	for {
		i = strings.Index(template, "{")
		if i < 0 {
			break
		}
		p = append(p, template[:i]...)
		template = template[i+1:]
		i = strings.Index(template, "}")
		if s, ok := match[template[:i]]; ok {
			p = append(p, s...)
		} else {
			j, _ := strconv.Atoi(template[:i])
			if j >= len(subs) {
				p = append(p, []byte("Missing")...)
			} else {
				p = append(p, subs[j]...)
			}
		}
		template = template[i+1:]
	}
	p = append(p, template...)
	return string(p)
}

// Reverse s string, support unicode
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

func MysqlFilter(key string) string {
	key = strings.Replace(key, "'", "\\'", -1)
	key = strings.Replace(key, "\"", "\\\"", -1)
	return key
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

// PKCS7Padding pads as prescribed by the PKCS7 standard
func PKCS7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS7UnPadding unpads as prescribed by the PKCS7 standard
func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > aes.BlockSize || unpadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	pad := src[len(src)-unpadding:]
	for i := 0; i < unpadding; i++ {
		if pad[i] != byte(unpadding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

//3des ecb PKCS7
func EncryptDesECB(ori, key []byte) ([]byte, error) {

	block, err := des.NewTripleDESCipher(key[:24])
	if err != nil {
		return nil, err
	}
	ori = PKCS7Padding(ori)
	blockMode := NewECBEncrypter(block)
	crypted := make([]byte, len(ori))
	blockMode.CryptBlocks(crypted, ori)
	return crypted, nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

//返回ECB方式的加密器
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("dec_ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("dec_ecb: output smaller than input")
	}

	for len(src) > 0 {
		// Write to the dst
		x.b.Encrypt(dst[:x.blockSize], src[:x.blockSize])

		// Move to the next block
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func uniqid() string {
	r.Seed(time.Now().UnixNano())
	hash := md5.New()
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i := 0; i < 128; i++ {
		io.WriteString(hash, string(chars[r.Intn(len(chars))]))
	}

	return string(hash.Sum(nil))
}

func Generate() string {
	str := uniqid()

	return fmt.Sprintf("%x-%x-%x-%x-%x", str[0:4], str[4:6], str[6:8], str[8:10], str[10:16])
}

// 按字节截取字符串 utf-8不乱码
func SubstrByByte(str string, length int) string {
	bs := []byte(str)[:length]
	bl := 0
	for i := len(bs) - 1; i >= 0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i]&252 == 252:
				cl = 6
			case bs[i]&248 == 248:
				cl = 5
			case bs[i]&240 == 240:
				cl = 4
			case bs[i]&224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}

func SubString(str string, begin, length int) string {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)
	endstr := ""
	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	} else {
		endstr = ""
	}
	// 返回子串
	return string(rs[begin:end]) + endstr
}

// camel2Underline:snake string, XxYy to xx_yy , XxYY to xx_yy
func Camel2Underline(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// underline2Camel camel string, xx_yy to XxYy 大驼峰
func Underline2Camel(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 构造缓存 key
// key 自定义key
// from pc/wap
func GetKey(siteId string, siteIndexId string, key string) string {
	return fmt.Sprintf("%s%s%s", siteId, siteIndexId, key)
}

// 转换detail
func DetailConversion(fcId, playId int, playDetails string) string {
	type config struct {
		Id      []int    `json:"id"`
		Pos     []string `json:"pos"`
		PlayIds []int    `json:"playIds"`
	}

	configs := []config{
		config{
			// 时时彩/排列五/11选5
			Id: []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
				27, 28, 29, 30, 31, 32, 33, 34, 70, 107, 108, 109, 112, 113, 115, 116, 117, 118, 121},
			Pos:     []string{"万", "千", "百", "十", "个"},
			PlayIds: []int{81, 94, 100, 109, 126},
		}, config{
			// 福彩3D/上海时时乐/必赢3D/排列三
			Id:      []int{9, 67, 68, 69},
			Pos:     []string{"百", "十", "个"},
			PlayIds: []int{109, 169},
		}, config{
			// 七星彩
			Id:      []int{71},
			Pos:     []string{"千", "百", "十", "个"},
			PlayIds: []int{220, 222, 224},
		}, config{
			// 幸运赛车
			Id:      []int{76, 77, 78, 79, 80, 81},
			Pos:     []string{"冠", "亚", "季", "第四", "第五", "第六", "第七", "第八", "第九", "第十"},
			PlayIds: []int{174, 252, 253},
		}, config{
			// 快乐十分
			Id:      []int{55, 56, 57, 58, 59, 60, 61, 62, 63, 119},
			Pos:     []string{"第一", "第二", "第三", "第四", "第五", "第六", "第七", "第八"},
			PlayIds: []int{233, 236, 234},
		},
	}

	details := strings.Split(playDetails, "|")
	var detailss []string
	var ok1, ok2 bool
	for _, v := range configs {
		for _, vv := range v.Id {
			if vv == fcId {
				ok1 = true
			}
		}
		for _, vvv := range v.PlayIds {
			if vvv == playId {
				ok2 = true
			}
		}
		if ok1 && ok2 {
			if len(v.Pos) == len(details) {
				for k, p := range v.Pos {
					for l := range details {
						if details[k] != "" {
							_ = l
							detailss = append(detailss, p+":"+details[k])
							break
						}
					}
				}
			}
		}
	}

	if len(detailss) == 0 {
		return playDetails
	}

	return strings.Join(detailss, "|")
}

// 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}

// PasswordEncrypt
func NewPasswordEncrypt(account string, password string) string {
	return Md5(account + Md5(Md5(password)) + REAR_STRINGS)
}
