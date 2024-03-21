package bolt

type mapFunc[E any, L any] func(E) L

func Map[E any, L any](s []E, f mapFunc[E,L]) []L {
    result := make([]L, len(s))
    for i := range s {
        result[i] = f(s[i])
    }
    return result
}
