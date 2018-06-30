#include <stdio.h>
#include <stdlib.h>

int main(void)
{
    char *data;
    long m, n;
    printf("Content-type: text/html\n\n");
    printf("<TITLE>Add Result</TITLE>");
    printf("<H3>Add Result</H3>");

    data = getenv("QUERY_STRING");
    if(data == NULL)
        printf("<P>Don't transfer data or transfer error");
    else if(sscanf(data,"m=%ld&n=%ld",&m,&n)!=2)
        printf("<P>Error, invalid format, data have to number");
    else
        printf("<P>%ld and %ld result: %ld", m, n, m + n);
    return 0;
}