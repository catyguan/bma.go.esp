package httpmux4glua

import (
	"bmautil/valutil"
	"esp/glua"
	"net/http"
	"strings"
)

type GLuaHttpReq struct {
	glua.ContextLuaInfo
	timeout int
}

type Dispatcher func(req *http.Request, path string) (*GLuaHttpReq, error)

func CommonDispatcher(req *http.Request, path string) (*GLuaHttpReq, error) {
	r := new(GLuaHttpReq)

	r.Script = strings.Replace(path, "/", ".", -1)
	r.FuncName = strings.Replace(path, "/", "_", -1)
	r.timeout = valutil.ToInt(req.FormValue("_to"), 0)

	return r, nil
}
