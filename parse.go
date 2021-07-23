package mediaindex

import (
	"errors"
	"fmt"

	"github.com/suconghou/mediaindex/ebml"
	"github.com/suconghou/mediaindex/sidx"
)

// ParseMp4 info, start 最小值为 indexEndOffset+1 , end 最大值为文件大小(sidx内部自动读出来的)
func ParseMp4(data []byte, indexEndOffset uint64) (res map[int][2]uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	parser := sidx.NewParser(data)
	info := parser.Parse(uint32(indexEndOffset))
	res = map[int][2]uint64{}
	for i, item := range info.References {
		res[i] = [2]uint64{uint64(item.StartRange), uint64(item.EndRange)}
	}
	return
}

// ParseWebM info,start最小值为indexEndOffset+1,end最大值为文件大小;下一个块的start是等于上一个块的end
func ParseWebM(data []byte, indexEndOffset uint64, totalSize uint64) (res map[int][2]uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	parser := ebml.NewParser(data)
	info := parser.Parse()
	res = map[int][2]uint64{}
	if len(info) > 0 && len(info[0].Children) > 0 {
		var arr = []uint64{}
		for _, item := range info[0].Children {
			if item.ID != "bb" {
				continue
			}
			var (
				CueTrackPositions  = item.Children[1]
				CueClusterPosition = CueTrackPositions.Children[1]
				n                  = ebml.VarNum(CueClusterPosition.Value)
			)
			arr = append(arr, n)
		}
		var l = len(arr) - 1
		if l < 0 {
			return
		}
		var segmentStart = indexEndOffset - arr[0] + 1
		var segmentEnd = totalSize
		for i, item := range arr {
			var start = item + segmentStart
			var end uint64
			if i < l {
				end = arr[i+1] + segmentStart
			} else {
				// last item,range end is its length
				end = segmentEnd
			}
			res[i] = [2]uint64{start, end}
		}
	}
	return
}
