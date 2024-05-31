# Diceroller

[![License](https://img.shields.io/badge/Licence-GNU%20GPL%20v3-blue)](COPYING)
[![Go Report Card](https://goreportcard.com/badge/github.com/vaughany/diceroller)](https://goreportcard.com/report/github.com/vaughany/diceroller)
[![Go Test](https://github.com/vaughany/diceroller/actions/workflows/go.yml/badge.svg)](https://github.com/vaughany/diceroller/actions/workflows/go.yml)
[![CodeQL](https://github.com/vaughany/diceroller/actions/workflows/codeql.yml/badge.svg)](https://github.com/vaughany/diceroller/actions/workflows/codeql.yml)
![Coverage](https://img.shields.io/endpoint?url=https%3A%2F%2Fgist.githubusercontent.com%2Fvaughany%2Fcc4ca9197c72abf20858c15ea662adf6%2Fraw%2Fdiceroller-go-coverage.json)
![CoverageDetails](https://img.shields.io/endpoint?url=https%3A%2F%2Fgist.githubusercontent.com%2Fvaughany%2Fcc4ca9197c72abf20858c15ea662adf6%2Fraw%2Fdiceroller-go-tests.json)

[Diceroller](https://github.com/vaughany/diceroller): A Go module to parse and simulate rolling dice for TTRPGs.

Written by Paul Vaughan, [github.com/vaughany](https://github.com/vaughany).

```go
roll, _ := diceroller.Roll("2d6", "4d4", "5d10+10")
fmt.Printf("%#v\n", roll) // []int{8, 11, 38}
```


## Overview

Discover dice rolls in strings such as `"roll 2d6 and 4d4 please"` with the `Parse` function.

Perform discovered rolls such as `2d6` and `4d4` and get results in various formats with the `Roll` functions.

Get the roll's results returned as nicely-formatted strings with the `Prettify` functions.


## Installation

```bash
go get -u github.com/vaughany/diceroller
```


## Importing

```bash
import (
    "github.com/vaughany/diceroller"
)
```


## Examples of Use

**Note:** all error handling has been removed for brevity.


### Parsing

`Parse()`: Parse one or more strings and return the discovered 'rolls'. A 'roll' is a string in the format `nDn+n` (or `-n`), e.g. `2d6`, `4d4+4`, `1 D4 -1`.  This can then be passed into one of the `Roll` functions.

```go
parse, _ := diceroller.Parse("roll 4d6 +4 please", "and 1d6")
fmt.Printf("%#v\n", parse)
// []string{"4d6+4", "1d6"}
```

**Note:** You can put multiple rolls in a string and this package will attempt to parse them, but be sure to separate them with somthing other than white space, e.g. `"1d6, 2d8"` (comma) or `"1d6 and 2d8"` (the word 'and') are both acceptable. `"1d6 2d8"` will be parsed as `1d62`.


### Rolling 

`RollOne()`: Roll one dice and return the total as an int.

```go
rollOne, _ := diceroller.RollOne("5d10+10")
fmt.Printf("%#v\n", rollOne)
// 37
```


`Roll()`: Roll one or more dice and return the totals as a slice of ints.

```go
roll, _ := diceroller.Roll("2d6", "4d4", "5d10+10")
fmt.Printf("%#v\n", roll)
// []int{2, 6, 44}
```


`RollTotal()`: Roll one or more dice and return the sum total as an int.

```go
rollTotal, _ := diceroller.RollTotal("2d6", "4d4", "5d10+10")
fmt.Printf("%#v\n", rollTotal)
// 65
```


`RollDetails()`: Roll one or more dice and return all the details as a slice of `DiceRoll` structs:

```go
type DiceRoll struct {
	DiscoveredRoll string // The 'nDn+n'-format string we've discovered and are processing.
	Faces          int    // How many faces our dice has: 4, 6, 8, 10, 12 and 20 are common, but we can handle up to 99,999.
	Rolls          int    // How many times we're going to roll the above dice.
	Modifier       int    // A '+n' or '-n' modifier to add to the total, or 0.
	Results        []int  // Each roll, for the curious.
	Total          int    // Total of all rolls.
}
```

**Note:** Notice in the below example, the /2, a typo, is ignored. This is why the 'discoverd' roll is also returned as it may differ from what was passed in.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2", "4d8/2")
fmt.Printf("%#v\n", rollDetails)
// []diceroller.DiceRoll{diceroller.DiceRoll{DiscoveredRoll:"3d6-2", Faces:6, Rolls:3, Modifier:-2, Results:[]int{2, 2, 1}, Total:3}, diceroller.DiceRoll{DiscoveredRoll:"4d8", Faces:8, Rolls:4, Modifier:0, Results:[]int{2, 3, 3, 7}, Total:15}}
```


### Prettifying

For all `Prettify...()` functions, the modifier is omitted if it is zero, and appears in brackets if present.

Additionally, if only one dice is rolled and there's no modifier, the equals sign `=` and total are omitted, because e.g. `6 = 6` is redundant and ugly.

`Prettify()`: Prettify takes in details of one or more rolls and outputs a slice of 'pretty' strings.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2", "4d8/2")
prettify := diceroller.Prettify(rollDetails)
fmt.Printf("%#v\n", prettify)
// []string{"2 + 5 + 3 (-2) = 8", "7 + 3 + 7 + 7 = 24"}
```


`PrettifyFull()`: PrettifyFull takes in details of one or more rolls and outputs a slice of 'pretty' strings, including the discovered roll.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2", "4d8/2")
prettifyFull := diceroller.PrettifyFull(rollDetails)
fmt.Printf("%#v\n", prettifyFull)
// []string{"3d6-2: 2 + 3 + 4 (-2) = 7", "4d8: 5 + 3 + 5 + 5 = 18"}
```


`PrettifyOne()`: PrettifyOne takes in details of one roll and outputs a 'pretty' string.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2")
prettifyOne := diceroller.PrettifyOne(rollDetails[0])
fmt.Printf("%#v\n", prettifyOne)
// "6 + 4 + 5 (-2) = 13"
```


`PrettifyOneFull()`: PrettifyOneFull takes in details of one roll and outputs a 'pretty' string, including the discovered roll.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2")
prettifyOneFull := diceroller.PrettifyOneFull(rollDetails[0])
fmt.Printf("%#v\n", prettifyOneFull)
// "3d6-2: 6 + 2 + 1 (-2) = 7"
```


`PrettifyHTML()`: PrettifyHTML takes in details of one or more rolls and outputs a 'pretty' string with HTML.

```go
rollDetails, _ := diceroller.RollDetails("3d6-2")
prettifyHTML := diceroller.PrettifyHTML(rollDetails)
fmt.Printf("%#v\n", prettifyHTML)
// []string{"<strong>5 + 6 + 2 (-2) = 11</strong>"}
```


`PrettifyHTMLFull()`: PrettifyHTMLFull takes in details of one or more rolls and outputs a 'pretty' string with HTML, including the discovered roll.
```go
rollDetails, _ := diceroller.RollDetails("1d20", "3d6-2")
prettifyHTMLFull := diceroller.PrettifyHTMLFull(rollDetails)
fmt.Printf("%#v\n", prettifyHTMLFull)
// []string{"<strong>1d20:</strong> <em>19</em>", "<strong>3d6-2:</strong> <em>3 + 6 + 3 (-2) = 10</em>"}
```


## Full Example

Below is a full example. Error handling has been removed for brevity.

```go
rollStrings := []string{
    "Roll 4d6 +4 and 8 D4 please.",
    "roll 4 d4 and 4 d6",
    "1d20+3",
}

parsed, _ := diceroller.Parse(rollStrings...)
details, _ := diceroller.RollDetails(parsed...)
pretty := diceroller.PrettifyFull(details)

fmt.Printf("%#v\n", rollStrings)
for i, p := range pretty {
    fmt.Printf("Roll %d: %s.\n", i+1, p)
}
```

Output:

```
[]string{"Roll 4d6 +4 and 8 D4 please.", "roll 4 d4 and 4 d6", "1d20+3"}
Roll 1: 4d6+4: 5 + 3 + 2 + 4 (+4) = 18.
Roll 2: 8d4: 2 + 3 + 1 + 3 + 2 + 4 + 1 + 2 = 18.
Roll 3: 4d4: 3 + 4 + 2 + 2 = 11.
Roll 4: 4d6: 6 + 5 + 1 + 1 = 13.
Roll 5: 1d20+3: 18 (+3) = 21.
```


## History

See the [commit history](https://github.com/vaughany/diceroller/commits/main/) or `git log` for full details.

* **v0.1.8 (2024-05-31):** Added `PrettifyHTML...` functions
* **v0.1.7 (2024-05-29):** Changed the output logic around rolling one dice without modifiers so that we don't get a '6 = 6' output, which is both redundant and ugly.
* **v0.1.6 (2024-05-19):** Another badge.
* **v0.1.5 (2024-05-19):** Added some badges to the readme.
* **v0.1.4 (2024-05-18):** More linting changes.
* **v0.1.3 (2024-05-18):** Handling details and errors returned by the benchmarking functions, in order to keep the linters happy.
* **v0.1.2 (2024-05-18):** Added tests.
* **v0.1.1 (2024-05-13):** Added the readme.
* **v0.1.0 (2024-05-13):** Initial release.


## Contributing

Want to contribute?  [Raise an issue](https://github.com/vaughany/diceroller/issues/new), or [fork the repo](https://github.com/vaughany/diceroller/fork) and submit a pull request. :)


## Licence

[diceroller](https://github.com/vaughany/diceroller) © 2024 by [Paul Vaughan](https://github.com/vaughany) is licensed under the [GNU General Public License v3.0](https://choosealicense.com/licenses/gpl-3.0/).
