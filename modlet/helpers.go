/*
Copyright © 2025 Donovan C. Young <dyoung522@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/
package modlet

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/donovanmods/7dtd-modtools/lib/logger"
)

// HasPrefix returns true if s starts with prefix
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix returns true if s ends with suffix
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Match returns true if s matches the regex pattern
func Match(s, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		logger.Error("invalid regex pattern %q: %v", pattern, err)
		return false
	}
	return re.MatchString(s)
}

// NotMatch returns true if s does NOT match the regex pattern
func NotMatch(s, pattern string) bool {
	return !Match(s, pattern)
}

// MultValue multiplies a numeric value (or CSV of values) by factor, rounding to int
func MultValue(value string, factor float64) string {
	if value == "" {
		return "0"
	}

	parts := strings.Split(value, ",")
	results := make([]string, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			logger.Error("invalid number %q: %v", part, err)
			results[i] = part
			continue
		}
		result := math.Round(num * factor)
		results[i] = strconv.FormatFloat(result, 'f', 0, 64)
	}

	return strings.Join(results, ",")
}

// ProbMult multiplies a probability value by factor, capping at 1.0
func ProbMult(value string, factor float64) string {
	if value == "" {
		return "0"
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Error("invalid probability %q: %v", value, err)
		return value
	}

	result := num * factor
	if result >= 1.0 {
		return "1"
	}

	// Round to avoid floating point precision issues
	result = math.Round(result*1000000) / 1000000
	formatted := strconv.FormatFloat(result, 'f', -1, 64)
	return formatted
}
