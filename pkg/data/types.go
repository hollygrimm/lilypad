package data

// used by resource providers to describe their resources
// use by job offers to describe their requirements
// when used by resource providers - these are absolute values
// when used by job offers - these are minimum requirements
type Spec struct {
	// Milli-GPU
	// Whilst it's unlikely that partial GPU's make sense
	// let's not use a float and fix the precision to 1/1000
	GPU int `json:"gpu"`

	// Milli-CPU
	CPU int `json:"cpu"`

	// Megabytes
	RAM int `json:"ram"`
}

type ModuleInputs map[string]string

// this is what is loaded from the template file in the git repo
type Module struct {
	// the min spec that this module requires
	// e.g. does this module need a GPU?
	// the module file itself will contain this spec
	// and so the module will need to be downloaded
	// and executed for this spec to be known
	Spec Spec `json:"spec"`
}

// describes a workload to be run
// this pins a go-template.yaml file
// that is a bacalhau job spec
type ModuleConfig struct {

	// used for the shortcuts
	// this is in the modules package
	// where we keep a map of named modules
	// and their versions onto the
	// repo, hash and path below
	Name    string `json:"name"`
	Version string `json:"version"`

	// needs to be a http url for a git repo
	// we must be able to clone it without credentials
	Repo string `json:"repo"`
	// the git hash to pin the module
	// we will 'git checkout' this hash
	Hash string `json:"hash"`
	// once the checkout has been done
	// this is the path to the module template
	// within the repo
	Path string `json:"path"`
}

type Result struct {
	// this is the cid of the result where ID is set to empty string
	ID     string `json:"id"`
	DealID string `json:"deal_id"`
	// the CID of the actual results
	DataID           string `json:"results_id"`
	InstructionCount uint64 `json:"instruction_count"`
}

type PricingMode string

// MarketPrice means - get me the best deal
// job creators will do this by default i.e. "just buy me the cheapest"
// FixedPrice means - take it or leave it
// resource creators will do this by default i.e. "this is my price"
const (
	MarketPrice PricingMode = "MarketPrice"
	FixedPrice  PricingMode = "FixedPrice"
)

// represents the cost of a job
type Pricing struct {
	Mode                      PricingMode `json:"mode"`
	InstructionPrice          uint64      `json:"instruction_price"`
	Timeout                   uint64      `json:"timeout"`
	TimeoutCollateral         uint64      `json:"timeout_collateral"`
	PaymentCollateral         uint64      `json:"payment_collateral"`
	ResultsCollateralMultiple uint64      `json:"results_collateral_multiple"`
	MediationFee              uint64      `json:"mediation_fee"`
}

// posted to the solver by a job creator
type JobOffer struct {
	// this is the cid of the job offer where ID is set to empty string
	ID string `json:"id"`
	// the address of the job creator
	JobCreator string `json:"job_creator"`
	// this is the CID of the Module description
	ModuleID string `json:"module_id"`
	// the actual module that is being offered
	// this must hash to the ModuleID above
	Module ModuleConfig `json:"module"`
	// the user inputs to the module
	// these values will power the go template
	Inputs ModuleInputs `json:"inputs"`
	// the offered price
	Pricing Pricing `json:"pricing"`
}

// posted to the solver by a resource provider
type ResourceOffer struct {
	// this is the cid of the resource offer where ID is set to empty string
	ID string `json:"id"`
	// the address of the job creator
	ResourceProvider string `json:"resource_provider"`
	// allows a resource provider to manage multiple offers
	// that are essentially the same
	Index int `json:"index"`
	// the spec being offered
	Spec Spec `json:"spec"`
	// the module ID's that this resource provider can run
	// an empty list means ALL modules
	Modules []string `json:"modules"`
	// the default pricing for this resource offer
	// i.e. this is for any module
	DefaultPricing Pricing `json:"default_pricing"`
	// the pricing for each module
	// this allows a resource provider to charge more
	// for certain modules
	ModulePricing map[string]Pricing `json:"module_pricing"`
}

// generated by the solver
type Match struct {
	// this is the cid of the resource offer where ID is set to empty string
	ID string `json:"id"`

	// how long the resource offer is valid for
	Timeout uint64 `json:"timeout"`

	ResourceOffer string `json:"resource_offer"`
	JobOffer      string `json:"job_offer"`
}

// represents the cost of a job
// this is both the bid and ask in the 2 sided marketplace
// job creators will attach this to their job offers
// and resource providers will attach this to their resource offers
// the solvers job is to propose the most efficient match
// for the job creator
type Deal struct {
	// this is the cid of the deal where ID is set to empty string
	ID               string        `json:"id"`
	ResourceProvider string        `json:"resource_provider"`
	JobCreator       string        `json:"job_creator"`
	JobOffer         JobOffer      `json:"job_offer"`
	ResourceOffer    ResourceOffer `json:"resource_offer"`
	// this is the agreed upon price
	Pricing Pricing `json:"pricing"`
}
