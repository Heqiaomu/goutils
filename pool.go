package util

type ConnRes interface{}
type Factory func() (ConnRes, error)
type Destory func(ConnRes)

type Pool struct {
	conns   chan ConnRes
	factory Factory
	destory Destory
}

func NewPool(factory Factory, destroy Destory, cap int) *Pool {
	return &Pool{
		conns:   make(chan ConnRes, cap),
		factory: factory,
		destory: destroy,
	}
}

func (p *Pool) DestoryPool() {

	close(p.conns)

	for k := range p.conns {
		p.destory(k)
	}
}

func (p *Pool) new() (ConnRes, error) {
	return p.factory()
}

func (p *Pool) close(conn ConnRes) {
	p.destory(conn)
}

func (p *Pool) Get() (conn ConnRes) {
	select {
	case conn = <-p.conns:
		{
		}
	default:
		conn, _ = p.new()
	}
	return
}

func (p *Pool) Put(conn ConnRes) {
	select {
	case p.conns <- conn:
		{
		}
	default:
		p.close(conn)
		conn = nil
	}
}
