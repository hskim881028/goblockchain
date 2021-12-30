package blockchain

type fakeDB struct {
	fakeGetBlock func() []byte
	fakeGetChain func() []byte
}

func (f fakeDB) GetBlock(hash string) []byte {
	return f.fakeGetBlock()
}
func (f fakeDB) GetChain() []byte {
	return f.fakeGetChain()
}
func (fakeDB) PutBlock(hash string, data []byte) {}
func (fakeDB) DeleteAllBlocks()                  {}
func (fakeDB) PutChain(data []byte)              {}
