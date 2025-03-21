package db

import (
	"fmt"
	"log"
	"net"
	"io"
	"os"
	"path/filepath"

	"github.com/hirochachacha/go-smb2"
	infisical "github.com/infisical/go-sdk"
	"github.com/99designs/gqlgen/graphql"

	"backoffice/helpers"
)

// SMBConfig holds the configuration for connecting to the SMB server
type SMBConfig struct {
	Server   string // e.g., "192.168.1.100"
	Share    string // e.g., "shared-folder"
	User     string
	Password string
}

// SMBClient handles the SMB connection and operations
type SMBClient struct {
	conn net.Conn
	sess *smb2.Session
	fs   *smb2.Share
}

func SMBConfigure(secrets []infisical.Secret) (*SMBConfig, error) {
	config := &SMBConfig{
		Share: "cdn", // Default value
	}
	
	secretKeys := map[string]*string {
		"OMV_HOST":     &config.Server,
		"OMV_CDN_USER": &config.User,
		"OMV_CDN_PASS": &config.Password,
	}
	
	for _, secret := range secrets {
		if ptr, exists := secretKeys[secret.SecretKey]; exists {
			*ptr = secret.SecretValue
		}
	}

	// Validate required fields
	if config.Server == "" || config.User == "" || config.Password == "" {
		return nil, fmt.Errorf("Failed to retrieve OMV secrets")
	}

	log.Println("Succesfully create config for SMB connection")

	return config, nil
}

// SMBConnect establishes a new SMB connection and returns a client instance
func SMBConnect(config *SMBConfig) (*SMBClient, error) {
	client := &SMBClient{}

	// Dial the SMB server
	var err error
	client.conn, err = net.Dial("tcp", fmt.Sprintf("%s:445", config.Server))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMB server: %w", err)
	}

	// Authenticate with SMB
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.User,
			Password: config.Password,
		},
	}

	client.sess, err = d.Dial(client.conn)
	if err != nil {
		client.conn.Close() // Close connection if auth fails
		return nil, fmt.Errorf("failed to authenticate to SMB server: %w", err)
	}

	// Mount the SMB share
	client.fs, err = client.sess.Mount(config.Share)
	if err != nil {
		client.sess.Logoff()
		client.conn.Close()
		return nil, fmt.Errorf("failed to mount SMB share: %w", err)
	}

	log.Println("Connected to SMB share")

	return client, nil
}

// createRemoteDir ensures the directory exists in the SMB share
func (c *SMBClient) smbCreateRemoteDir(dirPath string) error {
	_, err := c.fs.Stat(dirPath)
	if os.IsNotExist(err) {
		return c.fs.MkdirAll(dirPath, 0755)
	}
	return err
}

// UploadFiles uploads multiple files to the SMB share
func (c *SMBClient) SMBFileUpload(files []*graphql.Upload, remoteDir string, dirPrefix string) ([]*string, error) {
	var uploadedFilePaths []*string

	// Ensure remote directory exists
	err := c.smbCreateRemoteDir(dirPrefix+remoteDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote directory: %w", err)
	}

	for _, file := range files {
		if file == nil {
			continue
		}

		// Define the remote file path
		remoteFileName := filepath.Base(file.Filename)
		remoteFilePath := filepath.Join(dirPrefix, remoteDir, remoteFileName)

		log.Println("Image being prepared for uplaod: " + remoteFilePath)

		// Check if the file already exists
		_, err := c.fs.Stat(remoteFilePath)
		if err == nil {
			// File already exists, skip uploading but add to results
			uploadedFilePaths = append(uploadedFilePaths, &remoteFilePath)
			log.Printf("File already exists, skipping upload: %s", remoteFilePath)
			continue
		} else if !os.IsNotExist(err) {
			// Other error while checking file existence
			return nil, fmt.Errorf("error checking file existence: %w", err)
		}

		// Open the uploaded file
		uploadFile := file.File
		// Check if uploadFile has a Close() method before deferring
		if closer, ok := uploadFile.(io.Closer); ok {
			defer closer.Close()
		}

		// Create file on SMB share
		remoteFile, err := c.fs.Create(remoteFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file on SMB share: %w", err)
		}
		defer remoteFile.Close()

		// Copy file data
		_, err = io.Copy(remoteFile, uploadFile)
		if err != nil {
			return nil, fmt.Errorf("failed to write file to SMB share: %w", err)
		}

		log.Println(remoteFilePath)

		// Store uploaded file path
		uploadedFilePaths = append(uploadedFilePaths, &remoteFilePath)
	}

	log.Println("Files uploaded successfully")

	return uploadedFilePaths, nil
}

// UpdateFiles uploads multiple files to the SMB share
func (c *SMBClient) SMBFileUpdate(currentFiles []*string, newFiles []*graphql.Upload, deletedFiles []*string, dir string) ([]*string, error) {
	// Ensure remote directory exists
	_, err := c.fs.Stat(dir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %w", err)
	}

	// Handle Current Files
	var uploadedFilePaths = currentFiles

	// Handle Deleted Files
	for _, pathPtr := range deletedFiles {
		if pathPtr == nil {
			continue
		}

		path := *pathPtr

		// Check if the file exists
		_, err := c.fs.Stat(path)
		if err == nil {
			log.Println(path + "does not exist, continuing")
			continue
		} else if !os.IsNotExist(err) {
			// Other error while checking file existence
			return nil, fmt.Errorf("error checking file existence: %w", err)
		}

		err = c.fs.Remove(path)
		if err != nil {
			return nil, fmt.Errorf("error removing file: %w", err)
		}

		uploadedFilePaths = helpers.RemoveByValue(uploadedFilePaths, &path)
	}

	// Handle New Files
	for _, file := range newFiles {
		if file == nil {
			continue
		}

		// Define the remote file path
		remoteFileName := filepath.Base(file.Filename)
		remoteFilePath := filepath.Join(dir, remoteFileName)

		log.Println("Image being prepared for uplaod: " + remoteFilePath)

		// Check if the file already exists
		_, err := c.fs.Stat(remoteFilePath)
		if err == nil {
			// File already exists, skip uploading but add to results
			uploadedFilePaths = append(uploadedFilePaths, &remoteFilePath)
			log.Printf("File already exists, skipping upload: %s", remoteFilePath)
			continue
		} else if !os.IsNotExist(err) {
			// Other error while checking file existence
			return nil, fmt.Errorf("error checking file existence: %w", err)
		}

		// Open the uploaded file
		uploadFile := file.File
		// Check if uploadFile has a Close() method before deferring
		if closer, ok := uploadFile.(io.Closer); ok {
			defer closer.Close()
		}

		// Create file on SMB share
		remoteFile, err := c.fs.Create(remoteFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file on SMB share: %w", err)
		}
		defer remoteFile.Close()

		// Copy file data
		_, err = io.Copy(remoteFile, uploadFile)
		if err != nil {
			return nil, fmt.Errorf("failed to write file to SMB share: %w", err)
		}

		log.Println(remoteFilePath)

		// Store uploaded file path
		uploadedFilePaths = append(uploadedFilePaths, &remoteFilePath)
	}

	log.Println("Files updated successfully")

	return uploadedFilePaths, nil
}

// SMBRemoveDirRecursive removes a directory and all its contents
func (c *SMBClient) SMBRemoveDirRecursive(dirPath string) error {
	log.Println("Deleting " + dirPath)
	err := c.fs.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the SMB connection
func (c *SMBClient) SMBClose() {
	if c.fs != nil {
		c.fs.Umount()
	}
	if c.sess != nil {
		c.sess.Logoff()
	}
	if c.conn != nil {
		c.conn.Close()
	}

	log.Println("SMB Share cleaned up")
}