# Version number reflecting the contract between Go and C.
GATE_RUNTIME_BINARY_SUFFIX	:= 0

GATE_SANDBOX			:= 1

CPPFLAGS += -DGATE_RUNTIME_BINARY_SUFFIX=$(GATE_RUNTIME_BINARY_SUFFIX) -DGATE_SANDBOX=$(GATE_SANDBOX)
