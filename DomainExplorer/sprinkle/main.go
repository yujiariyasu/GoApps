package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{
	otherWord,
	otherWord,
	otherWord,
	otherWord,
	otherWord + "app",
	otherWord + "site",
	otherWord + "time",
	"get" + otherWord,
	"go" + otherWord,
	"lets " + otherWord,
}

func main() {
	// 現在の時刻から乱数の元(シード)を作成
	rand.Seed(time.Now().UTC().UnixNano())
	// 標準入力のストリームkからデータを読み込むように指定したbufio.Scannerオブジェクトを生成
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		// transformsのうちの1つをtに代入
		t := transforms[rand.Intn(len(transforms))]
		// 出力
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
