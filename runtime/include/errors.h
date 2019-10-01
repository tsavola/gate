// Code generated by internal/runtime-err-gen.  DO NOT EDIT.

// Executor
#define ERR_CONT_EXEC_EXECUTOR 10
#define ERR_EXEC_PRCTL_NOT_DUMPABLE 11
#define ERR_EXEC_SETRLIMIT_DATA 12
#define ERR_EXEC_FCNTL_GETFD 13
#define ERR_EXEC_FCNTL_CLOEXEC 14
#define ERR_EXEC_SIGMASK 16
#define ERR_EXEC_KILL 17
#define ERR_REAP_WAITPID 18
#define ERR_EXEC_PPOLL 19
#define ERR_EXEC_RECVMSG 20
#define ERR_EXEC_SEND 21
#define ERR_EXEC_VFORK 22
#define ERR_EXEC_MSG_CTRUNC 23
#define ERR_EXEC_CMSG_LEVEL 24
#define ERR_EXEC_CMSG_TYPE 25
#define ERR_EXEC_CMSG_LEN 26
#define ERR_EXEC_CMSG_NXTHDR 27
#define ERR_EXEC_SENDBUF_OVERFLOW_CMSG 28
#define ERR_EXEC_SENDBUF_OVERFLOW_REAP 29
#define ERR_EXEC_KILLBUF_OVERFLOW 30
#define ERR_EXEC_DEADBUF_OVERFLOW 31
#define ERR_EXEC_KILLMSG_PID 32
#define ERR_EXEC_PERSONALITY_ADDR_NO_RANDOMIZE 33
#define ERR_EXEC_PRLIMIT 34
#define ERR_EXEC_SETRLIMIT_STACK 36
#define ERR_EXEC_PAGESIZE 37
#define ERR_REAP_SENTINEL 43
#define ERR_EXEC_NODE_ALLOC 44
#define ERR_EXEC_BRK 45
#define ERR_EXEC_MAP_REMOVE 46
#define ERR_REAP_WRITEV 47
#define ERR_REAP_WRITE_ALIGN 48
#define ERR_EXEC_MAP_PID 49
#define ERR_EXEC_MAP_INSERT 50
#define ERR_EXEC_OP 51
#define ERR_EXEC_THREAD_ATTR 52
#define ERR_EXEC_THREAD_CREATE 53
#define ERR_EXEC_SIGACTION 54
#define ERR_EXEC_PRLIMIT_CPU 56
#define ERR_EXEC_FORK_SENTINEL 57
#define ERR_EXEC_KILL_SENTINEL 58
#define ERR_EXEC_MSG_LEN 60
#define ERR_EXEC_CMSG_OP_MISMATCH 62
#define ERR_EXEC_ID_RANGE 63
#define ERR_EXEC_RAISE 64
#define ERR_EXEC_NO_NEW_PRIVS 65
#define ERR_EXEC_CLEAR_CAPS 66
#define ERR_EXEC_PROCSTAT_OPEN 67
#define ERR_EXEC_PROCSTAT_READ 68
#define ERR_EXEC_PROCSTAT_PARSE 69
#define ERR_EXEC_CLOSE 70

// Process
#define ERR_RT_RECVFROM 4
#define ERR_RT_WRITE 5
#define ERR_RT_DEBUG 6
#define ERR_RT_MPROTECT 7
#define ERR_RT_MREMAP 8
#define ERR_RT_CLOCK_GETTIME 9
#define ERR_SENTINEL_SIGNAL_HANDLER 10
#define ERR_SENTINEL_CLOSE 11
#define ERR_SENTINEL_PAUSE 12
#define ERR_EXECHILD_DUP2 13
#define ERR_EXECHILD_EXEC_LOADER 14
#define ERR_LOAD_SETRLIMIT_NOFILE 19
#define ERR_LOAD_SETRLIMIT_NPROC 20
#define ERR_LOAD_PRCTL_NOT_DUMPABLE 21
#define ERR_LOAD_PERSONALITY_DEFAULT 22
#define ERR_LOAD_READ_INFO 23
#define ERR_LOAD_MAGIC_1 24
#define ERR_LOAD_MAGIC_2 25
#define ERR_LOAD_MMAP_VECTOR 26
#define ERR_LOAD_MPROTECT_VECTOR 27
#define ERR_LOAD_MMAP_TEXT 28
#define ERR_LOAD_MMAP_STACK 29
#define ERR_LOAD_MMAP_HEAP 30
#define ERR_LOAD_CLOSE_STATE 31
#define ERR_LOAD_MUNMAP_STACK 32
#define ERR_LOAD_SIGACTION 33
#define ERR_LOAD_MUNMAP_LOADER 34
#define ERR_LOAD_SECCOMP 35
#define ERR_LOAD_ARG_ENV 36
#define ERR_LOAD_NO_VDSO 37
#define ERR_LOAD_FCNTL_OUTPUT 39
#define ERR_LOAD_MPROTECT_HEAP 40
#define ERR_LOAD_CLOSE_TEXT 41
#define ERR_LOAD_SETPRIORITY 42
#define ERR_SENTINEL_PRCTL_PDEATHSIG 44
#define ERR_SENTINEL_SIGMASK 45
#define ERR_LOAD_NO_CLOCK_GETTIME 47
#define ERR_LOAD_READ_TEXT 48
#define ERR_LOAD_MREMAP_HEAP 49
