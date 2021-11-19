#pragma once

typedef struct _point
{
    double x;
    double y;
} point;

point create_point(double x, double y);

double distance(const point *p1, const point *p2);
