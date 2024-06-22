package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suyog1pathak/services/migration"
	"github.com/suyog1pathak/services/pkg/config"
	"github.com/suyog1pathak/services/pkg/datastore"
	S "github.com/suyog1pathak/services/pkg/server"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

var mysqlContainer *mysql.MySQLContainer
var s *gin.Engine
var err error

func performInfraSetup() context.Context {
	log.Default().Printf("Setting up test infra..")
	ctx := context.Background()
	mysqlContainer, err = mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0.32"),
		mysql.WithDatabase("TestDb"),
		mysql.WithPassword("TestDbPassword"),
		mysql.WithUsername("TestDbUser"),
	)

	if err != nil {
		log.Fatalf("Failed to run mysql cluster: %v", err)
	}

	err = mysqlContainer.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start mysql cluster: %v", err)
	}

	portWithType, _ := mysqlContainer.MappedPort(ctx, "3306")
	_, _ = mysqlContainer.PortEndpoint(ctx, nat.Port("3306/tcp"), "mysql_port_name")
	strPort := strings.Replace(string(portWithType), "/tcp", "", -1)

	hostIp, err := mysqlContainer.ContainerIP(ctx)
	if err != nil {
		log.Fatalf("error getting host IP of mysql container %v", err)
	}

	fmt.Println("DEBUG------>", hostIp, " ", strPort)
	os.Setenv("APP_DB_HOST", hostIp)
	config.GetConfig()
	config.Data.Db.Host = "0.0.0.0"
	config.Data.Db.Port, _ = strconv.Atoi(strPort)
	config.Data.Db.User = "TestDbUser"
	config.Data.Db.Password = "TestDbPassword"
	config.Data.Db.Name = "TestDb"

	return ctx
}

func teardownInfraSetup(ctx context.Context) {
	log.Default().Printf("Tearing down test infra..")
	if err := mysqlContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err.Error())
	}
}

func RunMigrations() {
	log.Println("running migrations.")
	db, err := datastore.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to create db connection: %v", err)
	}
	migration.RunMigrations(db)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	ctx := performInfraSetup()
	RunMigrations()
	// execute tests
	exitCode := m.Run()
	//
	teardownInfraSetup(ctx)
	os.Exit(exitCode)
}

var TestUrlPrefix = "/api/v1"

func makeRequest(method, url, body string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, TestUrlPrefix+url, bytes.NewBuffer([]byte(body)))

	response := httptest.NewRecorder()
	S.InitRouter().ServeHTTP(response, request)
	return response
}

func getRespBodyBytes(writer *httptest.ResponseRecorder) []byte {
	resp := writer.Result()
	body, _ := io.ReadAll(resp.Body)
	return body
}

// start test cases from here.
func TestShouldCheckGetAllServiceResponse(t *testing.T) {
	response := makeRequest("GET", "/services", "")
	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
}
