package optimizationpatterns

import (
	"reflect"
	"runtime"
	"syscall"
	"unsafe"
)

func pointer8GB() {
	//fmt.Printf("PID: %d\n", os.Getpid())
	a := make([]*int, 1e8)

	for i := 0; i < 10; i++ {
		//start := time.Now()
		runtime.GC()
		//fmt.Printf("GC took %s\n", time.Since(start))
	}

	runtime.KeepAlive(a)
}

func value8GB() {
	//fmt.Printf("PID: %d\n", os.Getpid())
	a := make([]int, 1e8)

	for i := 0; i < 10; i++ {
		//start := time.Now()
		runtime.GC()
		//fmt.Printf("GC took %s\n", time.Since(start))
	}

	runtime.KeepAlive(a)
}

func osDirectlyMapMemory8GB() {
	var example *int
	slice := makeSlice(1e8, unsafe.Sizeof(example))
	a := *(*[]*int)(unsafe.Pointer(&slice))

	for i := 0; i < 10; i++ {
		//start := time.Now()
		runtime.GC()
		//fmt.Printf("GC took %s\n", time.Since(start))
	}

	runtime.KeepAlive(a)
}

func makeSlice(len int, eltsize uintptr) reflect.SliceHeader {
	fd := -1
	data, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0, // address
		uintptr(len)*eltsize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
		uintptr(fd), // No file descriptor
		0,           // offset
	)
	if errno != 0 {
		panic(errno)
	}

	return reflect.SliceHeader{
		Data: data,
		Len:  len,
		Cap:  len,
	}
}
