package cdn

import (
	"fmt"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"os"
	"sync"
)

const (
	POST     = "https://ayang.ink/posts/"
	CATEGORY = "https://ayang.ink/categories/"
	TAG      = "https://ayang.ink/tags/"
)

var (
	// 文章跟新时，需要更改其他页面
	updatePath = []string{POST, CATEGORY, TAG}
)

type Cdn interface {
	Add(urlEncode bool, urls ...string)
	Update(urlEncode bool, urls ...string)
}

type Client struct {
	once   sync.Once
	client *cdn.Client
}

func (c *Client) initClient() {
	if c == nil {
		panic("未初始化 Client")
	}
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	secretKey := os.Getenv("CDN_SECRET_KEY")
	credential := common.NewCredential(
		"AKID4xjHckecMu9D4g4FEzWh4E9qNglbDbqs",
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	//实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cdn.NewClient(credential, "", cpf)
	c.client = client
}

// Add 网站 post 文章
func (c *Client) Add(urlEncode bool, urls ...string) {
	c.once.Do(c.initClient)
	// 缓存预热
	c.pushUrlsCache(urls, urlEncode)
	// 刷新目录，包括 TAG 等
	c.purgePathCache(updatePath, urlEncode)
}

// Update 更该文章
func (c *Client) Update(urlEncode bool, urls ...string) {
	c.once.Do(c.initClient)
	// 只刷新这篇文章
	c.purgePathCache(urls, urlEncode)
}

// purgePathCache 刷新目录
func (c *Client) purgePathCache(urls []string, urlEncode bool) {
	urlsPoint := strToStrPoint(urls)
	boolPoint := new(bool)
	if urlEncode {
		*boolPoint = true
	} else {
		*boolPoint = false
	}
	str := new(string)
	*str = "delete"
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPurgePathCacheRequest()
	request.Paths = append(request.Paths, urlsPoint...)
	request.UrlEncode = boolPoint
	request.FlushType = str

	// 返回的resp是一个PurgePathCacheResponse的实例，与请求对象对应
	response, err := c.client.PurgePathCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}

// pushUrlsCache 预热 url
func (c *Client) pushUrlsCache(urls []string, urlEncode bool) {
	urlsPoint := strToStrPoint(urls)
	boolPoint := new(bool)
	if urlEncode {
		*boolPoint = true
	} else {
		*boolPoint = false
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPushUrlsCacheRequest()

	request.Urls = append(request.Urls, urlsPoint...)
	request.UrlEncode = boolPoint

	// 返回的resp是一个PurgePathCacheResponse的实例，与请求对象对应
	response, err := c.client.PushUrlsCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}

// strToStrPoint 不知道为啥腾讯云全是指针
func strToStrPoint(str []string) (strPoint []*string) {
	strPoint = make([]*string, len(str))
	for i, _ := range str {
		strPoint[i] = &str[i]
	}
	return
}
