## go-unsafe

Runtime interface type cast (and modification?)

since compiled types are live on binary .rodata, it is not writable*, we need to rebuild the type on runtime, calculate the hashcode (fnv1), assign methods etc etc.

interface represents (aligned 8*2):
- rtype ptr, where used as type reference.
- word/ptr data, where used as reference to data.


interface method implementation (dict.type|interface):
- another rtype, contains functions/methods coresponding to itab.
- referencing primitive type on `(*rtype).ptrdata`.


runtime inspect type modification:
- example

generic pointer casting/reinterpretation:
- example:
    ```go
    var (
        a = new(uint64)
        b *int64
    )
    CastPtr(&b, a)
    ```
    ```go
    var a = new(uint64)
    b := ReinterpretPtr[int64](a) // b *int64
    ```
    or if you wish to
    ```go
    b := ReinterpretPtr[uint64](0x600000) // *uint64
    ```