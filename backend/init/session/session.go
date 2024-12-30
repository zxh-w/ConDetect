package session

import (
	"ConDetect/backend/global"
	"ConDetect/backend/init/session/psession"
)

func Init() {
	global.SESSION = psession.NewPSession(global.CACHE)
	global.LOG.Info("init session successfully")
}
