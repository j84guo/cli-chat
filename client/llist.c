#include <stdio.h>
#include <stdlib.h>
#include "llist.h"

int llist_init(llist_t *list)
{
    if (!list)
        return -1;

    list->head = NULL;
    list->tail = NULL;
    list->size = 0;

    return 0;
}

int llist_dtry(llist_t *list)
{
    if (!list)
        return -1;

    lnode_t *node = list->head, *next;
    while (node) {
        next = node->next;
        free(node);
        node = next;
    }

    return 0;
}

int llist_addf(llist_t *list, void *data)
{
    if (!list || !data)
        return -1;

    lnode_t *node = malloc(sizeof(lnode_t));
    if (!node)
        return -1;

    node->data = data;
    node->next = list->head;
    node->prev = NULL;

    list->head = node;
    if (++list->size == 1)
        list->tail = node;

    return 0;
}

int llist_addl(llist_t *list, void *data)
{
    if (!list || !data)
        return -1;

    lnode_t *node = malloc(sizeof(lnode_t));
    if(!node)
        return -1;

    node->data = data;
    node->next = NULL;
    node->prev = list->tail;

    list->tail = node;
    if (++list->size == 1)
        list->head = node;

    return 0;
}

void *llist_remf(llist_t *list)
{
    if (!list || !list->head)
        return NULL;

    lnode_t *node = list->head;
    void *data = node->data;

    list->head = node->next;
    free(node);
    if (--list->size == 0)
        list->tail = NULL;

    return data;
}

void *llist_reml(llist_t *list)
{
    if (!list || !list->tail)
        return NULL;

    lnode_t *node = list->tail;
    void *data = node->data;

    list->tail = node->prev;
    free(node);
    if (--list->size == 0)
        list->head = NULL;

    return data;
}
