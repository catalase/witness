import sympy
from sympy.core.cache import cacheit

sympy.var("n")

# 1 ?댁긽 u ????섏뿬 n(n+1) ????긽 ?몄옄濡?媛吏먯씠 利앸챸?섏뿀??
def powseries(u):
    if u == 0:
        return n

    return n * (n + 1) * _powseries(u)

@cacheit
def _powseries(u):
    if u == 1:
        return sympy.Rational(1, 2)

    if u == 2:
        return (2 * n + 1) / 6

    z = ((u + 2) * (n + 1) + 1) * _powseries(u - 1) - (1 + n) * _powseries(u - 2)
    z /= u + 3

    return sympy.factor(z)

##def main():
##    import io
##    import clipboard # pyperclip
##
##    def ternary(lst):
##        for i in range(0, len(lst), 3):
##            yield lst[i:i + 3]
##
##    sep = '\\left('
##    w = io.StringIO()
##
##    for i in range(0, 5):
##        print('\\begin{align}', file=w)
##
##        eq = sympy.latex(powseries(i))
##        x, _, y = eq.rpartition(sep)
##
##        print('\\sum_{k=1}^{k^%d} = ' % i, x, sep, '& ', file=w)
##
##        it = ternary(y.split('+'))
##        print('+'.join(next(it)), file=w)
##
##        if i == 0:
##            print(')', file=w)
##
##        for part in it:
##            print('&', '+', '+'.join(part), '\\\\', file=w)
##
##        print(r'\\\\', file=w)
##        print('\\end{align}', file=w)
##
##    clipboard.copy(w.getvalue())
##
##if __name__ == '__main__':
##    main()