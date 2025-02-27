package app

import (
	"context"
	"testing"

	"github.com/rovany706/url-shortener/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetFullURL(t *testing.T) {
	tests := []struct {
		name          string
		existingLinks map[string]string
		shortID       string
		wantFullURL   string
		wantOk        bool
	}{
		{
			name: "shortID exists",
			existingLinks: map[string]string{
				"id1": "http://example.com/",
			},
			shortID:     "id1",
			wantFullURL: "http://example.com/",
			wantOk:      true,
		},
		{
			name: "shortID don't exist",
			existingLinks: map[string]string{
				"id1": "http://example.com/",
			},
			shortID:     "testid",
			wantFullURL: "",
			wantOk:      false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repository := mock.NewMockRepository(ctrl)
			repository.EXPECT().GetFullURL(gomock.Any(), tt.shortID).Return(tt.wantFullURL, tt.wantOk).AnyTimes()

			app := NewURLShortenerApp(repository)

			fullLink, ok := app.GetFullURL(ctx, tt.shortID)

			if tt.wantOk {
				assert.True(t, ok)
				assert.Equal(t, tt.wantFullURL, fullLink)
				return
			}

			assert.False(t, ok)
		})
	}
}

func TestGetShortID(t *testing.T) {
	tests := []struct {
		name        string
		fullURL     string
		wantShortID string
		wantErr     bool
	}{
		{
			name:        "valid url",
			fullURL:     "http://example.com/123",
			wantShortID: "488575e6",
			wantErr:     false,
		},
		{
			name:        "invalid url",
			fullURL:     "http,,:example.com",
			wantShortID: "",
			wantErr:     true,
		},
		{
			name:        "long url",
			fullURL:     "http://www.reallylong.link/rll/PdB6vKJe9sq9MxnLiUHZFFxPCdzt9O_ntvESFf0KbV88N6Embw9LfhRReYAs5lyTCL2JDDBQl2o4k0v374xQXyu7NGJ3_PklOemPuIy00IHZ8LVU4oSGdDc0gWa5vkI3kWS0SdG9yvXlxFQGXkl1WvIqFpTsNHpsuVu0IR0lnr33OPN3FHA0DDBAu8NjYezil5QKFAmQj42LRrAwt1P2IHT2wYPibP3xEYPGDZwgOmnHHxLaxWZyqWqNK6XaYDkonoPurIAUdWkAa5/az78mtNMOweGvefwksumLGc1/dGPXapKdfn/yL7Il_FpzQ7F3NhZoZiUB6sqgQIxoFHnV7eZX7BFpw5_pfSClYb6KV6rXPfTj1V2wy2QY4ACK94iniPq64PDEOlqqPiZF5LXKNxkt2_tESUHGsTLn6v5Hy1PP73hlqW1J1N7M62YtwbaCmbcyFbgxkpd80ZfrtPQlug16F190FUEDQC4JiYJevy9UhR1XIZ4rCPiIxaBuG/IyaBUhozzg_KfA2zYfL03r44sAL27R8T1EAh27mpVa2RRWSpVDdDhDoFiaYYEt6WHVNL6cSxFa8UP1eSfmENyTev1O_6Mo7iKWUZW0Js__N_FimlFI3PqA18CQZQBjpstupLPYRA0XxaeeBA2CNSelRHylvy_ttdRfsd8x3cRoN1kkcbwBFSFq8EqA4CxgKvrVCcnQAYXQPhfFIPBPO1K7FbQ7RaU7RfsbgdXhfg1hhiUEIM1qzJ992MsA0v9H4IIlRTC3CspA0y08oFefus7Ik1IBo3wzqS8ta3CpYpbgr98kxsmPKoPqb/mRFopiWDufm3KYzLzHEUDwBeGG8_191cehhiCc8QMwo11bbiEQBUlOm9JpSdEWyUtwpi6qHmMmuE50IqKO6twlz11UHo7/ToUIfSrfszJyjkkiWGi9D/Tbw1j16n1qCy3yfGcBPstWICWW1hT0lG7OjvPRIxyB7FtFInDjdeMvySTtyEITz6J660rZVsYrPrLrPZzx7ETdX5EkeV5Gc8ssUr8W2sxPHT5btpZmB/k3C4sjt1pzmDVS2ISzpYopu1Gs_d4hCQUXqJFzrchOx6YCLQoSxqtpr9M1NVwCApDzNCvCFE/dpj22Mx6AjE2oHkX6tqEbOM615WKck2QJncDY7hFPUHeBvb1Vu1_dXmufVkIM2ph6UC9k42LrTo59aftx3PvKq6FRG8gOy/9WOWkdaUFx8ZeMe6xGF8i6Pxq2rBMUASlXzd29fxCPtFvlRKnP29JuAWF6lGWA_rEaU7VQ2lMy2zXkohScTmcQGD4XYeeIxEj1B9WS7SC1Vif_pvggW6JJtJCDDlpv7lGSR6uv33z1cmzJI5DqUk9A2knGnvxphkTT2GuPgzblXrQUM687DzKMEkFIqGZfeYaoFu5SpxTtdFYhQj5krIMygOaNQUIIItaEz1l7CDU1Q37as02M9PU2VRwHUSDkkQ7H6YEripVFTtwh8f2KjyRAGKnUwOwHZX2aG1ohI7KZQBrL2xEpeB5HJ52P7q8WVl6tW6T3Xn63OWOSn3em6jwUsEqC9FepDIOZOAeX8VyBhlRI4ZRFyacfO_GppiRgt1WY9N31ap79nPr78jYuTxlabUQf4AX4G1VAaTtWXYBSZYF428x6FgOiGo0Ko8wAUIr2S5Z1Ax3s4LCSKFLm5/pYYA7otqjgvSMwGU6TIqZMfTX3GNCs5clBDNzkBZGT3_Ap4Ps1itN3A10hqMGY9abCmRGtK570Arw2xc4duWnb2p5vavQRGVv0r5ZV/PAARq5BAAxCwwcHdCmmCZFImc5XC3M7TSU2oseG8G3bR2kHLbhKw/IXIDjxP/KLY_6icL_gGyMKzy1C/HZmmfUJK983EoHuMANHIuhaRFStAgVPA82cXvEx2ANYuuzFRrCJChRbWnZUZHXnun2CWclM_yEeenqgU89CGzDSFk7ikrt",
			wantShortID: "4c95793e",
			wantErr:     false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repository := mock.NewMockRepository(ctrl)
			repository.EXPECT().SaveEntry(gomock.Any(), gomock.Any(), gomock.Any(), tt.fullURL).Return(nil).AnyTimes()

			app := NewURLShortenerApp(repository)
			shortID, err := app.GetShortID(ctx, 1, tt.fullURL)

			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.wantShortID, shortID)
				return
			}

			assert.Error(t, err)
		})
	}
}
