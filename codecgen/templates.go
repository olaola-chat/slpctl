package codecgen

var templateV2 = `package codec  
  
import (  
    "context"   
	"database/sql"   
	"fmt"   
	"time"  
	"%mode/app/dao"   
	"%mode/app/pb"   
	"%mode/library/go2cache" 
	"%mode/library"   
	"google.golang.org/protobuf/proto"
)  
var (  
    %pbNameRedisCodec *go2cache.Server   
	expiredTtl%pbNameSeconds     = int64(%ttl)
)  
  
const (  
    tableKey%pbName = "table.key.%tableKey.%id.%d"
)  
  
func init() {  
	expiredTime := time.Hour    
	if expiredTtl%pbNameSeconds > 0 {   
		expiredTime = time.Duration(expiredTtl%pbNameSeconds) * time.Second  
	}   
	%pbNameRedisCodec = go2cache.NewOnlyRedisServer(library.%redis-db, &%lowerNameCodec{}, go2cache.WithTtl(expiredTime))
	TableCodecMap["%tableKey"] = %pbNameRedisCodec
}  
  
type %lowerNameCodec struct {  
}  
  
func (b %lowerNameCodec) Pt() proto.Message {  
    return &pb.Entity%pbName{}
}  
  
// Key 生成缓存key  
func (b %lowerNameCodec) Key(key uint32) string {  
    if key == 0 {     
		return ""   
	}   
	return fmt.Sprintf(tableKey%pbName, key)
}  
  
// Pk 根据proto数据，获取主键信息  
func (b %lowerNameCodec) Pk(data proto.Message) uint32 {  
    if entity, ok := data.(*pb.Entity%pbName); ok {    
		id:= entity.%entity_id
		return uint32(id)
	}    
	return 0
}  
  
func (b %lowerNameCodec) One(ctx context.Context, key uint32, data proto.Message) error {  
    //排除大字段，description  
    err := dao.%pbName.Ctx(ctx).Where("%id = ?", key).Struct(data)  
	if err != nil && err != sql.ErrNoRows {      
		return err   
	}    
	return nil
}  
  
func (b %lowerNameCodec) FindAll(ctx context.Context, keys []uint32, callback go2cache.Find2Item) error {  
    res, err := dao.%pbName.Ctx(ctx).Where("%id in (?)", keys).FindAll()   
	if err != nil { 
		return err   
	}   
	for _, item := range res {     
		callback(item)  
	}  
	return nil
}  
`
