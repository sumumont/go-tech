/* ******************************************************************************
* 2019 - present Contributed by Apulis Technology (Shenzhen) Co. LTD
*
* This program and the accompanying materials are made available under the
* terms of the MIT License, which is available at
* https://www.opensource.org/licenses/MIT
*
* See the NOTICE file distributed with this work for additional
* information regarding copyright ownership.
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
* WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
* License for the specific language governing permissions and limitations
* under the License.
*
* SPDX-License-Identifier: MIT
******************************************************************************/
package configs

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-tech/internal/logging"
)

var appConfig *AppConfig

type AppConfig struct {
	Debug      bool
	HTTPServer HTTPServer
	DbConfig   DbConfig
	HttpClient HttpClient
}
type HTTPServer struct {
	Network string
	Addr    string
}

type DbConfig struct {
	ServerType      string
	Username        string
	Password        string
	Host            string
	Port            int
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	Debug           bool
	Sslmode         string
	DbPing          int
}

type HttpClient struct {
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	TimeoutSeconds      int
}

func InitConfig(ctx context.Context) error {
	logging.Info(ctx).Msg("reading config")
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("configs")

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	appConfig = &AppConfig{}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		logging.Info(ctx).Str("fs.event", in.Name).Msg("config.yaml changed")
		err := v.Unmarshal(&appConfig)
		if err != nil {
			logging.Error(ctx, err).Send()
		}
	})
	err = v.Unmarshal(&appConfig)
	if err != nil {
		return err
	}
	return nil
}

func Conf() *AppConfig {
	return appConfig
}
