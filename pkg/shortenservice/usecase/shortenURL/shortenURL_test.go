package shortenURL_test

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/mocks"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/usecase/shortenURL"
	"reflect"
	"testing"
)

func TestShortenURLUseCase_Find(t *testing.T) {
	redirectRepoStub := mocks.NewRedirectReposiotyStub()
	shortenURLUsecase := shortenURL.NewShortenURLUseCase(redirectRepoStub)
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Redirect
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Found URL from code",
			args{"hXH2eyrZg"},
			&model.Redirect{
				Code:     "hXH2eyrZg",
				URL:      "https://github.com",
				CreateAt: 1,
			},
			false,
		},
		{
			"Not URL from code",
			args{"notfound"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shortenURLUsecase.Find(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenURLUseCase.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortenURLUseCase.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenURLUseCase_Store(t *testing.T) {
	redirectRepoStub := mocks.NewRedirectReposiotyStub()
	shortenURLUsecase := shortenURL.NewShortenURLUseCase(redirectRepoStub)
	redirect := &model.Redirect{
		URL: "https://google.com",
	}

	invalidRedirect := &model.Redirect{
		URL: "https/google.com",
	}

	fmt.Printf("%v\n", redirect)

	type args struct {
		redirect *model.Redirect
	}
	type wants struct {
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		want wants
	}{
		// TODO: Add test cases.
		{
			"Store successfully",
			args{redirect},
			wants{
				wantErr: false,
			},
		},
		{
			"Store unsuccessfully, invalid URL",
			args{invalidRedirect},
			wants{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := shortenURLUsecase.Store(tt.args.redirect); (err != nil) != tt.want.wantErr {
				t.Errorf("ShortenURLUseCase.Store() error = %v, wantErr %v", err, tt.want.wantErr)
			}

			fmt.Printf("%v\n", redirect)
		})
	}
}
