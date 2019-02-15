rm alluxiosc.so
gcc -D_GNU_SOURCE -fPIC -shared -O2 alluxiosc.c -o alluxiosc.so -ldl
