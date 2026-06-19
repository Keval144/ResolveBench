package dns

type Provider struct {
	Name         string
	PrimaryDNS   string
	SecondaryDNS string
	UseCase      string
	Description  string
}

var Providers = []Provider{
	{
		Name:         "Cloudflare",
		PrimaryDNS:   "1.1.1.1",
		SecondaryDNS: "1.0.0.1",
		UseCase:      "Speed and Privacy",
		Description:  "Cloudflare DNS focuses on fast resolution and strong privacy protections with DNSSEC support.",
	},
	{
		Name:         "Google DNS",
		PrimaryDNS:   "8.8.8.8",
		SecondaryDNS: "8.8.4.4",
		UseCase:      "Reliability",
		Description:  "Google Public DNS offers high reliability and global infrastructure with low latency worldwide.",
	},
	{
		Name:         "Quad9",
		PrimaryDNS:   "9.9.9.9",
		SecondaryDNS: "149.112.112.112",
		UseCase:      "Security",
		Description:  "Quad9 blocks known malicious domains, providing security against phishing and malware threats.",
	},
	{
		Name:         "AdGuard DNS",
		PrimaryDNS:   "94.140.14.14",
		SecondaryDNS: "94.140.15.15",
		UseCase:      "Ad Blocking",
		Description:  "AdGuard DNS filters out ads, trackers, and malicious websites at the DNS level system-wide.",
	},
	{
		Name:         "OpenDNS",
		PrimaryDNS:   "208.67.222.222",
		SecondaryDNS: "208.67.220.220",
		UseCase:      "Content Filtering",
		Description:  "OpenDNS provides customizable content filtering and phishing protection for families and businesses.",
	},
}

var BenchDomains = []string{
	"github.com",
	"google.com",
	"reddit.com",
	"wikipedia.org",
	"cloudflare.com",
	"microsoft.com",
	"amazon.com",
	"stackoverflow.com",
	"youtube.com",
	"openai.com",
}

var NetworkTargets = []string{
	"1.1.1.1",
	"8.8.8.8",
	"9.9.9.9",
	"cloudflare.com",
	"google.com",
}
