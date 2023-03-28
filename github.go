package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

// Fetch repo langs, returns nil if theres an error
func repoLangFetch(r respRepo, headers http.Header) map[string]int {
	langs := map[string]int{}
	if r.IsFork {
		return map[string]int{}
	}

	uri := fmt.Sprintf("https://api.github.com/repos/%s/languages", r.FullName)
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header = headers
	resp, _ := http.DefaultClient.Do(req)

	err := json.NewDecoder(resp.Body).Decode(&langs)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if langs["HTML"] != 0 && (langs["TypeScript"]+langs["JavaScript"] != 0) && langs["Svelte"] == 0 {
		tsTot := float64(langs["TypeScript"])
		jsTot := float64(langs["JavaScript"])

		langs["JavaScript"] = int(math.Round(jsTot * 0.2))
		langs["TypeScript"] = int(math.Round(tsTot * 0.2))
		langs["React (JS)"] = int(math.Round(jsTot * 0.3))
		langs["React (TS)"] = int(math.Round(tsTot * 0.3))
	}

	return langs
}
