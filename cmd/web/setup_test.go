// Run this before our test run. It won't run in our compile program.

package main

import (
	"net/http"
	"os"
	"testing"
)

/*
It is sometimes necessary for a test program to do extra setup or teardown before or after testing. It is also sometimes necessary for a test to control which code runs on the main thread. To support these and other cases, if a test file contains a function:

func TestMain(m *testing.M)
then the generated test will call TestMain(m) instead of running the tests directly.
*/

// middleware test or main test
func TestMain(m *testing.M) {
	// do something inside this funciton "TestMain", run the test "m.Run", and then exit "os.Exit"

	os.Exit(m.Run()) // stop after finished
}

type myHandler struct{}

// create the handler tird to myHandler that can satisfy ServeHTTP
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only implemented, not return anything yet
}
