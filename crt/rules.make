prog.bc: $(OBJECTS)
	$(LLVMLINK) -o $@ $(CRTOBJECTS) $(OBJECTS)

%.bc: %.c
	$(CLANG) $(CPPFLAGS) $(CFLAGS) -c -o $@ $*.c

%.bc: %.cpp
	$(CLANGPP) $(CPPFLAGS) $(CFLAGS) $(CXXFLAGS) -include $(CRTDIR)/main.hpp -c -o $@ $*.cpp
