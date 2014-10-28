local pform = httpserv.post()
go.debug("upload", "post= %v", pform)

local f,fn,sz = httpserv.formFile("pic")
go.debug("upload", "file=%v, fileName=%v, fileSize=%v", f, fn, sz)

return {
	Content="ok"
}