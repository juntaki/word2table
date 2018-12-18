package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
	"golang.org/x/text/unicode/norm"
)

func main() {
	dic, err := tokenizer.NewUserDic("./dic.csv")
	if err != nil {
		panic(err)
	}
	t := tokenizer.New()
	t.SetUserDic(dic)

	fp, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	headerMap := make(map[string]struct{})

	for scanner.Scan() {
		input := scanner.Text()
		input = string(norm.NFKC.Bytes([]byte(input)))

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

			headerMap[token.Surface] = struct{}{}
		}
	}

	header := []string{}
	for k := range headerMap {
		header = append(header, k)
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

		output := make([]string, len(header))
		tokens := t.Tokenize(input)
		for _, token := range tokens {
			if token.Class == tokenizer.DUMMY {
				continue
			}
			if token.Features()[1] == "空白" {
				continue
			}

			for i, h := range header {
				if h == token.Surface {
					output[i] = "1"
				}
			}
		}
		fmt.Printf("\"%s\",", input)
		fmt.Println(strings.Join(output, ","))
	}
}
