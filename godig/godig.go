/*
**********************************************************
* godig.go
* --------------------------------------------------------
* Clone of dig (not really) from gnu written in Go.
* TODO: Reduce bloat.
* - Parse, print answer(s).
* - Use TTL to determine if domain was cached.
* - Check TLD in processName().
* - PTR records for reverse lookup (ipv4 -> name).
* - User Server input
http://brianc2788.github.io/
**********************************************************
*/
package main

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

/* Constants for domain server options. */
const (
	RESOLVER_URL_SEP         = ":"
	RESOLVER_PORT            = "53"
	RESOLVER_IPV4_GOOGLE     = "8.8.8.8"
	RESOLVER_IPV4_CLOUDFLARE = "1.1.1.1"
	RESOLVER_IPV4_LOCAL      = ""
)

/* MAIN */
func main() {
	argVars := os.Args
	argCount := len(argVars)

	for nArg := 1; nArg < argCount; nArg++ {
		nameIn := argVars[nArg]
		fqname := processName(nameIn)
		rMsg := new(dns.Msg)

		rMsg, merr := getRecordA(fqname)
		if merr != nil {
			panic(merr)
		}

		fmt.Printf("### INPUT ###\n")
		fmt.Printf("name: %s\n", argVars[nArg])
		fmt.Printf("fqdn: %s\n", fqname)
		for _, q := range rMsg.Question {
			fmt.Printf("payload: %s\n", q.Name)
		}
		fmt.Printf("\n")
	}
}

/* Check user's domain names. */
func processName(TestName string) string {
	if !dns.IsFqdn(TestName) {
		return dns.Fqdn(TestName)
	}
	return TestName
}

/* Concatonate domain name variables. */
func catDnsAddr() string {
	fullAddr := RESOLVER_IPV4_GOOGLE + RESOLVER_URL_SEP + RESOLVER_PORT
	return fullAddr
}

/* Resolve names via DNS, print to stdout. */
func getRecordA(dn string) (*dns.Msg, error) {
	var nMsg dns.Msg
	nMsg.SetQuestion(dn, dns.TypeA)
	_, err := dns.Exchange(&nMsg, catDnsAddr())

	if err != nil {
		return nil, fmt.Errorf("Error during dns exchange.\nRuntime Error:\n%s", err)
	}

	return &nMsg, nil
}
