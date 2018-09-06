package main

import(
	"fmt"
	"flag"
	"os"
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

func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}

	log.Println("データベースに接続します...")
	db, err != mgo.Dial("localhost")
	if err != nil {
		fatal(err)
		return
	}

	defer func() {
		log.Println("データベース接続を閉じます...")
		db.Close()
	}
}