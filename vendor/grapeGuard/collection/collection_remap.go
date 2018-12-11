package collection

import (
	"sync"
	"sync/atomic"
	"time"
	remap "grapeGuard/remaps"
	proto "grapeLoggers/protos"

	"context"

	log "grapeLoggers/clientv1"
)

var (
	remapHot sync.Map

	comitTicker = time.NewTicker(time.Minute) // 提交周期
)

type RemapCollection struct {
	// 哪个HOST
	RemapFrom string

	// 流量计算
	RecvBytes int64
	SendBytes int64

	// 请求数
	RqTotal int64
}

func (c *RemapCollection) Reset() {
	atomic.StoreInt64(&c.RqTotal,0)
	atomic.StoreInt64(&c.RecvBytes,0)
	atomic.StoreInt64(&c.SendBytes,0)
}

func CollectRemapRequst(item *HostQueueItem) {
	remapItem := remap.SearchRemapNonReq(item.Host,item.Path,item.IsSSL)
	if remapItem == nil {
		return // 不处理
	}

	val,ok := remapHot.Load(remapItem.FromSrc)
	if !ok {
		newData := &RemapCollection{
			RemapFrom:remapItem.FromSrc,
			RecvBytes:item.ReqBodySize,
			SendBytes:0,
			RqTotal:1,
		}

		remapHot.Store(remapItem.FromSrc,newData)
		return
	}

	Data := val.(*RemapCollection)

	atomic.AddInt64(&Data.RecvBytes,item.ReqBodySize)
	atomic.AddInt64(&Data.RqTotal,1)
}

func CollectRemapResp(item *HostQueueItem) {
	remapItem := remap.SearchRemapNonReq(item.Host,item.Path,item.IsSSL)
	if remapItem == nil {
		return // 不处理
	}

	val,ok := remapHot.Load(remapItem.FromSrc)
	if !ok {
		newData := &RemapCollection{
			RemapFrom:remapItem.FromSrc,
			RecvBytes:0,
			SendBytes:item.RespBodySize,
			RqTotal:0,
		}

		remapHot.Store(remapItem.FromSrc,newData)
		return
	}

	Data := val.(*RemapCollection)
	atomic.AddInt64(&Data.SendBytes,item.RespBodySize)
}

func prcoAutoRemapCommit() {
	for {
		select {
		case <- comitTicker.C:

			newReq := &proto.RemapCommitReq{}
			timesp := time.Now().Unix()
			remapHot.Range(func(key, value interface{}) bool {
				item := value.(*RemapCollection)

				newReq.Data = append(newReq.Data,&proto.RemapItem{
					FromSrc:item.RemapFrom,
					SendBytes:item.SendBytes,
					RecvBytes:item.RecvBytes,
					Rqtotal:int32(item.RqTotal),
					Time:timesp,
				})

				item.Reset()
				return true
			})

			_,err := loggerCli.RemapCollCommit(context.Background(),newReq)
			if err != nil {
				log.WithField("type", "collect").Error("提交采集数据异常:", err)
			}
		}
	}
}