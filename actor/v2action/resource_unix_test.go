// +build !windows

package v2action_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "code.cloudfoundry.org/cli/actor/v2action"
	"code.cloudfoundry.org/cli/actor/v2action/v2actionfakes"
	"code.cloudfoundry.org/ykk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource Actions", func() {
	var (
		actor                     Actor
		fakeCloudControllerClient *v2actionfakes.FakeCloudControllerClient
	)

	BeforeEach(func() {
		fakeCloudControllerClient = new(v2actionfakes.FakeCloudControllerClient)
		actor = NewActor(fakeCloudControllerClient, nil)
	})

	Describe("ZipResources", func() {
		var (
			srcDir string

			resultZip  string
			resources  []Resource
			executeErr error
		)

		BeforeEach(func() {
			var err error
			srcDir, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())

			subDir := filepath.Join(srcDir, "level1", "level2")
			err = os.MkdirAll(subDir, 0777)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(filepath.Join(subDir, "tmpFile1"), []byte("why hello"), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(filepath.Join(srcDir, "tmpFile2"), []byte("Hello, Binky"), 0751)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(filepath.Join(srcDir, "tmpFile3"), []byte("Bananarama"), 0655)
			Expect(err).ToNot(HaveOccurred())

			resources = []Resource{
				{Filename: "level1"},
				{Filename: "level1/level2"},
				{Filename: "level1/level2/tmpFile1"},
				{Filename: "tmpFile2"},
				{Filename: "tmpFile3"},
			}
		})

		JustBeforeEach(func() {
			resultZip, executeErr = actor.ZipResources(srcDir, resources)
		})

		AfterEach(func() {
			err := os.RemoveAll(srcDir)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when zipping on UNIX", func() {
			It("zips the directory and keeps the file permissions", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(resultZip).ToNot(BeEmpty())
				zipFile, err := os.Open(resultZip)
				Expect(err).ToNot(HaveOccurred())
				defer zipFile.Close()

				zipInfo, err := zipFile.Stat()
				Expect(err).ToNot(HaveOccurred())

				reader, err := ykk.NewReader(zipFile, zipInfo.Size())
				Expect(err).ToNot(HaveOccurred())

				Expect(reader.File).To(HaveLen(5))
				Expect(reader.File[0].Name).To(Equal("level1"))
				Expect(reader.File[1].Name).To(Equal("level1/level2"))
				Expect(reader.File[2].Name).To(Equal("level1/level2/tmpFile1"))
				Expect(reader.File[3].Name).To(Equal("tmpFile2"))
				Expect(reader.File[4].Name).To(Equal("tmpFile3"))

				Expect(reader.File[2].Mode()).To(Equal(os.FileMode(0644)))
				Expect(reader.File[3].Mode()).To(Equal(os.FileMode(0751)))
				Expect(reader.File[4].Mode()).To(Equal(os.FileMode(0655)))
			})
		})
	})
})
