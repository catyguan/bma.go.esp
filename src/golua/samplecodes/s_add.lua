local a, b
print(_REQUEST)
a = types.int(_REQUEST["a"])
b = types.int(_REQUEST["b"])
c = a+b
print("中文")
return c