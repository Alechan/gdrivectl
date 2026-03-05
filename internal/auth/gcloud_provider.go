package auth

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Alechan/gdrivectl/internal/fail"
)

type GcloudTokenProvider struct {
	bin string
}

func NewGcloudTokenProvider(bin string) *GcloudTokenProvider {
	return &GcloudTokenProvider{bin: bin}
}

func (p *GcloudTokenProvider) AccessToken(ctx context.Context) (string, error) {
	if strings.TrimSpace(p.bin) == "" {
		return "", fail.NewConfig("gcloud binary path is empty", "set --gcloud-bin or GDRIVECTL_GCLOUD_BIN")
	}
	cmd := exec.CommandContext(ctx, p.bin, "auth", "print-access-token")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s := strings.TrimSpace(string(out))
		if strings.Contains(strings.ToLower(s), "scope") || strings.Contains(strings.ToLower(s), "insufficient") {
			return "", fail.NewScope("insufficient auth scope", "run: gcloud auth login --enable-gdrive-access --update-adc")
		}
		if strings.Contains(strings.ToLower(s), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "executable file not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no such file or directory") {
			return "", fail.NewConfig("gcloud binary not found", "set --gcloud-bin to your gcloud executable")
		}
		return "", fail.NewAuth("unable to get access token", "run: gcloud auth login --enable-gdrive-access --update-adc")
	}
	token := strings.TrimSpace(string(out))
	if token == "" {
		return "", fail.NewAuth("received empty access token", "run: gcloud auth login --enable-gdrive-access --update-adc")
	}
	return token, nil
}
