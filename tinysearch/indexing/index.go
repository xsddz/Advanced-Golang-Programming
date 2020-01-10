package indexing

import (
	"encoding/json"
	"fmt"
	"log"
	"tinysearch/document"
	"tinysearch/tokenization"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
)

// IndexDocument 在redis中建立指定文档的倒排索引
func IndexDocument(rconn redis.Conn, docInfo *document.Document) {
	fmt.Println("rindexing document:", docInfo.Title)

	// 1. 获取文档内容Tokenize后的单词集合
	words := tokenization.Tokenize(docInfo.Content)
	fmt.Println("new tokens:", words)
	if len(words) < 1 {
		return
	}

	// 2. 构建重新Tokenize后被删除的单词
	// 3. 构建重新Tokenize后需新增的单词
	var delwords, addwords []string
	oldwordsJSON, err := redis.String(rconn.Do("GET", util.MakeDocWordsIndexRedisKey(docInfo.Title)))
	fmt.Println("old tokens:", oldwordsJSON)
	if err != nil {
		addwords = words
	} else {
		var oldwords []string
		err := json.Unmarshal([]byte(oldwordsJSON), &oldwords)
		if err != nil {
			log.Fatal("JSON deocde old words error", err)
		}
		delwords = diffWords(oldwords, words)
		addwords = diffWords(words, oldwords)
	}
	fmt.Println("tokens need to del:", delwords)
	fmt.Println("tokens need to add:", addwords)

	// 4. 删除删除单词的倒排索引
	for _, token := range delwords {
		rconn.Send("SREM", util.MakeDocWordsRIndexRedisKey(token), docInfo.Title)
	}
	// 5. 建立新增单词的倒排索引
	for _, token := range addwords {
		rconn.Send("sadd", util.MakeDocWordsRIndexRedisKey(token), docInfo.Title)
	}
	// 6. 建立 文档-单词集合 索引
	wordsJSON, err := json.Marshal(words)
	if err != nil {
		log.Fatal("JSON encode words error", err)
	}
	rconn.Send("SET", util.MakeDocWordsIndexRedisKey(docInfo.Title), string(wordsJSON))
	rconn.Flush()
	v, err := rconn.Receive()
	if err != nil {
		log.Fatal("Redis pipeline error", err)
	}

	fmt.Println(v)
}

func diffWords(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return
}
