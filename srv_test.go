package paymail

import (
	"fmt"
	"net"
	"testing"
)

// TestClient_GetSRVRecord will test the method GetSRVRecord()
func TestClient_GetSRVRecord(t *testing.T) {
	// t.Parallel() (turned off - race condition)

	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		service       string
		protocol      string
		domainName    string
		expectedError bool
		expectedNil   bool
	}{
		{DefaultServiceName, DefaultProtocol, "domain.com", true, true},
		{"", "", "domain.com", true, true},
		{"", DefaultProtocol, "domain.com", true, true},
		{DefaultServiceName, "", "domain.com", true, true},
		{"", "", "", true, true},
		{DefaultServiceName, DefaultProtocol, "", true, true},
		{"bogus", DefaultProtocol, "moneybutton.com", true, true},
		{DefaultServiceName, DefaultProtocol, "moneybutton.com", false, false},
		{DefaultServiceName, DefaultProtocol, "relayx.io", false, false},
		{DefaultServiceName, DefaultProtocol, "mypaymail.co", false, false},
	}

	// Test all
	for _, test := range tests {
		if srv, err := client.GetSRVRecord(test.service, test.protocol, test.domainName); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and error not expected but got: %s", t.Name(), test.service, test.protocol, test.domainName, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and error was expected", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv should have not been nil", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv != nil && test.expectedNil {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv should have been nil", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv != nil && srv.Port == 0 {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv port was empty", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv != nil && srv.Priority == 0 {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv priority was empty", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv != nil && srv.Weight == 0 {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv weight was empty", t.Name(), test.service, test.protocol, test.domainName)
		} else if srv != nil && len(srv.Target) == 0 {
			t.Errorf("%s Failed: [%s] [%s] [%s] inputted and srv target was empty", t.Name(), test.service, test.protocol, test.domainName)
		}
	}
}

// ExampleClient_GetSRVRecord example using GetSRVRecord()
func ExampleClient_GetSRVRecord() {
	client, _ := NewClient(nil, nil)
	srv, _ := client.GetSRVRecord(DefaultServiceName, DefaultProtocol, "moneybutton.com")
	fmt.Printf("port: %d priority: %d weight: %d target: %s", srv.Port, srv.Priority, srv.Weight, srv.Target)
	// Output:port: 443 priority: 1 weight: 10 target: www.moneybutton.com
}

// BenchmarkClient_GetSRVRecord benchmarks the method GetSRVRecord()
func BenchmarkClient_GetSRVRecord(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_, _ = client.GetSRVRecord(DefaultServiceName, DefaultProtocol, "moneybutton.com")
	}
}

// TestClient_ValidateSRVRecord will test the method ValidateSRVRecord()
func TestClient_ValidateSRVRecord(t *testing.T) {

	// t.Parallel() (turned off - race condition)

	client, err := newTestClient()
	if err != nil {
		t.Fatalf("error loading client: %s", err.Error())
	}

	// Create the list of tests
	var tests = []struct {
		srv           *net.SRV
		port          uint16
		priority      uint16
		weight        uint16
		expectedError bool
	}{
		{&net.SRV{Target: "domain.com", Port: DefaultPort, Priority: DefaultPriority, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, false},
		{&net.SRV{Target: "domain.com", Port: DefaultPort, Priority: DefaultPriority, Weight: DefaultWeight}, 0, 0, 0, false},
		{&net.SRV{Target: "", Port: DefaultPort, Priority: DefaultPriority, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{&net.SRV{Target: "domain", Port: DefaultPort, Priority: DefaultPriority, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{&net.SRV{Target: "domain.com", Port: 123, Priority: DefaultPriority, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{&net.SRV{Target: "domain.com", Port: DefaultPort, Priority: 123, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{&net.SRV{Target: "domain.com", Port: DefaultPort, Priority: DefaultPriority, Weight: 123}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{&net.SRV{Target: "baddomain10901919.com", Port: DefaultPort, Priority: DefaultPriority, Weight: DefaultWeight}, DefaultPort, DefaultPriority, DefaultWeight, true},
		{nil, DefaultPort, DefaultPriority, DefaultWeight, true},
	}

	// Test all
	for _, test := range tests {
		if err := client.ValidateSRVRecord(test.srv, test.port, test.priority, test.weight); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] [%d] inputted and error not expected but got: %s", t.Name(), test.srv, test.port, test.priority, test.weight, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] [%d] [%d] [%d] inputted and error was expected", t.Name(), test.srv, test.port, test.priority, test.weight)
		}
	}
}

// ExampleClient_ValidateSRVRecord example using ValidateSRVRecord()
func ExampleClient_ValidateSRVRecord() {
	client, _ := NewClient(nil, nil)
	err := client.ValidateSRVRecord(&net.SRV{
		Target:   "moneybutton.com",
		Port:     DefaultPort,
		Priority: 1,
		Weight:   DefaultWeight,
	}, DefaultPort, DefaultPriority, DefaultWeight)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	// Output:error: srv priority 1 does not match 10
}

// BenchmarkClient_ValidateSRVRecord benchmarks the method ValidateSRVRecord()
func BenchmarkClient_ValidateSRVRecord(b *testing.B) {
	client, _ := NewClient(nil, nil)
	for i := 0; i < b.N; i++ {
		_ = client.ValidateSRVRecord(&net.SRV{
			Target:   "moneybutton.com",
			Port:     DefaultPort,
			Priority: DefaultPriority,
			Weight:   DefaultWeight,
		}, DefaultPort, DefaultPriority, DefaultWeight)
	}
}
