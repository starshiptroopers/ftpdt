// Copyright 2020 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Implements a simple random UID generator
//The UID length, chars and format can be customized
//Math.rand is using as a random generator, seed is initialized with time.Now().UnixNano()

/*
	usage:


 */
package uid

import (
	"errors"
	"math/rand"
	"regexp"
	"time"
)

type UID interface {
	New() string
	Validate(string) (string, error)
	Validator() string
}

//simple short uid generator
type simpleRandomUID struct {
	alfa      		string
	format    		string
	validator 		string
	validatorRgxp 	*regexp.Regexp
	randomGenerator *rand.Rand
}


//Configuration descriptor for UID generator
type Cfg struct {
	Alfa      string  //The chars used in the uid generation, for example "1234567890"
	Format    string  //uid format, every X is replaced with a random generated char, for example "XXX-XXXXXX-XXX"
	Validator string  //uid validation regexp, without (), for example "[0-9]{3}-[0-9]{6}-[0-9]{3}"
}

func NewGenerator(c *Cfg) (UID) {

	//default UID format
	if c == nil {
		c = &Cfg{
			Alfa:      "1234567890",
			Format:    "XXX-XXXXXX-XXX",
			Validator: "[0-9]{3}-[0-9]{6}-[0-9]{3}",
		}
	}
	return &simpleRandomUID{
		c.Alfa,
		c.Format,
		c.Validator,
		regexp.MustCompile("(" + c.Validator + ")"),
		rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// generates a random string according the format
func (s *simpleRandomUID) New() string {
	size := len(s.format)
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		if s.format[i] == 'X' {
			buf[i] = s.alfa[s.randomGenerator.Intn(len(s.alfa))]
		} else {
			buf[i] = s.format[i]
		}
	}
	return string(buf)
}

func (s *simpleRandomUID) SetAlfa(alfa string) {
	s.alfa = alfa
}

func (s *simpleRandomUID) SetFormat(format string, validatingRegexp string) {
	s.format = format
	s.validator = validatingRegexp
}

func (s *simpleRandomUID) Validate(str string) (string, error) {
	matches := s.validatorRgxp.FindStringSubmatchIndex(str)
	if len(matches) > 0 {
		return string(s.validatorRgxp.ExpandString(nil, "$1", str, matches)), nil
	}
	return "", errors.New("uid isn't found in the string")
}

func (s *simpleRandomUID) Validator() string {
	return s.validator
}