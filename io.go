//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package epgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// genarateBucketNames takes a bucket size and returns an array of bucket names
func generateBucketNames(size int) []string {
	var nameList []string
	for i := 0; i < size; i++ {

	}
	return nameList
}

// countToBucketID turns a count assigns it to a letter sequence (e.g. 0-999 is aa, 1000 - 1999 is ab, etc)
func countToBucketID(i int, bucketNames []string) string {
	bucketsize := len(bucketNames)
	// Calculate bucket number
	bucketIndex := i % bucketsize
	return bucketNames[bucketIndex]
}

// generateCombinations from an alphabet and length
//
// this function is based on gist at https://play.golang.org/p/0bWDCibSUJ
func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c) // Once the iteration function is finished, we close the channel

		addLetter(c, "", alphabet, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// addLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
//
// this function is based on gist at https://play.golang.org/p/0bWDCibSUJ
func addLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

func GenerateBucketNames(alphabet string, length int) []string {
	l := []string{}
	for combo := range generateCombinations(alphabet, length) {
		if len(combo) == length {
			l = append(l, combo)
		}
	}
	return l
}

func IntToBucketName(val int, buckets []string) string {
	size := len(buckets)
	return buckets[(val % size)]
}

// WriteJSON write out a data structure as a JSON file.
func WriteJSON(data interface{}, dir string, fname string) error {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return fmt.Errorf("WriteJSON() mkdir %s/%s, %s", dir, fname, err)
	}
	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("WriteJSON() JSON encode %s/%s, %s", dir, fname, err)
	}
	return ioutil.WriteFile(path.Join(dir, fname), src, 0664)
}
