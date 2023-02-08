package S3utils

import (
	"context"
	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	csvutils "github.com/alessiosavi/GoGPUtils/csv"
	fileutils "github.com/alessiosavi/GoGPUtils/files"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// FIXME: create test method that create the S3 bucket for test the methods

func TestListBucketObject(t *testing.T) {
	type args struct {
		bucket string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket: "aws-website-tbimg-pdds3-backup",
			},
			want:    []string{"bbl-positional-conf-margins.csv", "bbl-positional-conf-moq.csv", "size_distribution.py"},
			wantErr: false,
		},
		{
			name: "ko",
			args: args{
				bucket: "this-bucket-does-not-exists",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListBucketObjects(tt.args.bucket, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ListBucketObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListBucketObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetObject(t *testing.T) {
	type args struct {
		bucket   string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{

		{
			name: "ok",
			args: args{
				bucket:   "my-bucket-test-s3",
				fileName: "empty.txt",
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetObject(tt.args.bucket, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyObject(t *testing.T) {
	type args struct {
		bucket       string
		bucketTarget string
		key          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket:       "my-bucket-test-s3",
				bucketTarget: "my-bucket-test-s3-copy",
				key:          "covid.csv",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyObject(tt.args.bucket, tt.args.bucketTarget, tt.args.key, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("CopyObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectExists(t *testing.T) {
	type args struct {
		bucket string
		key    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				bucket: "aws-glue-scripts-796325849317-eu-west-1",
				key:    "FabricaLabAdmin/prova",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ko1",
			args: args{
				bucket: "my-bucket-test-s3",
				key:    "this-file-does-not-exists",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "ko2",
			args: args{
				bucket: "this-bucket-does-not-exists",
				key:    "this-file-does-not-exists",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ObjectExists(tt.args.bucket, tt.args.key)
			if got != tt.want {
				t.Errorf("ObjectExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHead(t *testing.T) {
	if head, err := S3Client.HeadObject(context.Background(), &s3.HeadObjectInput{Bucket: aws.String("prod-lambda-asset"), Key: aws.String("describe-jobs.zip")}); err != nil {
		return
	} else {
		t.Log(helper.MarshalIndent(head))
	}
}

func TestIsDifferent2(t *testing.T) {
	diff := IsDifferent("qa-lambda-asset", "prod-lambda-asset", "go-am-parser.zip", "test.zip")
	if diff != true {
		t.Error()
	}
}

func TestSyncBucket(t *testing.T) {
	bucket, err := SyncBucket("prod-demand-planning-forecast-temp", "", "qa-demand-planning-forecast-temp")
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(bucket))
}

func TestIsDifferent(t *testing.T) {
	object, err := ListBucketObjects("qa-lambda-asset", "")
	if err != nil {
		t.Error(err)
	}
	object = arrayutils.RemoveStrings(object, []string{"orchestrator/", "orchestrator/entrypoint.zip"})
	log.Println(object)
	now := time.Now()
	for _, lambda := range object {
		IsDifferent("qa-lambda-asset", "prod-lambda-asset", lambda, lambda)
	}
	log.Printf("Executed time: %s\n", time.Since(now))

	now = time.Now()
	for _, lambda := range object {
		IsDifferentLegacy("qa-lambda-asset", "prod-lambda-asset", lambda, lambda)
	}
	log.Printf("Executed time: %s\n", time.Since(now))
}

func BenchmarkIsDifferent(b *testing.B) {
	object, err := ListBucketObjects("qa-lambda-asset", "")
	if err != nil {
		b.Error(err)
	}
	object = arrayutils.RemoveStrings(object, []string{"orchestrator/", "orchestrator/entrypoint.zip"})
	for i := 0; i < b.N; i++ {
		for _, lambda := range object {
			IsDifferent("qa-lambda-asset", "prod-lambda-asset", lambda, lambda)
		}
	}
}

func BenchmarkIsDifferentLegacy(b *testing.B) {
	object, err := ListBucketObjects("qa-lambda-asset", "")
	if err != nil {
		b.Error(err)
	}
	object = arrayutils.RemoveStrings(object, []string{"orchestrator/", "orchestrator/entrypoint.zip"})
	for i := 0; i < b.N; i++ {
		for _, lambda := range object {
			IsDifferentLegacy("qa-lambda-asset", "prod-lambda-asset", lambda, lambda)
		}
	}
}

func TestGetAfterDate(t *testing.T) {
	dates, err := GetAfterDate("prod-data-lake-bucket", "input/CENTRIC/upload/", time.Date(2021, 12, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	}
	log.Println(helper.MarshalIndent(dates))

	for _, date := range dates {
		object, err := GetObject("prod-data-lake-bucket", date)
		if err != nil {
			panic(err)
		}
		if err = os.WriteFile("/tmp/centric/data/Style/"+path.Base(date), object, 0600); err != nil {
			panic(err)
		}
	}
}

func TestList(t *testing.T) {

	details, err := ListBucketObjectsDetails("bbl-prod-input-for-data-lake", "")
	if err != nil {
		return
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].LastModified.After(*details[j].LastModified)

	})
	os.WriteFile("/tmp/list.json", []byte(helper.MarshalIndent(details)), 0600)
}

func TestSyncAfterDate2(t *testing.T) {
	type args struct {
		bucket    string
		prefix    string
		localPath string
		date      time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				bucket:    "prod-data-lake-bucket",
				prefix:    "input/CEGID/ftp-base/history/",
				localPath: "/tmp/CEGID/",
				date:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncAfterDate(tt.args.bucket, tt.args.prefix, tt.args.localPath, tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("SyncBeforeDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncAfterDate(t *testing.T) {
	type args struct {
		bucket    string
		prefix    string
		localPath string
		date      time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				bucket:    "prod-data-lake-bucket",
				prefix:    "input/WAC/",
				localPath: "/tmp/WAC/CEGID",
				date:      time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "OK",
			args: args{
				bucket:    "prod-data-lake-bucket",
				prefix:    "input/SAP/upload/WAC/",
				localPath: "/tmp/WAC/SAP",
				date:      time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncAfterDate(tt.args.bucket, tt.args.prefix, tt.args.localPath, tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("SyncBeforeDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBetweenDate(t *testing.T) {
	type args struct {
		bucket string
		prefix string
		start  time.Time
		stop   time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				bucket: "prod-data-lake-bucket",
				prefix: "input/WAC",
				start:  time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
				stop:   time.Date(2022, 9, 10, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "OK",
			args: args{
				bucket: "prod-data-lake-bucket",
				prefix: "input/SAP/upload/WAC/",
				start:  time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
				stop:   time.Date(2022, 9, 10, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBetweenDate(tt.args.bucket, tt.args.prefix, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBetweenDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, object := range got {
				getObject, err := GetObject(tt.args.bucket, object)
				if err != nil {
					panic(err)
				}
				p := path.Join("/tmp/wac", path.Dir(object))
				if !fileutils.IsDir(p) {
					if err := os.MkdirAll(p, 0755); err != nil {
						return
					}
				}
				if err = os.WriteFile(path.Join(p, path.Base(object)), getObject, 0755); err != nil {
					panic(err)
				}
			}
		})
	}
}

func TestImageMorello(t *testing.T) {
	file, err := os.ReadFile("/tmp/data.csv")
	if err != nil {
		panic(err)
	}
	_, csv, err := csvutils.ReadCSV(file, ',')
	if err != nil {
		panic(err)
	}
	var imagesTarget []string

	for i := range csv {
		imagesTarget = append(imagesTarget, csv[i][0])
	}
	log.Println(imagesTarget)

	objects, err := ListBucketObjects("thom-browne-images", "TB_Master_Images")
	if err != nil {
		panic(err)
	}

	bar := progressbar.Default(int64(len(imagesTarget)))
	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)
	for i := range objects {
		if arrayutils.ContainStrings(imagesTarget, objects[i]) {
			bar.Describe(objects[i])
			guard <- struct{}{} // would block if guard channel is already filled
			go func() {
				object, err := GetObject("thom-browne-images", objects[i])
				if err != nil {
					panic(err)
				}
				if err = os.MkdirAll(path.Join("/tmp/images", path.Dir(objects[i])), 0755); err != nil {
					panic(err)
				}
				if err = os.WriteFile(path.Join("/tmp/images", objects[i]), object, 0755); err != nil {
					panic(err)
				}
				bar.Add(1)
				<-guard
			}()
		}
	}
	bar.Close()
}

func TestLoadMorello(t *testing.T) {
	ordered, err := fileutils.ListFilesOrdered("/tmp/images")
	if err != nil {
		panic(err)
	}
	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)
	bar := progressbar.Default(int64(len(ordered)))
	defer bar.Close()
	for i := range ordered {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(i int) {
			bar.Describe(ordered[i])
			file, err := os.ReadFile(ordered[i])
			if err != nil {
				panic(err)
			}
			bar.Describe(ordered[i])
			if err = PutObject("aws-website-tbimg-pdds3", path.Join("images/", strings.ReplaceAll(path.Base(ordered[i]), " ", "")), file); err != nil {
				panic(err)
			}
			bar.Add(1)
			<-guard
		}(i)
	}
}

func TestVito(t *testing.T) {
	objects, err := ListBucketObjects("prod-data-lake-bucket", "input/CEGID/")
	if err != nil {
		panic(err)
	}
	for _, o := range objects {
		if strings.Contains(o, "CL1_SALES_CN-01232023") {
			log.Println(o)
		}
	}
}
