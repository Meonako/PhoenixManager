package party

type TARGET int

const (
	FrontierClash = iota
	VoidRift
	WorldBoss
	Raid

	Len
)

func (t TARGET) String() string {
	switch t {
	case Raid:
		return "Raid"
	case VoidRift:
		return "Void Rift"
	case WorldBoss:
		return "World Boss"
	default:
		return "Frontier Clash"
	}
}

func (t TARGET) GetMaxPlayer() int {
	if t == Raid {
		return 8
	}

	return 4
}

func GetEnum(i string) int {
	switch i {
	case "VOID-RIFT":
		return VoidRift
	case "WORLD-BOSS":
		return WorldBoss
	case "RAID":
		return Raid
	default:
		return FrontierClash
	}
}
