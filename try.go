package try

type Catcher interface {
	Catch(fn func(e any))
	throw(e any)
	setDone()
}

type _catcher struct {
	IsPanic    bool
	PanicError any

	done  bool
	catch func(e any)
}

func (c *_catcher) setDone() {
	c.done = true
	if c.catch != nil && c.IsPanic {
		c.catch(c.PanicError)
	}
}

func (c *_catcher) throw(e any) {
	c.IsPanic = true
	c.PanicError = e
	c.setDone()
}

func Run(fn func()) (c Catcher) {
	c = &_catcher{}
	defer func() {
		if r := recover(); r != nil {
			c.throw(r)
		}
	}()

	fn()
	c.setDone()
	return
}

func (c *_catcher) Catch(fn func(e any)) {
	c.catch = fn
	if c.done {
		c.setDone()
	}
}

func RunInNewThread(fn func()) (c Catcher) {

	c = &_catcher{}

	go func(cc Catcher) {
		defer func() {
			if r := recover(); r != nil {
				cc.throw(r)
			}
		}()
		fn()
		cc.setDone()
	}(c)

	return c
}
