#ifndef LIBXMTDMLS_H
#define LIBXMTDMLS_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct ValidationResult {
  bool ok;
  char *message;
} ValidationResult;

struct ValidationResult validate_inbox_id_key_package_ffi(const unsigned char *data_ptr,
                                                          size_t data_len);

void free_c_string(char *s);

#endif  /* LIBXMTDMLS_H */
