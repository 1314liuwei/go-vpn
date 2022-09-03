package protocol

type IPRaw struct {
	VersionAndHeaderLength      byte
	DifferentiatedServicesField [2]byte
	TotalLength                 [2]byte
	Identification              [2]byte
	Flags                       [2]byte
	TimeToLive                  byte
	Protocol                    byte
	HeaderChecksum              [2]byte
	SourceAddr                  [4]byte
	DestAddr                    [4]byte
}

func parseIPRawPacket(packet IPRaw) {

}
