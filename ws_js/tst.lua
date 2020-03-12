local start = os.date()
print(start)

function fabio(i)
    if i==0 then return 0 end
    if i==1 then return 1 end
    return fabio(i-1)+fabio(i-2)
end
fabio(36)
local start = os.date()
print(start)