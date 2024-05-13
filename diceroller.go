/*
 *   Copyright 2024 Paul Vaughan
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package diceroller

import (
	"fmt"
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DiceRoll struct {
	DiscoveredRoll string // The 'nDn+n'-format string we've discovered and are processing.
	Faces          int    // How many faces our dice has: 4, 6, 8, 10, 12 and 20 are common, but we can handle up to 99,999.
	Rolls          int    // How many times we're going to roll the above dice.
	Modifier       int    // A '+n' or '-n' modifier to add to the total, or 0.
	Results        []int  // Each roll, for the curious.
	Total          int    // Total of all rolls.
}

var (
	// This is the regex used to locate the e.g. 1d6, 2D8+2 rolls. It allows 5-digit numbers (bit daft but whatever).
	diceRollRegex = regexp.MustCompile(`(\d{1,5})[dD](\d{1,5})([\+-]\d{1,5})?`)

	// Pairs of strings: replace spaces, tabs and line endings with nothing.
	inputReplacer = strings.NewReplacer(" ", "", "\t", "", "\n", "")

	// Deterministic random source.
	// random = rand.New(rand.NewPCG(42, 1024))
	// Random random source.
	random = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
)

/*
 * RollOne accepts one string in the correct 'nDn+n' format and returns an int sum of the rolls.
 * e.g. RollOne("2d6") // 7
 */
func RollOne(input string) (int, error) {
	roll, err := roll(input)

	return roll.Total, err
}

/*
 * Roll accepts one or more strings in the correct 'nDn+n' format and returns []int with the totals.
 * e.g. Roll("2d6", "2d8") // []int{7, 12}
 */
func Roll(input ...string) (output []int, err error) {
	var dr DiceRoll

	for _, in := range input {
		dr, err = roll(in)
		if err != nil {
			return
		}

		output = append(output, dr.Total)
	}

	return
}

/*
 * RollTotal accepts one or more strings in the correct 'nDn+n' format and returns an int sum of the rolls.
 * e.g. RollTotal("2d6", "2d8") // 19
 */
func RollTotal(input ...string) (output int, err error) {
	var dr DiceRoll

	for _, in := range input {
		dr, err = roll(in)
		if err != nil {
			return
		}

		output += dr.Total
	}

	return
}

/*
 * RollDetails accepts one or more strings in the correct 'nDn+n' format and returns structs with details of the roll, modifier, and total.
 * e.g. RollTotal("2d6") // [{2d6 [5 2] 0 7}]
 *                       // []diceroller.DiceRoll{diceroller.DiceRoll{DiscoveredRoll:"2d6", Rolls:[]int{5, 2}, Modifier:0, Total:7}}
 */
func RollDetails(input ...string) (output []DiceRoll, err error) {
	var dr DiceRoll

	for _, in := range input {
		dr, err = roll(in)
		if err != nil {
			return
		}

		output = append(output, dr)
	}

	return
}

/*
 * Parse takes in one or more strings and returns a slice of strings containing the discovered dice rolls.
 */
func Parse(input ...string) (output []string, err error) {
	for _, in := range input {
		output = append(output, parse(in)...)
	}

	return
}

/*
 * Prettify takes in a slice of DiceRoll structs and returns a slice of strings with the rolls displayed nicely.
 * e.g. []string{"1 + 2 + 3 + 4 = 10"}
 */
func Prettify(input []DiceRoll) (output []string) {
	output = make([]string, len(input))

	for i, in := range input {
		output[i] = prettify(in, false)
	}

	return
}

/*
 * PrettifyFull takes in a slice of DiceRoll structs and returns a slice of strings with the discovered roll and rolls displayed nicely.
 * e.g. []string{"4d4: 1 + 2 + 3 + 4 = 10"}
 */
func PrettifyFull(input []DiceRoll) (output []string) {
	output = make([]string, len(input))

	for i, in := range input {
		output[i] = prettify(in, true)
	}

	return
}

/*
 * Prettify takes in a DiceRoll struct and returns a string with the rolls displayed nicely.
 * e.g. "1 + 2 + 3 + 4 = 10"
 */
func PrettifyOne(input DiceRoll) (output string) {
	return prettify(input, false)
}

/*
 * Prettify takes in a DiceRoll struct and returns a string with the discovered roll and rolls displayed nicely.
 * e.g. "4d4: 1 + 2 + 3 + 4 = 10"
 */
func PrettifyOneFull(input DiceRoll) (output string) {
	return prettify(input, true)
}

/*
 * roll takes one string in the 'nDn+n' format and rolls that size/face dice that many times, returning a DiceRoll struct with the details.
 */
func roll(input string) (output DiceRoll, err error) {
	// Split the string up into it's component parts.
	result := diceRollRegex.FindStringSubmatch(input)

	// We return the 'discovered' roll so the user knows what we saw.
	// This is important as if we try to process e.g. '2d6/2' (a typo: instead of '2d6+2'),
	//   we'll *actually* be processing '2d6', with no modifier, and the user might not be expecting this.
	output.DiscoveredRoll = result[0]

	// Converting strings to ints.
	output.Rolls, err = strconv.Atoi(result[1])
	if err != nil {
		return
	}

	output.Faces, err = strconv.Atoi(result[2])
	if err != nil {
		return
	}

	// If the modifier's *length* is greater than 0, not if the modifier is greater than zero.
	if len(result[3]) > 0 {
		output.Modifier, err = strconv.Atoi(result[3])
		if err != nil {
			return
		}
	}

	// Pre-allocate the Rolls slice.
	output.Results = make([]int, output.Rolls)

	// Simulate a number of dice being rolled.
	for times := 0; times < output.Rolls; times++ {
		// Roll one dice.
		rolled := random.IntN(output.Faces) + 1

		output.Results[times] = rolled
		output.Total += rolled
	}

	output.Total += output.Modifier

	return
}

/*
 * parse takes in one string and uses a regex to find dice rolls and returns any and all as a slice of strings.
 */
func parse(input string) []string {
	return diceRollRegex.FindAllString(inputReplacer.Replace(input), -1)
}

/*
 * prettify is the function that builds the string from the available data.
 */
func prettify(input DiceRoll, full bool) (output string) {
	var (
		totalsStr = make([]string, len(input.Results))
		total     int
	)

	if full {
		output += strings.ToLower(input.DiscoveredRoll) + ": "
	}

	for i, v := range input.Results {
		totalsStr[i] = strconv.Itoa(v)
		total += v
	}

	output += strings.Join(totalsStr, " + ")

	switch {
	case input.Modifier > 0:
		output += fmt.Sprintf(" (+%d)", input.Modifier)
	case input.Modifier < 0:
		output += fmt.Sprintf(" (-%.0f)", math.Abs(float64(input.Modifier)))
	}

	output += fmt.Sprintf(" = %d", total+input.Modifier)

	return
}
