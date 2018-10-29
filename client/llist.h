typedef struct lnode_t lnode_t;
typedef struct llist_t llist_t;

struct lnode_t {
    void *data;
    lnode_t *next;
    lnode_t *prev;
};

struct llist_t {
    lnode_t *head;
    lnode_t *tail;
    int size;
};

int llist_init(llist_t *list);

int llist_dtry(llist_t *list);

int llist_addf(llist_t *list, void *data);

int llist_addl(llist_t *list, void *data);

void *llist_remf(llist_t *list);

void *llist_reml(llist_t *list);
