package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type ResponseData struct {
	Code          int      `json:"code"`
	Msg           string   `json:"msg"`
	DomainName    string   `json:"domain_name"`
	Status        string   `json:"status"`
	ContentLength string   `json:"centent_length"`
	RawCode       string   `json:"raw_code"`
	Data          []string `json:"data"` // 使用[]string来表示字符串切片
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func f(n *html.Node) *ResponseData {
	resp := ResponseData{
		Code: 0,
		Msg:  "查询成功",
		Data: []string{},
	}
	fmt.Println("标签为：")
	fmt.Println(n.Data)
	//if (n.Type == html.ElementNode) && (n.Data == "link" || n.Data == "script" || n.Data == "dev" || n.Data == "body" || n.Data == "html" || n.Data == "head" || n.Data == "meta" || n.Data == "image") {
	if n.Type == html.ElementNode && n.Data != "a" {
		fmt.Println("n.Attr:")
		fmt.Println(n.Attr)
		for _, a := range n.Attr {
			//if a.Key == "href" {
			//	parsedURL, err := url.Parse(a.Val)
			//	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
			//		resp.Data = append(resp.Data, parsedURL.String())
			//	}
			//}
			//if a.Key == "itemtype" {
			//	parsedURL, err := url.Parse(a.Val)
			//	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
			//		resp.Data = append(resp.Data, parsedURL.String())
			//	}
			//}
			//if a.Key == "src" {
			//	parsedURL, err := url.Parse(a.Val)
			//	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
			//		resp.Data = append(resp.Data, parsedURL.String())
			//	}
			//}
			parsedURL, err := url.Parse(a.Val)
			if err != nil {
				fmt.Println("url.Parse(a.Val)：报错")
				fmt.Println(err)
				resp.Code = -1
				return &resp
			}

			//if n.Data == "link" {
			//	fmt.Println("a.Val")
			//	fmt.Println(a.Val)
			//	fmt.Println("url.Parse(a.Val)")
			//	fmt.Println(url.Parse(a.Val))
			//	fmt.Println("parsedURL.Scheme")
			//	fmt.Println(parsedURL.Scheme)
			//	fmt.Println("parsedURL.Host")
			//	fmt.Println(parsedURL.Host)
			//	fmt.Println("parsedURL.String()")
			//	fmt.Println(parsedURL.String())
			//	host := parsedURL.String()
			//	fmt.Println(strings.HasPrefix(host, "http://"))
			//	fmt.Println(strings.HasPrefix(host, "https://"))
			//	fmt.Println(strings.HasPrefix(host, "//"))
			//}

			if err == nil && parsedURL.Host != "" {
				host := parsedURL.String()
				if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") || strings.HasPrefix(host, "//") {
					resp.Data = append(resp.Data, host)

				}
			}
		}
		// 如果是 script 标签，检查标签内的文本内容
		//if n.Data == "script" {
		//	var scriptText strings.Builder
		//	html.Render(&scriptText, n)
		//	re := regexp.MustCompile(`http?://[^\s]+`)
		//	matches := re.FindAllString(scriptText.String(), -1)
		//	for _, match := range matches {
		//		//links = append(links, match)
		//		resp.Data = append(resp.Data, match)
		//	}
		//}
		//// 如果是 script 标签，检查标签内的文本内容
		//if n.Data == "script" {
		//	var scriptText strings.Builder
		//	html.Render(&scriptText, n)
		//	re := regexp.MustCompile(`https?://[^\s]+`)
		//	matches := re.FindAllString(scriptText.String(), -1)
		//	for _, match := range matches {
		//		//links = append(links, match)
		//		resp.Data = append(resp.Data, match)
		//	}
		//}
		//
		//// 匹配标签之间文字信息里的链接
		//if n.Type == html.TextNode {
		//	text := n.Data
		//	re := regexp.MustCompile(`https?://[^\s]+`)
		//	matches := re.FindAllString(text, -1)
		//	for _, match := range matches {
		//		parsedURL, err := url.Parse(match)
		//		if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		//			resp.Data = append(resp.Data, parsedURL.String())
		//		}
		//	}
		//}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childResp := f(c)
		resp.Data = append(resp.Data, childResp.Data...)
	}
	return &resp
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	defer recoverFromPanic()
	errMsg := ResponseData{
		Code: -1,
		Msg:  "ok",
		Data: []string{},
	}

	name = strings.ReplaceAll(name, " ", "")
	_, err := url.ParseRequestURI(name)
	if err != nil {
		errMsg.Msg = "url规范错误"
		jsonResp, _ := json.Marshal(errMsg)
		return fmt.Sprintf("%s", jsonResp)
	}
	if name == "" {
		errMsg.Msg = "域名不能为空"
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)
	}
	// 定义一个基本的 URL 正则表达式
	// 这个正则表达式用于检查 URL 的常见格式，可能不覆盖所有的情况
	pattern := `^(http|https)://[a-zA-Z0-9\.-]+\.[a-zA-Z]{2,}(?:[/?#][^\s]*)?$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(name) {
		errMsg.Msg = "url不符合规范"
		jsonResp, _ := json.Marshal(errMsg)
		return fmt.Sprintf("%s", jsonResp)
	}
	name = strings.ReplaceAll(name, " ", "")
	_, err = url.ParseRequestURI(name)
	if err != nil {
		errMsg.Msg = "url规范错误"
		jsonResp, _ := json.Marshal(errMsg)
		return fmt.Sprintf("%s", jsonResp)
	}
	if name == "" {
		errMsg.Msg = "域名不能为空"
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)
	}
	rs := strings.Contains(name, "http")
	if !rs {
		errMsg.Msg = "输入必须包含请求协议：http/https"
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)
	}
	fmt.Println("输入域名：" + name)
	// 创建一个带有超时设置的 HTTP 客户端
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间为 5 秒
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", name, nil)
	if err != nil {
		if err != nil {
			errMsg.Code = -1
			errMsg.Msg = "创建请求失败"
			jsonResp, _ := json.Marshal(errMsg)
			return fmt.Sprintf("%s", jsonResp)
		}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Connection", "keep-alive")
	// 设置模拟浏览器的请求头

	resp, err := client.Do(req)
	if err != nil {

		// 检查是否是超时错误
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			errMsg.Code = 0
			errMsg.Msg = "请求超时"
			errMsg.Status = "请求超时"
			errMsg.ContentLength = "0"
			jsonResp, _ := json.Marshal(errMsg)

			return fmt.Sprintf("%s", jsonResp)
		}
		errMsg.Code = 0
		errMsg.Msg = "请求错误"
		errMsg.Status = "请求错误"
		errMsg.ContentLength = "0"
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)
	}
	defer resp.Body.Close()

	fmt.Println("返回数据了")
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg.Msg = fmt.Sprintf("读取响应体错误：%v", err)
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)
	}

	// 解析HTML内容
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		errMsg.Msg = fmt.Sprintf("解析 HTML 错误：%v", err)
		jsonResp, _ := json.Marshal(errMsg)

		return fmt.Sprintf("%s", jsonResp)

	}

	result := f(doc)

	// 提取二级域名
	newData := make([]string, 0)

	//fmt.Println("result.Data")
	//fmt.Println(result.Data)
	//fmt.Println("result.Code")
	//fmt.Println(result.Code)
	//if len(result.Data) == 0 {
	if len(result.Data) == 0 || result.Code == 0 {
		// 创建一个 Chrome 浏览器上下文
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		// 设置超时
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		// 存储 href 属性值的切片
		var hrefs []string
		var res string

		// 使用 Chrome 浏览器获取网页中的所有 href 属性值
		err := chromedp.Run(ctx,
			chromedp.Navigate(name),
			chromedp.OuterHTML("html", &res),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &hrefs),
		)
		if err != nil {
			fmt.Println("获取网页内容失败: %v", err)
		}
		// 计算内容长度（千字节），保留两位小数
		contentLengthKB := float64(len(res)) / 1024
		formattedContentLength := fmt.Sprintf("%.2f", contentLengthKB)
		//fmt.Println("result.RawCode-chromedp")
		//fmt.Println(res)
		// 打印所有 href 属性值
		//fmt.Println("所有 href 属性值:")
		result.RawCode = res
		result.ContentLength = formattedContentLength
		result.DomainName = name
		if len(hrefs) > 0 {
			result.Status = "200 OK"
		}

		for _, href := range hrefs {
			result.Data = append(result.Data, href)
		}

		//
		//// 打印网页内容的前1000个字符
		//fmt.Println("网页内容前1000个字符:")
		//fmt.Println(res[:1000])
	}
	for _, urlStr := range result.Data {
		//fmt.Println("url.Parse(urlStr)")
		//fmt.Println(url.Parse(urlStr))
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}

		secondLevelDomain, err := getSecondLevelDomain(u.Hostname())

		if err != nil {
			fmt.Println(err)
		}

		// 删除指定的子字符串
		modifiedString := strings.Replace(u.Hostname(), secondLevelDomain, "", -1)

		// 查找最后一个点的位置
		lastDotIndex := strings.LastIndex(modifiedString, ".")

		if lastDotIndex != -1 {

			// 如果没有找到点，则返回空字符串
			// 计算字符串中点字符（'.'）的数量
			count := strings.Count(modifiedString, ".")
			finalURL := ""
			if count >= 2 {
				// 查找倒数第二个点的位置
				secondLastDotIndex := strings.LastIndex(modifiedString[:lastDotIndex], ".")
				if secondLastDotIndex != -1 {

					// 从倒数第二个点之后的位置提取子字符串
					finalURL = modifiedString[secondLastDotIndex+1:] + secondLevelDomain

				}
			} else {
				finalURL = modifiedString + secondLevelDomain
			}

			// 定义要删除的字符集合
			toRemove := []string{"\"", "'", ";"}

			// 遍历每一个要删除的字符，并进行替换
			for _, char := range toRemove {
				finalURL = strings.ReplaceAll(finalURL, char, "")
			}
			if len(finalURL) < 20 {
				newData = append(newData, finalURL)

			}

		}

	}

	//处理主域名问题------start

	u, err := url.Parse(name)
	if err != nil {
		fmt.Println(err)
	}

	secondLevelDomain, err := getSecondLevelDomain(u.Hostname())

	if err != nil {
		fmt.Println(err)
	}

	// 删除指定的子字符串
	modifiedString := strings.Replace(u.Hostname(), secondLevelDomain, "", -1)

	// 查找最后一个点的位置
	lastDotIndex := strings.LastIndex(modifiedString, ".")
	finalURL := ""

	if lastDotIndex != -1 {

		// 如果没有找到点，则返回空字符串
		// 计算字符串中点字符（'.'）的数量
		count := strings.Count(modifiedString, ".")
		if count >= 2 {
			// 查找倒数第二个点的位置
			secondLastDotIndex := strings.LastIndex(modifiedString[:lastDotIndex], ".")
			if secondLastDotIndex != -1 {

				// 从倒数第二个点之后的位置提取子字符串
				finalURL = modifiedString[secondLastDotIndex+1:] + secondLevelDomain

			}
		} else {
			finalURL = modifiedString + secondLevelDomain
		}

		// 定义要删除的字符集合
		toRemove := []string{"\"", "'", ";"}

		// 遍历每一个要删除的字符，并进行替换
		for _, char := range toRemove {
			finalURL = strings.ReplaceAll(finalURL, char, "")
		}

		fmt.Println("finalURL")
		fmt.Println(finalURL)
		if len(finalURL) < 20 {
			newData = append(newData, finalURL)

		}

	}

	//处理主域名---end
	result.Data = newData
	result.Data = append([]string{finalURL}, result.Data...)

	// 去重操作
	uniqueData := make([]string, 0)
	seen := make(map[string]struct{})
	for _, item := range result.Data {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			uniqueData = append(uniqueData, item)
		}
	}

	// 要添加的元素
	//newElement := "first"

	// 向 Data 切片的第一个位置添加新元素
	//response.Data = append([]string{newElement}, response.Data...)
	result.Data = uniqueData
	if len(result.Data) == 0 {
		result.Code = 0
		result.Msg = "没有查到数据"
		// errMsg.Data 是空的
	}
	if result.Status == "" {
		result.Status = resp.Status
	}

	if result.ContentLength == "" {
		// 计算内容长度（千字节）
		contentLengthKB := float64(len(body)) / 1024

		// 格式化内容长度为保留两位小数的字符串
		formattedContentLength := fmt.Sprintf("%.2f", contentLengthKB)
		result.ContentLength = formattedContentLength
	}

	if result.RawCode == "" {
		result.RawCode = string(body)

	}

	jsonResp, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("JSON 序列化错误：%v", err)
	}
	//fmt.Sprintf("%s", jsonResp)
	return fmt.Sprintf("%s", jsonResp)
}

func getSecondLevelDomain(host string) (string, error) {
	domain, err := publicsuffix.PublicSuffix(host)
	if !err {
		return "", nil
	}
	return domain, nil
}

func recoverFromPanic() {
	if r := recover(); r != nil {

		fmt.Sprintf("捕获到 panic: %v\n", r)
		// 在这里可以记录错误日志，进行清理工作等
	}
}
