package operation

import (
	"log"
	"regexp"
	"strings"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
)

var queryParseReg = regexp.MustCompile(`[+-]?[\S]{2,}`)

// SearchAndZSort 搜索并利用zset进行排序分页
func SearchAndZSort(rconn redis.Conn, query string, id string, ttl int, update int, vote int, start int, num int, desc bool) (int, []string, string) {
	stat, _ := redis.Bool(rconn.Do("EXPIRE", id, ttl))
	if stat == false {
		id = ""
	}

	if len(id) == 0 {
		id = Search(rconn, query, ttl)

		keyWeight := map[string]int{
			id: 0,
			util.MakeDocUpdateRedisKey(): update,
			util.MakeDocVoteRedisKey():   vote,
		}
		id = ZSetIntersect(rconn, keyWeight, ttl, true)
	}

	rconn.Send("MULTI")
	rconn.Send("ZCARD", id)
	if desc {
		rconn.Send("ZREVRANGE", id, start, start+num-1)
	} else {
		rconn.Send("ZRANGE", id, start, start+num-1)
	}
	res, err := redis.Values(rconn.Do("EXEC"))
	if err != nil {
		log.Fatal("Redis pipeline error", err)
	}

	var count int
	var items []string
	for _, item := range res {
		switch item.(type) {
		case int, int64:
			count, _ = redis.Int(item, nil)
		case []interface{}:
			items, _ = redis.Strings(item, nil)
		}
	}
	return count, items, id
}

/*
Search 搜索入口

搜索词规则如下：
+ expre: word [+word1] [-word2] 执行word和word2的并集，然后和word2执行差集
+ expre1 [expre2] 执行expre1和expre2的交集
*/
func Search(rconn redis.Conn, query string, ttl int) string {
	wanted, unwanted := reParse(query)

	var toIntersect []string
	for _, words := range wanted {
		if len(words) > 1 {
			toIntersect = append(toIntersect, SetUnion(rconn, util.MakeDocWordsRIndexRedisKeys(words), ttl, true))
		} else {
			toIntersect = append(toIntersect, util.MakeDocWordsRIndexRedisKey(words[0]))
		}
	}

	var intersectResult string
	if len(toIntersect) > 1 {
		intersectResult = SetIntersect(rconn, toIntersect, ttl, true)
	} else {
		intersectResult = toIntersect[0]
	}

	if len(unwanted) > 0 {
		toDifference := []string{intersectResult}
		toDifference = append(toDifference, util.MakeDocWordsRIndexRedisKeys(unwanted)...)
		return SetDifference(rconn, toDifference, ttl, true)
	}

	return intersectResult
}

// reParse 解析搜索词
func reParse(query string) ([][]string, []string) {
	var unwanted, current []string
	var wanted [][]string

	words := queryParseReg.FindAllString(query, -1)
	for _, word := range words {
		prefix := word[:1]
		if strings.ContainsAny(prefix, "+-") {
			word = word[1:]
		} else {
			prefix = ""
		}

		if prefix == "-" {
			unwanted = append(unwanted, word)
			continue
		}

		if (len(current) > 0) && (len(prefix) == 0) {
			wanted = append(wanted, current)
			current = []string{}
		}
		current = append(current, word)
	}

	if len(current) > 0 {
		wanted = append(wanted, current)
	}

	return wanted, unwanted
}
