package alpm

/*
#include <alpm.h>
void go_alpm_set_logging(alpm_handle_t *handle);
*/
import "C"

var DefaultLogLevel = LogWarning

func DefaultLogCallback(lvl uint, s string) {
	if lvl <= DefaultLogLevel {
		print("go-alpm: ", s)
	}
}

var log_callback = DefaultLogCallback

//export LogCallback
func LogCallback(level uint, cstring *C.char) {
	log_callback(level, C.GoString(cstring))
}

func (h *Handle) SetLogCallback(cb func(uint, string)) {
	log_callback = cb
	C.go_alpm_set_logging(h.ptr)
}
