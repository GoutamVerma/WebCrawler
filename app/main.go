package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	graph "github.com/GoutamVerma/WebCrawler/crawler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

var (
	config        = &tls.Config{InsecureSkipVerify: true}
	transport     = &http.Transport{TLSClientConfig: config}
	netClient     = &http.Client{Transport: transport}
	graphMap      = graph.NewGraph()
	lastCrawlTime time.Time
)

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/crawl", crawlHandler)

	e.Logger.Fatal(e.Start(":1234"))
}

func crawlHandler(c echo.Context) error {
	hasCrawled := make(map[string]bool)
	urlQueue := make(chan string)
	crawledURLs := []string{}
	stopCrawling := false
	baseUrl := c.QueryParam("url")
	if baseUrl == "" {
		return c.String(http.StatusBadRequest, "URL is missing")
	}

	deepStr := c.QueryParam("deep")
	maxDeep := 0
	if deepStr != "" {
		var err error
		maxDeep, err = strconv.Atoi(deepStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid value for deep")
		}
	}
	crawledURLs = []string{}

	go func() {
		urlQueue <- baseUrl
	}()

	timeout := viper.GetInt("timeout")
	timer := time.NewTimer(time.Duration(timeout) * time.Second)

	for {
		select {
		case href := <-urlQueue:
			if !hasCrawled[href] {
				crawlLink(href, baseUrl, &urlQueue, &crawledURLs, &hasCrawled)
				lastCrawlTime = time.Now()
			}

			if maxDeep > 0 && len(crawledURLs) >= maxDeep {
				stopCrawling = true
			}

		case <-timer.C:
			if time.Since(lastCrawlTime) >= 10*time.Second {
				stopCrawling = true
			}
		}

		if stopCrawling {
			root := buildSiteMap(crawledURLs)

			hostname, err := extractHostname(baseUrl)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to extract hostname")
			}

			sitemap := fmt.Sprintf("%s\n%s", hostname, printSiteMap(root, 1))
			return c.String(http.StatusOK, sitemap)
		}

		timer.Reset(10 * time.Second)
	}
}

func crawlLink(baseHref string, baseUrl string, urlQueue *chan string, crawledURLs *[]string, hasCrawled *map[string]bool) {
	if !isSameDomain(baseHref, baseUrl) {
		return
	}
	if strings.Contains(baseHref, "?") {
		return
	}
	graphMap.AddVertex(baseHref)
	(*hasCrawled)[baseHref] = true
	*crawledURLs = append(*crawledURLs, baseHref)
	resp, err := netClient.Get(baseHref)
	checkErr(err)
	defer resp.Body.Close()

	links, err := graph.All(resp.Body)
	checkErr(err)

	for _, l := range links {
		if l.Href == "" {
			continue
		}
		fixedUrl := toFixedUrl(l.Href, baseHref)
		if baseHref != fixedUrl {
			graphMap.AddEdge(baseHref, fixedUrl)
		}
		go func(url string) {
			*urlQueue <- url
		}(fixedUrl)
	}
}

func isSameDomain(href, baseUrl string) bool {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}

	parsedURL, err := url.Parse(href)
	if err != nil {
		return false
	}

	return parsedURL.Hostname() == base.Hostname()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func toFixedUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil || uri.Scheme == "mailto" || uri.Scheme == "tel" {
		return base
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

type siteMapNode struct {
	Path     string
	Children map[string]*siteMapNode
}

func newSiteMapNode(path string) *siteMapNode {
	return &siteMapNode{
		Path:     path,
		Children: make(map[string]*siteMapNode),
	}
}

func buildSiteMap(urls []string) *siteMapNode {
	root := newSiteMapNode("")
	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}
		pathComponents := strings.Split(u.Path, "/")
		node := root
		for _, component := range pathComponents {
			if component == "" {
				continue
			}
			if _, ok := node.Children[component]; !ok {
				node.Children[component] = newSiteMapNode(component)
			}
			node = node.Children[component]
		}
	}
	return root
}

func printSiteMap(node *siteMapNode, indent int) string {
	result := ""
	keys := make([]string, 0, len(node.Children))
	for key := range node.Children {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		child := node.Children[key]
		result += fmt.Sprintf("%s-%s\n", strings.Repeat(" ", indent*2), child.Path)
		result += printSiteMap(child, indent+1)
	}
	return result
}

func extractHostname(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	hostname := parsedURL.Hostname()
	return hostname, nil
}
