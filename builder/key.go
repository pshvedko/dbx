package builder

import (
	"github.com/pshvedko/dbx/filter"
)

type Key struct {
	Builder
}

func (k *Key) AppendValue(t filter.Type, v any) (int, error) {
	var n1, n2 int
	var err error
	switch t {
	case filter.TRUE, filter.FALSE:
	default:
		n1, err = k.WriteString(t.String())
		if err != nil {
			return n1, err
		}
	}
	z := k.Size()
	f := k.Value(v)
	switch x := f.(type) {
	case Keyword:
		return n1 + 1, k.WriteByte(x[0])
	default:
		if x != nil && z == k.Size() {
			n2, err = k.AppendFormat(x)
			return n1 + n2, err
		}
		return n1, nil
	}
}

func (k *Key) AppendColumn(t string, c string, m map[string]int) (int, error) {
	return k.AppendInt(m[c])
}

func (k *Key) AppendVia(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('&')
	}
	return 1, k.WriteByte('|')
}

func (k *Key) AppendParenthesis(b bool) (int, error) {
	if b {
		return 1, k.WriteByte('[')
	}
	return 1, k.WriteByte(']')
}

func (k *Key) Straight(j filter.Projector, u bool, x filter.Filter) (err error) {
	return Straight(k, j, u, x)
}

func (k *Key) Conjunct(j filter.Projector, u bool, ff []filter.Filter) (err error) {
	_, err = k.AppendParenthesis(true)
	if err != nil {
		return
	}
	for i, f := range ff {
		if i > 0 {
			_, err = k.AppendVia(u)
			if err != nil {
				return
			}
		}
		err = f.To(k, j)
		if err != nil {
			return
		}
	}
	_, err = k.AppendParenthesis(false)
	return
}

type Origin struct {
	File string
	Line int
}

type Hash struct {
	Origin
	static string
}

type Keeper interface {
	Get(Hash) (string, bool)
	Put(Hash, string)
}

type Cache struct {
	Origin
	Keeper
}

func (c Cache) IsCacheEnabled() bool {
	return c.Keeper != nil && c.Origin.File != "" && c.Origin.Line != 0
}

func (k *Key) WriteRange(i, z int) (n int, err error) {
	b := byte('-')
	switch i - z {
	case 0:
	case 1:
		b = ','
		fallthrough
	default:
		err = k.WriteByte(b)
		if err != nil {
			return
		}
		n, err = k.AppendInt(i - 1)
		n++
	}
	return
}

func (k *Key) WriteIndices(j filter.Fielder, f Fielder) (err error) {
	y := true
	z := 0
	for i, n := range j.Names() {
		if !f.Returned(n) {
			if !y {
				_, err = k.WriteRange(i, z)
				if err != nil {
					return
				}
				y = true
			}
			continue
		}
		if y {
			if z > 0 {
				err = k.WriteByte(',')
				if err != nil {
					return
				}
			}
			_, err = k.AppendInt(i)
			if err != nil {
				return
			}
			z = i + 1
			y = false
		}
	}
	if !y {
		_, err = k.WriteRange(j.Len(), z)
		if err != nil {
			return
		}
	}
	return
}

//
//func WritePresentRanges(idx []int) {
//	y := true // разрешено печатать индексы
//	z := 0    // индекс начала интервала
//	for i, x := range idx {
//		if x == 0 { // поле не в списке, пропускаем
//			if !y { // было запрещено печатать,
//				y = true // разрешаем печатать
//				switch i - z {
//				case 0: // такое невозможно!
//					panic(i)
//				case 1: // напечатан только один, добавить нечего
//				case 2: // напечатано два, между соседними "-" не влезет
//					println(",")
//					println(i - 1) // конец интервала
//				default: // больше двух, ставим "-" закрывая интервал
//					println("-")
//					println(i - 1) // конец интервала
//				}
//				println(",") // закрываем интервал
//			}
//			// тут было разрешено печатать, но нечего
//			continue
//		}
//
//		if y { // можно печатать индексы
//			println(i) // печатаем начало интервала
//			z = i      // запоминаем начало интервала
//			y = false  // запретим печатать индексы
//		}
//		// тут нельзя печатать индексы
//	}
//
//	if !y { // вышли из цикла и запрещено печатать
//		println("-")
//	}
//
//	println("===", fmt.Sprint(idx))
//}
//
//func (c *Constructor) WriteFSM(idx []int) {
//	y := true // разрешено печатать начало (мы "вне" интервала)
//	z := 0    // индекс начала интервала
//
//	for i, x := range idx {
//		if x == 0 { // поле не активно
//			if !y { // мы были внутри интервала, надо его закрыть
//				c.closeRange(z, i-1)
//				y = true
//			}
//			continue
//		}
//
//		// Поле активно
//		if y {
//			if z > 0 || !c.first { // если не самое первое поле вообще
//				c.WriteByte(',')
//			}
//			c.writeInteger(i) // печатаем начало
//			z = i
//			y = false // "заходим" в интервал, запрещаем печатать новые начала
//		}
//	}
//
//	if !y { // закрываем хвост, если цикл кончился на активном поле
//		c.closeRange(z, len(idx)-1)
//	}
//}
//
//func (c *Constructor) closeRange(start, end int) {
//	if end > start {
//		if end-start == 1 {
//			c.WriteByte(',') // или '-', если хочешь 0-1 вместо 0,1
//		} else {
//			c.WriteByte('-')
//		}
//		c.writeInteger(end)
//	}
//}
//
// ...
//if !y { // Мы вышли из цикла, находясь "внутри" интервала
//    c.WriteByte('-')
//    // Мы не пишем i, просто оставляем дефис как маркер "до конца"
//}
// ... внутри цикла ...
//if !y {
//    // Мы закончили цикл на активном поле.
//    // Закрываем интервал последним реальным индексом (i-1)
//    if (i-1) > z {
//        if (i-1) - z == 1 {
//            c.WriteByte(',') // или '-' по вкусу
//        } else {
//            c.WriteByte('-')
//        }
//        c.writeInteger(i-1)
//    }
//}
