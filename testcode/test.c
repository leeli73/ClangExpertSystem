#include <stdlib.h>
#include <stdio.h>
#include <string.h>
 
void init_buffer(void)
{
    char filename[128] = {""};
    printf("test of buffer\n");
    dump(filename, 128);
    char unused_buffer[7*1024*1024] = {0};
    char unused_buffer1[1*1024*1024] = {0};
    strcpy(unused_buffer1, "hello");
}
 
void out_of_array(void)
{
    int foo[2];
    int aaa = 250;
    foo[0] = 1;
    foo[1] = 2;
    foo[2] = 3;
    foo[3] = 4;
    foo[4] = 5;
    foo[5] = 6;
    printf("%d %d \n", foo[0], foo[1]);
}
 
#include <sys/types.h>
#include <dirent.h>
void open_not_close()
{
    char* p = new char[100];
    strcpy(p, "hello");
    printf("p:%s\n", p);
    FILE* fp = NULL;
    fp = fopen("aaa", "a");
    if (fp)
    {
        return;
    }
    fclose(fp);
    DIR *dir = NULL;
    dir = opendir("./");
    
}
 
int main(void)
{
    int ret = 0;
    foo();
    init_buffer();
    out_of_array();
    open_not_close()
    return 0;
}