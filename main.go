package main

func (redis *RedisPokemonCacheServiceImpl) addNode(pokemon *Pokemon) (bool, *Node) {
	node := new(Node)
	node.data = *pokemon
	if redis.end != nil {
		redis.end.next = node
		node.previous = redis.end
		redis.end = node
		redis.size++
		return true, node
	}
	redis.end = node
	redis.start = node
	redis.size++
	return true, node
}

func (redis *RedisPokemonCacheServiceImpl) deleteNode(node *Node) {
	if redis.start == node {
		redis.start = redis.start.next
		redis.start.previous = nil
	} else if redis.end == node {
		redis.end = redis.end.previous
		redis.end.next = nil
	} else {
		node.previous.next = node.next
		node.next.previous = node.previous
	}
	redis.size--
}

type Node struct {
	data     Pokemon
	next     *Node
	previous *Node
}

type Pokemon struct {
	name string
	//id        UUID
	Type      string
	height    float64
	weight    float64
	abilities []string
}

type PokemonCacheService interface {
	set(name string, pokemon *Pokemon)
	get(name string) *Pokemon
	delete(name string) bool
}

type RedisPokemonCacheServiceImpl struct {
	redisMap  map[string]*Node
	threshold int
	size      int
	start     *Node
	end       *Node
}

func (redis *RedisPokemonCacheServiceImpl) set(name string, pokemon *Pokemon) {

	if data, found := redis.redisMap[name]; found {
		redis.deleteNode(data)
		status, node := redis.addNode(pokemon)
		if status {
			redis.redisMap[name] = node
		}
	} else {
		if redis.size >= redis.threshold {
			delete(redis.redisMap, redis.start.data.name)
			redis.deleteNode(redis.start)

		}
		result, node := redis.addNode(pokemon)
		if result {
			redis.redisMap[name] = node
		}

	}

}

func (redis *RedisPokemonCacheServiceImpl) get(name string) *Pokemon {
	if node, found := redis.redisMap[name]; found {
		status, newNode := redis.addNode(&node.data)
		if status {
			redis.deleteNode(node)
		}
		redis.redisMap[name] = newNode
		return &newNode.data
	} else {
		return &Pokemon{}
	}
}

func (redis *RedisPokemonCacheServiceImpl) delete(name string) bool {
	if node, found := redis.redisMap[name]; found {
		redis.deleteNode(node)
		delete(redis.redisMap, name)
		return true
	}
	return false
}

type PokemonCache struct {
	cacheService PokemonCacheService
}

func main() {
	cache := new(PokemonCache)
	cache.cacheService = &RedisPokemonCacheServiceImpl{
		redisMap:  make(map[string]*Node),
		threshold: 5,
		size:      0,
	}
	pokemon1 := &Pokemon{name: "Sourav", Type: "A", abilities: []string{"A", "B", "C"}}
	pokemon2 := &Pokemon{name: "Soma", Type: "B", abilities: []string{"A", "B", "C"}}
	pokemon3 := &Pokemon{name: "Dhriti", Type: "C", abilities: []string{"A", "B", "C"}}
	pokemon4 := &Pokemon{name: "Subham", Type: "D", abilities: []string{"A", "B", "C"}}
	pokemon5 := &Pokemon{name: "Sakaldeep", Type: "E", abilities: []string{"A", "B", "C"}}
	pokemon6 := &Pokemon{name: "Sita", Type: "F", abilities: []string{"A", "B", "C"}}

	cache.cacheService.set(pokemon1.name, pokemon1)
	cache.cacheService.set(pokemon2.name, pokemon2)
	cache.cacheService.set(pokemon3.name, pokemon3)
	cache.cacheService.set(pokemon4.name, pokemon4)
	cache.cacheService.set(pokemon1.name, pokemon4)
	cache.cacheService.delete(pokemon2.name)
	cache.cacheService.set(pokemon5.name, pokemon5)
	cache.cacheService.set(pokemon6.name, pokemon6)
}
