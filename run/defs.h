// some of these are also defined in defs.go

#define GATE_RODATA_ADDR 0x10000
#define GATE_STACK_PAGES 3 // minimum workable value, determined on Linux 4.2
#define GATE_LOADER_FD   2
#define GATE_MAPS_FD     3
#define GATE_WAKEUP_FD   4

#define GATE_ABI_VERSION     0
#define GATE_MAX_PACKET_SIZE 0x10000
