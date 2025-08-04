package termuxlaunchfix

import "fmt"

func makeStartScript(execPath, prootPath string) string {
	return fmt.Sprintf(`#!/bin/bash
if [ ! -f %s.bin ]; then
  mv %s %s.bin
  echo '#!/bin/bash' > %s
  echo '%s -b $PREFIX/etc/resolv.conf:/etc/resolv.conf %s.bin $@' >> %s
  chmod +x %s
fi
%s $@
`,
		execPath,
		execPath,
		execPath,
		execPath,
		prootPath,
		execPath,
		execPath,
		execPath,
		execPath)
}
