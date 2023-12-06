/*
 * CBOR Example code in C
 * Copyright 2023 John Douglas Pritchard.  All rights reserved.
 * 
 * RFC8949 Appendix C Pseudocode
 * Copyright (c) 2020 IETF Trust and Bormann and Hoffman.  All rights reserved.
 * Carsten Bormann, Universität Bremen TZI
 * Paul Hoffman, ICANN
 *
 * Redistribution and use in source and binary forms, with
 * or without modification, are permitted provided that the
 * following conditions are met:
 *
 * 1. Redistributions of source code must retain the above
 * copyright notice, this list of conditions and the
 * following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the
 * above copyright notice, this list of conditions and the
 * following disclaimer in the documentation and/or other
 * materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS “AS IS” AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
 * PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
 * OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
 * DAMAGE.
 */
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <unistd.h>
/*
 * CPL/LIBC
 */
#define null 0
#define true 1
#define false 0

typedef uint8_t bool_t;
typedef uint8_t tag_t;
typedef uint8_t byte;
typedef int fd_t;

const fd_t stdfd_in = 0;
const fd_t stdfd_out = 1;
const fd_t stdfd_err = 3;
/*
 * CBOR
 */
#define stop 0xFF

tag_t well_formed_indefinite(tag_t, bool_t);

byte* take(size_t n){
  byte *bary = calloc(n,1);
  if (null != bary){
    ssize_t rd = read(stdfd_in,bary,n);
    if (rd == n){

      return bary;
    }
    else {
      free(bary);
    }
  }
  return null;
}
uint64_t uint64(byte *bary){
  if (null != bary){
    uint64_t a = bary[0];
    uint64_t b = bary[1];
    uint64_t c = bary[2];
    uint64_t d = bary[3];
    uint64_t e = bary[4];
    uint64_t f = bary[5];
    uint64_t g = bary[6];
    uint64_t h = bary[7];

    uint64_t hi = ((a << 24)|(b << 16)|(c << 8)|(d));
    uint64_t lo = ((e << 24)|(f << 16)|(g << 8)|(h));

    uint64_t value = ((hi << 32)|(lo));
  
    free(bary);
    return value;
  }
  else {
    return 0;
  }
}
uint32_t uint32(byte *bary){
  if (null != bary){
    uint32_t a = bary[0];
    uint32_t b = bary[1];
    uint32_t c = bary[2];
    uint32_t d = bary[3];

    uint32_t value = ((a << 24)|(b << 16)|(c << 8)|(d));

    free(bary);
    return value;
  }
  else {
    return 0;
  }
}
uint16_t uint16(byte *bary){
  if (null != bary){
    uint16_t a = bary[0];
    uint16_t b = bary[1];

    uint16_t value = ((a << 8)|(b));

    free(bary);
    return value;
  }
  else {
    return 0;
  }
}
uint8_t uint8(byte *bary){
  if (null != bary){
    uint8_t value = bary[0];

    free(bary);
    return value;
  }
  else {
    return 0;
  }
}
void fail(){
  exit(1);
}

tag_t well_formed(bool_t breakable) {
  // process initial bytes
  uint8_t ib = uint8(take(1));
  uint8_t mt = (ib >> 5);
  uint8_t ai = (ib & 0x1F);
  switch (ai) {
  case 24:{
    uint8_t val = uint8(take(1));
    // process content
    switch (mt) {
      // case 0, 1, 7 do not have content; just use val
    case 2: case 3:{
      void *utfstr = take(val);
      if (null != utfstr){
	free(utfstr);
      }
      break; // bytes/UTF-8
    }
    case 4:{
      uint8_t i;
      for (i = 0; i < val; i++) well_formed(false);
      break;
    }
    case 5:{
      uint8_t i;
      for (i = 0; i < val*2; i++) well_formed(false);
      break;
    }
    case 6:
      well_formed(false);
      break;     // 1 embedded data item
    case 7:
      if (ai == 24 && val < 32) fail(); // bad simple
    }
    break;
  }
  case 25:{
    uint16_t val = uint16(take(2));
    // process content
    switch (mt) {
      // case 0, 1, 7 do not have content; just use val
    case 2: case 3:
      take(val);
      break; // bytes/UTF-8
    case 4:{
      uint16_t i;
      for (i = 0; i < val; i++) well_formed(false);
      break;
    }
    case 5:{
      uint16_t i;
      for (i = 0; i < val*2; i++) well_formed(false);
      break;
    }
    case 6:
      well_formed(false);
      break;     // 1 embedded data item
    case 7:
      if (ai == 24 && val < 32) fail(); // bad simple
    }
    break;
  }
  case 26:{
    uint32_t val = uint32(take(4));
    // process content
    switch (mt) {
      // case 0, 1, 7 do not have content; just use val
    case 2: case 3:
      take(val);
      break; // bytes/UTF-8
    case 4:{
      uint16_t i;
      for (i = 0; i < val; i++) well_formed(false);
      break;
    }
    case 5:{
      uint16_t i;
      for (i = 0; i < val*2; i++) well_formed(false);
      break;
    }
    case 6:
      well_formed(false);
      break;     // 1 embedded data item
    case 7:
      if (ai == 24 && val < 32) fail(); // bad simple
    }
    break;
  }
  case 27:{
    uint64_t val = uint64(take(8));
    // process content
    switch (mt) {
      // case 0, 1, 7 do not have content; just use val
    case 2: case 3:
      take(val);
      break; // bytes/UTF-8
    case 4:{
      uint64_t i;
      for (i = 0; i < val; i++) well_formed(false);
      break;
    }
    case 5:{
      uint64_t i;
      for (i = 0; i < val*2; i++) well_formed(false);
      break;
    }
    case 6:
      well_formed(false);
      break;     // 1 embedded data item
    case 7:
      if (ai == 24 && val < 32) fail(); // bad simple
    }
    break;
  }
  case 28: case 29: case 30:
    fail();
  case 31:
    return well_formed_indefinite(mt, breakable);
  }
  return mt;                    // definite-length data item
}

tag_t well_formed_indefinite(tag_t mt, bool_t breakable) {
  tag_t it;
  switch (mt) {
  case 2: case 3:
    while (stop != (it = well_formed(true))){
      if (it != mt)           // need definite-length chunk
	fail();               //    of same type
    }
    break;
  case 4:
    while (stop != well_formed(true));
    break;
  case 5:
    while (stop != well_formed(true)) well_formed(false);
    break;
  case 7:
    if (breakable)
      return stop;              // signal break out
    else
      fail();              // no enclosing indefinite
  default:
    fail();            // wrong mt
  }
  return 99;                    // indefinite-length data item
}
/*
 * Compiles with "clang -o wellformed wellformed.c".
 */
int main(int argc, char **argv){

  tag_t m = well_formed(false);

  if (stop == m){

    printf("tag <stop>\n");

  } else if (null == m){

    printf("tag <null>\n");

  } else {

    printf("tag 0x%x\n",m);
  }
  return 0;
}
