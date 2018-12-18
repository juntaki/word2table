package main

import (
	"fmt"

	"github.com/ikawaha/kagome/tokenizer"
	"golang.org/x/text/unicode/norm"
)

func main() {
	input := "ＡＢＣ  -1ﾍﾟﾝｷﾞﾝ"
	input = string(norm.NFKC.Bytes([]byte(input)))

	dic, err := tokenizer.NewUserDic("./dic.csv")
	if err != nil {
		panic(err)
	}
	t := tokenizer.New()
	t.SetUserDic(dic)

	output := []string{}
	tokens := t.Tokenize(input)
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		if token.Features()[1] == "空白" {
			continue
		}
		output = append(output, token.Surface)
	}

	fmt.Println(output)
}
