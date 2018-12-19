package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
	lsd "github.com/mattn/go-lsd"
	"golang.org/x/text/unicode/norm"
)

func userdic() tokenizer.UserDic {
	fp, err := os.Open("udic.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	r := tokenizer.UserDicRecords{}
	for scanner.Scan() {
		token := scanner.Text()
		r = append(r, tokenizer.UserDicRecord{
			Text:   token,
			Tokens: []string{token},
			Yomi:   []string{""},
			Pos:    "名詞",
		})
	}
	udic, err := r.NewUserDic()
	if err != nil {
		panic(err)
	}

	return udic
}

func main() {
	dic, err := tokenizer.NewUserDic("./dic.csv")
	if err != nil {
		panic(err)
	}
	t := tokenizer.New()
	t.SetUserDic(dic)
	t.SetUserDic(userdic())

	fp, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	headerMap := make(map[string]int)

	re := strings.NewReplacer(
		"\"", "",
		"\n", "",
		"\r", "",
	)

	for scanner.Scan() {
		input := scanner.Text()
		input = string(norm.NFKC.Bytes([]byte(input)))
		input = re.Replace(input)

		output := []string{}
		tokens := t.Tokenize(input)
		for _, token := range tokens {
			if token.Class == tokenizer.DUMMY {
				continue
			}
			if token.Features()[1] == "空白" || token.Features()[1] == "接尾" {
				continue
			}
			if token.Features()[0] == "助詞" ||
				token.Features()[0] == "助動詞" ||
				token.Features()[0] == "記号" {
				continue
			}
			output = append(output, token.Surface)

			headerMap[token.Surface] += 1
		}
	}

	header := []string{}
	for k, v := range headerMap {
		if v > 5 {
			header = append(header, k)
		}
	}

	sort.Slice(header, func(i int, j int) bool {
		return header[i] > header[j]
	})
	fmt.Printf("\"input\",")
	fmt.Println(strings.Join(header, ","))

	fp.Seek(0, 0)
	scanner = bufio.NewScanner(fp)
	for scanner.Scan() {
		input := scanner.Text()
		input = string(norm.NFKC.Bytes([]byte(input)))
		input = re.Replace(input)

		output := make([]string, len(header))
		tokens := t.Tokenize(input)
		for _, token := range tokens {
			if token.Class == tokenizer.DUMMY {
				continue
			}
			if token.Features()[1] == "空白" || token.Features()[1] == "接尾" {
				continue
			}
			if token.Features()[0] == "助詞" ||
				token.Features()[0] == "助動詞" ||
				token.Features()[0] == "記号" {
				continue
			}

			for i, h := range header {
				if lsd.StringDistance(h, token.Surface) == 1 && output[i] == "" && len(h) > 4 {
					output[i] = "1"
				}
				if h == token.Surface {
					output[i] = "0"
				}
			}
		}
		fmt.Printf("\"%s\",", input)
		fmt.Println(strings.Join(output, ","))
	}
}
