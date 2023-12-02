# gds GoLang data structure

Generics implementation for data structures in Go.

### Usage

#### Tree

###### Red Black Tree 

[Red–black tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree)

```go
package mypkg

import (
	"fmt"
	"github.com/viant/gds/tree/redblack"
)


func ExampleNewTree() {
	tree := redblack.NewTree[int, redblack.Integer]()
	tree.Insert(1)
	tree.Insert(101)
	tree.Insert(54)
	sorted := tree.InOrderTraversal()
	fmt.Printf("%v\n", sorted)
}

```


## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.
