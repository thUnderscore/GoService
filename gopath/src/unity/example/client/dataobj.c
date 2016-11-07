#include "dataobj.h"



void CallStringCallback(stringCallback callBack, void *p, int len){
    callBack(p, len);
}
