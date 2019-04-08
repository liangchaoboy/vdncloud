package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type top struct {
	url   string
	value int
}

type toplist []top

func (p toplist) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p toplist) Len() int           { return len(p) }
func (p toplist) Less(i, j int) bool { return p[i].value < p[j].value }

func sortMapByValue(m map[string]int) toplist {

	p := make(toplist, len(m))
	i := 0
	for k, v := range m {
		p[i] = top{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

func main() {

	ak := "账号 ak"
	sk := "账号 sk"

	// url 、访问次数
	var top map[string]int
	top = make(map[string]int)

	// url、 访问流量
	var topflow map[string]int
	topflow = make(map[string]int)

	urls := getUrl("2019-04-07", ak, sk)

	for _, j := range urls {

		for key, value := range getCount(j) {
			_, ok := top[key]

			if ok {
				top[key] = top[key] + value
			}

			top[key] = value
		}
	}

	for _, f := range urls {

		for key, value := range getFlow(f) {
			_, ok := topflow[key]

			if ok {
				topflow[key] = topflow[key] + value
			}

			topflow[key] = value
		}
	}

	for _, sortCount := range sortMapByValue(top) {
		fmt.Println(sortCount)
	}

	fmt.Println("\n")

	for _, sortflow := range sortMapByValue(topflow) {
		fmt.Println(sortflow)
	}

}

// 获取下载链接 date:2019-03-29
func getUrl(date string, ak string, sk string) []string {

	s := make([]string, 0)

	mac := qbox.NewMac(ak, sk)

	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)

	domain := "http://pp32edcpl.bkt.clouddn.com"
	limit := 1000
	bucketName := "vdncloud"
	timeStr := date
	prefix := "csv/cdn3rd10sjylive.voole.com/"
	preString := fmt.Sprintf("%s%s", prefix, timeStr)
	delimiter := ""
	//初始列举marker为空
	marker := ""

	entries, _, _, _, err := bucketManager.ListFiles(bucketName, preString, delimiter, marker, limit)

	if err != nil {

		fmt.Println("list error,", err)

	}

	//print entries
	for _, entry := range entries {

		str := entry.Key

		if strings.Contains(str, "part") {

			deadline := time.Now().Add(time.Second * 3600).Unix()

			getDownLoadUrl := storage.MakePrivateURL(mac, domain, entry.Key, deadline)

			s = append(s, getDownLoadUrl)

		}
	}
	return s
}

// 获取 url 和访问次数
func getCount(str string) map[string]int {

	var top map[string]int

	top = make(map[string]int)

	resp, err := http.Get(str)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, errs := ioutil.ReadAll(resp.Body)

	if errs != nil {
		// handle error
	}

	t := string(body)

	var i []string = strings.Split(t, "\n")

	for index := 0; index < len(i); index++ {
		if strings.Contains(i[index], "http") {

			k := strings.Split(i[index], ",")
			intk, _ := strconv.Atoi(k[2])
			top[k[0]] = intk
		}
	}

	return top

}

// 获取 url 和访问流量
func getFlow(str string) map[string]int {

	var top map[string]int

	top = make(map[string]int)

	resp, err := http.Get(str)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, errs := ioutil.ReadAll(resp.Body)

	if errs != nil {
		// handle error
	}

	t := string(body)

	var i []string = strings.Split(t, "\n")

	for index := 0; index < len(i); index++ {
		if strings.Contains(i[index], "http") {
			k := strings.Split(i[index], ",")
			intk, _ := strconv.Atoi(k[1])
			top[k[0]] = intk
		}
	}

	return top

}
