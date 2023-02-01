// Package opt implements "optional" values.
// An Opt can either hold an error or a value. The special error ErrEmpty signifies that the option is empty.
// Options are not "futures" or "promises" but hold a pre-computed value.
package opt
