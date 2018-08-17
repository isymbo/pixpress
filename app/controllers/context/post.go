package context

import (
	"github.com/isymbo/pixpress/app/models"
)

type Post struct {
	// AccessMode   models.AccessMode
	//IsWatching   bool
	Post  *models.Post
	Owner *models.User
	// BranchName string
	// TagName    string
	// TreePath   string
	// CommitID   string
	// RepoLink   string
}

// // IsOwner returns true if current user is the owner of repository.
// func (r *Repository) IsOwner() bool {
// 	return r.AccessMode >= models.ACCESS_MODE_OWNER
// }

// // IsAdmin returns true if current user has admin or higher access of repository.
// func (r *Repository) IsAdmin() bool {
// 	return r.AccessMode >= models.ACCESS_MODE_ADMIN
// }

// // IsWriter returns true if current user has write or higher access of repository.
// func (r *Repository) IsWriter() bool {
// 	return r.AccessMode >= models.ACCESS_MODE_WRITE
// }

// // HasAccess returns true if the current user has at least read access for this repository
// func (r *Repository) HasAccess() bool {
// 	return r.AccessMode >= models.ACCESS_MODE_READ
// }

// // [0]: issues, [1]: wiki
// func RepoAssignment(pages ...bool) macaron.Handler {
// 	return func(c *Context) {
// 		var (
// 			owner        *models.User
// 			err          error
// 			isIssuesPage bool
// 			isWikiPage   bool
// 		)

// 		if len(pages) > 0 {
// 			isIssuesPage = pages[0]
// 		}
// 		if len(pages) > 1 {
// 			isWikiPage = pages[1]
// 		}

// 		ownerName := c.Params(":username")
// 		repoName := strings.TrimSuffix(c.Params(":reponame"), ".git")
// 		refName := c.Params(":branchname")
// 		if len(refName) == 0 {
// 			refName = c.Params(":path")
// 		}

// 		// Check if the user is the same as the repository owner
// 		if c.IsLogged && c.User.LowerName == strings.ToLower(ownerName) {
// 			owner = c.User
// 		} else {
// 			owner, err = models.GetUserByName(ownerName)
// 			if err != nil {
// 				c.NotFoundOrServerError("GetUserByName", errors.IsUserNotExist, err)
// 				return
// 			}
// 		}
// 		c.Repo.Owner = owner
// 		c.Data["Username"] = c.Repo.Owner.Name

// 		repo, err := models.GetRepositoryByName(owner.ID, repoName)
// 		if err != nil {
// 			c.NotFoundOrServerError("GetRepositoryByName", errors.IsRepoNotExist, err)
// 			return
// 		}

// 		c.Repo.Repository = repo
// 		c.Data["RepoName"] = c.Repo.Repository.Name
// 		c.Data["IsBareRepo"] = c.Repo.Repository.IsBare
// 		c.Repo.RepoLink = repo.Link()
// 		c.Data["RepoLink"] = c.Repo.RepoLink
// 		c.Data["RepoRelPath"] = c.Repo.Owner.Name + "/" + c.Repo.Repository.Name

// 		// Admin has super access.
// 		if c.IsLogged && c.User.IsAdmin {
// 			c.Repo.AccessMode = models.ACCESS_MODE_OWNER
// 		} else {
// 			mode, err := models.AccessLevel(c.UserID(), repo)
// 			if err != nil {
// 				c.ServerError("AccessLevel", err)
// 				return
// 			}
// 			c.Repo.AccessMode = mode
// 		}

// 		// Check access
// 		if c.Repo.AccessMode == models.ACCESS_MODE_NONE {
// 			// Redirect to any accessible page if not yet on it
// 			if repo.IsPartialPublic() &&
// 				(!(isIssuesPage || isWikiPage) ||
// 					(isIssuesPage && !repo.CanGuestViewIssues()) ||
// 					(isWikiPage && !repo.CanGuestViewWiki())) {
// 				switch {
// 				case repo.CanGuestViewIssues():
// 					c.Redirect(repo.Link() + "/issues")
// 				case repo.CanGuestViewWiki():
// 					c.Redirect(repo.Link() + "/wiki")
// 				default:
// 					c.NotFound()
// 				}
// 				return
// 			}

// 			// Response 404 if user is on completely private repository or possible accessible page but owner doesn't enabled
// 			if !repo.IsPartialPublic() ||
// 				(isIssuesPage && !repo.CanGuestViewIssues()) ||
// 				(isWikiPage && !repo.CanGuestViewWiki()) {
// 				c.NotFound()
// 				return
// 			}

// 			c.Repo.Repository.EnableIssues = repo.CanGuestViewIssues()
// 			c.Repo.Repository.EnableWiki = repo.CanGuestViewWiki()
// 		}

// 		if repo.IsMirror {
// 			c.Repo.Mirror, err = models.GetMirrorByRepoID(repo.ID)
// 			if err != nil {
// 				c.ServerError("GetMirror", err)
// 				return
// 			}
// 			c.Data["MirrorEnablePrune"] = c.Repo.Mirror.EnablePrune
// 			c.Data["MirrorInterval"] = c.Repo.Mirror.Interval
// 			c.Data["Mirror"] = c.Repo.Mirror
// 		}

// 		gitRepo, err := git.OpenRepository(models.RepoPath(ownerName, repoName))
// 		if err != nil {
// 			c.ServerError(fmt.Sprintf("RepoAssignment Invalid repo '%s'", c.Repo.Repository.RepoPath()), err)
// 			return
// 		}
// 		c.Repo.GitRepo = gitRepo

// 		tags, err := c.Repo.GitRepo.GetTags()
// 		if err != nil {
// 			c.ServerError(fmt.Sprintf("GetTags '%s'", c.Repo.Repository.RepoPath()), err)
// 			return
// 		}
// 		c.Data["Tags"] = tags
// 		c.Repo.Repository.NumTags = len(tags)

// 		c.Data["Title"] = owner.Name + "/" + repo.Name
// 		c.Data["Repository"] = repo
// 		c.Data["Owner"] = c.Repo.Repository.Owner
// 		c.Data["IsRepositoryOwner"] = c.Repo.IsOwner()
// 		c.Data["IsRepositoryAdmin"] = c.Repo.IsAdmin()
// 		c.Data["IsRepositoryWriter"] = c.Repo.IsWriter()

// 		c.Data["DisableSSH"] = setting.SSH.Disabled
// 		c.Data["DisableHTTP"] = setting.Repository.DisableHTTPGit
// 		c.Data["CloneLink"] = repo.CloneLink()
// 		c.Data["WikiCloneLink"] = repo.WikiCloneLink()

// 		if c.IsLogged {
// 			c.Data["IsWatchingRepo"] = models.IsWatching(c.User.ID, repo.ID)
// 			c.Data["IsStaringRepo"] = models.IsStaring(c.User.ID, repo.ID)
// 		}

// 		// repo is bare and display enable
// 		if c.Repo.Repository.IsBare {
// 			return
// 		}

// 		c.Data["TagName"] = c.Repo.TagName
// 		brs, err := c.Repo.GitRepo.GetBranches()
// 		if err != nil {
// 			c.ServerError("GetBranches", err)
// 			return
// 		}
// 		c.Data["Branches"] = brs
// 		c.Data["BrancheCount"] = len(brs)

// 		// If not branch selected, try default one.
// 		// If default branch doesn't exists, fall back to some other branch.
// 		if len(c.Repo.BranchName) == 0 {
// 			if len(c.Repo.Repository.DefaultBranch) > 0 && gitRepo.IsBranchExist(c.Repo.Repository.DefaultBranch) {
// 				c.Repo.BranchName = c.Repo.Repository.DefaultBranch
// 			} else if len(brs) > 0 {
// 				c.Repo.BranchName = brs[0]
// 			}
// 		}
// 		c.Data["BranchName"] = c.Repo.BranchName
// 		c.Data["CommitID"] = c.Repo.CommitID

// 		c.Data["IsGuest"] = !c.Repo.HasAccess()
// 	}
// }

// func RequireRepoAdmin() macaron.Handler {
// 	return func(c *Context) {
// 		if !c.IsLogged || (!c.Repo.IsAdmin() && !c.User.IsAdmin) {
// 			c.NotFound()
// 			return
// 		}
// 	}
// }

// func RequireRepoWriter() macaron.Handler {
// 	return func(c *Context) {
// 		if !c.IsLogged || (!c.Repo.IsWriter() && !c.User.IsAdmin) {
// 			c.NotFound()
// 			return
// 		}
// 	}
// }
