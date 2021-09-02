package httpserver

var PiDummy = []PiData{
	{
		Name:     "dummyPI",
		Ip:       "10.0.0.2",
		LastSeen: "10 august 00:00:00",
		Services: []Services{
			{
				Name:    "FTP",
				Running: true,
			},
			{
				Name:    "SSH",
				Running: false,
			},
			{
				Name:    "RDP",
				Running: true,
			},
		},
	},
	{
		Name:     "dummyPI2",
		Ip:       "10.0.0.3",
		LastSeen: "10 august 01:00:00",
		Services: []Services{
			{
				Name:    "SSH",
				Running: false,
			},
		},
	},
}

var pd = PageData{
	PIs: PiDummy,
}

type PageData struct {
	Title string
	Loggedin bool
	Navbar bool
	PIs []PiData
}

type PiData struct {
	Name     string
	Ip       string
	LastSeen string
	Services []Services
}
type Services struct {
	Name    string
	Running bool
}
