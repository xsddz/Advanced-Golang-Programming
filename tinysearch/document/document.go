package document

import (
	"io/ioutil"
	"log"
	"os"
	"tinysearch/util"

	"github.com/gomodule/redigo/redis"
)

// Document 文档结构
type Document struct {
	Title   string
	Content string
}

// ReadDoc 读取指定标题的文档
func ReadDoc(file string) *Document {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Open file error:", err)
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Read file content error:", err)
	}

	return &Document{
		file,
		string(content),
	}
}

// VoteDoc 对文档进行投票
func (doc *Document) VoteDoc(rconn redis.Conn, num int) {
	rconn.Do("ZINCRBY", util.MakeDocVoteRedisKey(), num, doc.Title)
}

// UpdateDocUpdateTime 更新文档内容更新时间
func (doc *Document) UpdateDocUpdateTime(rconn redis.Conn, timestamp int64) {
	rconn.Do("ZADD", util.MakeDocUpdateRedisKey(), timestamp, doc.Title)
}
