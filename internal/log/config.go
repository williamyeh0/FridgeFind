package log

type Config struct {
	/*
		configure max size
	*/
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
