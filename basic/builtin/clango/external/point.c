#include "point.h"
#include <math.h>

// 创建 `point` 结构体变量
point create_point(double x, double y)
{
    point p = {.x = x, .y = y};
    return p;
}

// 计算两个 `point` 变量表示的点之间的距离
double distance(const point *p1, const point *p2)
{
    return sqrt(pow(p2->x - p1->x, 2) + pow(p2->y - p1->y, 2));
}
