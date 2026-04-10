package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const githubReleaseAPI = "https://api.github.com/repos/ShinoharaHaruna/cc-provider/releases/latest"

type githubRelease struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cc-provider to the latest release.",
	Long:  `Fetch the latest release from GitHub and replace the current binary.`,
	Run:   runUpgradeCmd,
}

func runUpgradeCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Checking for the latest release...")

	release, err := fetchLatestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release info: %v\n", err)
		os.Exit(1)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(Version, "v")

	if current == latest {
		fmt.Printf("Already on the latest version (%s).\n", Version)
		return
	}

	fmt.Printf("New version available: %s -> %s\n", Version, release.TagName)

	assetName := fmt.Sprintf("cc-provider-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	checksumName := assetName + ".sha256"

	var downloadURL, checksumURL string
	for _, a := range release.Assets {
		switch a.Name {
		case assetName:
			downloadURL = a.BrowserDownloadURL
		case checksumName:
			checksumURL = a.BrowserDownloadURL
		}
	}

	if downloadURL == "" {
		fmt.Fprintf(os.Stderr, "No asset found for %s/%s (looked for: %s)\n",
			runtime.GOOS, runtime.GOARCH, assetName)
		os.Exit(1)
	}

	// Get expected checksum
	var expectedSum string
	if checksumURL != "" {
		expectedSum, err = fetchChecksum(checksumURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch checksum: %v\n", err)
		}
	}

	// Determine the path of the current binary
	selfPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error locating current binary: %v\n", err)
		os.Exit(1)
	}
	selfPath, err = filepath.EvalSymlinks(selfPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving symlinks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloading %s...\n", assetName)
	newBinary, err := downloadAndExtract(downloadURL, expectedSum)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading update: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(newBinary)

	// Atomic replace: write to a temp file beside the target, then rename
	dir := filepath.Dir(selfPath)
	tmp, err := os.CreateTemp(dir, ".cc-provider-upgrade-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		os.Exit(1)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath) // cleaned up if rename fails

	src, err := os.Open(newBinary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening downloaded binary: %v\n", err)
		os.Exit(1)
	}
	if _, err = io.Copy(tmp, src); err != nil {
		src.Close()
		tmp.Close()
		fmt.Fprintf(os.Stderr, "Error writing new binary: %v\n", err)
		os.Exit(1)
	}
	src.Close()
	tmp.Close()

	if err = os.Chmod(tmpPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting permissions: %v\n", err)
		os.Exit(1)
	}

	if err = os.Rename(tmpPath, selfPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error replacing binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully upgraded to %s.\n", release.TagName)
	fmt.Println("Run 'cc-provider setup' if shell integration needs refreshing.")
}

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest(http.MethodGet, githubReleaseAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "cc-provider/"+Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func fetchChecksum(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// File format: "<hash>  filename" or just "<hash>"
	return strings.Fields(string(body))[0], nil
}

// downloadAndExtract downloads the tar.gz, verifies checksum, extracts the binary,
// and returns the path to a temp file containing the extracted binary.
func downloadAndExtract(url, expectedSum string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Buffer to temp file while computing hash
	tmp, err := os.CreateTemp("", "cc-provider-dl-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	h := sha256.New()
	if _, err = io.Copy(io.MultiWriter(tmp, h), resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	if expectedSum != "" {
		actual := hex.EncodeToString(h.Sum(nil))
		if !strings.EqualFold(actual, expectedSum) {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSum, actual)
		}
	}

	// Rewind and extract
	if _, err = tmp.Seek(0, io.SeekStart); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	gr, err := gzip.NewReader(tmp)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(tmp.Name())
			return "", err
		}

		// Extract only the binary (any file named "cc-provider" at any depth)
		if filepath.Base(hdr.Name) == "cc-provider" && hdr.Typeflag == tar.TypeReg {
			out, err := os.CreateTemp("", "cc-provider-new-*")
			if err != nil {
				os.Remove(tmp.Name())
				return "", err
			}
			if _, err = io.Copy(out, tr); err != nil { //nolint:gosec
				out.Close()
				os.Remove(out.Name())
				os.Remove(tmp.Name())
				return "", err
			}
			out.Close()
			os.Remove(tmp.Name())
			return out.Name(), nil
		}
	}

	os.Remove(tmp.Name())
	return "", fmt.Errorf("binary 'cc-provider' not found in archive")
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
