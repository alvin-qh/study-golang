#include "point.h"
#include <math.h>

point create_point(double x, double y)
{
    point p = {.x = x, .y = y};
    return p;
}

double distance(const point *p1, const point *p2)
{
    return sqrt(pow(p2->x - p1->x, 2) + pow(p2->y - p1->y, 2));
}
