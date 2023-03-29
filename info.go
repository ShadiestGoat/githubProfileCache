package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

type cacheT struct {
	*sync.RWMutex
	m map[langType]map[string]float64
}

var c = &cacheT{
	RWMutex: &sync.RWMutex{},
	m:       map[langType]map[string]float64{},
}

func updateCache() {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer " + os.Getenv("AUTHENTICATION"))

	page := 1

	allStuff := map[string]int{}
	lock := &sync.Mutex{}

	for {
		uri := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d", page)
		req, _ := http.NewRequest("GET", uri, nil)
		req.Header = headers
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			if resp != nil {
				fmt.Println(resp.StatusCode)
			}
			continue
		}

		p := []*respRepo{}
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil || resp.StatusCode != 200 {
			fmt.Println(err)
			if resp != nil {
				fmt.Println(resp.StatusCode)
			}
			continue
		}

		wg := &sync.WaitGroup{}

		for _, r := range p {
			wg.Add(1)
			go func(r respRepo) {
				defer wg.Done()

				langs := repoLangFetch(r, headers)

				for langs == nil {
					time.Sleep(30 * time.Minute)
					langs = repoLangFetch(r, headers)
				}

				for l, b := range langs {
					if l == "" {
						delete(langs, l)
					}
					lock.Lock()
					allStuff[l] += b
					lock.Unlock()
				}
			}(*r)
		}

		wg.Wait()
		if len(p) != 100 {
			break
		}

		page++
	}

	perc := map[langType]*typeInfo{}

	for lang, num := range allStuff {
		allTypes, ok := languageTypeInfo[lang]
		if !ok {
			continue
		}
		for _, t := range allTypes {
			v, ok := perc[t]
			if !ok {
				perc[t] = &typeInfo{
					m: map[string]int{},
				}
				v = perc[t]
			}
			v.total += num
			v.m[lang] += num
		}
	}

	vals := map[langType]map[string]float64{}

	for t, v := range perc {
		vals[t] = map[string]float64{}
		norm := languageNormalization[t]
		if norm == nil {
			norm = map[string]float64{}
		}

		for lang, n := range v.m {
			mult := norm[lang]
			if mult == 0 {
				mult = 1
			}
			n = int(math.Round(float64(n) * mult))

			size := 100 * 100.0
			vals[t][lang] = math.Round(size*float64(n)/float64(v.total)) / size
		}
	}

	c.Lock()
	defer c.Unlock()
	c.m = vals
}

func StartCacheLoop() {
	for {
		updateCache()
		time.Sleep(2 * time.Hour)
	}
}
