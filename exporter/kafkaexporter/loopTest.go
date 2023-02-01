package kafkaexporter

import (
	"fmt"
)

func main() {
	sd := NewStreamDataRecordMessage()

	sd.Data = &UnionNullMapArrayUnionStringNull{}
	sd.Data.MapArrayUnionStringNull = make(map[string][]*UnionStringNull)
	sd.Data.MapArrayUnionStringNull["1"] = append(sd.Data.MapArrayUnionStringNull["1"], &UnionStringNull{
		String: "hello1",
	})
	sd.Data.MapArrayUnionStringNull["1"] = append(sd.Data.MapArrayUnionStringNull["1"], &UnionStringNull{
		String: "world1",
	})
	sd.Data.MapArrayUnionStringNull["2"] = append(sd.Data.MapArrayUnionStringNull["2"], &UnionStringNull{
		String: "hello2",
	})
	sd.Data.MapArrayUnionStringNull["2"] = append(sd.Data.MapArrayUnionStringNull["2"], &UnionStringNull{
		String: "world2",
	})

	fmt.Println(sd.Data.MapArrayUnionStringNull["1"][0].String)
	fmt.Println(sd.Data.MapArrayUnionStringNull["1"][1].String)
	fmt.Println(sd.Data.MapArrayUnionStringNull["2"][0].String)
	fmt.Println(sd.Data.MapArrayUnionStringNull["2"][1].String)
}
