#include <stdio.h>
#include <stdlib.h>

int main()
{
    freopen("./html/post.html", "w", stdout);
    printf("<TITLE>Post</TITLE>");
    printf("<H3>Post</H3>");
    printf("<p>Post to <a href=\"cgi/post\">cgi/post</a> you'll see the change of this page.</p>");
    char *str, ch;
    str = getenv("CONTENT_LENGTH");
    int content_length;
    sscanf(str, "%d", &content_length);
    int i = 0;
    printf("<p>Post Content: ");
    while ((ch = getchar()) != EOF && i < content_length) {
        putchar(ch);
    }
    printf("</p>");
    return 0;
}