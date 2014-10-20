local vr

vr = json.encode({a=1,b="hello",c={d=true,e=[1,2,3,4]}})
print(json.decode(vr, true))

return vr