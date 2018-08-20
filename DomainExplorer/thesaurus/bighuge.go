package thesaurus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	// スライス型のsyns定義
	var syns []string
	// apiのresponse取得
	response, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKey + "/" + term + "/json")
	if err != nil {
		return syns, fmt.Errorf("bighuge: %qの類語検索に失敗しました: %v", term, err)
	}
	// synonyms型のdata定義
	var data synonyms
	defer response.Body.Close()
	fmt.Println(response)
	fmt.Println(response.Body)
	// レスポンスの本体をjson.NewDecoderメソッドに渡し、バイト列からsynonyms型へとデコードを行なって結果をdata変数にセットする。
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return syns, err
	}
	fmt.Println(data)
	syns = append(syns, data.Noun.Syn...)
	syns = append(syns, data.Verb.Syn...)
	return syns, nil
}
