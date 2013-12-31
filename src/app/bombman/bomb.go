package bombman

type bombInfo struct {
	seqid int
	MapPos
	who   int
	fired bool
}
