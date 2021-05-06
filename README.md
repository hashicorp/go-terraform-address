# tf-address

[![CircleCI](https://circleci.com/gh/hashicorp/go-terraform-address.svg?style=svg)](https://app.circleci.com/pipelines/github/hashicorp/go-terraform-address)

This package provides utilities for properly parsing Terraform addresses.

The parser is implemented using [pigeon][1], a [PEG][2] implementation.

## Addresses

Resource addresses are described by the [Terraform documentation][3].

Identifiers are described in the [Terraform Configuration Syntax document][4]

## Generating

If you change the peg, please regenerate the go code with:

```sh
go get -u github.com/mna/pigeon
go generate .
```

You may need to clean up the go.mod with:

```sh
go mod tidy
```

## Examples

```go
package main

import (
	"fmt"

	address "github.com/hashicorp/go-terraform-address"
)

func main() {

	a, err := address.NewAddress(`module.first.module.second["xyz"].resource.name[2]`)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(a.ModulePath))                  // 2
	fmt.Println(a.ModulePath[0].Name)               // "first"
	fmt.Println(a.ModulePath[1].Index.String())     // "xyz"
	fmt.Printf("%T\n", a.ModulePath[1].Index.Value) // string
	fmt.Println(a.ResourceSpec.Type)                // "resource"
	fmt.Println(a.ResourceSpec.Name)                // "name"
	fmt.Printf("%T\n", a.ResourceSpec.Index.Value)  // int
}
```


[1]: https://godoc.org/github.com/mna/pigeon
[2]: https://en.wikipedia.org/wiki/Parsing_expression_grammar
[3]: https://github.com/hashicorp/terraform/blob/ef071f3d0e49ba421ae931c65b263827a8af1adb/website/docs/internals/resource-addressing.html.markdown#resource-addressing
[4]: https://www.terraform.io/docs/configuration/syntax.html#identifiers
