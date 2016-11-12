#include <stdint.h>
#ifndef DATAOBJ_H
#define DATAOBJ_H


typedef void (*stringCallback)(void*, int);

extern void CallStringCallback(stringCallback, void*, int);

typedef struct GoStatisticTag
{	
	int64_t Interval;	
	int NumGoroutine;  
	uint64_t Alloc;
	uint64_t Mallocs;
	uint64_t Frees;
	uint64_t HeapAlloc;
	uint64_t StackInuse;
	uint64_t PauseTotalNs;
	uint64_t NumGC;			
} GoStatistic;

typedef struct ClientConnectorTag
{
	int counterValue;	 
	stringCallback log;
} ClientConnector;




#endif