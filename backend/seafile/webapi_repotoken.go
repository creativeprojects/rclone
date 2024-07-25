package seafile

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/rclone/rclone/backend/seafile/api"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/lib/rest"
)

func (f *Fs) getRepoInfo(ctx context.Context) (*api.RepoInfo, error) {
	// API Documentation
	// https://download.seafile.com/published/web-api/v2.1/library-api-tokens.md#user-content-Get%20Repo%20Info
	opts := rest.Opts{
		Method: "GET",
		Path:   APIv21RepoToken + "repo-info/",
	}
	result := &api.RepoInfo{}
	var resp *http.Response
	var err error
	err = f.pacer.Call(func() (bool, error) {
		resp, err = f.srv.CallJSON(ctx, &opts, nil, &result)
		return f.shouldRetry(ctx, resp, err)
	})
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 401 || resp.StatusCode == 403 {
				return nil, fs.ErrorPermissionDenied
			}
			if resp.StatusCode == 404 {
				return nil, fs.ErrorDirNotFound
			}
		}
		return nil, fmt.Errorf("failed to get library information: %w", err)
	}
	result.RepoName = f.opt.Enc.ToStandardName(result.RepoName)
	return result, nil
}

func (f *Fs) getDirectoryEntriesWithRepoToken(ctx context.Context, dirPath string, recursive bool) ([]api.DirEntryRepoToken, error) {
	// API Documentation
	// https://download.seafile.com/published/web-api/v2.1/library-api-tokens.md#user-content-List%20Items%20in%20Directory
	dirPath = path.Join("/", dirPath)

	recursiveFlag := "0"
	if recursive {
		recursiveFlag = "1"
	}
	opts := rest.Opts{
		Method: "GET",
		Path:   APIv21RepoToken + "/dir/",
		Parameters: url.Values{
			"recursive": {recursiveFlag},
			"path":      {f.opt.Enc.FromStandardPath(dirPath)},
		},
	}
	result := &api.DirEntriesRepoToken{}
	var resp *http.Response
	var err error
	err = f.pacer.Call(func() (bool, error) {
		resp, err = f.srv.CallJSON(ctx, &opts, nil, &result)
		return f.shouldRetry(ctx, resp, err)
	})
	if err != nil {
		if resp != nil {
			if resp.StatusCode == 401 || resp.StatusCode == 403 {
				return nil, fs.ErrorPermissionDenied
			}
			if resp.StatusCode == 404 {
				return nil, fs.ErrorDirNotFound
			}
			if resp.StatusCode == 440 {
				// Encrypted library and password not provided
				return nil, fs.ErrorPermissionDenied
			}
		}
		return nil, fmt.Errorf("failed to get directory contents: %w", err)
	}

	// Clean up encoded names
	for index, fileInfo := range result.Entries {
		fileInfo.Name = f.opt.Enc.ToStandardName(fileInfo.Name)
		fileInfo.Path = f.opt.Enc.ToStandardPath(fileInfo.Path)
		result.Entries[index] = fileInfo
	}
	return result.Entries, nil
}
