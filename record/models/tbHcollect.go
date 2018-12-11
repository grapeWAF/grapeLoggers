package model

type TbHcollect struct {
	Idx           int64   `xorm:"not null pk autoincr BIGINT(20)"`
	System        string
	Byterecv      int
	Bytesend      int
	Procmemper    float32
	Cputotal      float32
	Proccpuper    float32
	Hdfree        int64
	Hdtotal       int64
	Cpuper        float32
	Hostuuid      string
	Hostname      string
	Sysmemfree    int64
	Timestamp     int64
	Remote        string
	Sysmemper     float32
	Sysmemused    int64
	Cputype       string
	Hdper         float32
	Totalrq       int
	Totalbytesend int64
	Totalbyterecv int64
}
