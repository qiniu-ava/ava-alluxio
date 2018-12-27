#define _GNU_SOURCE_

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <dlfcn.h>
#include <string.h>
#include <limits.h>
#include <mntent.h>
#include <sys/types.h>
#include <sys/stat.h>

#define MNT_MAX 100
char *seed = ".tmp.ava.alluxiosc.tmp"; 			// keep same with server
char *cache_seed = ".tmp.cache.tmp.ava.alluxiosc.tmp"; 	// keep same with server

typedef struct ava_struct {
  char *mnt_slots[MNT_MAX];
  int mnt_idx;
  int debug;
  int query;
  int cache;
} ava_struct;

ava_struct ava;

FILE* (*my_fopen)(const char *filename, const char* mode);
int  (*my_open)(const char *filename, int flags);
int  (*my_open64)(const char *filename, int flags);
 
int (*my__xstat)(int ver, const char *path, struct stat *buf);
int (*my__xstat64)(int ver, const char *path, struct stat64 *buf);

void __attribute__ ((constructor)) init(void){
  FILE *file= setmntent("/proc/mounts", "r");
  struct mntent *ent;

  ava.debug = (getenv("alluxiosc_debug") != NULL);
  ava.query = (getenv("alluxiosc_query") != NULL);
  ava.cache = (getenv("alluxiosc_cache") != NULL);
  if (ava.cache) seed = cache_seed;
  ava.mnt_idx = 0;

  my_fopen = dlsym(RTLD_NEXT, "fopen");
  my_open = dlsym(RTLD_NEXT, "open");
  my_open64 = dlsym(RTLD_NEXT, "open64");
  my__xstat = dlsym(RTLD_NEXT, "__xstat");
  my__xstat64 = dlsym(RTLD_NEXT, "__xstat64");

  if (file== NULL) {
    perror("setmntent");
    return;
  }

  while (NULL != (ent = getmntent(file))) {
    if (0 == strcmp("alluxio-fuse", ent->mnt_fsname)) {
      ava.mnt_slots[ava.mnt_idx++] = strdup(ent->mnt_dir);
    }
  }
  endmntent(file);
}

int is_alluxio_file(const char *fullpath)
{
  for (int i = 0; i < ava.mnt_idx; i++) {
    char *path = strstr(fullpath, ava.mnt_slots[i]);
    if (path == fullpath) {
      return 1;
    }
  }
  return 0;
}

int get_sc(const char *filename, char *sc, int size, int *af) 
{
  char buf[PATH_MAX + 1], *path = NULL;
  int rc = 0, i = 0;
  FILE *f = NULL;

  *af = 0;
  if (NULL == realpath(filename, buf)) return 0;
  if (!is_alluxio_file(buf)) return 0;
  *af = 1;

  strcat(buf, seed);
  
  if ((f = my_fopen(buf, "r")) && fgets(sc, size, f)) {
    i = strlen(sc);
    if (sc[i - 1] == '\n') sc[i - 1] = '\0';
    rc = ('/' == sc[0]);
  }
  if (f) fclose(f);
  return rc;
}

int __xstat(int ver, const char *path, struct stat *buf)
{
  char full[PATH_MAX + 1];
  if (ava.query) {
    char sc[PATH_MAX + 1];
    int af, rc = get_sc(path, sc, sizeof(sc), &af);
    printf("query=%s\n", sc);
    return my__xstat(ver, path, buf);
  }

  if (NULL == realpath(path, full)) return my__xstat(ver, path, buf);  // error out
  if (!is_alluxio_file(full)) return my__xstat(ver, path, buf);

  if (ava.debug) fprintf(stderr, "--- async cache %s\n", full);
  strcat(full, seed);
  return my__xstat(ver, full, buf);
}

int __xstat64(int ver, const char *path, struct stat64 *buf)
{
  char full[PATH_MAX + 1];
  if (ava.query) {
    char sc[PATH_MAX + 1];
    int af, rc = get_sc(path, sc, sizeof(sc), &af);
    printf("query=%s\n", sc);
    return my__xstat64(ver, path, buf);
  }

  if (NULL == realpath(path, full)) return my__xstat64(ver, path, buf);  // error out
  if (!is_alluxio_file(full)) return my__xstat64(ver, path, buf);

  if (ava.debug) fprintf(stderr, "--- async cache64 %s\n", full);
  strcat(full, seed);
  return my__xstat64(ver, full, buf);
}

FILE* fopen(const char* filename, const char* mode){
  char sc[PATH_MAX + 1];
  memset(sc, 0, sizeof(sc));
  int af, rc = get_sc(filename, sc, sizeof(sc), &af);
  if (ava.debug && af) fprintf(stderr, "--- fopen filename=%s, rc=%d, sc=%s\n", filename, rc, sc);
  FILE *f = rc ? my_fopen(sc, mode) : NULL;
  return f ? f : my_fopen(filename, mode);
}

int open(const char *filename, int flags)
{
  char sc[PATH_MAX + 1];
  memset(sc, 0, sizeof(sc));
  int af, rc = get_sc(filename, sc, sizeof(sc), &af);
  if (ava.debug && af) fprintf(stderr, "--- open filename=%s, rc=%d, sc=%s\n", filename, rc, sc);
  int fd = rc ? my_open(sc, flags) : 0;
  return fd ? fd : my_open(filename, flags);
}

int open64(const char *filename, int flags)
{
  char sc[PATH_MAX + 1];
  memset(sc, 0, sizeof(sc));
  int af, rc = get_sc(filename, sc, sizeof(sc), &af);
  if (ava.debug && af) fprintf(stderr, "--- open64 filename=%s, rc=%d, sc=%s\n", filename, rc, sc);
  int fd = rc ?  my_open64(sc, flags) : 0; 
  return fd ? fd : my_open64(filename, flags);
}

