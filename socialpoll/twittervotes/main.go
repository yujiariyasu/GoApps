package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDBにダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}
func closedb() {
	db.Close()
	log.Println("データベース接続が閉じられました")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	// 投票を待ち、それをNSQにパブリッシュする。votesチャネルが閉じられたらfor文を抜けてstopChanにシグナルを送る
	go func() {
		// 以下のfor文では、チャネルであるvotesから定期的に値を読みだしている。
		// votesチャネルを継続的にチェックしているが、チャネルを閉じることでループを終了させることができる。
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // 投票内容をパブリッシュ
		}
		log.Println("Publisher: 停止中です")
		pub.Stop()
		log.Println("Publisher: 停止しました")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {
	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	// signalChanにシグナルが送られたら、stopフラグをtrueにしてstopChanにstruct{}{}を送信し、closeConnを呼ぶ
	// stopChanが受信すると、startTwitterStream内の無名関数を出る。その際にstoppedChanにシグナル送る
	// 上記のような終了の仕組みを以下のgoroutineでセットしておく
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("停止します....")
		stopChan <- struct{}{}
		closeConn()
	}()
	// プログラムを終了しようとした時にsignalChanにシグナルを送るように設定
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// MongoDBに接続
	if err := dialdb(); err != nil {
		log.Fatalln("MongoDBへのダイヤルに失敗しました:", err)
	}
	defer closedb()
	// 処理を開始
	votes := make(chan string)                  // 投票結果のためのチャネル
	publisherStoppedChan := publishVotes(votes) // パブリッシャーをスタートさせている
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	// 単に一文ごとにcloseConnを呼び出して接続を切断する
	// こうすることでreadFromTwitterが再び呼ばれ、選択肢を最新のものに保てる
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			// 2つのgoroutineが同じ変数にアクセスする時は、競合を避けるためにstoplockを利用
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStoppedChan // twitterStoppedChanからデータを読み込むまでは実行をブロック
	close(votes)
	<-publisherStoppedChan
}
