package blog

import server "blog/server/system"

func init() {
	server.RegisterInit(initOrderApi, &initApi{})
}
