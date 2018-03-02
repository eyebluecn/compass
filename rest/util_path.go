package rest

import (
	"os"
	"path/filepath"
)

//判断文件或文件夹是否已经存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取该应用可执行文件的位置。
//例如：C:\Users\lishuang\AppData\Local\Temp
func GetHomePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

//获取前端静态资源的位置。如果你在开发模式下，可以将这里直接返回compass/build下面的html路径。
//例如：C:/Users/lishuang/AppData/Local/Temp/html
func GetHtmlPath() string {

	homePath := GetHomePath()
	filePath := homePath + "/html"
	exists, err := PathExists(filePath)
	if err != nil {
		panic("判断上传文件是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(filePath, 0777)
		if err != nil {
			panic("创建上传文件夹时出错！")
		}
	}

	return filePath
}

//如果文件夹存在就不管，不存在就创建。 例如：/var/www/matter
func MakeDirAll(dirPath string) string {

	exists, err := PathExists(dirPath)
	if err != nil {
		panic("判断文件是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(dirPath, 0666)
		if err != nil {
			panic("创建文件夹时出错！")
		}
	}

	return dirPath
}


//获取配置文件存放的位置
//例如：C:\Users\lishuang\AppData\Local\Temp/conf
func GetConfPath() string {

	homePath := GetHomePath()
	filePath := homePath + "/conf"
	exists, err := PathExists(filePath)
	if err != nil {
		panic("判断日志文件夹是否存在时出错！")
	}
	if !exists {
		err = os.MkdirAll(filePath, 0666)
		if err != nil {
			panic("创建日志文件夹时出错！")
		}
	}

	return filePath
}

