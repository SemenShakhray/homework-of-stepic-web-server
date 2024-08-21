package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"hw3/model"
)

func FastSearch(out io.Writer) {
	const filePath string = "./data/users.txt"
	var foundUsers strings.Builder
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	seenBrowsers := make(map[string]struct{})
	uniqueBrowsers := 0
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		br := model.Brows{}
		br.UnmarshalJSON([]byte(line))
		isAndroid := false
		isMSIE := false

		for _, browserRaw := range br.Browsers {
			if ok := strings.Contains(browserRaw, "Android"); ok {
				isAndroid = true
				if _, ok := seenBrowsers[browserRaw]; !ok {
					seenBrowsers[browserRaw] = struct{}{}
					uniqueBrowsers++
				}
			}
			if ok := strings.Contains(browserRaw, "MSIE"); ok {
				isMSIE = true
				if _, ok := seenBrowsers[browserRaw]; !ok {
					seenBrowsers[browserRaw] = struct{}{}
					uniqueBrowsers++
				}
			}
		}
		if !(isAndroid && isMSIE) {
			continue
		}
		user := &model.User{}
		user.UnmarshalJSON([]byte(line))
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		fmt.Fprintf(&foundUsers, "[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
