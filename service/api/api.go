package api

//import packages and http configs
import (   
	"context"
	"time"

	goHttp "net/http"

	"github.com/pterm/pterm"
	httpIface "github.com/taubyte/http"
	http "github.com/taubyte/http/basic"
	"github.com/taubyte/http/options"
	"github.com/taubyte/tau/libdream/common"
	"github.com/taubyte/tau/libdream/services"
)

type multiverseService struct { 	
	//responsible for setting up various routes @http_routes.go from github.com/taubyte/http
	rest httpIface.Service		
	common.Multiverse
}

//initialize and start Dreamland service
func BigBang() error {   
	var err error

	//Creating multiverse instance
	srv := &multiverseService{
		Multiverse: services.NewMultiVerse(),    
	}

	//Setup rest service and handle it's errors
	srv.rest, err = http.New(srv.Context(), options.Listen(common.DreamlandApiListen), options.AllowedOrigins(true, []string{".*"})) 
	if err != nil {		
		return err	
	}

	//Set up HTTP routes for the rest service and start it
	srv.setUpHttpRoutes().Start()

	//cancel requests that exceeds 10seconds after starting http route
	waitCtx, waitCtxC := context.WithTimeout(srv.Context(), 10*time.Second)
	defer waitCtxC()

	//wait until dreamland is ready 
	for {
		select {
		case <-waitCtx.Done():
			return waitCtx.Err()
		case <-time.After(100 * time.Millisecond):
			if srv.rest.Error() != nil {
				pterm.Error.Println("Dreamland failed to start")
				return srv.rest.Error()
			}
			_, err := goHttp.Get("http://" + common.DreamlandApiListen)
			if err == nil {
				pterm.Info.Println("Dreamland ready")
				return nil
			}
		}
	}
}
