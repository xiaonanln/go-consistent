from __future__ import print_function
import sys
from cffi import FFI

is_64b = sys.maxsize > 2**32

ffi = FFI()
if is_64b: ffi.cdef("typedef long GoInt;\n")
else:      ffi.cdef("typedef int GoInt;\n")

ffi.cdef("""
typedef struct {
    void* data;
    GoInt len;
    GoInt cap;
} GoSlice;

typedef struct {
    const char *data;
    GoInt len;
} GoString;

void Add(GoString host);
void Remove(GoString host);
GoString Hash(GoString key);
void SetReplica(GoInt replica);
""")

lib = ffi.dlopen("./consistent.so")
print(lib.Add)
# print("awesome.Add(12,99) = %d" % lib.Add(12,99))
# print("awesome.Cosine(1) = %f" % lib.Cosine(1))
#
# data = ffi.new("GoInt[]", [74,4,122,9,12])
# nums = ffi.new("GoSlice*", {'data':data, 'len':5, 'cap':5})
# lib.Sort(nums[0])
# print("awesome.Sort(74,4,122,9,12) = %s" % [
#     ffi.cast("GoInt*", nums.data)[i]
#     for i in range(nums.len)])
#
# data = ffi.new("char[]", b"Hello Python!")
# msg = ffi.new("GoString*", {'data':data, 'len':13})
# print("log id %d" % lib.Log(msg[0]))