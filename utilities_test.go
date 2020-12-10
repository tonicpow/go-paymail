package paymail

import (
	"fmt"
	"testing"
	"time"
)

// TestSanitizePaymail will test the method SanitizePaymail()
func TestSanitizePaymail(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		input           string
		expectedAlias   string
		expectedDomain  string
		expectedPaymail string
	}{
		{"test@domain.com", "test", "domain.com", "test@domain.com"},
		{"TEST@domain.com", "test", "domain.com", "test@domain.com"},
		{"TEST@Domain.com", "test", "domain.com", "test@domain.com"},
		{"TEST@DomaiN.COM", "test", "domain.com", "test@domain.com"},
		{"@DomaiN.COM", "", "domain.com", ""},
		{"test@", "test", "", ""},
		{"test@domain", "test", "domain", "test@domain"},
		{"domain.com", "domain.com", "", ""},
		{"1337@Moneybutton.com", "1337", "moneybutton.com", "1337@moneybutton.com"},
	}

	for _, test := range tests {
		if outputAlias, outputDomain, outputPaymail := SanitizePaymail(test.input); outputAlias != test.expectedAlias {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expectedAlias, outputAlias)
		} else if outputDomain != test.expectedDomain {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expectedDomain, outputDomain)
		} else if outputPaymail != test.expectedPaymail {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expectedPaymail, outputPaymail)
		}
	}
}

// ExampleSanitizePaymail example using SanitizePaymail()
//
// See more examples in /examples/
func ExampleSanitizePaymail() {
	alias, domain, address := SanitizePaymail("Paymail@Domain.com")
	fmt.Printf("alias: %s domain: %s address: %s", alias, domain, address)
	// Output:alias: paymail domain: domain.com address: paymail@domain.com
}

// BenchmarkSanitizePaymail benchmarks the method SanitizePaymail()
func BenchmarkSanitizePaymail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = SanitizePaymail("paymail@domain.com")
	}
}

// TestValidatePaymail will test the method ValidatePaymail()
func TestValidatePaymail(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		input         string
		expectedError bool
	}{
		{"test@domain.com", false},
		{"test@example", true},
		{"@example", true},
		{"example", true},
		{"test@.", true},
		{"test@.com", true},
		{"test@.com..", true},
		{"test@.com.", true},
	}

	for _, test := range tests {
		if err := ValidatePaymail(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
	}
}

// ExampleValidatePaymail example using ValidatePaymail()
//
// See more examples in /examples/
func ExampleValidatePaymail() {
	err := ValidatePaymail("bad@paymail")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("example failed")
	}
	// Output:paymail address failed format validation: email is not a valid address format
}

// BenchmarkValidatePaymail benchmarks the method ValidatePaymail()
func BenchmarkValidatePaymail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidatePaymail("paymail@domain.com")
	}
}

// TestValidateDomain will test the method ValidateDomain()
func TestValidateDomain(t *testing.T) {

	t.Parallel()

	var tests = []struct {
		input         string
		expectedError bool
	}{
		{"domain", true},
		{"example.com", false},
		{"google.com", false},
		{"test@domain.com", true},
		{"example..", true},
		{"example.c", false},
	}

	for _, test := range tests {
		if err := ValidateDomain(test.input); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.input, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.input)
		}
	}
}

// ExampleValidateDomain example using ValidateDomain()
//
// See more examples in /examples/
func ExampleValidateDomain() {
	err := ValidateDomain("domain")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("example failed")
	}
	// Output:domain name is invalid: domain
}

// BenchmarkValidateDomain benchmarks the method ValidateDomain()
func BenchmarkValidateDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateDomain("domain.com")
	}
}

// TestConvertHandle will test the method ConvertHandle()
func TestConvertHandle(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		input    string
		beta     bool
		expected string
	}{
		{"$mr-z", false, "mr-z@handcash.io"},
		{"$MR-Z", false, "mr-z@handcash.io"},
		{"invalid$mr-z", false, "invalid$mr-z"},
		{"$", false, "@handcash.io"},
		{"$", true, "@beta.handcash.io"},
		{"1handle", false, "handle@relayx.io"},
		{"1337@moneybutton.com", false, "1337@moneybutton.com"},
		{"1337", false, "337@relayx.io"},
		{"1PN7K19Jmj7QQCpUg37WHpSRUw5gKhJVRa", false, "1PN7K19Jmj7QQCpUg37WHpSRUw5gKhJVRa"},
		{"$misterz", true, "misterz@beta.handcash.io"},
	}

	for _, test := range tests {
		if output := ConvertHandle(test.input, test.beta); output != test.expected {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		}
	}
}

// ExampleConvertHandle example using the method ConvertHandle()
//
// See more examples in /examples/
func ExampleConvertHandle() {
	paymail := ConvertHandle("$mr-z", false)
	fmt.Println(paymail)
	// Output:mr-z@handcash.io
}

// BenchmarkConvertHandle benchmarks the method ConvertHandle()
func BenchmarkConvertHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertHandle("$mr-z", false)
	}
}

// TestValidateTimestamp will test the method ValidateTimestamp()
func TestValidateTimestamp(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		timestamp     string
		expectedError bool
	}{
		{"", true},
		{"0000-00-00T00:00:00Z", true},
		{"12345", true},
		{"2017", true},
		{"2018-01-01", true},
		{"2020-04-09 12:00", true},
		{"2020-04-09 12:00:00", true},
		{"2020-04-09 12:00B", true},
		{"2020-04-09T12:00:00", true},
		{"abcdef", true},
		{time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339), false},
		{time.Now().UTC().Add(-118 * time.Second).Format(time.RFC3339), false},
		{time.Now().UTC().Add(-122 * time.Second).Format(time.RFC3339), true},
		{time.Now().UTC().Add(-3 * time.Minute).Format(time.RFC3339), true},
		{time.Now().UTC().Add(-4 * time.Minute).Format(time.RFC3339), true},
		{time.Now().UTC().Add(1 * time.Minute).Format(time.RFC3339), false},
		{time.Now().UTC().Add(122 * time.Second).Format(time.RFC3339), true},
		{time.Now().UTC().Add(3 * time.Minute).Format(time.RFC3339), true},
		{time.Now().UTC().Add(4 * time.Minute).Format(time.RFC3339), true},
		{time.Now().UTC().Format(time.RFC3339), false},
	}

	for _, test := range tests {
		if err := ValidateTimestamp(test.timestamp); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error not expected but got: %s", t.Name(), test.timestamp, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] inputted and error was expected", t.Name(), test.timestamp)
		}
	}
}

// ExampleValidateTimestamp example using the method ValidateTimestamp()
//
// See more examples in /examples/
func ExampleValidateTimestamp() {
	err := ValidateTimestamp(time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	} else {
		fmt.Printf("timestamp was valid!")
	}
	// Output:timestamp was valid!
}

// BenchmarkValidateTimestamp benchmarks the method ValidateTimestamp()
func BenchmarkValidateTimestamp(b *testing.B) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	for i := 0; i < b.N; i++ {
		_ = ValidateTimestamp(timestamp)
	}
}
