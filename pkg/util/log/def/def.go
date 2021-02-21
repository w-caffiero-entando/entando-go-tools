package def

import (
    "entando_go_tools/pkg/util/log"
)

// The default Logger
var base = log.Logger{
    LogVisibleFromDepth: log.Business,
    DefaultLoggingDepth: log.Business,
}

var (
    BizLog     = base.At(log.Business)
    DomLog     = base.At(log.Domain)
    GenLog     = base.At(log.Generic)
    SysLog     = base.At(log.System)
    AbyssalLog = base.At(log.AbyssalCode)
)
