package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var fatalErr error

// エラーが発生した際にlog.Fatalやos.Exitを呼び出すと、deferしておいたコードは実行されない。
// エラーが発生した際に以下のfatalメソッドを呼ぶようにするとdefer文も実行される。
// ※エラーが発生したらエラーを返し最後にmainでlog.Fatalを呼び出すような書き方でもOK。こちらの方がdeferを素直に使える。
func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

const updateDuration = 1 * time.Second

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()

	// DB接続
	log.Println("データベースに接続します...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}

	defer func() {
		log.Println("データベース接続を閉じます...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	// マップとロック(sinc.Mutex)はよく使われる組み合わせ
	// 複数のgoroutineが1つのマップにアクセスできる場合に、同時に読み書きを行なってマップが破壊されるのを防ぐ
	var countsLock sync.Mutex
	var counts map[string]int

	log.Println("NSQに接続します...")
	// NSQのvotesトピックを監視するオブジェクトをqにセット
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}

	// この書き方で、votes上でメッセージが受信されるたびに呼び出される
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		fatal(err)
		return
	}

	log.Println("NSQ上での投票を待機します...")
	var updater *time.Timer
	// time.Afterfuncを呼び出すと、引数として指定された関数を一定時間後に自身のgoroutineの中で実行する
	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("新しい投票はありません。データベースの更新をスキップします")
		} else {
			log.Println(counts)
			ok := true
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if _, err := pollData.UpdateAll(sel, up); err != nil {
					log.Println("更新に失敗しました：", err)
					ok = false
					continue
				}
				counts[option] = 0
			}
			if ok {
				log.Println("データベースの更新が完了しました。")
				counts = nil // 得票数をリセット
			}
		}
		// Resetを呼び出すと、同じ手順が再び行われる。つまり更新のためのコードが定期的に繰り返し実行される。
		updater.Reset(updateDuration)
	})
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-termChan:
			updater.Stop()
			q.Stop()
		case <-q.StopChan:
			// 完了
			return
		}
	}
}
