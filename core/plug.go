package core

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"plugin"

	"github.com/google/gopacket"
	hp "github.com/sinuxlee/go-sniffer/plugSrc/http"
	kafka "github.com/sinuxlee/go-sniffer/plugSrc/kafka"
	mongodb "github.com/sinuxlee/go-sniffer/plugSrc/mongodb"
	mssql "github.com/sinuxlee/go-sniffer/plugSrc/mssql"
	mysql "github.com/sinuxlee/go-sniffer/plugSrc/mysql"
	redis "github.com/sinuxlee/go-sniffer/plugSrc/redis"
)

// PlugInterface All internal plug-ins must implement this interface
type PlugInterface interface {
	// ResolveStream - entry
	ResolveStream(net gopacket.Flow, transport gopacket.Flow, r io.Reader)

	// BPFFilter     - set BPF, like: mysql(tcp and port 3306)
	BPFFilter() string

	// SetFlag       - plug-in params
	SetFlag([]string)

	// Version       - plug-in version
	Version() string
}

type Plug struct {
	dir           string
	ResolveStream func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
	BPF           string

	InternalPlugList map[string]PlugInterface
	ExternalPlugList map[string]ExternalPlug
}

type ExternalPlug struct {
	Name          string
	Version       string
	ResolvePacket func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
	BPFFilter     func() string
	SetFlag       func([]string)
}

func NewPlug() (*Plug, error) {
	var p Plug
	p.dir, _ = filepath.Abs("./plug/")
	p.LoadInternalPlugList()
	err := p.LoadExternalPlugList()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Plug) LoadInternalPlugList() {
	p.InternalPlugList = map[string]PlugInterface{
		"http":    hp.NewInstance(),
		"redis":   redis.NewInstance(),
		"mysql":   mysql.NewInstance(),
		"mongodb": mongodb.NewInstance(),
		"kafka":   kafka.NewInstance(),
		"mssql":   mssql.NewInstance(),
	}
}

func (p *Plug) LoadExternalPlugList() error {
	dir, err := os.ReadDir(p.dir)
	if err != nil {
		return err
	}

	p.ExternalPlugList = make(map[string]ExternalPlug)
	for _, fi := range dir {
		if fi.IsDir() || path.Ext(fi.Name()) != ".so" {
			continue
		}

		plug, err := plugin.Open(p.dir + "/" + fi.Name())
		if err != nil {
			return err
		}

		versionFunc, err := plug.Lookup("Version")
		if err != nil {
			return err
		}

		setFlagFunc, err := plug.Lookup("SetFlag")
		if err != nil {
			return err
		}

		BPFFilterFunc, err := plug.Lookup("BPFFilter")
		if err != nil {
			return err
		}

		ResolvePacketFunc, err := plug.Lookup("ResolvePacket")
		if err != nil {
			return err
		}

		version := versionFunc.(func() string)()
		p.ExternalPlugList[fi.Name()] = ExternalPlug{
			ResolvePacket: ResolvePacketFunc.(func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)),
			SetFlag:       setFlagFunc.(func([]string)),
			BPFFilter:     BPFFilterFunc.(func() string),
			Version:       version,
			Name:          fi.Name(),
		}
	}

	return nil
}

func (p *Plug) ChangePath(dir string) {
	p.dir = dir
}

func (p *Plug) PrintList() {
	// Print Internal Plug
	for inPlugName := range p.InternalPlugList {
		fmt.Println("internal plug : " + inPlugName)
	}

	// split
	fmt.Println("-- --- --")

	// print External Plug
	for exPlugName := range p.ExternalPlugList {
		fmt.Println("external plug : " + exPlugName)
	}
}

func (p *Plug) SetOption(plugName string, plugParams []string) error {
	fmt.Println("internalPlug", plugName)

	// Load Internal Plug
	if internalPlug, ok := p.InternalPlugList[plugName]; ok {
		p.ResolveStream = internalPlug.ResolveStream
		internalPlug.SetFlag(plugParams)
		p.BPF = internalPlug.BPFFilter()

		return nil
	}

	// Load External Plug
	plug, err := plugin.Open("./plug/" + plugName)
	if err != nil {
		return err
	}

	resolvePacket, err := plug.Lookup("ResolvePacket")
	if err != nil {
		return err
	}

	setFlag, err := plug.Lookup("SetFlag")
	if err != nil {
		return err
	}

	BPFFilter, err := plug.Lookup("BPFFilter")
	if err != nil {
		return err
	}

	p.ResolveStream = resolvePacket.(func(net gopacket.Flow, transport gopacket.Flow, r io.Reader))
	setFlag.(func([]string))(plugParams)
	p.BPF = BPFFilter.(func() string)()

	return nil
}
