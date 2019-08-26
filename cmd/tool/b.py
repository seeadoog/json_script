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
st=time.time()
for i in range(20000):

    a=[0,1,2,3,4,5,6,7,8,9]
    bearsh(a,8)
end=time.time()
print(end-st)