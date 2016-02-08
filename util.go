package kupak

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func fetchURL(url string) ([]byte, error) {
	if strings.HasPrefix(strings.ToLower(url), "http://") ||
		strings.HasPrefix(strings.ToLower(url), "https://") {
		c := &http.Client{}
		resp, err := c.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return ioutil.ReadFile(url)
}

func joinURL(baseUrl string, url string) string {
	return path.Join(path.Dir(baseUrl), url)
}

func getMapChild(keys []string, m map[string]interface{}) (interface{}, error) {
	var innerMap map[string]interface{}
	var v interface{}
	var has, ok bool
	for i := range keys {
		if innerMap == nil {
			innerMap = m
		} else {
			innerMap, ok = v.(map[string]interface{})
			if !ok {
				return nil, errors.New("key not found " + keys[i])
			}
		}
		v, has = innerMap[keys[i]]
		if !has {
			return nil, errors.New("key not found " + keys[i])
		}
	}
	return v, nil
}

func mergeStringMaps(a map[string]string, b map[string]string) map[string]string {
	if a == nil {
		a = make(map[string]string)
	}
	out := make(map[string]string)
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
