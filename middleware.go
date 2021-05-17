package main

import (
	// "encoding/json"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	// "io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// bodyCacheWriter is used to cache responses in gin.
type bodyCacheWriter struct {
	gin.ResponseWriter
	cache      *cache.Cache
	requestURI string
}

// Write a JSON response to gin and cache the response.
func (w bodyCacheWriter) Write(b []byte) (int, error) {
	// Write the response to the cache only if a success code
	status := w.Status()
	if 200 <= status && status <= 299 {
		w.cache.Set(w.requestURI, b, cache.DefaultExpiration)
	}

	// Then write the response to gin
	return w.ResponseWriter.Write(b)
}

func CacheInit() (*cache.Cache, error) {
	// var rs map[string]cache.Item

	// if _, err := os.Stat(config.App.Cache.Path); os.IsNotExist(err) {
	// 	return cache.New(time.Duration(config.App.Cache.Expiration)*time.Hour, time.Duration(config.App.Cache.Cleanup)*time.Hour), nil
	// }

	// data, err := os.ReadFile(config.App.Cache.Path)
	// if err != nil {
	// 	return nil, err
	// }

	// err = json.Unmarshal(data, &rs)
	// if err != nil {
	// 	return nil, err
	// }

	// return cache.NewFrom(time.Duration(config.App.Cache.Expiration)*time.Hour, time.Duration(config.App.Cache.Cleanup)*time.Hour, rs), nil
	cache := cache.New(time.Duration(config.App.Cache.Expiration)*time.Hour, time.Duration(config.App.Cache.Cleanup)*time.Hour)
	if _, err := os.Stat(config.App.Cache.Path); err == nil {
		err := cache.LoadFile(config.App.Cache.Path)
		if err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
	} else {
		return nil, err
	}

	return cache, nil
}

func CacheSave(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// items := cache.Items()

		// data, err := json.Marshal(items)
		// if err != nil {
		// 	util.Error(c, err)
		// 	return
		// }

		// err = ioutil.WriteFile("cache.json", data, 0644)
		// if err != nil {
		// 	util.Error(c, err)
		// 	return
		// }
		cache.SaveFile(config.App.Cache.Path)
		util.Success(c, nil)
	}
}

// CacheCheck sees if there are any cached responses and returns
// the cached response if one is available.
func CacheCheck(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the ignoreCache parameter
		ignoreCache := strings.ToLower(c.Query("ignoreCache")) == "true"

		// See if we have a cached response
		response, exists := cache.Get(c.Request.RequestURI)
		if !ignoreCache && exists {
			// If so, use it

			if string(response.([]byte)) == `{"error":null,"response":{"games":null},"status":"OK"}` {
				// If not, pass our cache writer to the next middleware
				bcw := &bodyCacheWriter{cache: cache, requestURI: c.Request.RequestURI, ResponseWriter: c.Writer}
				c.Writer = bcw
				c.Next()
			} else {
				log.Println(c.Request.RequestURI, "is cached")
				c.Data(200, "application/json", response.([]byte))
				c.Abort()
			}
		} else {
			// If not, pass our cache writer to the next middleware
			bcw := &bodyCacheWriter{cache: cache, requestURI: c.Request.RequestURI, ResponseWriter: c.Writer}
			c.Writer = bcw
			c.Next()
		}
	}
}
