package tokenization

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const tokenizeAPITpl = "http://api.pullword.com/get.php?source=%s&param1=1&param2=0&json=1"

type tokenizeAPIResult struct {
	T string `json:"t"`
}

// Tokenize tokenize指定内容
func Tokenize(content string) (tokens []string) {
	temp := map[string]struct{}{}
	ps := strings.Split(content, "\n")
	for _, p := range ps {
		p = strings.Trim(p, " \n#")
		if len(p) < 1 {
			continue
		}

		// tokenize content
		url := fmt.Sprintf(tokenizeAPITpl, p)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Get pullword fail", err)
		}
		defer resp.Body.Close()

		// fmt.Println(url)
		// c, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(c))

		// 将 json 响应解码到结构类型
		var tokenInfoList []tokenizeAPIResult
		err = json.NewDecoder(resp.Body).Decode(&tokenInfoList)
		if err != nil {
			log.Fatal("Deocode json error", err)
		}

		// fmt.Println(tokenInfoList)

		// 收集 token
		for _, tokenInfo := range tokenInfoList {
			if _, ok := temp[tokenInfo.T]; !ok {
				temp[tokenInfo.T] = struct{}{}
				tokens = append(tokens, tokenInfo.T)
			}
		}

		// break
	}
	return
}
