// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cover is a program for analyzing the coverage profiles generated by
'go test -coverprofile=cover.out'.

Cover is also used by 'go test -cover' to rewrite the source code with
annotations to track which parts of each function are executed.
It operates on one Go source file at a time, computing approximate
lib block information by studying the source. It is thus more portable
than binary-rewriting coverage tools, but also a little less capable.
For instance, it does not probe inside && and || expressions, and can
be mildly confused by single statements with multiple function literals.

For usage information, please see:
	go help testflag
	go tool cover -help
*/
package main
