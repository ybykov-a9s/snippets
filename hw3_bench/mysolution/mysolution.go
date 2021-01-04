package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

// FastSearch version of Search function is down here
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// fileContents, err := ioutil.ReadAll(file) //!!! MEM alloc 9.74 Mb
	// if err != nil {
	// 	panic(err)
	// }

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	// lines := strings.Split(string(fileContents), "\n") //!!! MEM alloc 2.25 Mb
	lines := bufio.NewScanner(file)
	users := make([]map[string]interface{}, 0)
	// for _, line := range lines {
	for lines.Scan() {
		user := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal([]byte(lines.Text()), &user) // !!! 330ms CPU //!!! MEM alloc 9.18 Mb
		if err != nil {
			panic(err)
		}
		users = append(users, user)

	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			// if ok, err := regexp.MatchString("Android", browser); ok && err == nil { // !!! 250ms CPU //!!! MEM alloc 31.23 Mb
			if ok = strings.Contains(browser, "Android"); ok && err == nil {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}
			//if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil { // !!! 290 ms CPU //!!! MEM alloc 20.75 Mb
			if ok = strings.Contains(browser, "MSIE"); ok && err == nil {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email) //!!! MEM alloc 0.79 Mb
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}
