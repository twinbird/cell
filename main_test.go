package main

import "testing"

func TestSimpleNumberExpression(t *testing.T) {
	con := &ExecContext{}
	con.code = "1"
	run(con)
	if con.exitCode != 1 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 1, con.exitCode)
	}
}
