#include <inttypes.h>
#include <stdlib.h>
#include <stdio.h>

// Read an integer from stdin
int64_t read_int() {
  int64_t i;
  scanf("%" SCNd64, &i);
  return i;
}

int64_t add_five(int64_t x) {
  return x + 5;
}

void print_int(int64_t x) {
  printf("OOOOOh here I am!!");
  printf("%" PRId64, x);
}
