# Go Utils
Utilities and packages containing helpful and shared code!

## Usage
```
git module add git@github.com:uptimize/go-utils.git pkg
```
```golang
// get root path
import "<project-name>/pkg/util"

path := util.FromRootPath("binaries")
```
```golang
// json to bytes
import "<project-name>/pkg/util"

data := util.MustJsonToBytes(jsonData)
```
