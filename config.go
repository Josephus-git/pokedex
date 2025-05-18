package main

type config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Result   []data `json:"results"`
}
type data struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
