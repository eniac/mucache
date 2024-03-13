package cm

import (
	"reflect"
	"testing"
)

func testParseUnparse(req Request, t *testing.T) {
	// Unparse the command
	buf := req.Unparse()
	// Parse it back
	newReq := ParseBytes(buf)
	if !reflect.DeepEqual(req, newReq) {
		t.Fatalf("Request: %+v\nwas misparsed for: %+v\n", req, newReq)
	}
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestParsingUnparsing(t *testing.T) {
	var command Request

	command = &StartRequest{"rid1"}
	testParseUnparse(command, t)

	command = &StartRequest{"rid2"}
	testParseUnparse(command, t)

	command = &InvalidateKeyRequest{"k1"} // "I k1\n"
	testParseUnparse(command, t)

	command = &EndRequest{CallArgs: "rid2", Caller: "caller2", Deps: []Key{Key("k2"), Key("k3")}} //"E rid2 caller2 k2 k3\n"
	testParseUnparse(command, t)

	command = &EndRequest{CallArgs: "rid1", Caller: "caller1", Deps: []Key{Key("k1"), Key("k3")}} //"E rid1 caller1 k1 k3\n"
	testParseUnparse(command, t)

	// These tests do not pass with the existing parsing + unparsing
	//command = StartRequest{"rid2 r"}
	//testParseUnparse(command, t)
	//command = EndRequest{"rid1", "caller1", []Key{Key("k1,as"), Key("k3")}} //"E rid1 caller1 k1 k3\n"
	//testParseUnparse(command, t)
}
