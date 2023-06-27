#pragma once

// 定义结构体
typedef struct _point
{
    double x;
    double y;
} point;

// 创建 `point` 结构体变量
point create_point(double x, double y);

// 计算两个 `point` 变量表示的点之间的距离
double distance(const point *p1, const point *p2);
