package storageserver

import (
	"fmt"
	"sort"
	"sync"
)

type Pool struct {
	servers        []*Server
	fileserversmap map[string][]*Server
	filespace      map[string]int
	len            int
	mu             *sync.RWMutex
}

// Append appends a server to the pool.
//
// The function takes a pointer to a Pool struct and a pointer to a Server struct as input parameters.
// It adds the server to the list of servers in the pool and increments the length of the pool.
// The function returns a slice of pointers to Server structs, which is the updated list of servers in the pool.
func (p *Pool) Append(s *Server) []*Server {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.servers = append(p.servers, s)
	p.len++
	return p.servers
}

// LeastCompleted returns a slice of *Server that represents the servers in the Pool
// sorted in ascending order based on the space reserved.
//
// No parameters.
// Returns a slice of *Server.
func (p *Pool) LeastCompleted() []*Server {
	p.mu.Lock()
	defer p.mu.Unlock()

	sort.Slice(p.servers, func(i, j int) bool {
		return p.servers[i].SpaceReserved() < p.servers[j].SpaceReserved()
	})

	return p.servers
}

// Len returns the length of the Pool.
//
// No parameters.
// Returns an integer.
func (p *Pool) Len() int {
	return p.len
}

func (p *Pool) SetFileServers(userId int, filename string, ss []*Server) {
	p.fileserversmap[fmt.Sprintf("%d_%s", userId, filename)] = ss
}

func (p *Pool) GetFileServers(filename string) []*Server {
	return p.fileserversmap[filename]
}

// NewPool creates a new Pool instance.
//
// It takes a slice of Server pointers as a parameter and returns a Pool instance.
func NewPool(ss []*Server) Pool {
	return Pool{
		servers: ss,
		len:     len(ss),
		mu:      &sync.RWMutex{},
	}
}
