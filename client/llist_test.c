#include <stdio.h>
#include "llist.h"

void test()
{
    struct llist_t list;
    llist_init(&list);

    int a = 0, b = 1, c = 2;
    llist_addf(&list, &a);
    llist_addf(&list, &b);
    llist_addf(&list, &c);

    int *p = llist_remf(&list);
    printf("%d\n", *p);
    printf("size: %d\n", list.size);

    p = llist_remf(&list);
    printf("%d\n", *p);    
    printf("size: %d\n", list.size);

    p = llist_remf(&list);
    printf("%d\n", *p);
    printf("size: %d\n", list.size);
}

int main()
{
    test();
}
