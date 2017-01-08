#include <stdint.h>

#include <gate.h>

long workaround;
void (*indirection)(const gate_op_header *);

namespace {

void implementation(const gate_op_header *p)
{
	gate_send_packet(p);
}

template <int PayloadSize>
class Packet
{
	Packet(const Packet &) = delete;
	Packet &operator=(const Packet &) = delete;

public:
	enum {
		header_size  = 8,
		payload_size = PayloadSize,
		size         = header_size + payload_size,
	};

	Packet()
	{
		for (int i = 0; i < header_size; i++)
			buf[i] = 0;
	}

	char *payload()
	{
		return buf + header_size;
	}

	const gate_op_header *op_data(enum gate_op_code code, uint16_t flags = 0)
	{
		gate_op_header *header = reinterpret_cast<gate_op_header *> (buf);
		header->size = size;
		header->code = code;
		header->flags = flags;
		return header;
	}

private:
	char buf[size];
};

} // namespace

int main()
{
	indirection = implementation;

	char str[] = "hello world\n";
	Packet<sizeof (str) - 1> p;

	for (int i = 0; i < p.payload_size; i++)
		p.payload()[i] = str[i];

	if (p.size > gate_max_packet_size)
		return 1;

	indirection(p.op_data(GATE_OP_CODE_ORIGIN));

	return 0;
}
