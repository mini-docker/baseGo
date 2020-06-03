package server

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"model"
	"net/http"
	"red_admin/app/middleware/validate"
	"strings"
	"time"
	"unsafe"
)

const (
	defaultExpire        = 60 * time.Second
	defaultClearInterval = 10 * time.Minute
	defaultContentType   = MIMEApplicationJSONCharsetUTF8
	defaultMarkID        = "#"
	defaultEmptyParams   = "#"
	cacheKeyHashPattern  = "%s"
	cacheKeyPattern      = "api:cache:%s:%s:%s:%s"
	ctxWriteBackKey      = "api:cache:back:reply"

	// 头部协议标志
	protocolMarkTextHtml   = "1"
	protocolMarkTextPlain  = "2"
	protocolMarkJson       = "3"
	protocolMarkXml        = "4"
	protocolMarkJavaScript = "5"
	protocolMarkProtobuf   = "6"

	// 错误协议标志
	protocolMarkError = "e"
)

var (
	// 全局接口缓存接口
	apiCache ApiCache
	// 全局缓存控制开关
	cacheSwitch bool
	// 协议标志映射的HTTP头信息
	protocolMarksHttpMapping = map[string]string{
		protocolMarkTextHtml:   MIMETextHTMLCharsetUTF8,
		protocolMarkTextPlain:  MIMETextPlainCharsetUTF8,
		protocolMarkJson:       MIMEApplicationJSONCharsetUTF8,
		protocolMarkXml:        MIMEApplicationXMLCharsetUTF8,
		protocolMarkJavaScript: MIMEApplicationJavaScriptCharsetUTF8,
		protocolMarkProtobuf:   MIMEApplicationProtobuf,
	}
)

// 接口缓存接口
// 为了方便后面使用其他的缓存进行扩展
type ApiCache interface {
	// 实例化缓存对象
	// defExpire:默认缓存过期时间
	// clearInterval:清除失效缓存的间隔时间
	New(defExpire, clearInterval time.Duration)
	// 读取缓存
	// key:缓存键
	Read(key string) (interface{}, bool)
	// 写缓存
	// key:缓存键
	// value:缓存值
	// expire:指定该缓存的单独的过期时间
	Write(key string, value interface{}, expire time.Duration)
}

// 缓存选项
type cacheOptions struct {
	defExpire     time.Duration
	clearInterval time.Duration
	cacheSwitch   bool
}

type CacheOptions func(*cacheOptions)

// 设置默认的缓存超时时间
func WithDefaultExpire(defExpire time.Duration) CacheOptions {
	return func(o *cacheOptions) {
		o.defExpire = defExpire
	}
}

// 设置清除过期缓存的时间
func WithClearInterval(clearInterval time.Duration) CacheOptions {
	return func(o *cacheOptions) {
		o.clearInterval = clearInterval
	}
}

// 设置缓存开关
func WithCacheSwitch(on bool) CacheOptions {
	return func(o *cacheOptions) {
		o.cacheSwitch = on
	}
}

// 初始化接口缓存
func ApiCacheInit(cache ApiCache, opts ...CacheOptions) {
	if cache == nil {
		panic("the api cache implement instance must not be empty")
	}
	o := &cacheOptions{
		// 默认开启缓存
		cacheSwitch: true,
	}
	if opts != nil {
		for _, opt := range opts {
			opt(o)
		}
	}
	if o.defExpire <= 0 {
		o.defExpire = defaultExpire
	}
	if o.clearInterval <= 0 {
		o.clearInterval = defaultClearInterval
	}
	apiCache = cache
	cacheSwitch = o.cacheSwitch
	apiCache.New(o.defExpire, o.clearInterval)
}

// 选项
type options struct {
	expire      time.Duration // 缓存过期时间
	contentType string        // 返回头信息
	handler     HandlerFunc   // 控制器
}

type Options func(o *options)

// 设置缓存过期选项
func Expire(expire time.Duration) Options {
	return func(o *options) {
		o.expire = expire
	}
}

// 设置处理器选项
func Handler(h HandlerFunc) Options {
	return func(o *options) {
		o.handler = h
	}
}

// 缓存包装
type cacheWrapper struct {
	cache ApiCache // go-cache缓存对象
	opt   *options // 缓存选项
}

// 获取包装后的HandlerFunc
func ApiCacheWrapper(opts ...Options) HandlerFunc {
	if apiCache == nil {
		panic("the api cache is not init, please do the fun:'ApiCacheInit(...)'")
	}
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.handler == nil {
		panic("the handle func must not be empty")
	}
	if o.contentType == "" {
		o.contentType = defaultContentType
	}
	if o.expire <= 0 {
		o.expire = defaultExpire
	}
	wrapper := &cacheWrapper{cache: apiCache, opt: o}
	return wrapper.handle
}

// 执行处理器方法
func (cw *cacheWrapper) handle(ctx Context) error {

	// 缓存开关
	// 关闭开关时不做任何缓存处理
	if !cacheSwitch {
		return cw.opt.handler(ctx)
	}

	r := ctx.Request()
	w := ctx.Response()

	var params string
	method := r.Method
	path := r.URL.Path

	// 获取标志ID
	// 默认使用用户的sessionId
	// 如果没有的话说明是没有授权的路由,无需标志ID,使用默认符号#号代替
	// 对获取到的sessionID需要处理,将uuid部分去除,只保留用户ID和设备ID用":"连接
	var markID string
	user := ctx.User()
	if user == nil || user.Info() == nil {
		markID = defaultMarkID
	} else {
		session, ok := user.Info().(*model.AdminSession)
		if !ok || session == nil || session.SessionId == "" {
			markID = defaultMarkID
		} else {
			bf := bytes.Buffer{}
			elems := strings.Split(session.SessionId, "_")
			for i, elem := range elems {
				if i == 0 {
					continue
				}
				if i == len(elems)-1 {
					bf.WriteString(elem)
				} else {
					bf.WriteString(elem)
					bf.WriteString(":")
				}
			}
			markID = bf.String()
		}
	}

	// GET参数在路径上
	// POST参数在body中
	switch method {
	case GET:
		params = r.URL.RawQuery
	case POST:
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewHTTPError(http.StatusInternalServerError)
		} else {
			r.Body = ioutil.NopCloser(bytes.NewReader(bs))
		}
		params = B2S(bs)
	}

	// 开始命中本地缓存
	// 如果命中,直接返回缓存数据
	// 如果没有命中,则执行控制器方法,并将控制器方法中返回数据回写缓存
	cacheKey := cacheKey(markID, method, path, params)
	cacheValue, hit := cw.cache.Read(cacheKey)
	// 命中缓存,开始解析缓存并返回
	if hit {
		b, ok := cacheValue.([]byte)
		if !ok {
			return NewHTTPError(http.StatusInternalServerError)
		}
		log.Println("命中本地缓存, KEY:", cacheKey, ", VALUE:", B2S(b))
		// 存储的数据中前1位是一个协议位
		// 主要是区分返回的一个头部信息
		// 协议为:获取出的数据的第一位为返回头"ContentType"类型,string类型取值为"1","2","3","4","5","6";[]byte类型取值为[49],[50],[51],[52],[53],[54]
		// 其string值对应的"ContentType"类型为：
		// 		"1":"text/html; charset=UTF-8" 即echo中的 MIMETextHTMLCharsetUTF8
		// 		"2":"text/plain; charset=UTF-8" 即echo中的 MIMETextPlainCharsetUTF8
		// 		"3":"application/json; charset=UTF-8" 即echo中的 MIMEApplicationJSONCharsetUTF8
		// 		"4":"application/xml; charset=UTF-8" 即echo中的 MIMEApplicationXMLCharsetUTF8
		//		"5":"application/javascript; charset=UTF-8" 即echo中的 MIMEApplicationJavaScriptCharsetUTF8
		// 		"6":"application/protobuf" 即echo中的 MIMEApplicationProtobuf
		// 同时这个协议位除了表示返回头信息之外,同时还表示缓存的错误信息的类型,string类型取值为"e";[]byte取值为[101]
		if b != nil && len(b) > 0 {
			// 协议位
			protocolBit := B2S(b[0:1])
			// 原始数据
			replies := b[1:]
			// 先在头协议中查找
			contentType, ok := protocolMarksHttpMapping[protocolBit]
			if !ok {
				// 没有找到协议位在头中的定义，再查找错误信息
				switch protocolBit {
				case protocolMarkError:
					vErr := &validate.Err{}
					err := json.Unmarshal(replies, vErr)
					if err != nil {
						return NewHTTPError(http.StatusInternalServerError)
					}
					accept := r.Header.Get(HeaderAccept)
					switch {
					case strings.HasPrefix(accept, MIMEApplicationJSON):
						return ctx.JSON(http.StatusBadRequest, vErr)
					case strings.HasPrefix(accept, MIMEApplicationProtobuf):
						return ctx.Pb(http.StatusBadRequest, &validate.Error{Code: int64(vErr.Code), Msg: vErr.Msg})
					default:
						return ctx.JSON(http.StatusBadRequest, vErr)
					}
				}
			} else {
				w.Header().Set(HeaderContentType, contentType)
				w.WriteHeader(http.StatusOK)
				w.Write(replies)
				return nil
			}
		}
	}

	// 没有命中缓存,执行处理函数
	// 并将处理数据进行处理后存入缓存
	err := cw.opt.handler(ctx)
	if err != nil {
		// 不是服务器错误时,需要将错误也一起缓存,避免发生缓存穿透
		// 当前业务返回的错误都是使用"validate.Err"
		// 所以这里只需要缓存"validate.Err"的错误信息
		// 错误的协议标志为"e"
		switch err.(type) {
		case validate.Err, *validate.Err:
			vErrBs, _ := json.Marshal(err)
			writeBuffer := bytes.NewBuffer(S2B(protocolMarkError))
			writeBuffer.Write(vErrBs)
			cw.cache.Write(cacheKey, writeBuffer.Bytes(), cw.opt.expire)
			log.Println("缓存没有命中,回写数据, KEY:", cacheKey, ", VALUE:", writeBuffer.String())
		}
		return err
	}

	// 没有返回错误时提取返回数据
	value := ctx.Get(ctxWriteBackKey)
	if value == nil {
		return nil
	}

	reply, ok := value.([]byte)
	if !ok {
		return nil
	}

	// 开始回写缓存
	// 回写时需要加上协议位
	var contentTypeByte []byte
	contentType := w.Header().Get(HeaderContentType)
	for protocol, ctype := range protocolMarksHttpMapping {
		if contentType == ctype {
			contentTypeByte = S2B(protocol)
			break
		}
	}
	if len(contentTypeByte) == 0 {
		contentTypeByte = S2B(protocolMarkJson)
	}
	writeBuffer := bytes.NewBuffer(contentTypeByte)
	writeBuffer.Write(reply)
	cw.cache.Write(cacheKey, writeBuffer.Bytes(), cw.opt.expire)
	log.Println("缓存没有命中,回写数据, KEY:", cacheKey, ", VALUE:", writeBuffer.String())
	return nil
}

// 将指定的数据进行MD5混淆得到缓存KEY
func cacheKey(markId, method, path, params string) string {
	var hashStr string
	if params == "" {
		hashStr = defaultEmptyParams
	} else {
		hash := md5.New()
		hash.Write(S2B(params))
		hashed := hash.Sum(nil)
		hashStr = hex.EncodeToString(hashed)
	}
	return fmt.Sprintf(cacheKeyPattern, method, path, markId, hashStr)
}

// 字符串转数组，提升性能
func S2B(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 数组转字符串，提升性能
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
