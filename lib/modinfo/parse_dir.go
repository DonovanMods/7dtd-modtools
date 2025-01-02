/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package modinfo

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ParseOpts struct {
	Directory string
	Verbosity int
}

func ParseDir(opts ParseOpts) *ModInfos {
	var directory = opts.Directory
	var verbosity = opts.Verbosity

	modInfos := Make(0)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(info.Name()) == "modinfo.xml" {
			data, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}

			var modInfo ModInfo
			modInfo.SetPath(filepath.Dir(path))

			err = xml.Unmarshal(data, &modInfo)
			if err != nil {
				log.Printf("Error unmarshalling XML from file %s: %v", path, err)
				return nil
			}

			if verbosity >= 2 {
				log.Printf("Read modinfo %q as %q\n", modInfo.Filename(), modInfo.NameStr())
			}

			*modInfos = modInfos.Add(&modInfo)
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %s: %v", directory, err)
	}

	return modInfos
}
