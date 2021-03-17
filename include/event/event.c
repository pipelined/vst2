#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include "include/vst.h"

typedef struct Event
{
	int32_t type;
	int32_t byteSize;
	int32_t deltaFrames;	///< sample frames related to the current block start sample position
	int32_t flags;			///< none defined yet (should be zero)
	int32_t dumpBytes;		///< byte size of sysexDump
	int64_t resvd1;		///< zero (Reserved for future use)
	char* sysexDump;		///< sysex dump
	int64_t resvd2;		///< zero (Reserved for future use)
} Event;

Events* newEvents(int32_t numEvents) {
	Events *e = malloc(sizeof(Events));
	e->numEvents = numEvents;
	e->events = malloc(sizeof(e->events) * numEvents);
	return e;
}

void setEvent(Events *events, void *event, int32_t pos) {
	events->events[pos] = event;
}

void* getEvent(Events *events, int32_t pos) {
	return events->events[pos];
}