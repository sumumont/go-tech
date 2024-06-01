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
package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-tech/common"
	"go-tech/internal/configs"
	"go-tech/internal/dao/model"
	"go-tech/internal/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

var (
	db  *gorm.DB
	dsn string
)

func InitDb(ctx context.Context) error {
	dbConf := configs.Conf().DbConfig

	var err error

	//if pg_passwd, exists := os.LookupEnv("POSTGRES_PASSWORD"); exists && !strings.HasPrefix(pg_passwd, "vault:") {
	//	dbConf.Password = pg_passwd
	//	logging.Debug().Str("POSTGRES_PASSWORD", pg_passwd).Send()
	//}
	switch dbConf.ServerType {
	case "postgres", "postgresql":
		// create database if not exists
		preDsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Sslmode)
		logging.Debug(ctx).Str("preDsn", strings.ReplaceAll(preDsn, fmt.Sprintf("password=%s", dbConf.Password), "")).Send()
		db, err = gorm.Open(postgres.Open(preDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		exit := 0
		res1 := db.Table("pg_database").Select("count(1)").Where("datname = ?", dbConf.Database).Scan(&exit)
		if res1.Error != nil {
			return res1.Error
		}

		if exit == 0 {
			logging.Info(ctx).Msgf("Trying to create database: %s", dbConf.Database)
			res2 := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConf.Database))
			if res2.Error != nil {
				return res2.Error
			}
		}

		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Database, dbConf.Sslmode)
		logging.Debug(ctx).Str("dsn", strings.ReplaceAll(dsn, fmt.Sprintf("password=%s", dbConf.Password), "")).Send()

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			//Logger: newLogger,
		})
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported database type")
	}
	SetConnPool(ctx, dbConf)
	if dbConf.Debug {
		db = db.Debug()
	}
	return initTables(ctx)
}

func initTables(ctx context.Context) error {

	modelTypes := []interface{}{
		&model.Book{},
	}

	for _, modelType := range modelTypes {
		err := autoMigrateTable(ctx, modelType)
		if err != nil {
			return err
		}
	}
	return nil
}

func autoMigrateTable(ctx context.Context, modelType interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(modelType))
	modelName := val.Type().Name()

	logging.Info(ctx).Msgf("Migrating Table of %s ...", modelName)

	err := db.AutoMigrate(modelType)
	if err != nil {
		return err
	}
	return nil
}

func SetConnPool(ctx context.Context, dbConf configs.DbConfig) {
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	if dbConf.MaxIdleConns > 0 {
		sqlDb.SetMaxIdleConns(dbConf.MaxIdleConns)
	}
	if dbConf.MaxOpenConns > 0 {
		sqlDb.SetMaxOpenConns(dbConf.MaxOpenConns)
	}
	if dbConf.ConnMaxLifetime > 0 {
		sqlDb.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifetime) * time.Second)

	}
	if dbConf.ConnMaxIdleTime > 0 {
		sqlDb.SetConnMaxIdleTime(time.Duration(dbConf.ConnMaxIdleTime) * time.Second)
	}
	data, _ := json.Marshal(sqlDb.Stats())
	logging.Info(ctx).Str("db.stats", string(data)).Send()
	go func() {
		var count = 0
		for {
			time.Sleep(time.Duration(common.Max(10, dbConf.DbPing)) * time.Second)
			if _, err := sqlDb.Exec("create table if not exists t1(c1 int)"); err != nil {
				logging.Error(ctx, err).Msg("create test table error")
				sqlDb.Close()
				if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
					logging.Error(ctx, err).Msg("re-connect database error")
				}
			}
			sqlDb, _ = db.DB()
			err := sqlDb.Ping()
			if err != nil {
				logging.Error(ctx, err).Msg("db ping error")
			}
			count++
			if count%10 == 0 { //每10个周期打印一次心跳
				logging.Debug(ctx).Int("count", count).Msg("db ping!")
			}
			if count > 10000 {
				count = 0
			}
		}
	}()
}
