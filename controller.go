package neuron

import (
	"github.com/neuronlabs/errors"

	"github.com/neuronlabs/neuron-core/class"
	"github.com/neuronlabs/neuron-core/config"
	"github.com/neuronlabs/neuron-core/controller"
)

// Controller gets default controller.
func Controller() *controller.Controller {
	return controller.Default()
}

// InitController initializes default controller for given configuration.
// Returns error if the default controller is already defined.
// In order to force initializing new controller set the controller.DefaultController to nil.
func InitController(cfg *config.Controller) (err error) {
	if controller.DefaultController != nil {
		return errors.NewDetf(class.InternalInitController, "default controller already defined")
	}
	controller.DefaultController, err = controller.New(cfg)
	return err
}

// NewController creates and initializes controller for provided config.
func NewController(cfg *config.Controller) (*controller.Controller, error) {
	return controller.New(cfg)
}

// RegisterModels registers all models into default controller.
// This function requires repositories to be registered before.
// Returns error if the model was already registered.
func RegisterModels(models ...interface{}) error {
	return controller.Default().RegisterModels(models...)
}

// RegisterRepository registers repository into default controller.
// Returns error if the repository with given name was already registered.
// By default the first registered repository is set to the default repository
// for all models that doesn't define their repositories, unless default
// controller's config DisallowDefaultRepository is set to false.
func RegisterRepository(name string, repo *config.Repository) error {
	return controller.Default().RegisterRepository(name, repo)
}
