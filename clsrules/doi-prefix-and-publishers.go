//
// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
package clsrules

var (
	doiPrefix = map[string]string{
		"10.1103":  "American Physical Society",
		"10.1063":  "American Institute of Physics",
		"10.1039":  "Royal Society of Chemistry",
		"10.1242":  "Company of Biologists",
		"10.1073":  "PNAS",
		"10.1109":  "IEEE",
		"10.2514":  "AIAA",
		"10.1029":  "AGU (pre-Wiley hosting)",
		"10.1093":  "MNRAS",
		"10.1046":  "Geophysical Journal International",
		"10.1175":  "American Meteorological Society",
		"10.1083":  "Rockefeller University Press",
		"10.1084":  "Rockefeller University Press",
		"10.1085":  "Rockefeller University Press",
		"10.26508": "Rockefeller University Press",
		"10.1371":  "PLOS",
		"10.5194":  "European Geosciences Union",
		"10.1051":  "EDP Sciences",
		"10.2140":  "Mathematical Sciences Publishers",
		"10.1074":  "ASBMB",
		"10.1091":  "ASCB",
		"10.1523":  "Society for Neuroscience",
		"10.1101":  "Cold Spring Harbor",
		"10.1128":  "American Society for Microbiology",
		"10.1115":  "ASME",
		"10.1061":  "ASCE",
		"10.1038":  "Nature",
		"10.1126":  "Science",
		"10.1021":  "American Chemical Society",
		"10.1002":  "Wiley",
		"10.1016":  "Elsevier",
	}
)
