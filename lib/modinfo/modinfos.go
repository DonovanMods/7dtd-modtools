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

// ModInfos is a map of ModInfo keyed to the directory it was found in
type ModInfos []*ModInfo

// New creates an empty GameXMLs map and returns a pointer to it
func NewModInfos() *ModInfos {
	modInfos := make(ModInfos, 0, 1)

	return &modInfos
}

func (M *ModInfos) Add(modInfo *ModInfo) ModInfos {
	return append(*M, modInfo)
}

func (M *ModInfos) Find(searchStr string) (*ModInfo, bool) {
	if searchStr == "" {
		return nil, false
	}

	for _, modInfo := range *M {
		switch searchStr {
		case modInfo.Name.Value:
			return modInfo, true
		case modInfo.DisplayName.Value:
			return modInfo, true
		case modInfo.Filename():
			return modInfo, true
		case modInfo.Path():
			return modInfo, true
		}
	}

	return nil, false
}
