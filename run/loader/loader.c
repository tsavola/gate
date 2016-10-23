#include <stddef.h>
#include <stdint.h>

#include <sys/mman.h>
#include <sys/prctl.h>
#include <sys/syscall.h>
#include <sys/types.h>
#include <linux/prctl.h>
#include <linux/seccomp.h>

#include "../defs.h"

#define xstr(s) str(s)
#define str(s)  #s

__attribute__ ((noreturn))
static void sys_exit(int status)
{
	asm volatile (
		" syscall \n"
		" int3    \n"
		:
		: "a" (SYS_exit), "D" (status)
	);
	__builtin_unreachable();
}

static ssize_t sys_read(int fd, void *buf, size_t count)
{
	ssize_t retval;

	asm volatile (
		"syscall"
		: "=a" (retval)
		: "a" (SYS_read), "D" (fd), "S" (buf), "d" (count)
		: "cc", "rcx", "r11", "memory"
	);

	return retval;
}

static void *sys_mmap(void *addr, size_t length, int prot, int flags, int fd, off_t offset)
{
	void *retval;

	register void *rdi asm ("rdi") = addr;
	register size_t rsi asm ("rsi") = length;
	register int rdx asm ("rdx") = prot;
	register int r10 asm ("r10") = flags;
	register int r8 asm ("r8") = fd;
	register off_t r9 asm ("r9") = offset;

	asm volatile (
		"syscall"
		: "=a" (retval)
		: "a" (SYS_mmap), "r" (rdi), "r" (rsi), "r" (rdx), "r" (r10), "r" (r8), "r" (r9)
		: "cc", "rcx", "r11", "memory"
	);

	return retval;
}

static int sys_mprotect(void *addr, size_t len, int prot)
{
	int retval;

	asm volatile (
		"syscall"
		: "=a" (retval)
		: "a" (SYS_mprotect), "D" (addr), "S" (len), "d" (prot)
		: "cc", "rcx", "r11", "memory"
	);

	return retval;
}

static int read_full(void *buf, size_t size)
{
	for (size_t pos = 0; pos < size; ) {
		ssize_t len = sys_read(0, buf + pos, size - pos);
		if (len < 0)
			return -1;
		pos += len;
	}

	return 0;
}

static void trap_handler(int num)
{
	if (num != 0)
		num += 100;
	sys_exit(num);
}

__attribute__ ((noreturn))
static void enter(uint64_t page_size, void *text_ptr, void *memory_ptr, void *init_memory_limit, void *grow_memory_limit, void *stack_ptr, void *stack_limit)
{
	register void *rax asm ("rax") = stack_ptr;
	register void (*rdx)(int) asm ("rdx") = trap_handler;
	register void *rcx asm ("rcx") = grow_memory_limit;
	register uint64_t rsi asm ("rsi") = GATE_STACK_PAGES * page_size;
	register uint64_t r11 asm ("r11") = page_size;
	register void *r12 asm ("r12") = text_ptr;
	register void *r13 asm ("r13") = stack_limit;
	register void *r14 asm ("r14") = memory_ptr;
	register void *r15 asm ("r15") = init_memory_limit;

	asm volatile (
		// MMX registers
		"        movq    %%rdx, %%mm0                            \n"
		"        movq    %%rcx, %%mm1                            \n"
		// replace stack
		"        mov     %%rsp, %%rdi                            \n"
		"        mov     %%rax, %%rsp                            \n"
		// unmap old stack (fails if stack pointer wasn't somewhere in the last frame)
		"        dec     %%r11                                   \n"
		"        add     %%r11, %%rdi                            \n"
		"        not     %%r11                                   \n"
		"        and     %%r11, %%rdi                            \n"
		"        sub     %%rsi, %%rdi                            \n"
		"        mov     $"xstr(SYS_munmap)", %%eax              \n"
		"        syscall                                         \n"
		"        mov     $41, %%edi                              \n"
		"        test    %%rax, %%rax                            \n"
		"        jne     .exit                                   \n"
		// enable seccomp
		"        mov     $"xstr(PR_SET_SECCOMP)", %%edi          \n"
		"        mov     $"xstr(SECCOMP_MODE_STRICT)", %%esi     \n"
		"        xor     %%edx, %%edx                            \n"
		"        xor     %%r10, %%r10                            \n"
		"        xor     %%r8, %%r8                              \n"
		"        mov     $"xstr(SYS_prctl)", %%eax               \n"
		"        syscall                                         \n"
		"        mov     $42, %%edi                              \n"
		"        test    %%rax, %%rax                            \n"
		"        jne     .exit                                   \n"
		// clear unused registers
		"        xor     %%edx, %%edx                            \n"
		"        xor     %%ecx, %%ecx                            \n"
		"        xor     %%ebx, %%ebx                            \n"
		"        xor     %%ebp, %%ebp                            \n"
		"        xor     %%esi, %%esi                            \n"
		"        xor     %%r8, %%r8                              \n"
		"        xor     %%r9, %%r9                              \n"
		"        xor     %%r10, %%r10                            \n"
		"        xor     %%r11, %%r11                            \n"
		// 0 = no resume
		"        xor     %%eax, %%eax                            \n"
		// skip trap code
		"        mov     %%r12, %%rdx                            \n"
		"        add     $16, %%rdx                              \n"
		"        jmp     *%%rdx                                  \n"
		".exit:                                                  \n"
		"        mov     $"xstr(SYS_exit)", %%rax                \n"
		"        syscall                                         \n"
		"        int3                                            \n"
		:
		: "r" (rax), "r" (rdx), "r" (rcx), "r" (rsi), "r" (r11), "r" (r12), "r" (r13), "r" (r14), "r" (r15)
	);
	__builtin_unreachable();
}

static int main(void)
{
	struct __attribute__ ((packed)) {
		uint32_t page_size;
		uint32_t rodata_size;
		uint32_t text_size;
		uint32_t data_size;
		uint32_t memory_offset;
		uint32_t init_memory_size;
		uint32_t grow_memory_size;
		uint32_t stack_size;
	} info;

	if (read_full(&info, sizeof (info)) != 0)
		return 20;

	if (info.rodata_size > 0) {
		void *ptr = sys_mmap((void *) GATE_RODATA_ADDR, info.rodata_size, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS|MAP_FIXED, -1, 0);
		if (ptr != (void *) GATE_RODATA_ADDR)
			return 21;

		if (read_full(ptr, info.rodata_size) != 0)
			return 24;

		if (sys_mprotect(ptr, info.rodata_size, PROT_READ) != 0)
			return 25;
	}

	void *text_ptr = sys_mmap(NULL, info.text_size, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0);
	if (text_ptr == MAP_FAILED)
		return 21;

	if (read_full(text_ptr, info.text_size) != 0)
		return 22;

	if (sys_mprotect(text_ptr, info.text_size, PROT_EXEC) != 0)
		return 23;

	size_t globals_memory_size = info.memory_offset + info.grow_memory_size;
	void *memory_ptr = NULL;

	if (globals_memory_size > 0) {
		void *ptr = sys_mmap(NULL, globals_memory_size, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS|MAP_NORESERVE, -1, 0);
		if (ptr == MAP_FAILED)
			return 21;

		if (read_full(ptr, info.data_size) != 0)
			return 26;

		memory_ptr = ptr + info.memory_offset;
	}

	void *init_memory_limit = memory_ptr + info.init_memory_size;
	void *grow_memory_limit = memory_ptr + info.grow_memory_size;

	void *stack_limit = sys_mmap(NULL, info.stack_size, PROT_READ|PROT_WRITE, MAP_GROWSDOWN|MAP_PRIVATE|MAP_ANONYMOUS|MAP_NORESERVE|MAP_STACK, -1, 0);
	if (stack_limit == MAP_FAILED)
		return 34;

	void *stack_ptr = stack_limit + info.stack_size;

	enter(info.page_size, text_ptr, memory_ptr, init_memory_limit, grow_memory_limit, stack_ptr, stack_limit);
}

__attribute__ ((noreturn))
void _start(void)
{
	sys_exit(main());
}
