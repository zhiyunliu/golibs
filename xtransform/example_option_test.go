package xtransform_test

import (
	"fmt"
	"github.com/zhiyunliu/golibs/xtransform"
)

func ExampleTranslate_withOptions() {
	data := map[string]interface{}{
		"name": "World",
		"item": "Go",
		"lang": "Golang",
	}

	// Example 1: Default behavior (backward compatibility) - all modes enabled
	template1 := "Hello @name, welcome to @{item}!"
	result1 := xtransform.Translate(template1, data)
	fmt.Println("Default (all modes):", result1)

	// Example 2: Only @{...} mode enabled
	template2 := "Hello @name, welcome to @{lang}!"
	result2 := xtransform.Translate(template2, data, xtransform.WithAtBraceMode())
	fmt.Println("Only @{...} mode:", result2)

	// Example 3: @... and {...} modes enabled
	template3 := "Hello @name, learn {lang} today!"
	result3 := xtransform.Translate(template3, data, xtransform.WithAtMode(), xtransform.WithBraceMode())
	fmt.Println("@... and {...} modes:", result3)

	// Example 4: Only {...} mode enabled
	template4 := "Welcome to {lang} programming!"
	result4 := xtransform.Translate(template4, data, xtransform.WithBraceMode())
	fmt.Println("Only {...} mode:", result4)

	// Output:
	// Default (all modes): Hello World, welcome to Go!
	// Only @{...} mode: Hello @name, welcome to Golang!
	// @... and {...} modes: Hello World, learn Golang today!
	// Only {...} mode: Welcome to Golang programming!
}