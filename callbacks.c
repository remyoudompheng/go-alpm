#include <stdio.h>
#include <stdarg.h>
#include <alpm.h>
#include "_cgo_export.h"

void go_alpm_log_cb(alpm_loglevel_t level, const char *fmt, va_list arg) {
  char *s = malloc(128);
  if (s == NULL) return;
  int length = vsnprintf(s, 128, fmt, arg);
  if (length > 128) {
    length = (length + 16) & ~0xf;
    s = realloc(s, length);
  }
  if (s != NULL) {
		LogCallback(level, s);
		free(s);
  }
}

void go_alpm_set_logging(alpm_handle_t *handle) {
  alpm_option_set_logcb(handle, go_alpm_log_cb);
}

