package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"hw3/model"
)

func FastSearch(out io.Writer) {
	const filePath string = "./data/users.txt"
	var foundUsers strings.Builder

	reader, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	seenBrowsers := make(map[string]struct{})
	uniqueBrowsers := 0
	lines := strings.Split(string(reader), "\n")
	for i, line := range lines {
		user := model.User{}
		user.UnmarshalJSON([]byte(line))
		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
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
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		fmt.Fprintf(&foundUsers, "[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func main() {
	FastSearch(os.Stdout)
}
