local a, b, c
a = types.int(_REQUEST,"a",0)
b = types.int(_REQUEST,"b",0)
c = a+b
return "result = "..types.string(c)