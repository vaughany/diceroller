/*
 * diceroller: A Go module to parse and simulate rolling dice for TTRPGs.
 * Copyright (C) 2024 Paul Vaughan, github.com/vaughany.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package diceroller

import (
	"math/rand/v2"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	random = rand.New(rand.NewPCG(42, 1024))

	os.Exit(m.Run())
}

type rollOneTest struct {
	got  string
	want int
}

var rollOneTests = []rollOneTest{
	{"2d6", 9},
	{"4d4+4", 16},
}

// TestRollOne calls diceroller.RollOne with one valid dice roll string (e.g. '2d6'), checking for valid return values.
func TestRollOne(t *testing.T) {
	for _, test := range rollOneTests {
		output, err := RollOne(test.got)

		if !reflect.DeepEqual(output, test.want) || err != nil {
			t.Errorf("have %v, wanted %v, err %v", output, test.want, err)
		}
	}
}

// BenchmarkRollOne benchmarks diceroller.RollOne.
func BenchmarkRollOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = RollOne(rollOneTests[1].got)
	}
}

type rollTest struct {
	got  []string
	want []int
}

var rollTests = []rollTest{
	{[]string{"2d6", "4d4+4"}, []int{7, 15}},
}

// TestRoll calls diceroller.RollOne with one valid dice roll string (e.g. '2d6'), checking for valid return values.
func TestRoll(t *testing.T) {
	for _, test := range rollTests {
		output, err := Roll(test.got...)

		if !reflect.DeepEqual(output, test.want) || err != nil {
			t.Errorf("have %v, wanted %v, err %v", output, test.want, err)
		}
	}
}

// BenchmarkRoll benchmarks diceroller.Roll.
func BenchmarkRoll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Roll(rollTests[0].got...)
	}
}

type rollTotalTest struct {
	got  []string
	want int
}

var rollTotalTests = []rollTotalTest{
	{[]string{"2d6", "4d4+4"}, 23},
}

// TestRollTotal calls diceroller.RollTotal with one valid dice roll string (e.g. '2d6'), checking for valid return values.
func TestRollTotal(t *testing.T) {
	for _, test := range rollTotalTests {
		output, err := RollTotal(test.got...)

		if !reflect.DeepEqual(output, test.want) || err != nil {
			t.Errorf("have %v, wanted %v, err %v", output, test.want, err)
		}
	}
}

// BenchmarkRollTotal benchmarks diceroller.RollTotal.
func BenchmarkRollTotal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = RollTotal(rollTotalTests[0].got...)
	}
}

type rollDetailsTest struct {
	got  []string
	want []DiceRoll
}

var rollDetailsTests = []rollDetailsTest{
	{[]string{"2d6", "4d4+4"}, []DiceRoll{{"2d6", 6, 2, 0, []int{1, 2}, 3}, {"4d4+4", 4, 4, 4, []int{3, 2, 3, 4}, 16}}},
}

// TestRollDetails calls diceroller.RollDetails with one or more valid dice roll string (e.g. '2d6'), checking for valid return values.
func TestRollDetails(t *testing.T) {
	for _, test := range rollDetailsTests {
		output, err := RollDetails(test.got...)

		if !reflect.DeepEqual(output, test.want) || err != nil {
			t.Errorf("have %v, wanted %v, err %v", output, test.want, err)
		}
	}
}

// BenchmarkRollTotal benchmarks diceroller.RollTotal.
func BenchmarkRollDetails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = RollDetails(rollTotalTests[0].got...)
	}
}

type parseTest struct {
	got  string
	want []string
}

var parseTests = []parseTest{
	{"2d6", []string{"2d6"}},
	{"2 d 6", []string{"2d6"}},
	{"roll a 2d6 please", []string{"2d6"}},
	{"roll a 2 d 6 please", []string{"2d6"}},
	{"roll a 2d6 and 2 d 8 please", []string{"2d6", "2d8"}},
	{"2D6 and 2 D 8", []string{"2D6", "2D8"}},

	{"0d6", []string{"0d6"}},                              // Valid regex match, dumb roll.
	{"try rolling 99999d99999?", []string{"99999d99999"}}, // Who knows what people will try to roll...

	{"roll 2d6+2", []string{"2d6+2"}},
	{"roll 2 d 6 + 2", []string{"2d6+2"}},
	{"then roll 2d8-3", []string{"2d8-3"}},
	{"then roll 2 d 8 - 3", []string{"2d8-3"}},

	{"2\td\t6", []string{"2d6"}},
	{"2\nd\n6", []string{"2d6"}},
	{`2
d
 6`, []string{"2d6"}},
	{"2\td\t6\t+\t2", []string{"2d6+2"}},
	{"2\nd\n6\n+\n2", []string{"2d6+2"}},
	{`2
d
 6
+
2 `, []string{"2d6+2"}},

	{"So 1d6 +2 of something and 2d8-3 harmless something else and 3D12+0 whatever of 8d10+20 nope.", []string{"1d6+2", "2d8-3", "3D12+0", "8d10+20"}},
}

// TestParse calls diceroller.Parse with many strings, checking for valid return values.
func TestParse(t *testing.T) {
	for _, test := range parseTests {
		output, err := Parse(test.got)

		if !reflect.DeepEqual(output, test.want) || err != nil {
			t.Errorf("have %v, wanted %v, err %v", output, test.want, err)
		}
	}
}

// BenchmarkParse benchmarks diceroller.Parse.
func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse(parseTests[2].got)
	}
}

type prettifyTest struct {
	got  []DiceRoll
	want []string
}

var prettifyTests = []prettifyTest{
	{[]DiceRoll{{"2d6", 6, 2, 0, []int{1, 2}, 3}, {"4d4+4", 4, 4, 4, []int{3, 2, 3, 4}, 16}}, []string{"1 + 2 = 3", "3 + 2 + 3 + 4 (+4) = 16"}},
	{[]DiceRoll{{"1d4-1", 1, 4, -1, []int{2}, 1}}, []string{"2 (-1) = 1"}},
}

// TestPrettify calls diceroller.Prettify with one or more valid DiceRoll structs, checking for valid return values.
func TestPrettify(t *testing.T) {
	for _, test := range prettifyTests {
		output := Prettify(test.got)

		if !reflect.DeepEqual(output, test.want) {
			t.Errorf("have %v, wanted %v", output, test.want)
		}
	}
}

// BenchmarkPrettify benchmarks diceroller.Prettify.
func BenchmarkPrettify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Prettify(prettifyTests[0].got)
	}
}

var prettifyFullTests = []prettifyTest{
	{[]DiceRoll{{"2d6", 6, 2, 0, []int{1, 2}, 3}, {"4d4+4", 4, 4, 4, []int{3, 2, 3, 4}, 16}}, []string{"2d6: 1 + 2 = 3", "4d4+4: 3 + 2 + 3 + 4 (+4) = 16"}},
}

// TestPrettifyFull calls diceroller.PrettifyWide with one or more valid DiceRoll structs, checking for valid return values.
func TestPrettifyFull(t *testing.T) {
	for _, test := range prettifyFullTests {
		output := PrettifyFull(test.got)

		if !reflect.DeepEqual(output, test.want) {
			t.Errorf("have %v, wanted %v", output, test.want)
		}
	}
}

// BenchmarkPrettifyFull benchmarks diceroller.PrettifyWide.
func BenchmarkPrettifyFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PrettifyFull(prettifyTests[0].got)
	}
}

type prettifyOneTest struct {
	got  DiceRoll
	want string
}

var prettifyOneTests = []prettifyOneTest{
	{DiceRoll{"2d6", 6, 2, 0, []int{1, 2}, 3}, "1 + 2 = 3"},
}

// TestPrettifyOne calls diceroller.PrettifyOne with one valid DiceRoll structs, checking for valid return values.
func TestPrettifyOne(t *testing.T) {
	for _, test := range prettifyOneTests {
		output := PrettifyOne(test.got)

		if !reflect.DeepEqual(output, test.want) {
			t.Errorf("have %v, wanted %v", output, test.want)
		}
	}
}

// BenchmarkPrettifyOne benchmarks diceroller.PrettifyOne.
func BenchmarkPrettifyOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PrettifyOne(prettifyOneTests[0].got)
	}
}

var prettifyOneFullTests = []prettifyOneTest{
	{DiceRoll{"2d6", 6, 2, 0, []int{1, 2}, 3}, "2d6: 1 + 2 = 3"},
}

// TestPrettifyOneFull calls diceroller.PrettifyOneFull with one valid DiceRoll structs, checking for valid return values.
func TestPrettifyOneFull(t *testing.T) {
	for _, test := range prettifyOneFullTests {
		output := PrettifyOneFull(test.got)

		if !reflect.DeepEqual(output, test.want) {
			t.Errorf("have %v, wanted %v", output, test.want)
		}
	}
}

// BenchmarkPrettifyOneFull benchmarks diceroller.PrettifyOneFull.
func BenchmarkPrettifyOneFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PrettifyOneFull(prettifyOneFullTests[0].got)
	}
}
