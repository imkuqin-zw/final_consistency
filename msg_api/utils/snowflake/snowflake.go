package snowflake

import (
	"go.uber.org/zap"
	"shop/basic/config"
	z "shop/plugins/zap"
	"sync"
	"time"
)

const (
	ts_mask         = 0x1FFFFFFFFFF // 41bit
	center_id_mask  = 0x1F          // 5bit
	machine_id_mask = 0x1F          // 5bit
	sn_mask         = 0xFFF         // 12bit
)

var (
	once sync.Once
	s    *Snowflake
	log  *z.Logger
)

type Cfg struct {
	DataCenterId uint64 `json:"data_center_id"`
}

type Snowflake struct {
	mu           sync.Mutex
	dataCenterId uint64
	workId       uint64
	sequence     uint64
	lastTs       int64
}

/**
GetUID 得到全局唯一ID int64类型
首位0(1位) + 毫秒时间戳(41位) + 数据中心标识(5位) + 工作机器标识(5位) + 自增id(12位)
时间可以保证400年不重复
数据中心和机器标识一起标识节点，最多支持1024个节点
每个节点每一毫秒能生成最多4096个id
63      62            21            16        11       0
+-------+-------------+-------------+---------+--------+
| 未使用 | 毫秒级时间戳 | 数据中心标识 | 工作机器 | 自增id |
+-------+-------------+-------------+---------+--------+
*/
func (s *Snowflake) GetUUID() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.ts()
	if t != s.lastTs {
		s.sequence = 0
		s.lastTs = t
	} else {
		s.sequence = (s.sequence + 1) & sn_mask
		if s.sequence == 0 {
			t = s.waitMs(s.lastTs)
			s.lastTs = t
		}
	}
	ms := (uint64(t) & ts_mask) << 22
	center := (s.dataCenterId & center_id_mask) << 17
	work := (s.workId & machine_id_mask) << 12
	return ms | center | work | s.sequence
}

func (s *Snowflake) ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (s *Snowflake) waitMs(lastTs int64) int64 {
	t := s.ts()
	for t <= lastTs {
		t = s.ts()
	}
	return t
}

func Init(workId uint64) {
	once.Do(func() {
		log = z.GetLogger()
		initSnowflake(workId)
	})

}

func initSnowflake(workId uint64) {
	c := config.C()
	cfg := &Cfg{}
	if err := c.Path("snowflake", &cfg); err != nil {
		log.Panic(err.Error())
	}
	s := Snowflake{
		dataCenterId: cfg.DataCenterId,
		workId:       workId,
		sequence:     0,
		lastTs:       0,
	}
	if s.dataCenterId > center_id_mask {
		log.Panic("data center id too big", zap.Uint64("date_center_id", s.dataCenterId))
	}
	if s.workId > machine_id_mask {
		log.Panic("work id too big", zap.Uint64("work_id", s.workId))
	}
}

func GetSnowkflak() *Snowflake {
	return s
}
