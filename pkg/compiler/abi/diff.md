This ARM64 assembly code demonstrates the difference between passing large structs by value versus by pointer in Go. Here's what's happening:

## Function Call by Value (lines 280-288)
```assembly
MOVD	main.largeStruct-8032(SP), R0    ; Load struct address from stack
MOVD	R0, 8024(RSP)                    ; Copy struct to call stack
CALL	main.processByValue(SB)          ; Call function
```

**Key observations:**
- The large struct (7984 bytes based on earlier assembly) must be **copied entirely** onto the stack
- `8024(RSP)` indicates a large stack offset where the struct copy is placed
- This is expensive because it copies ~8KB of data

## Function Call by Pointer (lines 292-296)
```assembly
MOVD	main.pointerOfLargeStruct-16056(SP), R0  ; Load pointer from stack
CALL	main.processByPointer(SB)                ; Call function
```

**Key observations:**
- Only the **pointer** (8 bytes) is loaded into R0
- No struct copying occurs - just passing the memory address
- Much more efficient for large data structures

## Performance Impact
The assembly clearly shows why passing large structs by pointer is more efficient:
- **By value**: Requires copying the entire struct (7984 bytes)
- **By pointer**: Only passes an 8-byte memory address

The `PCDATA` instructions are debugging metadata for the garbage collector and don't affect the core logic.

This demonstrates Go's calling convention where large structs passed by value incur significant copying overhead, while pointer passing is constant-time regardless of struct size.