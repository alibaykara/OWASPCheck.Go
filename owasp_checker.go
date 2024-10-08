package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ASCII Art for OWASP Checker
const ASCII_ART = `
         _.-=-._ 
       o~'  '  'o
       \    OWASP  /
        '..____..-' 
   Designed by Ali Baykara
`

// Risk ratings
var riskRatings = map[string]string{
	"A01:2021 - Broken Access Control":                    "high",
	"A02:2021 - Cryptographic Failures":                   "high",
	"A03:2021 - Injection":                                "high",
	"A04:2021 - Insecure Design":                          "medium",
	"A05:2021 - Security Misconfiguration":                "medium",
	"A06:2021 - Vulnerable and Outdated Components":       "medium",
	"A07:2021 - Identification and Authentication Failures": "high",
	"A08:2021 - Software and Data Integrity Failures":     "high",
	"A09:2021 - Security Logging and Monitoring Failures": "low",
	"A10:2021 - Server-Side Request Forgery (SSRF)":       "high",
}

// OWASPTester struct to hold target URL
type OWASPTester struct {
	targetURL string
}

// Function to check for broken access control
func (tester OWASPTester) checkBrokenAccessControl() string {
	resp, err := http.Get(tester.targetURL + "/admin")
	if err == nil && resp.StatusCode == 200 {
		return "A01:2021 - Broken Access Control"
	}
	return ""
}

// Function to check for injection vulnerabilities
func (tester OWASPTester) checkInjection() string {
	payload := "username=admin' OR '1'='1&password=password"
	resp, err := http.Post(tester.targetURL+"/login", "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if strings.Contains(string(body), "Welcome") || resp.StatusCode == 200 {
			return "A03:2021 - Injection"
		}
	}
	return ""
}

// Function to check for security misconfiguration
func (tester OWASPTester) checkSecurityMisconfiguration() string {
	resp, err := http.Get(tester.targetURL + "/phpinfo.php")
	if err == nil && resp.StatusCode == 200 {
		return "A05:2021 - Security Misconfiguration"
	}
	return ""
}

// Function to run all tests and return found vulnerabilities
func (tester OWASPTester) runTests() []string {
	vulnerabilities := []string{}

	// Running each check and appending results if any
	if vuln := tester.checkBrokenAccessControl(); vuln != "" {
		vulnerabilities = append(vulnerabilities, vuln)
	}
	if vuln := tester.checkInjection(); vuln != "" {
		vulnerabilities = append(vulnerabilities, vuln)
	}
	if vuln := tester.checkSecurityMisconfiguration(); vuln != "" {
		vulnerabilities = append(vulnerabilities, vuln)
	}

	return vulnerabilities
}

// Function to get risk rating for a vulnerability
func getRiskRating(vulnerability string) string {
	risk, exists := riskRatings[vulnerability]
	if !exists {
		return "unknown"
	}
	return risk
}

// Main function to run the OWASP checker
func main() {
	// Parsing command-line arguments
	target := flag.String("t", "", "Target URL to test")
	outputFile := flag.String("o", "", "Output file to save the results")
	flag.Parse()

	if *target == "" || *outputFile == "" {
		fmt.Println("Usage: owasp_checker -t <target_url> -o <output_file>")
		return
	}

	// Add "http://" if the URL does not have a scheme
	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "http://" + *target
	}

	// Instantiate OWASPTester
	tester := OWASPTester{targetURL: *target}

	// Run tests
	vulnerabilities := tester.runTests()

	// Prepare the output
	output := fmt.Sprintf("%s\n\nTarget: %s\nResults:\n", ASCII_ART, *target)

	// Add found vulnerabilities to the output
	if len(vulnerabilities) > 0 {
		for _, vuln := range vulnerabilities {
			risk := getRiskRating(vuln)
			output += fmt.Sprintf("  - Vulnerability Found: %s - Risk Level: %s\n", vuln, risk)
		}
	} else {
		output += "  No vulnerabilities found.\n"
	}

	// Write output to file
	err := ioutil.WriteFile(*outputFile, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Print the ASCII art and results to the terminal
	fmt.Println(output)

	// Print completion message
	fmt.Printf("Scan completed. Results saved to %s\n", *outputFile)
}
