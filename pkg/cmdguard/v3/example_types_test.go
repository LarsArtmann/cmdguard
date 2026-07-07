package v3_test

import (
	"fmt"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func ExampleParsePort() {
	p, err := v3.ParsePort("8080")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(p.Int(), p.IsValid())
	// Output: 8080 true
}

func ExampleParsePort_named() {
	p, err := v3.ParsePort("https")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(p.Int())
	// Output: 443
}

func ExampleParseEmail() {
	e, err := v3.ParseEmail("user@example.com")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(e.String())
	// Output: user@example.com
}

func ExampleParseURL() {
	u, err := v3.ParseURL("https://example.com/path")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(u.String())
	// Output: https://example.com/path
}

func ExampleParseDuration() {
	d, err := v3.ParseDuration("1h30m")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(d.String(), d.Seconds())
	// Output: 1h30m0s 5400
}

func ExampleParseHostPort() {
	hp, err := v3.ParseHostPort("localhost:8080")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(hp.String())
	// Output: localhost:8080
}

func ExampleParseFilePath() {
	fp, err := v3.ParseFilePath("/tmp/data.log", false)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(fp.String() != "")
	// Output: true
}

func ExampleParseLogLevel() {
	l, err := v3.ParseLogLevel("debug")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println(l.String())
	// Output: debug
}
