package configdomain

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// PartialConfig contains configuration data as it is stored in the local or global Git configuration.
type PartialConfig struct {
	Aliases                  Aliases
	ContributionBranches     *gitdomain.LocalBranchNames
	GitHubToken              *GitHubToken
	GitLabToken              *GitLabToken
	GitUserEmail             *string
	GitUserName              *string
	GiteaToken               *GiteaToken
	HostingOriginHostname    *HostingOriginHostname
	HostingPlatform          *HostingPlatform
	Lineage                  *Lineage
	MainBranch               *gitdomain.LocalBranchName
	ObservedBranches         *gitdomain.LocalBranchNames
	Offline                  *Offline
	ParkedBranches           *gitdomain.LocalBranchNames
	PerennialBranches        *gitdomain.LocalBranchNames
	PerennialRegex           *PerennialRegex
	PushHook                 *PushHook
	PushNewBranches          *PushNewBranches
	ShipDeleteTrackingBranch *ShipDeleteTrackingBranch
	SyncBeforeShip           *SyncBeforeShip
	SyncFeatureStrategy      *SyncFeatureStrategy
	SyncPerennialStrategy    *SyncPerennialStrategy
	SyncUpstream             *SyncUpstream
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{ //nolint:exhaustruct
		Aliases: Aliases{},
	}
}
