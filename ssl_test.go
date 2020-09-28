package paymail

import (
	"fmt"
	"testing"
)

// TestClient_CheckSSL will test the method CheckSSL()
func TestClient_CheckSSL(t *testing.T) {
	// t.Parallel() (turned off - race condition)

	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		host          string
		expectedValid bool
		expectedError bool
	}{
		{"google.com", true, false},
		{"google", false, true},
		{"", false, true},
		{"domaindoesntexistatall101910.co", false, true},
		{"moneybutton.com", true, false},
	}

	// Test all
	for _, test := range tests {
		if valid, err := client.CheckSSL(test.host); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.host, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.host)
		} else if valid != test.expectedValid {
			t.Errorf("%s Failed: [%s] inputted and valid was not as expected", t.Name(), test.host)
		}
	}
}

// ExampleClient_CheckSSL example using CheckSSL()
//
// See more examples in /examples/
func ExampleClient_CheckSSL() {
	client, _ := NewClient(nil, nil)
	valid, _ := client.CheckSSL("moneybutton.com")
	if valid {
		fmt.Printf("valid SSL certificate found for: %s", "moneybutton.com")
	} else {
		fmt.Printf("invalid SSL certificate found for: %s", "moneybutton.com")
	}

	// Output:valid SSL certificate found for: moneybutton.com
}

// BenchmarkClient_CheckSSL benchmarks the method CheckSSL()
func BenchmarkClient_CheckSSL(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.CheckSSL("moneybutton.com")
	}
}
