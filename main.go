package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	lru "github.com/hashicorp/golang-lru"
)

// CacheEntry represents a key-value entry in the cache
type CacheEntry struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}

// LRUCache represents the LRU cache
type LRUCache struct {
	cache *lru.Cache
}

// NewLRUCache initializes a new LRUCache instance
func NewLRUCache(maxEntries int) *LRUCache {
	cache, _ := lru.New(maxEntries)
	return &LRUCache{cache: cache}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (string, error) {
	entry, ok := c.cache.Get(key)
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	cacheEntry := entry.(CacheEntry)
	if time.Now().After(cacheEntry.ExpiresAt) {
		c.cache.Remove(key)
		return "", fmt.Errorf("key expired")
	}

	return cacheEntry.Value, nil
}

// Set value in the cache
func (c *LRUCache) Set(key, value string, ttl time.Duration) {
	expiration := time.Now().Add(ttl)
	c.cache.Add(key, CacheEntry{Value: value, ExpiresAt: expiration})
}

var cache *LRUCache

// Function HandleGetCache handles GET requests to retrieve a value from the cache
func HandleGetCache(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]

	value, err := cache.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

// Function handles POST requests to set a value in the cache
func HandleSetCache(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}

	// Log the raw request body
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Error reading request body: %v\n", err)
		return
	}
	trimmedBody := strings.TrimSpace(string(rawBody)) // Trim whitespace
	fmt.Println("Raw Request Body:", trimmedBody)

	err = json.Unmarshal([]byte(trimmedBody), &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	// Extract the key from the JSON payload
	key := requestData.Key

	// Set the value in the cache
	cache.Set(key, requestData.Value, time.Duration(requestData.TTL)*time.Second)
	fmt.Fprintf(w, "Value set successfully for key: %s\n", key)
}

func main() {
	cache = NewLRUCache(1024)

	router := mux.NewRouter()
	router.HandleFunc("/cache/{key}", HandleGetCache).Methods("GET")
	router.HandleFunc("/cache", HandleSetCache).Methods("POST")

	fmt.Println("LRU Cache server started at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
