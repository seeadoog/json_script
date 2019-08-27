import time;
def bearsh(a,n):
    lo=0
    hi=len(a)
    while lo<=hi:
        mid=int((lo+hi)/2)
        if a[mid]==n:
            return mid
        if a[mid] < n:
            lo=mid
        else:
            hi=mid+1

    return -1


def getMap():
    a={'name':'lixang','beach':'dfdsf'}
    a['hello']=50
st=time.time()
for i in range(100000):
    getMap()
end=time.time()
print(end-st)