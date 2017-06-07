package algorithm


 // 一财一刻、两财一刻、三财一刻检测
// 财神归位检测
// 升序的slice、一定为胡牌型
// todo 两财一刻和财神归位同时存在的情况
func guiwei(cs []byte, ch, ps, ks []uint32, wildcard byte) (value int64) {
	caishenCount := 0
	for _, v := range cs {
		if v == WILDCARD {
			caishenCount ++
		}
	}

	if caishenCount == 0 {
		return 0
	}

	list := make([]byte, len(cs))
	copy(list, cs)
	cs = list

	// 复原财神
	for k, v := range cs {
		if v == WILDCARD {
			cs[k] = wildcard
		}
	}
	Sort(cs, 0, len(cs)-1)

	count := caishenCount
	for i := 0; i < caishenCount; i++ {
		for k, v := range cs {
			if v == wildcard {
				//glog.Errorf("==============%+x %d %d %x",cs,caishenCount,k,v)
				// 这里不加入13不靠检测，13不靠内部带有财神归位检测，13不靠没有n财n刻
				if existHu3n2(cs, wildcard) ||
					existLuanFeng(getHandPongKong(cs, ch, ps, ks, wildcard)) > 0 ||
					exist7pair(cs) > 0 {
					if count == 1 {
						//glog.Errorf("==============%+x %d %d %x",cs,caishenCount,k,v)
						value = HU_GUI_WEI1
						return
					} else if count == 2 {
						//glog.Errorf("==============%+x %d %d %x",cs,caishenCount,k,v)
						value = HU_GUI_WEI2
						return
					} else if count == 3 {
						//glog.Errorf("==============%+x %d %d %x",cs,caishenCount,k,v)
						value = HU_GUI_WEI3
						return
					}
				}
				cs[k] = WILDCARD
				Sort(cs, 0, len(cs)-1)
				count --
				break
			}
		}
	}

	return
}

func ExistNCaiNKe(cs []byte, ch, ps, ks []uint32, wildcard byte) int64 {
	// 计算财神的数量
	caishenCount := 0
	for _, v := range cs {
		if v == WILDCARD {
			caishenCount ++
		}
	}

	if caishenCount == 0 {
		return 0
	}

	if len(cs) >= 5 {
		// 3财1刻
		if caishenCount == 3 {
			list := make([]byte, 0, len(cs))
			for _, v := range cs {
				if v != WILDCARD {
					list = append(list, v)
				}
			}
			Sort(list, 0, len(list)-1)
			if existHu3n2(list, wildcard)  ||
				existLuanFeng(getHandPongKong(list, ch, ps, ks, wildcard)) > 0 {
				return HU_CAI_3
			}
		}

		baiCount := 0
		for _, v := range cs {
			if v == BAI {
				baiCount ++
			}
		}

		// 2财1刻
		if baiCount >= 1 && caishenCount >= 2 {
			list := make([]byte, 0, len(cs))
			bai := 0
			cai := 0
			for _, v := range cs {
				if bai < 1 && v == BAI {
					bai ++
					continue
				}

				if cai < 2 && v == WILDCARD {
					cai ++
					continue
				}

				list = append(list, v)
			}

			Sort(list, 0, len(list)-1)
			if  existHu3n2(list, wildcard) ||
				existLuanFeng(getHandPongKong(list, ch, ps, ks, wildcard)) > 0 {
				value := guiwei(list, ch, ps, ks, wildcard)
				return HU_CAI_2 | value
			}
		}
		// 1财1刻
		if baiCount >= 2 && caishenCount >= 1 {

			list := make([]byte, 0, len(cs))
			bai := 0
			cai := 0
			for _, v := range cs {
				if bai < 2 && v == BAI {
					bai ++
					continue
				}

				if cai < 1 && v == WILDCARD {
					cai ++
					continue
				}

				list = append(list, v)
			}

			Sort(list, 0, len(list)-1)
			if  existHu3n2(list, wildcard) ||
				existLuanFeng(getHandPongKong(list, ch, ps, ks, wildcard)) > 0 {
				value := guiwei(list, ch, ps, ks, wildcard)
				return HU_CAI_1 | value
			}
		}
	}
	//glog.Errorf("==============%+x",cs)
	return guiwei(cs, ch, ps, ks, wildcard)
}
