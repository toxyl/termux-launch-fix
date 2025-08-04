# Termux Launch Fix

A Go utility that automatically fixes `os.Args` when running Go applications in Termux on Android. It will also create a start script for the application that uses `proot` to fix problems with DNS resolution.

## `os.Args`

When running Go applications in Termux on Android `os.Args` has an additional, preceding argument. `os.Args[0]` seems to be the path entered by the user and `os.Args[1]` is the absolute path of the binary one would expect in `os.Args[0]`. 

When added as a blank import this utility will automatically strip the first argument from `os.Args` so your application can keep using it like on any other OS.

## DNS Resolution

If your application uses anything that requires DNS lookups it will fail in Termux on Android because Go can't find `/etc/resolv.conf` and falls back to a local resolver. That, however, doesn't exist, causing all DNS queries to fail.

The fix is to use `proot` to launch the application with the correct binding for `resolv.conf`:
```sh
proot -b $PREFIX/etc/resolv.conf:/etc/resolv.conf <path to app>
```

Unfortunately we can't automatically run the application with `proot` from the Go app. Instead this utility will create a start script in the user's home directory (`$HOME/<app name>.proot`) that can be used to launch the app through `proot`. In my use-case this is perfectly fine, but your mileage may vary. Feel free to fork and make your own customizations.

## Usage
Just add a blank import to the Go file serving as entry point:
```go
package main

import (
	"fmt"

	_ "github.com/toxyl/termux-launch-fix"  
)

func main() {
	fmt.Println("The blank import is all you need.")
}
```

## Exit Codes

- `100`: `proot` not found - install with `pkg install proot`
- `101`: Problem resolving home directory
- `102`: Problem writing start script
