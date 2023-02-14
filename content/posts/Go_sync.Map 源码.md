---
# 主页简介
# summary: ""
# 文章副标题
# subtitle: ""
# 作者信息
# author: ""
# authorLink: ""
# authorEmail: ""
# description: ""
# keywords: ""
# license: ""
# images: []
# 文章的特色图片
# featuredImage: ""
# 用在主页预览的文章特色图片
# featuredImagePreview: ""
# password:加密页面内容的密码，详见 主题文档 - 内容加密
# message:  加密提示信息，详见 主题文档 - 内容加密
# linkToMarkdown: true
# 上面一般不用动
title: "sync.Map 源码"
date: 2023-02-14T00:16:04+08:00
lastmod: 2022-02-14T00:16:04+08:00
categories: ["Go"]
---

sync.Map 源码分析，参考 https://mp.weixin.qq.com/s/HFUyiS7OH0jIPDg4rEmcLw，自己做了一些总结。  

## map 非线程安全

Go 的数据结构 map 不是线程安全的，解决方案：  

1. 自己配一把锁（sync.Mutex），或者更加考究一点配一把读写锁（sync.RWMutex）。这种方案简约直接，但是缺点也明显，就是性能不会太高。

2. 使用Go语言在2017年发布的Go 1.9中正式加入了并发安全的字典类型sync.Map。

## java ConcurrentHashMap 实现思想

JDK1.7的实现：由一个 Segment 数组和多个 HashEntry 组成。基于分桶上锁（**分段锁**）的思想，每个 Segment 数组的元素都是一个 HashMap，**对每个 HashMap 单独上锁**。

{{< image src="/images/sync.Map 源码/map_jdk1.7.webp" width=80% height=80% caption="jdk1.7" >}}

JDK1.8的实现：**Hash 数组**+链表+红黑树。**直接对 hash 表的每个元素单独上锁**，用 CAS 实现。（本质也是**分段锁**的思想，只不过比 1.7 **锁的粒度会小很多**。）

{{< image src="/images/sync.Map 源码/map_jdk1.8.webp" width=70% height=70% caption="jdk1.8" >}}

## syn.map 简介

1. **适合场景**：  
   **适合读多，写少的场景**。因为写是必须加锁的。

2. **实现简介**
   1. 由两个 map 实现，read map 和 dirty map。value 是 *entry，entry 是对 pointer 的封装，对 entry.point 都是原子操作和 CAS；  
   2. read 只有**一部分数据**，而 dirty map 一定具有全部的数据；  
   3. read map **直接读取 entry**，但是 load entry.Pointer 用原子操作 `atomic.LoadPointer(&e.p)`，修改数据（store、delete）使用的是 **CAS 乐观锁**无限尝试直接替换 entry 里的 pointer；  
   4. dirty map 的所有操作都需要加**互斥锁**，无论读取还是存储。获取 entry 加锁，同时存储 entry.pointre 再用**原子操作**；  
   5. entry.pointer 有三种状态：normal，nil（中间状态逻辑删除，dirty 还有），expunged（表示 read=>dirty 流程已经扫描过该key，且 dirty 中没有该 key。这是最终状态，如果下次 read=> dirty 还是处于这个状态，则不会存入 dirty。那么在 下一次 dirty => read 将因为 dirty = nil 的原因被释放）；
   6. **store 流程**：
      1. 看 read map 有没有，有（可以是nil）且不是 expunged 则直接 store（entry.Point）；有但是是 expunged 状态则加锁后存入 dirty 中(复用entry，但是dirty 中新加 key)；
      2. 看 dirty 有没有，有直接 store（entry.Point）；
      3. dirty 也没有，**判断是否需要发生 read => dirty**；
      4. **创建 entry**，加入 dirty；
   7. **load流程**：  
      1. 在 read map 中找，找到且不处于 expunged 则返回；处于 expunged 说明不存在，直接返回 nil；
      2. 在 dirty 中找，miss++，如果 miss 达到阈值，**dirty => read**；
      3. dirty 中找到且不处于 nil 和 expunged 则返回，否则返回 nil。
   8. **delete 流程**：
      1. 在 read 中，设置 entry.point = nil，逻辑删除。
      2. 不在 read 中，直接删除 dirty 的 key（真正删除），且 miss++。
   9. **range 流程**：
      1. 如果 read 和 dirty 相同，即判断 read.amended == false，成立则遍历 read map；  
      2. 如果为 true，则说明 dirty 有独占的 key，发生 dirty => read，再遍历 read map；  
      也就是说**每次遍历一定是在 read map 中，这样可以不用加锁**。  
   10. **dirty=>read** 操作：直接 sync.map.m = dirty，dirty = nil；
   11. **read=>dirty** 操作：把 read 中非 nil 非 expunged 的数据存入 dirty 中，nil 的数据设置为 expunged；

3. **既然nil也表示标记删除，那么再设计出一个expunged的意义是什么？**  
   expunged是有存在意义的，它作为删除的最终状态（待释放），这样nil就可以作为一种中间状态。如果仅仅使用nil，那么，在read=>dirty重塑的时候，可能会出现如下的情况：  
   1. 如果nil在read浅拷贝至dirty（read=>dirty）的时候仍然保留entry的指针（即拷贝完成后，对应键值下read和dirty中都有对应键下entry e的指针，且e.p=nil）那么之后在dirty=>read升级key的时候对应entry的指针仍然会保留。那么最终**合集会越来越大，存在大量nil的状态**，永远无法得到清理的机会。   
   2. 如果nil在read浅拷贝时不进入dirty(read=>dirty)，那么之后store某个Key键的时候，可能会出现read和dirty不同步的情况，即此时read中包含dirty不包含的键，那么之后用dirty替换read的时候就会出现数据丢失的问题。  
   3. 如果nil在read浅拷贝时直接把read中对应键删除（从而避免了 2 中不同步的问题），但这又必须对read加锁，违背了read读写不加锁的初衷。

## syn.map 源码

可以直接 copy 到 idea 中，方便追踪和查看。

```go
package map

type Map struct {
	mu sync.Mutex

	// read contains .... 省略原版的注释
	// read map是被atomic包托管的，这意味着它本身Load是并发安全的（但是它的Store操作需要锁mu的保护）
	// read map中的entries可以安全地并发更新，但是对于expunged entry，在更新前需要经它unexpunge化并存入dirty
	//（这句话，在Store方法的第一种特殊情况中，使用e.unexpungeLocked处有所体现）
	read atomic.Value // readOnly

	// dirty contains .... 省略原版的注释
	// 关于dirty map必须要在锁mu的保护下，进行操作。它仅仅存储 non-expunged entries
	// 如果一个 expunged entries需要存入dirty，需要先进行unexpunged化处理
	// 如果dirty map是nil的，则对dirty map的写入之前，需要先根据read map对dirty map进行浅拷贝初始化
	dirty map[interface{}]*entry

	// misses counts .... 省略原版的注释
	// 每当读取的是时候，read中不存在，需要去dirty查看，miss自增，到一定程度会触发dirty=>read升级转储
	// 升级完毕之后，dirty置空 &miss清零 &read.amended置false
	misses int
}

// 这是一个被原子包atomic.Value托管了的结构，内部仍然是一个map[interface{}]*entry
// 以及一个amended标记位，如果为真，则说明dirty中存在新增的key，还没升级转储，不存在于read中
type readOnly struct {
	m       map[interface{}]*entry
	amended bool // true if the dirty map contains some key not in m.
}

// expunged is an arbitrary pointer that marks entries which have been deleted
// from the dirty map.
// 表示 该entry 已经被扫描过了，当然可能处于 expunged 过程中（即read=>dirty）或者 expunged 已经全部完成
// 为什么存在这一个状态，主要是因为 expunged 过程（即read=>dirty）是在 dirty 加锁，但是 read 却可以直接 store
// 当要 store 时 key 在 read 找得到，
// 1. 如果这个 key 还没扫描，不管怎么样，把 read.key 直接 store 没问题，后面扫描到会 copy 到 dirty
// 2. 但是如果这个 key 已经扫描完且为nil（按道理应该是expunged，但是我们只简单设为nil），直接 read.key 进行 store，是不正确的。要加入 dirty 中
var expunged = unsafe.Pointer(new(any))

// An entry is a slot in the map corresponding to a particular key.
// 这是一个容器，可以存储任意的东西，因为成员p是unsafe.Pointer(*interface{})
// sync.Map中的值都不是直接存入map的，都是在entry的包裹下存入的
type entry struct {
	// p points ....  省略原版的注释
	// entry的p可能的状态：
	// e.p == nil：entry已经被标记删除，不过此时还未经过read=>dirty重塑，此时可能仍然属于dirty（如果dirty非nil）
	// e.p == expunged：entry已经被标记删除，经过read=>dirty重塑，不属于dirty，仅仅属于read，下一次dirty=>read升级，会被彻底清理
	// e.p == 普通指针：此时entry是一个不同的存在状态，属于read，如果dirty非nil，也属于dirty
	p unsafe.Pointer // *interface{}
}

func newEntry(i any) *entry {
	return &entry{p: unsafe.Pointer(&i)}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool) {
	// Avoid locking if it's a clean hit.
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

// tryLoadOrStore atomically loads or stores a value if the entry is not
// expunged.
//
// If the entry is expunged, tryLoadOrStore leaves the entry unchanged and
// returns with ok==false.
func (e *entry) tryLoadOrStore(i any) (actual any, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expunged {
		return nil, false, false
	}
	if p != nil {
		return *(*any)(p), true, true
	}

	// Copy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expunged {
			return nil, false, false
		}
		if p != nil {
			return *(*any)(p), true, true
		}
	}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map) Range(f func(key, value any) bool) {
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.m satisfies that property without
	// requiring us to hold m.mu for a long time.
	read, _ := m.read.Load().(readOnly)
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		if read.amended {
			read = readOnly{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

// Store sets the value for a key.
func (m *Map) Store(key, value interface{}) {
	// 首先把readonly字段原子地取出来
	// 如果key在readonly里面，则先取出key对应的entry，然后尝试对这个entry存入value的指针
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	// 如果readonly里面（1）不存在key或者是（2）对应的key是被擦除掉了的（expunged），则继续。
	m.mu.Lock() // 上锁

	// 锁的惯用模式：再次检查readonly，防止在上锁前的时间缝隙出现存储。（上锁后就不可能有出现这种问题了，因为这个操作是加锁的）
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		// 这里有两种情况：
		// 1. 上面的时间缝隙里面，出现了key的存储过程（可能是normal值，也可能是expunge值）
		//    （1）此时先校验e.p，如果是 normal，说明read和dirty里都有相同的entry（因为read中有，dirty一定有。），则直接设置 entry.pointer
		//    （2）如果是expunge值，则说明dirty里面已经不存在key了，需要先在dirty里面种上key，然后设置entry
		// 2. read 中的key 的 entry是expunge的状态，说明read=>dirty 已经完成或正在发生，但该key一定已经被扫描，直接加入dirty 中（注意dirty中一定没有该key,复用entry，但是dirty 中新加 key）
		if e.unexpungeLocked() {
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty map and this entry is not in it.
			m.dirty[key] = e
		}
		e.storeLocked(&value) // 将 Pointer存入容器 entry 中
	} else if e, ok := m.dirty[key]; ok {
		// readonly里面不存在，则查看dirty里面是否存在
		// 如果dirty里面存在，则直接设置dirty的对应key
		e.storeLocked(&value)
	} else {
		// dirty里面也不存在（或者dirty为nil），则应该先设置在ditry里面
		// 此时要检查read.amended，如果为假（则 dirty 刚刚变成 read，dirty 为 nil） or 两者均是初始化状态）
		// 此时要在dirty里面设置新的key，需要确保dirty是初始化的且需要设置amended为true（表示自此dirty多出了一些独有key）
		if !read.amended {
			// 利用read重塑dirty！即 dirty 变 read 后，dirty 变空，然后再把 read normal 的数据存入 dirty
			m.dirtyLocked()
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
	}

	// 解锁
	m.mu.Unlock()
}

// 这是一个自旋乐观锁：只有key是非expunged的情况下，会得到set操作
func (e *entry) tryStore(i *interface{}) bool {
	for {
		// 注意要用原子加载哦，才能使取到的是最新值，而不是缓存值
		p := atomic.LoadPointer(&e.p)
		// 如果p是expunged就不可以了set了
		// 因为expunged状态是read独有的，这种情况下说明 read=>dirty 全部完成或者正在 read=>dirty 过程中，但是该key的entry已经扫描过了。
		// 此时要新增只能在dirty中，不能在read中
		if p == expunged {
			return false
		}
		// 如果非expunged，则说明是normal的entry或者nil的entry，可以直接替换
		// 1. normal 的entry，即使正在read=>dirty也无所谓，直接替换，因为 read 和 dirty 的 e 是指向同一个
		// 2. 如果处于 nil，说明不处于 read=>dirty的过程中或处于 read=> dirty过程中，但是未扫描到;也可以直接赋值，反正后面扫描到会加入 dirty 中
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

// 利用了go的CAS，如果e.p是 expunged，则将e.p置为空，从而保证她是read和dirty共有的
func (e *entry) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}

// 真正的set操作，从这里也可以看出来2点：1是set是原子的 2是封装的过程
func (e *entry) storeLocked(i *interface{}) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

// 利用read重塑dirty！即 dirty 变 read 后，dirty 变空，然后再把 read normal 的数据存入 dirty
// 如果dirty为nil，则利用当前的read来初始化dirty（包括read本身也为空的情况）
// 此函数是在锁的保护下进行，所以不用担心出现不一致
func (m *Map) dirtyLocked() {
	if m.dirty != nil {
		return
	}
	// 经过这么一轮操作:
	// dirty里面存储了全部的非expunged的entry
	// read里面存储了dirty的全集，以及所有expunged的entry
	// 且read中不存在e.p == nil的entry（已经被转成了expunged）
	read, _ := m.read.Load().(readOnly)
	m.dirty = make(map[interface{}]*entry, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() { // 只有 normal 的key，能够重塑到dirty里面
			m.dirty[k] = e // 这里是浅拷贝，也就是说 read 和 dirty 的 entry 是同一个
		}
	}
}

// 利用乐观自旋锁，
// 如果e.p是nil，尽量将e.p置为expunged
// 返回最终e.p是否是expunged
func (e *entry) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	// 有可能是 normal
	return p == expunged
}

func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	// 把readonly字段原子地取出来
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]

	// 如果readonly没找到，且dirty包含了read没有的key，则尝试去dirty里面找
	if !ok && read.amended {
		m.mu.Lock()
		// 锁的惯用套路
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// Regardless of ... 省略英文
			// 记录miss次数，并在满足阈值后，触发dirty=>map的升级
			m.missLocked()
		}
		m.mu.Unlock()
	}

	// readonly和dirty的key列表，都没找到，返回nil
	if !ok {
		return nil, false
	}

	// 找到了对应entry，随即取出对应的值
	return e.load()
}

// 自增miss计数器
// 如果增加到一定程度，dirty会升级成为readonly（dirty自身清空 & read.amended置为false）
func (m *Map) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	// 直接用dirty覆盖到了read上（那也就是意味着dirty的值是必然是read的父集合，当然这不包括read中的expunged entry）
	m.read.Store(readOnly{m: m.dirty}) // 这里有一个隐含操作，read.amended再次变成false
	m.dirty = nil
	m.misses = 0
}

// entry是一个容器，从entry里面取出实际存储的值（以指针提取的方式）
func (e *entry) load() (value interface{}, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged {
		return nil, false
	}
	return *(*interface{})(p), true
}

// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	m.LoadAndDelete(key)
}

// LoadAndDelete 删除的逻辑和Load的逻辑，基本上是一致的
func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// 不在 read 中
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		// 不在 read 中
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 删除 dirty，如果有的话
			delete(m.dirty, key)
			// Regardless of ...省略
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		// 在 read 和 dirty 中，e 是 read 的，设置 entry.p 为 nil
		return e.delete()
	}
	return nil, false
}

// 如果e.p == expunged 或者nil，则返回false
// 否则，设置e.p = nil，返回删除的值的指针
func (e *entry) delete() (value interface{}, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return nil, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*interface{})(p), true
		}
	}
}
```
## End