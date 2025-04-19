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
package modlet_test

import (
	"fmt"
	"testing"

	"github.com/donovanmods/7dtd-modtools/modlet"
)

func TestValidateTemplate(t *testing.T) {
	assert := setup(t)

	tests := []struct {
		tmplFile   string
		errorMsg   string
		wantErr    bool
		createFile bool
	}{
		{
			tmplFile:   fmt.Sprintf("%s/template.tmpl", testTMP),
			createFile: true,
			wantErr:    false,
		},
		{
			tmplFile:   fmt.Sprintf("%s/template.tmpl.gz", testTMP),
			createFile: true,
			wantErr:    false,
		},
		{
			tmplFile:   fmt.Sprintf("%s/template", testTMP),
			errorMsg:   "does not appear to be a valid template file",
			createFile: true,
			wantErr:    true,
		},
		{
			tmplFile:   fmt.Sprintf("%s/template.gz", testTMP),
			errorMsg:   "does not appear to be a valid template file",
			createFile: true,
			wantErr:    true,
		},
		{
			tmplFile:   "/foo/bar/template.tmpl",
			errorMsg:   "file does not exist",
			createFile: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tmplFile, func(t *testing.T) {
			if tt.createFile {
				if f, err := FS.Create(tt.tmplFile); err != nil {
					t.Fatalf("Error creating template file: %v", err)
				} else {
					_ = f.Close()
				}
			}

			err := modlet.ValidateTemplate(tt.tmplFile)

			if tt.wantErr {
				assert.ErrorContains(err, tt.errorMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}
