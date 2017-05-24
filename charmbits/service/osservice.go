package service

import (
	"github.com/juju/gocharm/vendored/service"
	"github.com/juju/gocharm/vendored/service/common"
	"github.com/juju/utils"
	"github.com/juju/utils/series"
	"gopkg.in/errgo.v1"
)

// OSServiceParams holds the parameters for
// creating a new service.
type OSServiceParams struct {
	// Name holds the name of the service.
	Name string

	// Description holds the description of the service.
	Description string

	// Output holds the file where output will be put.
	Output string

	// Exe holds the name of the executable to run.
	Exe string

	// Args holds any arguments to the executable,
	// which should be OK to to pass to the shell
	// without quoting.
	Args []string
}

// NewService is used to create a new service.
// It is defined as a variable so that it can be
// replaced for testing purposes.
var NewService func(p OSServiceParams) OSService = newService

func newService(p OSServiceParams) OSService {
	series, err := series.HostSeries()
	if err != nil {
		return errorService{errgo.Notef(err, "cannot get host OS series")}
	}
	s, err := service.NewService(p.Name, common.Conf{
		Desc: p.Description,
		// Bizarrely, ServiceBinary/ServiceArgs and ExecStart are non-orthogonal
		// and both need to be provided.
		ServiceBinary: p.Exe,
		ServiceArgs:   p.Args,
		ExecStart:     quoteCommand(p.Exe, p.Args),
		Logfile:       p.Output,
	}, series)
	if err != nil {
		return errorService{errgo.Notef(err, "cannot create service object")}
	}
	return s
}

func quoteCommand(cmd string, args []string) string {
	buf := []byte(utils.ShQuote(cmd))
	for _, arg := range args {
		buf = append(buf, ' ')
		buf = append(buf, []byte(utils.ShQuote(arg))...)
	}
	return string(buf)
}

type errorService struct {
	err error
}

func (e errorService) Start() error {
	return errgo.Mask(e.err)
}

func (e errorService) Stop() error {
	return errgo.Mask(e.err)
}

func (e errorService) Install() error {
	return errgo.Mask(e.err)
}

func (e errorService) Remove() error {
	return errgo.Mask(e.err)
}

func (e errorService) Running() (bool, error) {
	return false, errgo.Mask(e.err)
}
