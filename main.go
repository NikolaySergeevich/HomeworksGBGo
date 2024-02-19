package main

import (
	"fmt"
	"time"
	"unicode/utf8"
)

type Cache interface {
	Get(k string) (string, bool)
	Set(k, v string)
	Del(k string)
	DelAll()
	AutoDel(k string, delay time.Duration)
}

var _ Cache = (*cacheImpl)(nil)

// Доработает конструктор и методы кеша, так чтобы они соответствовали интерфейсу Cache
func newCacheImpl() *cacheImpl {
	cah := cacheImpl{make(map[string]string)}
	return &cah
}

type cacheImpl struct {
	cache map[string]string
}

func (c *cacheImpl) Get(k string) (string, bool) { // получить из кэша
	if val, ok := c.cache[k]; ok {
		return val, true
	}
	return "", false
}

func (c *cacheImpl) Set(k, v string) {
	if _, ok := c.cache[k]; !ok {
		c.cache[k] = v
	}
}
//удалаяет запись из кэша по ключу
func (c *cacheImpl) Del(k string) {
	delete(c.cache, k)
}
//чистит весь кэш
func(c *cacheImpl) DelAll(){
	clear(c.cache)
}
//удаляет запись из кэша через определённый промежуток времени с момента добавления
func(c *cacheImpl) AutoDel(k string, delay time.Duration){
	go func ()  {
		time.Sleep(delay)
		delete(c.cache, k)
	}()

}


func newDbImpl(cache Cache) *dbImpl {
	dbImpl := &dbImpl{cache: cache, dbs: map[string]string{"hello": "world", "test": "test"}}
	for key, val := range dbImpl.dbs {
		dbImpl.cache.Set(key, val)
		dbImpl.cache.AutoDel(key, 20*time.Second)
	}
	return dbImpl
}

type dbImpl struct {
	cache Cache
	dbs   map[string]string
}

func (d *dbImpl) Get(k string) (string, bool) {
	v, ok := d.cache.Get(k)
	if ok {
		return fmt.Sprintf("answer from cache: key: %s, val: %s", k, v), ok
	}

	if v, ok := d.dbs[k]; ok {
		d.cache.Set(k, v)//если дошли до этой записи,то в кэше нет так записи и добавляем её
		d.cache.AutoDel(k, 20*time.Second)
		return fmt.Sprintf("answer from dbs: key: %s, val: %s", k, v), ok
	}
	return "Нет такой записи.", false
}

//Добавление записи в БД
func (d *dbImpl) Set(k, v string) {
	d.dbs[k] = v
	d.cache.Set(k, v)
	d.cache.AutoDel(k, 20*time.Second)
}

//удаление записи из БД(и из кэша если есть)
func (d *dbImpl) Del(k string) (er error){
	if _, ok := d.cache.Get(k); ok {
		d.cache.Del(k)
	}
	if _, ok := d.dbs[k]; !ok {
		return &ErrorAbsentKey{}
	}
	delete(d.dbs, k)
	return nil
}


func main() {
	c := newCacheImpl()
	db := newDbImpl(c)
	fmt.Println(help())
	for {
		fmt.Println("Введите правильную команду:")
		if c, er := comand(); er != nil {
			fmt.Println(er, " Попробуй ещё")
			continue
		} else {
			switch c {
			case 'x':
				return
			case 's':
				fmt.Println("Введите информацию")
				if k, v, err := ScanWord(db); err != nil {
					fmt.Println(err, " Попробуй ещё")
					continue
				} else {
					db.Set(k, v)
				}
			case 'd':
				fmt.Println("Введите ключ записи, которую нужно удалить:")
				var key string
				fmt.Scan(&key)
				if er := db.Del(key); er != nil{
					fmt.Println(er, " Попробуй ещё.")
					continue
				}
				fmt.Println("Запись была успешно удалена")
			case 'c':
				db.cache.DelAll()
				fmt.Println("Кэш был успешно очистен")
			case 'g':
				fmt.Println("Введите ключ записи, которую нужно получить:")
				var k string
				fmt.Scan(&k)
				fmt.Println(db.Get(k))
			}
		}

	}

}
//для воода инфы(ключ, значение)
func ScanWord(db *dbImpl) (_, _ string, err error) {
	var k, v string
	fmt.Scan(&k, &v)
	if _, ok := db.dbs[k]; ok {
		err = &ErrorExistsKey{k}
		return
	}
	return k, v, nil
}
//для ввода команды
func comand() (_ rune, er error) {
	var c string
	fmt.Scan(&c)
	if utf8.RuneCountInString(c) > 1 {
		er = &ErrorExpOneLet{}
		return
	}
	return rune(c[0]), nil
}

func help() string {
	return "Эта программа иммитиркет какоё-то поведение БД и Кэша.\n" +
		"- Всё, что вы добавляете в БД, попадает в кэш и через определённое время пропадает из него.\n" +
		"- Вы можете добавлять данные в БД и получать их. При этом, если информация была в кэше, то она\n" +
		"вытащится из него. А если она к моменту запроса уже была удалена из Кэша, то информация вытащится из БД\n" +
		"и снова попадёт в Кэш."+
		"s - добавить запись в БД\n" +
		"g - получить данные по ключу\n"+
		"d - удалить запись из БД\n" +
		"с - очистить весь Кэш\n" +
		"x - закрыть программу."
}
